//go:build !windows
// +build !windows

package osutils

import (
	"os"
	"os/exec"
	"syscall"
)

// CmdExitCode returns the exit code of a command
func CmdExitCode(cmd *exec.Cmd) (code int) {
	return cmd.ProcessState.ExitCode()
}

func BashifyPathEnv(pathList string) (string, error) {
	return pathList, nil // already bashified
}

// InheritEnv returns a union of the given environment and os.Environ(). If the given environment
// and os.Environ() share any environment variables, the former's will be used over the latter's.
func InheritEnv(env map[string]string) map[string]string {
	for k, v := range EnvSliceToMap(os.Environ()) {
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	return env
}

// IsAccessDeniedError is primarily used to determine if an operation failed due to insufficient
// permissions (e.g. attempting to kill an admin process as a normal user)
func IsAccessDeniedError(err error) bool {
	errno, ok := err.(syscall.Errno)
	return ok && errno == syscall.EPERM
}
