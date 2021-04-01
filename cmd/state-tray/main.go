package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/getlantern/systray"
	"github.com/gobuffalo/packr"

	"github.com/ActiveState/cli/cmd/state-tray/internal/open"
	"github.com/ActiveState/cli/cmd/state-tray/pkg/autostart"
	"github.com/ActiveState/cli/internal/errs"
	"github.com/ActiveState/cli/internal/fileutils"
	"github.com/ActiveState/cli/internal/locale"
	"github.com/ActiveState/cli/internal/logging"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	err := run()
	if err != nil {
		logging.Error("Systray encountered an error: %v", errs.Join(err, ": "))
		os.Exit(1)
	}
}

func run() error {
	if os.Getenv("VERBOSE") == "true" {
		// Doesn't seem to work, I think the systray lib and its logging solution is interfering
		logging.CurrentHandler().SetVerbose(true)
	}

	stateSvcExe := filepath.Join(filepath.Dir(os.Args[0]), "state-svc")
	if runtime.GOOS == "windows" {
		stateSvcExe = stateSvcExe + ".exe"
	}
	if !fileutils.FileExists(stateSvcExe) {
		return errs.New("Could not find: %s", stateSvcExe)
	}

	cmd := exec.Command(stateSvcExe, "start")
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	if err := cmd.Start(); err != nil {
		return errs.Wrap(err, "Could not start %s", stateSvcExe)
	}

	box := packr.NewBox("../../assets")
	systray.SetIcon(box.Bytes("icon.ico"))
	systray.SetTooltip(locale.Tl("tray_tooltip", "ActiveState State Tool"))

	mAbout := systray.AddMenuItem(
		locale.Tl("tray_about_title", "About State Tool"),
		locale.Tl("tray_about_tooltip", "Information about the State Tool"),
	)

	systray.AddSeparator()

	mDoc := systray.AddMenuItem(
		locale.Tl("tray_documentation_title", "Documentation"),
		locale.Tl("tray_documentation_tooltip", "Open State Tool Docs"),
	)

	mPlatform := systray.AddMenuItem(locale.Tl("tray_platform_title", "ActiveState Platform"), "")
	mLearn := mPlatform.AddSubMenuItem(
		locale.Tl("tray_learn_title", "Learn"),
		locale.Tl("tray_learn_tooltip", "Open ActiveState blog"),
	)
	mSupport := mPlatform.AddSubMenuItem(
		locale.Tl("tray_support_title", "Support"),
		locale.Tl("tray_support_tooltip", "Open support page"),
	)
	mAccount := mPlatform.AddSubMenuItem(
		locale.Tl("tray_account_title", "Account"),
		locale.Tl("tray_account_tooltip", "Open your account page"),
	)

	systray.AddSeparator()
	mAutoStart := systray.AddMenuItem(locale.Tl("tray_autostart", "Start on Login"), "")
	if autostart.New().IsEnabled() {
		mAutoStart.Check()
	}
	systray.AddSeparator()

	mQuit := systray.AddMenuItem(locale.Tl("tray_exit", "Exit"), "")

	for {
		select {
		case <-mAbout.ClickedCh:
			logging.Debug("About event")
			err := open.Prompt("state --version")
			if err != nil {
				logging.Error("Could not open command prompt, got error: %v", err)
			}
		case <-mDoc.ClickedCh:
			logging.Debug("Documentation event")
			// Not implemented
		case <-mLearn.ClickedCh:
			logging.Debug("Learn event")
			// Not implemented
		case <-mSupport.ClickedCh:
			logging.Debug("Support event")
			// Not implemented
		case <-mAccount.ClickedCh:
			logging.Debug("Account event")
			// Not implemented
		case <-mAutoStart.ClickedCh:
			logging.Debug("Autostart event")
			as := autostart.New()
			var err error
			if as.IsEnabled() {
				logging.Debug("Disable")
				err = as.Disable()
				if err == nil {
					mAutoStart.Uncheck()
				}
			} else {
				logging.Debug("Enable")
				err = as.Enable()
				if err == nil {
					mAutoStart.Check()
				}
			}
			if err != nil {
				logging.Error("Could not toggle autostart tray: %v", errs.Join(err, ": "))
			}
		case <-mQuit.ClickedCh:
			logging.Debug("Quit event")
			systray.Quit()
		}
	}

	return nil
}

func onExit() {
	// Not implemented
}