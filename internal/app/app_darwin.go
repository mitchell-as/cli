package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveState/cli/internal/assets"
	"github.com/ActiveState/cli/internal/constants"
	"github.com/ActiveState/cli/internal/errs"
	"github.com/ActiveState/cli/internal/fileutils"
	"github.com/ActiveState/cli/internal/installation"
	"github.com/ActiveState/cli/internal/logging"
	"github.com/ActiveState/cli/internal/strutils"
)

const (
	execFileSource   = "exec.sh.tpl"
	launchFileSource = "com.activestate.platform.app.plist.tpl"
	iconFile         = "icon.icns"
	assetAppDir      = "placeholder.app"
)

func (a *App) install() error {
	// Create all of the necessary directories and files in a temporary directory
	// Then move the temporary directory to the final location which for macOS will be the Applications directory
	tmpDir, err := ioutil.TempDir("", fmt.Sprintf("%s-", a.Name))
	if err != nil {
		return errs.Wrap(err, "Could not create temporary directory")
	}
	defer os.RemoveAll(tmpDir)

	tmpAppPath := filepath.Join(tmpDir, fmt.Sprintf("%s.app", a.Name))
	if err := fileutils.Mkdir(tmpAppPath); err != nil {
		return errs.Wrap(err, "Could not create .app directory")
	}

	if err := fileutils.CopyFilesDirReader(assets.NewAssetsFS(), assetAppDir, tmpAppPath, assets.PlaceholderFileName); err != nil {
		return errs.Wrap(err, "Could not copy files from assets")
	}

	if err := a.createIcon(tmpAppPath); err != nil {
		return errs.Wrap(err, "Could not create icon")
	}

	if err := a.createExecFile(tmpAppPath); err != nil {
		return errs.Wrap(err, "Could not create exec file")
	}

	if err := a.createInfoFile(tmpAppPath); err != nil {
		return errs.Wrap(err, "Could not create info file")
	}

	installDir := os.Getenv(constants.AppInstallDirOverrideEnvVarName)
	if installDir == "" {
		installDir, err = installation.ApplicationInstallPath()
		if err != nil {
			return errs.Wrap(err, "Could not get installation path")
		}
	}

	if err := fileutils.MoveAllFiles(tmpDir, installDir); err != nil {
		return errs.Wrap(err, "Could not move .app to Applications directory")
	}

	return nil
}

func (a *App) createIcon(path string) error {
	icon, err := assets.ReadFileBytes(a.options.IconFileSource)
	if err != nil {
		return errs.Wrap(err, "Could not read asset")
	}

	if err = fileutils.WriteFile(filepath.Join(path, "Contents", "Resources", iconFile), icon); err != nil {
		return errs.Wrap(err, "Could not write icon file")
	}

	return nil
}

func (a *App) createExecFile(base string) error {
	path := filepath.Join(base, "Contents", "MacOS")
	asset, err := assets.ReadFileBytes(execFileSource)
	if err != nil {
		return errs.Wrap(err, "Could not read asset")
	}

	scriptFile := fmt.Sprintf("%s.sh", filepath.Base(a.Exec))

	content, err := strutils.ParseTemplate(
		string(asset),
		map[string]interface{}{
			"Exec": a.Exec,
			"Args": strings.Join(a.Args, " "),
		})
	if err != nil {
		return errs.Wrap(err, "Could not parse launch file source")
	}

	err = fileutils.WriteFile(filepath.Join(path, scriptFile), []byte(content))
	if err != nil {
		return errs.Wrap(err, "Could not write Info.plist file")
	}

	err = os.Chmod(filepath.Join(path, scriptFile), 0755)
	if err != nil {
		return errs.Wrap(err, "Could not make executable")
	}

	return nil
}

func (a *App) createInfoFile(base string) error {
	path := filepath.Join(base, "Contents")
	asset, err := assets.ReadFileBytes(launchFileSource)
	if err != nil {
		return errs.Wrap(err, "Could not read asset")
	}

	scriptFile := fmt.Sprintf("%s.sh", filepath.Base(a.Exec))

	content, err := strutils.ParseTemplate(
		string(asset),
		map[string]interface{}{
			"Exec":         scriptFile,
			"Icon":         a.options.IconFileName,
			"HideDockIcon": a.options.MacHideDockIcon,
			"IsGUIApp":     a.options.IsGUIApp,
		})
	if err != nil {
		return errs.Wrap(err, "Could not parse launch file source")
	}

	err = fileutils.WriteFile(filepath.Join(path, "Info.plist"), []byte(content))
	if err != nil {
		return errs.Wrap(err, "Could not write Info.plist file")
	}

	return nil
}

func (a *App) uninstall() error {
	defaultPath, err := installation.ApplicationInstallPath()
	if err != nil {
		return errs.Wrap(err, "Could not get installation path")
	}

	installDir := filepath.Join(defaultPath, fmt.Sprintf("%s.app", a.Name))
	if !fileutils.DirExists(installDir) {
		logging.Debug("Directory does not exist, nothing to do")
		return nil
	}

	if err := os.RemoveAll(installDir); err != nil {
		logging.Debug("Could not remove %s: %v", installDir, err)
		return errs.Wrap(err, "Could not remove .app from Applications directory")
	}

	return nil
}