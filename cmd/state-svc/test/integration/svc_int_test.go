package integration

import (
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"syscall"
	"testing"
	"time"

	svcApp "github.com/ActiveState/cli/cmd/state-svc/app"
	svcAutostart "github.com/ActiveState/cli/cmd/state-svc/autostart"
	"github.com/ActiveState/cli/internal/app"
	"github.com/ActiveState/cli/internal/condition"
	"github.com/ActiveState/cli/internal/constants"
	"github.com/ActiveState/cli/internal/exeutils"
	"github.com/ActiveState/cli/internal/fileutils"
	"github.com/ActiveState/cli/internal/osutils/autostart"
	"github.com/ActiveState/cli/internal/svcctl"
	"github.com/ActiveState/cli/internal/testhelpers/e2e"
	"github.com/ActiveState/cli/internal/testhelpers/tagsuite"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/suite"
)

type SvcIntegrationTestSuite struct {
	tagsuite.Suite
}

func (suite *SvcIntegrationTestSuite) TestStartStop() {
	// Disable test until we can fix console output on Windows
	// See issue here: https://activestatef.atlassian.net/browse/DX-1311
	suite.T().SkipNow()
	suite.OnlyRunForTags(tagsuite.Service)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	cp := ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("stop"))
	cp.ExpectExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("status"))
	cp.Expect("Service cannot be reached")
	cp.ExpectExitCode(1)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("start"))
	cp.Expect("Starting")
	cp.ExpectExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("status"))
	cp.Expect("Checking")

	// Verify the server is running on its reported port.
	cp.ExpectRe("Port:\\s+:\\d+\\s")
	portRe := regexp.MustCompile("Port:\\s+:(\\d+)")
	port := portRe.FindStringSubmatch(cp.TrimmedSnapshot())[1]
	_, err := net.Listen("tcp", "localhost:"+port)
	suite.Error(err)

	// Verify it created and wrote to its reported log file.
	cp.ExpectRe("Log:\\s+.+?\\.log")
	logRe := regexp.MustCompile("Log:\\s+(.+?\\.log)")
	logFile := logRe.FindStringSubmatch(cp.TrimmedSnapshot())[1]
	suite.True(fileutils.FileExists(logFile), "log file '"+logFile+"' does not exist")
	suite.True(len(fileutils.ReadFileUnsafe(logFile)) > 0, "log file is empty")

	cp.ExpectExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("stop"))
	cp.Expect("Stopping")
	cp.ExpectExitCode(0)
	time.Sleep(500 * time.Millisecond) // wait for service to stop

	// Verify the port is free.
	server, err := net.Listen("tcp", "localhost:"+port)
	suite.NoError(err)
	server.Close()
}

func (suite *SvcIntegrationTestSuite) TestSignals() {
	if condition.OnCI() {
		// https://activestatef.atlassian.net/browse/DX-964
		// https://activestatef.atlassian.net/browse/DX-980
		suite.T().Skip("Signal handling on CI is unstable and unreliable")
	}

	if runtime.GOOS == "windows" {
		suite.T().Skip("Windows does not support signal sending.")
	}

	suite.OnlyRunForTags(tagsuite.Service)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	// SIGINT (^C)
	cp := ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("foreground"))
	cp.Expect("Starting")
	time.Sleep(1 * time.Second) // wait for the service to start up
	cp.Signal(syscall.SIGINT)
	cp.Expect("caught a signal: interrupt")
	cp.ExpectNotExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("status"))
	cp.Expect("Service cannot be reached")
	cp.ExpectExitCode(1)

	sockFile := svcctl.NewIPCSockPathFromGlobals().String()
	suite.False(fileutils.TargetExists(sockFile), "socket file was not deleted")

	// SIGTERM
	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("foreground"))
	cp.Expect("Starting")
	time.Sleep(1 * time.Second) // wait for the service to start up
	cp.Signal(syscall.SIGTERM)
	suite.NotContains(cp.TrimmedSnapshot(), "caught a signal")
	cp.ExpectExitCode(0) // should exit gracefully

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("status"))
	cp.Expect("Service cannot be reached")
	cp.ExpectExitCode(1)

	suite.False(fileutils.TargetExists(sockFile), "socket file was not deleted")
}

func (suite *SvcIntegrationTestSuite) TestStartDuplicateErrorOutput() {
	// https://activestatef.atlassian.net/browse/DX-1136
	suite.OnlyRunForTags(tagsuite.Service)
	if runtime.GOOS == "windows" {
		suite.T().Skip("Windows doesn't seem to read from svc at the moment")
	}

	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	cp := ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("stop"))
	cp.ExpectExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("status"))
	cp.ExpectNotExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("start"))
	cp.ExpectExitCode(0)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("foreground"))
	cp.Expect("not start service: An existing")
	cp.ExpectExitCode(1)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("foreground", "test this"))
	cp.Expect("not start service (invoked by \"test this\"): An existing")
	cp.ExpectExitCode(1)

	cp = ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("stop"))
	cp.ExpectExitCode(0)
}

func (suite *SvcIntegrationTestSuite) TestSingleSvc() {
	suite.OnlyRunForTags(tagsuite.Service)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	ts.SpawnCmdWithOpts(ts.SvcExe, e2e.WithArgs("stop"))
	time.Sleep(2 * time.Second) // allow for some time to stop the existing available process

	oldCount := suite.GetNumStateSvcProcesses() // may be non-zero due to non-test state-svc processes (using different sock file)
	for i := 1; i <= 10; i++ {
		go ts.SpawnCmdWithOpts(ts.Exe, e2e.WithArgs("--version"))
		time.Sleep(50 * time.Millisecond) // do not spam CPU
	}
	time.Sleep(2 * time.Second) // allow for some time to spawn the processes

	for attempts := 100; attempts > 0; attempts-- {
		suite.T().Log("iters left:", attempts, "procs:", suite.GetNumStateSvcProcesses())
		if suite.GetNumStateSvcProcesses() == oldCount+1 {
			break
		}
		time.Sleep(2 * time.Second) // keep waiting
	}

	newCount := suite.GetNumStateSvcProcesses()
	if newCount > oldCount+1 {
		// We only care if we end up with more services than anticipated. We can actually end up with less than we started
		// with due to other integration tests not always waiting for state-svc to have fully shut down before running the next test
		suite.Fail(fmt.Sprintf("spawning multiple state processes should only result in one more state-svc process at most, newCount: %d, oldCount: %d", newCount, oldCount))
	}
}

func (suite *SvcIntegrationTestSuite) GetNumStateSvcProcesses() int {
	procs, err := process.Processes()
	suite.NoError(err)

	count := 0
	for _, p := range procs {
		if name, err := p.Name(); err == nil {
			name = filepath.Base(name) // just in case an absolute path is returned
			if svcName := constants.ServiceCommandName + exeutils.Extension; name == svcName {
				count++
			}
		}
	}

	return count
}

func (suite *SvcIntegrationTestSuite) TestAutostartConfigEnableDisable() {
	// Disable test for v0.36: https://activestatef.atlassian.net/browse/DX-1501.
	// This test should be re-enabled by https://activestatef.atlassian.net/browse/DX-1435.
	suite.T().SkipNow()

	suite.OnlyRunForTags(tagsuite.Service)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	app, err := app.New(constants.SvcAppName, ts.SvcExe, nil, svcApp.Options)
	suite.Require().NoError(err)
	enabled, err := autostart.IsEnabled(app.Exec, svcAutostart.Options)
	suite.Require().NoError(err)

	// Toggle it via state tool config.
	cp := ts.SpawnWithOpts(e2e.WithArgs("config", "set", constants.AutostartSvcConfigKey, strconv.FormatBool(!enabled)))
	cp.ExpectExitCode(0)
	suite.checkEnabled(app.Exec, svcAutostart.Options, !enabled)

	// Toggle it again via state tool config.
	cp = ts.SpawnWithOpts(e2e.WithArgs("config", "set", constants.AutostartSvcConfigKey, strconv.FormatBool(enabled)))
	cp.ExpectExitCode(0)
	suite.checkEnabled(app.Exec, svcAutostart.Options, enabled)
}

type autostartApp interface {
	IsAutostartEnabled() (bool, error)
}

func (suite *SvcIntegrationTestSuite) checkEnabled(exec string, opts autostart.Options, expect bool) {
	timeout := time.After(1 * time.Minute)
	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-timeout:
			suite.Fail("autostart has not been changed")
		case <-tick:
			toggled, err := autostart.IsEnabled(exec, opts)
			suite.Require().NoError(err)
			if suite.Equal(expect, toggled) {
				return
			}
		}
	}
}

func TestSvcIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(SvcIntegrationTestSuite))
}
