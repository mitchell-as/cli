package executors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	rt "runtime"
	"strings"

	"github.com/ActiveState/cli/internal/exeutils"
	"github.com/ActiveState/cli/internal/installation"
	"github.com/ActiveState/cli/internal/osutils"
	"github.com/ActiveState/cli/pkg/platform/runtime/envdef"
	"github.com/ActiveState/cli/pkg/platform/runtime/executors/execmeta"
	"github.com/go-openapi/strfmt"

	"github.com/ActiveState/cli/internal/errs"
	"github.com/ActiveState/cli/internal/fileutils"
	"github.com/ActiveState/cli/internal/locale"
	"github.com/ActiveState/cli/internal/logging"
)

// executorTarget tracks the target executable of the executor and is used to determine whether an existing
// executor needs to be updating.
// Update this if you want to blow away older targets (ie. you made changes to the template)
const executorTarget = "Target: "

type Targeter interface {
	CommitUUID() strfmt.UUID
	Name() string
	Owner() string
	Dir() string
	Headless() bool
}

type Executors struct {
	executorPath string // The location to store the executors

	altExecSrcPath string // Path to alternate executor for testing. Executor() will use global func if not set.
}

func New(executorPath string) *Executors {
	return &Executors{
		executorPath: executorPath,
	}
}

func (es *Executors) ExecutorSrc() (string, error) {
	if es.altExecSrcPath != "" {
		return es.altExecSrcPath, nil
	}
	return installation.ExecutorExec()
}

func (es *Executors) Apply(sockPath string, targeter Targeter, env map[string]string, exes envdef.ExecutablePaths) error {
	logging.Debug("Creating executors at %s, exes: %v", es.executorPath, exes)

	// We need to cover the use case of someone running perl.exe/python.exe
	// Proper fix scheduled here https://www.pivotaltracker.com/story/show/177845386
	if rt.GOOS == "windows" {
		for _, exe := range exes {
			if !strings.HasSuffix(exe, exeutils.Extension) {
				continue
			}
			exes = append(exes, exe+exeutils.Extension) // Double up on the ext so only the first on gets dropped
		}
	}

	if err := es.Clean(); err != nil {
		return errs.Wrap(err, "Could not clean up old executors")
	}

	if err := fileutils.MkdirUnlessExists(es.executorPath); err != nil {
		return locale.WrapError(err, "err_mkdir", "Could not create directory: {{.V0}}", es.executorPath)
	}

	t := execmeta.Target{}
	m := execmeta.New(sockPath, osutils.EnvMapToSlice(env), t, exes)
	if err := m.WriteToDisk(es.executorPath); err != nil {
		return err
	}

	executorExec, err := es.ExecutorSrc()
	if err != nil {
		return locale.WrapError(err, "err_state_exec")
	}

	for _, exe := range exes {
		if err := copyExecutor(es.executorPath, exe, executorExec); err != nil {
			return locale.WrapError(err, "err_createexecutor", "Could not create executor for {{.V0}}.", exe)
		}
	}

	return nil
}

func (es *Executors) Clean() error {
	if !fileutils.DirExists(es.executorPath) {
		return nil
	}

	files, err := ioutil.ReadDir(es.executorPath)
	if err != nil {
		return errs.Wrap(err, "Could not read dir: %s", es.executorPath)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(es.executorPath, file.Name())
		b, err := fileutils.ReadFile(filePath)
		if err != nil {
			return locale.WrapError(err, "err_cleanexecutor_noread", "Could not read potential executor file: {{.V0}}.", file.Name())
		}

		if !isOwnedByUs(b) {
			continue
		}

		if err := os.Remove(filePath); err != nil {
			return locale.WrapError(err, "err_cleanexecutor_remove", "Could not remove executor: {{.V0}}", file.Name())
		}
	}

	return nil
}

func IsExecutor(filePath string) (bool, error) {
	if fileutils.IsDir(filePath) {
		return false, nil
	}

	b, err := fileutils.ReadFile(filePath)
	if err != nil {
		return false, locale.WrapError(err, "err_cleanexecutor_noread", "Could not read potential executor file: {{.V0}}.", filePath)
	}
	return isOwnedByUs(b), nil
}

func isOwnedByUs(fileContents []byte) bool {
	return strings.Contains(string(fileContents), "state-exec") ||
		execmeta.IsMetaFile(fileContents) ||
		legacyIsOwnedByUs(fileContents)
}

func copyExecutor(destDir, exe, srcExec string) error {
	name := filepath.Base(exe)
	target := filepath.Clean(filepath.Join(destDir, name))

	if strings.HasSuffix(exe, exeutils.Extension+exeutils.Extension) {
		// This is super awkward, but we have a double .exe to temporarily work around an issue that will be fixed
		// more correctly here - https://www.pivotaltracker.com/story/show/177845386
		exe = strings.TrimSuffix(exe, exeutils.Extension)
	}

	logging.Debug("w/Creating executor for %s at %s", exe, target)

	if fileutils.TargetExists(target) {
		b, err := fileutils.ReadFile(target)
		if err != nil {
			return locale.WrapError(err, "err_createexecutor_exists_noread", "Could not create executor as target already exists and could not be read: {{.V0}}.", target)
		}
		if !isOwnedByUs(b) {
			return locale.WrapError(err, "err_createexecutor_exists", "Could not create executor as target already exists: {{.V0}}.", target)
		}
	}

	if err := fileutils.CopyFile(srcExec, target); err != nil {
		return locale.WrapError(err, "err_copyexecutor_fail", "Could not copy {{.V0}} to {{.V1}}", srcExec, target)
	}

	if err := os.Chmod(target, 0755); err != nil {
		return locale.WrapError(err, "err_setexecmode_fail", "Could not set mode of {{.V0}}", target)
	}

	return nil
}

// denoter constants are to ensure we clean up old executors, but are deprecated as of this comment
const (
	legacyExecutorDenoter = "!DO NOT EDIT! State Tool Executor !DO NOT EDIT!"
	legacyShimDenoter     = "!DO NOT EDIT! State Tool Shim !DO NOT EDIT!"
)

func legacyIsOwnedByUs(fileContents []byte) bool {
	s := string(fileContents)
	return strings.Contains(s, legacyExecutorDenoter) || strings.Contains(s, legacyShimDenoter)
}