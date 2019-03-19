package auth_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ActiveState/cli/internal/constants"
	"github.com/ActiveState/cli/internal/environment"
	"github.com/ActiveState/cli/internal/failures"
	secretsapi "github.com/ActiveState/cli/internal/secrets-api"
	secretsModels "github.com/ActiveState/cli/internal/secrets-api/models"
	"github.com/ActiveState/cli/internal/testhelpers/httpmock"
	"github.com/ActiveState/cli/internal/testhelpers/osutil"
	"github.com/ActiveState/cli/internal/testhelpers/secretsapi_test"
	authlet "github.com/ActiveState/cli/pkg/cmdlets/auth"
	"github.com/ActiveState/cli/pkg/platform/api"
	"github.com/ActiveState/cli/pkg/platform/api/models"
	"github.com/ActiveState/cli/pkg/platform/authentication"
	authCmd "github.com/ActiveState/cli/state/auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

var Command = authCmd.Command

func setup(t *testing.T) {
	failures.ResetHandled()
	authentication.Logout()
	secretsapi_test.InitializeTestClient("bearer123")

	root, err := environment.GetRootPath()
	assert.NoError(t, err, "Should detect root path")
	os.Chdir(filepath.Join(root, "test"))

	Cc := Command.GetCobraCmd()
	Cc.SetArgs([]string{})
	authCmd.Args.Token = ""

	authlet.OpenURI = func(uri string) error { return nil }
}

func setupUser() *models.UserEditable {
	testUser := &models.UserEditable{
		Username: "test",
		Email:    "test@test.tld",
		Password: "foo", // this matches the passphrase on testdata/self-private.key
		Name:     "Test User",
	}
	return testUser
}

func TestExecuteNoArgs(t *testing.T) {
	setup(t)

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	defer httpmock.DeActivate()

	httpmock.RegisterWithCode("POST", "/login", 401)

	var execErr error
	osutil.WrapStdinWithDelay(10*time.Millisecond, func() { execErr = Command.Execute() },
		// prompted for username and password only
		// 10ms delay between writes to stdin
		"baduser",
		"badpass",
	)

	assert.NoError(t, execErr, "Executed without error")
	assert.Error(t, failures.Handled(), "No failure occurred")
	assert.Nil(t, authentication.ClientAuth(), "Did not authenticate")
}

func TestExecuteNoArgsAuthenticated_WithExistingKeypair(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")
	httpmock.Register("GET", "/renew")

	fail := authentication.Get().AuthenticateWithModel(&models.Credentials{
		Username: user.Username,
		Password: user.Password,
	})
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	require.NoError(t, fail.ToError())

	assert.NoError(t, Command.Execute(), "Executed without error")
	assert.NoError(t, failures.Handled(), "No failure occurred")
}

func TestExecuteNoArgsLoginByPrompt_WithExistingKeypair(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")
	secretsapiMock.Register("GET", "/keypair")

	var execErr error
	osutil.WrapStdinWithDelay(10*time.Millisecond, func() { execErr = Command.Execute() },
		user.Username,
		user.Password)

	assert.NoError(t, execErr, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
}

func TestExecuteNoArgsLoginByPrompt_NoExistingKeypair(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")

	var bodyKeypair *secretsModels.KeypairChange
	var bodyErr error
	secretsapiMock.RegisterWithCode("GET", "/keypair", 404)
	secretsapiMock.RegisterWithResponder("PUT", "/keypair", func(req *http.Request) (int, string) {
		reqBody, _ := ioutil.ReadAll(req.Body)
		bodyErr = json.Unmarshal(reqBody, &bodyKeypair)
		return 204, "empty"
	})

	var execErr error
	osutil.WrapStdinWithDelay(10*time.Millisecond, func() { execErr = Command.Execute() },
		user.Username,
		user.Password)

	assert.NoError(t, execErr, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")

	require.NoError(t, bodyErr, "unmarshalling keypair save response")
	assert.NotZero(t, bodyKeypair.EncryptedPrivateKey, "published private key")
	assert.NotZero(t, bodyKeypair.PublicKey, "published public key")
}

func TestExecuteNoArgsLoginThenSignupByPrompt(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	var secondRequest bool
	httpmock.RegisterWithResponder("POST", "/login", func(req *http.Request) (int, string) {
		if !secondRequest {
			secondRequest = true
			return 401, "login"
		}
		return 200, "login"
	})
	httpmock.Register("GET", "/users/uniqueUsername/test")
	httpmock.Register("POST", "/users")

	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")

	var bodyKeypair *secretsModels.KeypairChange
	var bodyErr error
	secretsapiMock.RegisterWithCode("GET", "/keypair", 404)
	secretsapiMock.RegisterWithResponder("PUT", "/keypair", func(req *http.Request) (int, string) {
		reqBody, _ := ioutil.ReadAll(req.Body)
		bodyErr = json.Unmarshal(reqBody, &bodyKeypair)
		return 204, "empty"
	})

	var execErr error
	osutil.WrapStdinWithDelay(10*time.Millisecond, func() { execErr = Command.Execute() },
		// prompted for username and password
		user.Username,
		user.Password,
		// prompted to signup instead
		"yes",
		// enter new user details
		user.Password, // confirmation
		user.Name,
		user.Email,
	)

	assert.NoError(t, execErr, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")

	require.NoError(t, bodyErr, "unmarshalling keypair save response")
	assert.NotZero(t, bodyKeypair.EncryptedPrivateKey, "published private key")
	assert.NotZero(t, bodyKeypair.PublicKey, "published public key")
}

func TestExecuteSignup(t *testing.T) {
	setup(t)

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.Register("GET", "/users/uniqueUsername/test")
	httpmock.Register("POST", "/users")
	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")

	var bodyKeypair *secretsModels.KeypairChange
	var bodyErr error
	secretsapiMock.RegisterWithResponder("PUT", "/keypair", func(req *http.Request) (int, string) {
		reqBody, _ := ioutil.ReadAll(req.Body)
		bodyErr = json.Unmarshal(reqBody, &bodyKeypair)
		return 204, "empty"
	})

	user := setupUser()

	Cc := Command.GetCobraCmd()
	Cc.SetArgs([]string{"signup"})

	var execErr error
	osutil.WrapStdinWithDelay(10*time.Millisecond, func() { execErr = Command.Execute() },
		user.Username,
		user.Password,
		user.Password, // confirmation
		user.Name,
		user.Email,
	)

	assert.NoError(t, execErr, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")

	require.NoError(t, bodyErr, "unmarshalling keypair save response")
	assert.NotZero(t, bodyKeypair.EncryptedPrivateKey, "published private key")
	assert.NotZero(t, bodyKeypair.PublicKey, "published public key")
}

func TestExecuteToken(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")

	fail := authentication.Get().AuthenticateWithModel(&models.Credentials{
		Username: user.Username,
		Password: user.Password,
	})
	token := viper.GetString("apiToken")
	authentication.Logout()
	assert.NoError(t, fail.ToError(), "Executed without error")
	assert.Nil(t, authentication.ClientAuth(), "Not Authenticated")

	Cc := Command.GetCobraCmd()
	Cc.SetArgs([]string{token})

	err := Command.Execute()
	assert.NoError(t, err, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
}

func TestExecuteLogout(t *testing.T) {
	setup(t)
	defer osutil.RemoveConfigFile(constants.KeypairLocalFileName + ".key")
	osutil.CopyTestFileToConfigDir("self-private.key", constants.KeypairLocalFileName+".key", 0600)

	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")

	a := authentication.Get()
	fail := a.AuthenticateWithModel(&models.Credentials{
		Username: user.Username,
		Password: user.Password,
	})
	require.NoError(t, fail.ToError())
	assert.True(t, a.Authenticated(), "Authenticated")

	Cc := Command.GetCobraCmd()
	Cc.SetArgs([]string{"logout"})

	err := Command.Execute()
	assert.NoError(t, err, "Executed without error")
	assert.False(t, a.Authenticated(), "Not Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")

	pkstat, err := osutil.StatConfigFile(constants.KeypairLocalFileName + ".key")
	require.Nil(t, pkstat)
	assert.Regexp(t, "no such file or directory", err.Error())
}

func TestExecuteAuthWithTOTP_WithExistingKeypair(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.RegisterWithResponder("POST", "/login", func(req *http.Request) (int, string) {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "totp") {
			return 449, "login"
		}
		return 200, "login"
	})
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")
	secretsapiMock.Register("GET", "/keypair")

	var execErr error
	// \x04 is the equivalent of a ctrl+d, which tells the survey prompter to stop expecting
	// input for the specific field
	osutil.WrapStdinWithDelay(10*time.Millisecond,
		func() { execErr = Command.Execute() },
		user.Username, user.Password, "\x04")

	require.NoError(t, execErr, "Executed without error")
	assert.Nil(t, authentication.ClientAuth(), "Not Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
	failures.ResetHandled()

	osutil.WrapStdinWithDelay(10*time.Millisecond,
		func() { execErr = Command.Execute() },
		user.Username, user.Password, "foo")

	require.NoError(t, execErr, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
	failures.ResetHandled()
}

func TestExecuteAuthWithTOTP_NoExistingKeypair(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()
	defer failures.ResetHandled()

	httpmock.RegisterWithResponder("POST", "/login", func(req *http.Request) (int, string) {
		bodyBytes, _ := ioutil.ReadAll(req.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "totp") {
			return 449, "login"
		}
		return 200, "login"
	})
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")

	var bodyKeypair *secretsModels.KeypairChange
	var bodyErr error
	secretsapiMock.RegisterWithCode("GET", "/keypair", 404)
	secretsapiMock.RegisterWithResponder("PUT", "/keypair", func(req *http.Request) (int, string) {
		reqBody, _ := ioutil.ReadAll(req.Body)
		bodyErr = json.Unmarshal(reqBody, &bodyKeypair)
		return 204, "empty"
	})

	var execErr error
	// \x04 is the equivalent of a ctrl+d, which tells the survey prompter to stop expecting
	// input for the specific field
	osutil.WrapStdinWithDelay(10*time.Millisecond,
		func() { execErr = Command.Execute() },
		user.Username, user.Password, "\x04")

	require.NoError(t, execErr, "Executed without error")
	assert.Nil(t, authentication.ClientAuth(), "Not Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
	failures.ResetHandled()

	osutil.WrapStdinWithDelay(10*time.Millisecond,
		func() { execErr = Command.Execute() },
		user.Username, user.Password, "foo")

	require.NoError(t, execErr, "Executed without error")
	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")

	require.NoError(t, bodyErr, "unmarshalling keypair save response")
	assert.NotZero(t, bodyKeypair.EncryptedPrivateKey, "published private key")
	assert.NotZero(t, bodyKeypair.PublicKey, "published public key")
}

func TestUsernameValidator(t *testing.T) {
	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	defer httpmock.DeActivate()

	httpmock.Register("GET", "/users/uniqueUsername/test")

	err := authlet.UsernameValidator("test")
	assert.NoError(t, err, "Username is unique")

	httpmock.RegisterWithCode("GET", "/users/uniqueUsername/test", 400)

	err = authlet.UsernameValidator("test")
	assert.Error(t, err, "Username is not unique")
}

func TestRequireAuthenticationLogin(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")
	httpmock.Register("GET", "/renew")
	secretsapiMock.Register("GET", "/keypair")

	osutil.WrapStdinWithDelay(10*time.Millisecond, func() {
		authlet.RequireAuthentication("")
	}, "", user.Username, user.Password)

	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
}

func TestRequireAuthenticationLoginFail(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	defer httpmock.DeActivate()

	httpmock.Register("GET", "/users/uniqueUsername/test")
	httpmock.RegisterWithCode("POST", "/login", 401)

	var fail *failures.Failure
	osutil.WrapStdinWithDelay(10*time.Millisecond, func() {
		fail = authlet.RequireAuthentication("")
	}, "", user.Username, user.Password)

	assert.Nil(t, authentication.ClientAuth(), "Not Authenticated")
	require.Error(t, fail.ToError(), "Failure occurred")
	assert.Equal(t, authlet.FailNotAuthenticated.Name, fail.Type.Name)
}

func TestRequireAuthenticationSignup(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.Register("GET", "/users/uniqueUsername/test")
	httpmock.Register("POST", "/users")
	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")

	secretsapiMock.RegisterWithResponder("PUT", "/keypair", func(req *http.Request) (int, string) {
		return 204, "empty"
	})

	osutil.WrapStdinWithDelay(50*time.Millisecond, func() {
		authlet.RequireAuthentication("")
	}, terminal.KeyArrowDown, "", user.Username, user.Password, user.Password, user.Name, user.Email)

	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
}

func TestRequireAuthenticationSignupBrowser(t *testing.T) {
	setup(t)
	user := setupUser()

	httpmock.Activate(api.GetServiceURL(api.ServicePlatform).String())
	secretsapiMock := httpmock.Activate(secretsapi.DefaultClient.BaseURI)
	defer httpmock.DeActivate()

	httpmock.Register("POST", "/login")
	httpmock.Register("GET", "/apikeys")
	httpmock.Register("DELETE", "/apikeys/"+constants.APITokenName)
	httpmock.Register("POST", "/apikeys")
	httpmock.Register("GET", "/renew")
	secretsapiMock.Register("GET", "/keypair")

	var openURICalled bool
	authlet.OpenURI = func(uri string) error {
		openURICalled = true
		return nil
	}

	osutil.WrapStdinWithDelay(50*time.Millisecond, func() {
		authlet.RequireAuthentication("")
	}, terminal.KeyArrowDown, terminal.KeyArrowDown, "", user.Username, user.Password)

	assert.NotNil(t, authentication.ClientAuth(), "Authenticated")
	assert.NoError(t, failures.Handled(), "No failure occurred")
	assert.True(t, openURICalled, "OpenURI was called")
}