package progress

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/ActiveState/cli/internal/constants"
	"github.com/ActiveState/cli/internal/errs"
	"github.com/ActiveState/cli/internal/locale"
	"github.com/ActiveState/cli/internal/logging"
	"github.com/ActiveState/cli/internal/output"
	"github.com/ActiveState/cli/pkg/platform/runtime/artifact"
	"github.com/ActiveState/cli/pkg/platform/runtime/setup/events"
	"github.com/vbauerster/mpb/v7"
	"golang.org/x/net/context"
)

type step struct {
	name     string
	verb     string
	priority int
}

func (s step) String() string {
	return s.name
}

var (
	StepBuild    = step{"build", locale.T("building"), 10000} // the priority is high because the artifact progress bars need to fit in between the steps
	StepDownload = step{"download", locale.T("downloading"), 20000}
	StepInstall  = step{"install", locale.T("installing"), 30000}
)

type artifactStepID string

type artifactStep struct {
	artifactID artifact.ArtifactID
	step       step
}

func (a artifactStep) ID() artifactStepID {
	return artifactStepID(a.artifactID.String() + a.step.String())
}

type ProgressDigester struct {
	// The max width to use for the name entries of progress bars
	maxNameWidth int

	// Progress bars and spinners
	mainProgress *mpb.Progress
	buildBar     *bar
	downloadBar  *bar
	installBar   *bar
	solveSpinner *output.Spinner
	artifactBars map[artifactStepID]*bar

	// Artifact name lookup map
	artifactNames artifact.Named

	// Track the totals required as the bars for these are only initialized for the first artifact received, at which
	// time we won't have the totals unless we previously recorded them.
	buildsExpected    map[artifact.ArtifactID]struct{}
	downloadsExpected map[artifact.ArtifactID]struct{}
	installsExpected  map[artifact.ArtifactID]struct{}

	// Debug properties used to reduce the number of log entries generated
	dbgEventLog []string

	out output.Outputer

	// We use a mutex because whilst this package itself doesn't do any threading; its consumers do.
	mutex *sync.Mutex

	// The cancel function for the mpb package
	cancelMpb context.CancelFunc

	// Record whether changes were made
	changesMade bool
	// Record whether the runtime install was successful
	success bool
}

func NewProgressIndicator(w io.Writer, out output.Outputer) *ProgressDigester {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProgressDigester{
		mainProgress: mpb.NewWithContext(
			ctx,
			mpb.WithWidth(progressBarWidth),
			mpb.WithOutput(w),
			mpb.WithRefreshRate(refreshRate),
		),

		artifactNames: map[artifact.ArtifactID]string{},
		artifactBars:  map[artifactStepID]*bar{},

		cancelMpb:    cancel,
		maxNameWidth: maxNameWidth(),
		out:          out,

		mutex: &sync.Mutex{},
	}
}

func (p *ProgressDigester) Handle(ev events.Eventer) error {
	p.dbgEventLog = append(p.dbgEventLog, fmt.Sprintf("%T", ev))

	p.mutex.Lock()
	defer p.mutex.Unlock()

	switch v := ev.(type) {

	case events.Start:
		logging.Debug("Initialize Event: %#v", v)

		// Ensure Start event is first.. because otherwise the prints below will cause output to be malformed.
		if p.buildBar != nil || p.downloadBar != nil || p.installBar != nil || p.solveSpinner != nil {
			return errs.New("Received Start event after bars were already initialized, event log: %v", p.dbgEventLog)
		}

		// Report the log file we'll be using. This has to happen here and not in the BuildStarted even as there's no
		// guarantee that no downloads or installs might have triggered before BuildStarted, in which case there's
		// already progressbars being displayed which won't play nice with newly printed output.
		if v.RequiresBuild {
			p.out.Notice(locale.Tr("progress_build_log", v.LogFilePath))
		}

		p.artifactNames = v.ArtifactNames

		p.buildsExpected = artifact.ArtifactIDsToMap(v.ArtifactsToBuild)
		p.downloadsExpected = artifact.ArtifactIDsToMap(v.ArtifactsToDownload)
		p.installsExpected = artifact.ArtifactIDsToMap(v.ArtifactsToInstall)

		if len(v.ArtifactsToBuild)+len(v.ArtifactsToDownload)+len(v.ArtifactsToInstall) == 0 {
			p.out.Notice(locale.T("progress_nothing_to_do"))
		} else {
			p.changesMade = true
		}

	case events.Success:
		p.success = true

	case events.SolveStart:
		p.solveSpinner = output.StartSpinner(p.out, locale.T("progress_solve"), refreshRate)

	case events.SolveError:
		if p.solveSpinner == nil {
			return errs.New("SolveError called before solveBar was initialized")
		}
		p.solveSpinner.Stop(locale.T("progress_fail"))
		p.solveSpinner = nil

	case events.SolveSuccess:
		if p.solveSpinner == nil {
			return errs.New("SolveSuccess called before solveBar was initialized")
		}
		p.solveSpinner.Stop(locale.T("progress_success"))
		p.solveSpinner = nil

	case events.BuildSkipped:
		if p.buildBar != nil {
			return errs.New("BuildSkipped called, but buildBar was initialized.. this should not happen as they should be mutually exclusive")
		}

	case events.BuildStarted:
		if p.buildBar != nil {
			return errs.New("BuildStarted called after buildbar was already initialized")
		}
		p.buildBar = p.addTotalBar(locale.Tl("progress_building", "Building"), v.Artifacts, mpb.BarPriority(StepBuild.priority))

	case events.BuildSuccess:
		if p.buildBar == nil {
			return errs.New("BuildSuccess called before buildbar was initialized")
		}

	case events.BuildFailure:
		if p.buildBar == nil {
			return errs.New("BuildFailure called before buildbar was initialized")
		}
		p.buildBar.Abort(false) // mpb has been known to stick around after it was told not to

	case events.ArtifactBuildStarted:
		if p.buildBar == nil {
			return errs.New("ArtifactBuildStarted called before buildbar was initialized")
		}
		if _, ok := p.buildsExpected[v.ArtifactID]; !ok {
			return errs.New("ArtifactBuildStarted called for an artifact that was not expected: %s", v.ArtifactID.String())
		}

	case events.ArtifactBuildSuccess:
		if p.buildBar == nil {
			return errs.New("ArtifactBuildSuccess called before buildbar was initialized")
		}
		if _, ok := p.buildsExpected[v.ArtifactID]; !ok {
			return errs.New("ArtifactBuildSuccess called for an artifact that was not expected: %s", v.ArtifactID.String())
		}
		if p.buildBar.Current() == p.buildBar.total {
			return errs.New("Build bar is already complete, this should not happen")
		}
		delete(p.buildsExpected, v.ArtifactID)
		p.buildBar.Increment()

	case events.ArtifactDownloadStarted:
		if p.downloadBar == nil {
			p.downloadBar = p.addTotalBar(locale.Tl("progress_building", "Downloading"), int64(len(p.downloadsExpected)), mpb.BarPriority(StepDownload.priority))
		}
		if _, ok := p.downloadsExpected[v.ArtifactID]; !ok {
			return errs.New("ArtifactDownloadStarted called for an artifact that was not expected: %s", v.ArtifactID.String())
		}

		if err := p.addArtifactBar(v.ArtifactID, StepDownload, int64(v.TotalSize), true); err != nil {
			return errs.Wrap(err, "Failed to add or update artifact bar")
		}

	case events.ArtifactDownloadProgress:
		if err := p.updateArtifactBar(v.ArtifactID, StepDownload, v.IncrementBySize); err != nil {
			return errs.Wrap(err, "Failed to add or update artifact bar")
		}

	case events.ArtifactDownloadSkipped:
		if p.downloadBar == nil {
			return errs.New("ArtifactDownloadSkipped called before downloadBar was initialized")
		}
		delete(p.downloadsExpected, v.ArtifactID)
		p.downloadBar.Increment()

	case events.ArtifactDownloadSuccess:
		if p.downloadBar == nil {
			return errs.New("ArtifactDownloadSuccess called before downloadBar was initialized")
		}
		if _, ok := p.downloadsExpected[v.ArtifactID]; !ok {
			return errs.New("ArtifactDownloadSuccess called for an artifact that was not expected: %s", v.ArtifactID.String())
		}
		if err := p.dropArtifactBar(v.ArtifactID, StepDownload); err != nil {
			return errs.Wrap(err, "Failed to drop install bar")
		}
		if p.downloadBar.Current() == p.downloadBar.total {
			return errs.New("Download bar is already complete, this should not happen")
		}
		delete(p.downloadsExpected, v.ArtifactID)
		p.downloadBar.Increment()

	case events.ArtifactInstallStarted:
		if p.installBar == nil {
			p.installBar = p.addTotalBar(locale.Tl("progress_building", "Installing"), int64(len(p.installsExpected)), mpb.BarPriority(StepInstall.priority))
		}
		if _, ok := p.installsExpected[v.ArtifactID]; !ok {
			return errs.New("ArtifactInstallStarted called for an artifact that was not expected: %s", v.ArtifactID.String())
		}
		if err := p.addArtifactBar(v.ArtifactID, StepInstall, int64(v.TotalSize), true); err != nil {
			return errs.Wrap(err, "Failed to add or update artifact bar")
		}

	case events.ArtifactInstallSkipped:
		if p.installBar == nil {
			return errs.New("ArtifactInstallSkipped called before installBar was initialized")
		}
		delete(p.installsExpected, v.ArtifactID)
		p.installBar.Increment()

	case events.ArtifactInstallSuccess:
		if p.installBar == nil {
			return errs.New("ArtifactInstall[Skipped|Success] called before installBar was initialized")
		}
		if _, ok := p.installsExpected[v.ArtifactID]; !ok {
			return errs.New("ArtifactInstallSuccess called for an artifact that was not expected: %s", v.ArtifactID.String())
		}
		if err := p.dropArtifactBar(v.ArtifactID, StepInstall); err != nil {
			return errs.Wrap(err, "Failed to drop install bar")
		}
		if p.installBar.Current() == p.installBar.total {
			return errs.New("Install bar is already complete, this should not happen")
		}
		delete(p.installsExpected, v.ArtifactID)
		p.installBar.Increment()

	case events.ArtifactInstallProgress:
		if err := p.updateArtifactBar(v.ArtifactID, StepInstall, v.IncrementBySize); err != nil {
			return errs.Wrap(err, "Failed to add or update artifact bar")
		}

	}

	return nil
}

func (p *ProgressDigester) Close() error {
	mainProgressDone := make(chan struct{}, 1)
	go func() {
		p.mainProgress.Wait()
		mainProgressDone <- struct{}{}
	}()

	select {
	case <-mainProgressDone:
		break

	// Wait one second, which should be plenty as we're really just waiting for the last frame to render
	// If it's not done after 1 second it's unlikely it will ever be and it means it did not receive events in a way
	// that we can make sense of.
	case <-time.After(time.Second):
		p.cancelMpb() // mpb doesn't have a Close, just a Wait. We force it as we don't want to give it the opportunity to block.

		bars := map[string]*bar{
			"build bar":    p.buildBar,
			"download bar": p.downloadBar,
			"install bar":  p.installBar,
		}

		pending := 0
		debugMsg := []string{}
		for name, bar := range bars {
			debugMsg = append(debugMsg, fmt.Sprintf("%s is at %v", name, func() string {
				if bar == nil {
					return "nil"
				}
				if !bar.Completed() {
					pending++
				}
				return fmt.Sprintf("%d out of %d", bar.Current(), bar.total)
			}()))
		}

		logging.Debug(`
Timed out waiting for progress bars to close.
Progress bars status:
%s
Still expecting:
 - Builds: %v
 - Downloads: %v
 - Installs: %v
Event log:
%s`,
			strings.Join(debugMsg, "\n"),
			p.buildsExpected, p.downloadsExpected, p.installsExpected,
			strings.Join(p.dbgEventLog, " > "),
		)

		if pending > 0 {
			// We only error out if we determine the issue is down to one of our bars not completing.
			// Otherwise this is an issue with the mpb package which is currently a known limitation, end goal is to get rid of mpb.
			return locale.NewError("err_rtprogress_outofsync", "", constants.BugTrackerURL, logging.FilePath())
		}
	}

	// Success message. Can't happen in event loop as progressbar lib clears new lines when it closes.
	if p.success && p.changesMade {
		p.out.Notice(locale.T("progress_completed"))
	}

	// Blank line to separate progress from rest of output
	p.out.Notice("")

	return nil
}