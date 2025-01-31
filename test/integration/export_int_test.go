package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/ActiveState/cli/internal/testhelpers/e2e"
	"github.com/ActiveState/cli/internal/testhelpers/tagsuite"
	"github.com/stretchr/testify/suite"
)

type ExportIntegrationTestSuite struct {
	tagsuite.Suite
}

func (suite *ExportIntegrationTestSuite) TestExport_Export() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	cp := ts.Spawn("export", "recipe")
	cp.Expect("{\"camel_flags\":")
	cp.ExpectExitCode(0)
}

func (suite *ExportIntegrationTestSuite) TestExport_ExportArg() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	cp := ts.Spawn("export", "recipe")
	cp.Expect("{\"camel_flags\":")
	cp.ExpectExitCode(0)
}

func (suite *ExportIntegrationTestSuite) TestExport_ExportPlatform() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	cp := ts.Spawn("export", "recipe", "--platform", "linux")
	cp.Expect("{\"camel_flags\":")
	cp.ExpectExitCode(0)
}

func (suite *ExportIntegrationTestSuite) TestExport_InvalidPlatform() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	cp := ts.Spawn("export", "recipe", "--platform", "junk")
	cp.ExpectExitCode(1)
}

func (suite *ExportIntegrationTestSuite) TestExport_ConfigDir() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	cp := ts.Spawn("export", "config", "--filter", "junk")
	cp.ExpectExitCode(1)
}

func (suite *ExportIntegrationTestSuite) TestExport_Config() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	cp := ts.Spawn("export", "config")
	cp.Expect(`dir: `)
	cp.ExpectLongString(ts.Dirs.Config, time.Second)
	cp.ExpectExitCode(0)
}

func (suite *ExportIntegrationTestSuite) TestExport_Env() {
	suite.OnlyRunForTags(tagsuite.Export)
	ts := e2e.New(suite.T(), false)
	defer ts.Close()

	suite.PrepareActiveStateYAML(ts)
	asyData := fmt.Sprintf(`project: "https://platform.activestate.com/ActiveState-CLI/Export?branch=main&commitID=5397f645-da8a-4591-b106-9d7fa99545fe"`)
	ts.PrepareActiveStateYAML(asyData)
	cp := ts.SpawnWithOpts(
		e2e.WithArgs("export", "env"),
		e2e.AppendEnv("ACTIVESTATE_CLI_DISABLE_RUNTIME=false"),
	)
	cp.Expect(`PATH: `)
	cp.ExpectExitCode(0)

	suite.Assert().NotContains(cp.TrimmedSnapshot(), "ACTIVESTATE_ACTIVATED")
}

func TestExportIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ExportIntegrationTestSuite))
}

func (suite *ExportIntegrationTestSuite) PrepareActiveStateYAML(ts *e2e.Session) {
	asyData := `project: "https://platform.activestate.com/cli-integration-tests/Export"`
	ts.PrepareActiveStateYAML(asyData)
}
