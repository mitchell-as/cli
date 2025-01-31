package integration

import (
	"fmt"
	"testing"

	"github.com/ActiveState/cli/internal/testhelpers/e2e"
	"github.com/ActiveState/cli/internal/testhelpers/tagsuite"
	"github.com/stretchr/testify/suite"
)

type RefreshIntegrationTestSuite struct {
	tagsuite.Suite
}

func (suite *RefreshIntegrationTestSuite) TestRefresh() {
	suite.OnlyRunForTags(tagsuite.Refresh)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts, "ActiveState-CLI", "Branches", "main", "35af7414-b44b-4fd7-aa93-2ecad337ed2b")

	cp := ts.SpawnWithOpts(
		e2e.WithArgs("refresh"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.Expect("Setting Up Runtime")
	cp.Expect("Runtime updated")
	cp.ExpectExitCode(0)

	cp = ts.SpawnWithOpts(
		e2e.WithArgs("exec", "--", "python3", "-c", "import requests"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.Expect("ModuleNotFoundError")
	cp.ExpectExitCode(1)

	suite.PrepareActiveStateYAML(ts, "ActiveState-CLI", "Branches", "secondbranch", "46c83477-d580-43e2-a0c6-f5d3677517f1")
	cp = ts.SpawnWithOpts(
		e2e.WithArgs("refresh"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.Expect("Setting Up Runtime")
	cp.Expect("Runtime updated")
	cp.ExpectExitCode(0)

	cp = ts.SpawnWithOpts(
		e2e.WithArgs("exec", "--", "python3", "-c", "import requests"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.ExpectExitCode(0)

	cp = ts.Spawn("refresh")
	suite.Assert().NotContains(cp.TrimmedSnapshot(), "Setting Up Runtime", "Unchanged runtime should not refresh")
	cp.Expect("Runtime updated")
	cp.ExpectExitCode(0)
}

func (suite *RefreshIntegrationTestSuite) PrepareActiveStateYAML(ts *e2e.Session, username, project, branch, commitID string) {
	asyData := fmt.Sprintf(`project: "https://platform.activestate.com/%s/%s?branch=%s&commitID=%s"`, username, project, branch, commitID)
	ts.PrepareActiveStateYAML(asyData)
}

func TestRefreshIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RefreshIntegrationTestSuite))
}
