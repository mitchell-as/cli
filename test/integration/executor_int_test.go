package integration

import (
	"testing"

	"github.com/ActiveState/cli/internal/testhelpers/e2e"
	"github.com/ActiveState/cli/internal/testhelpers/tagsuite"
	"github.com/stretchr/testify/suite"
)

type ExecutorIntegrationTestSuite struct {
	tagsuite.Suite
}

func TestExecutorIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutorIntegrationTestSuite))
}

func (suite *ExecutorIntegrationTestSuite) TestExecutorForwards() {
	suite.OnlyRunForTags(tagsuite.Executor)

	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	cp := ts.SpawnWithOpts(
		e2e.WithArgs("checkout", "ActiveState-CLI/Python3"),
	)
	cp.Expect("Checked out project")
	cp.ExpectExitCode(0)

	cp = ts.SpawnWithOpts(
		e2e.WithArgs("shell", "ActiveState-CLI/Python3"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.Expect("Activated")
	cp.WaitForInput()

	cp.SendLine("python3 -c \"import sys; print(sys.copyright)\"")
	cp.Expect("ActiveState Software Inc.")

	cp.SendLine("exit")
	cp.Expect("Deactivated")
	cp.ExpectExitCode(0)
}

func (suite *ExecutorIntegrationTestSuite) TestExecutorExitCode() {
	suite.OnlyRunForTags(tagsuite.Executor)

	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	cp := ts.SpawnWithOpts(
		e2e.WithArgs("checkout", "ActiveState-CLI/Python3"),
	)
	cp.Expect("Checked out project")
	cp.ExpectExitCode(0)

	cp = ts.SpawnWithOpts(
		e2e.WithArgs("shell", "ActiveState-CLI/Python3"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.Expect("Activated")
	cp.WaitForInput()

	cp.SendLine("python3 -c \"exit(42)\"")

	cp.SendLine("exit")
	cp.ExpectExitCode(42)
}