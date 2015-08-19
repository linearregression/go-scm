package scm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/codeship/go-exec"
)

const (
	clonePath = "clone"
)

var (
	ValidationErrorTypeRequiredFieldMissing                         ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                          ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType ValidationErrorType = "SecurityNotImplementedForCheckoutOptionsType"

	errorSecurityNotImplementedForCheckoutOptionsType = errors.New("SecurityNotImplementedForCheckoutOptionsType")
)

type ValidationErrorType string

type ValidationError interface {
	error
	Type() ValidationErrorType
}

//go:generate gen-enumtype

// @gen-enumtype CheckoutOptions git 0
type GitCheckoutOptions struct {
	User            string
	UserName        string
	Email           string
	Host            string
	Path            string
	Branch          string
	CommitID        string
	CommitMessage   string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions github 1
type GithubCheckoutOptions struct {
	User            string
	UserName        string
	Email           string
	Repository      string
	Branch          string
	CommitID        string
	CommitMessage   string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions hg 2
type HgCheckoutOptions struct {
	User            string
	UserName        string
	Email           string
	Host            string
	Path            string
	ChangesetID     string
	CommitMessage   string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions bitbucketGit 3
type BitbucketGitCheckoutOptions struct {
	User            string
	UserName        string
	Email           string
	Repository      string
	Branch          string
	CommitID        string
	CommitMessage   string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions bitbucketHg 4
type BitbucketHgCheckoutOptions struct {
	User            string
	UserName        string
	Email           string
	Repository      string
	ChangesetID     string
	CommitMessage   string
	SecurityOptions SecurityOptions
}

// @gen-enumtype SecurityOptions ssh 0
type SSHSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}

// @gen-enumtype SecurityOptions accessToken 1
type AccessTokenSecurityOptions struct {
	AccessToken string
}

func ConvertCheckoutOptions(checkoutOptions CheckoutOptions) (*ExternalCheckoutOptions, error) {
	return convertCheckoutOptions(checkoutOptions)
}

func ConvertSecurityOptions(securityOptions SecurityOptions) (*ExternalSecurityOptions, error) {
	return convertSecurityOptions(securityOptions)
}

type ExternalCheckoutOptions struct {
	Type            string                   `json:"type,omitempty" yaml:"type,omitempty"`
	User            string                   `json:"user,omitempty" yaml:"user,omitempty"`
	UserName        string                   `json:"user_name,omitempty" yaml:"user_name,omitempty"`
	Email           string                   `json:"email,omitempty" yaml:"email,omitempty"`
	Host            string                   `json:"host,omitempty" yaml:"host,omitempty"`
	Path            string                   `json:"path,omitempty" yaml:"path,omitempty"`
	Repository      string                   `json:"repository,omitempty" yaml:"repository,omitempty"`
	Branch          string                   `json:"branch,omitempty" yaml:"branch,omitempty"`
	CommitID        string                   `json:"commit_id,omitempty" yaml:"commit_id,omitempty"`
	CommitMessage   string                   `json:"commit_message,omitempty" yaml:"commit_message,omitempty"`
	ChangesetID     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
	SecurityOptions *ExternalSecurityOptions `json:"security_options,omitempty" yaml:"security_options,omitempty"`
}

type ExternalSecurityOptions struct {
	Type                  string `json:"type,omitempty" yaml:"type,omitempty"`
	StrictHostKeyChecking bool   `json:"strict_host_key_checking,omitempty" yaml:"strict_host_key_checking,omitempty"`
	PrivateKey            string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	AccessToken           string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
}

func ConvertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	return convertExternalCheckoutOptions(externalCheckoutOptions)
}
func Checkout(
	checkoutOptions CheckoutOptions,
	absolutePath string,
) error {
	return checkout(
		checkoutOptions,
		absolutePath,
	)
}

// ***** PRIVATE *****

func convertCheckoutOptions(checkoutOptions CheckoutOptions) (*ExternalCheckoutOptions, error) {
	var externalCheckoutOptions *ExternalCheckoutOptions
	if switchErr := CheckoutOptionsSwitch(
		checkoutOptions,
		func(gitCheckoutOptions *GitCheckoutOptions) error {
			var externalSecurityOptions *ExternalSecurityOptions
			var err error
			if gitCheckoutOptions.SecurityOptions != nil {
				externalSecurityOptions, err = ConvertSecurityOptions(gitCheckoutOptions.SecurityOptions)
				if err != nil {
					return err
				}
			}
			externalCheckoutOptions = &ExternalCheckoutOptions{
				Type:            "git",
				User:            gitCheckoutOptions.User,
				UserName:        gitCheckoutOptions.UserName,
				Email:           gitCheckoutOptions.Email,
				Host:            gitCheckoutOptions.Host,
				Path:            gitCheckoutOptions.Path,
				Branch:          gitCheckoutOptions.Branch,
				CommitID:        gitCheckoutOptions.CommitID,
				CommitMessage:   gitCheckoutOptions.CommitMessage,
				SecurityOptions: externalSecurityOptions,
			}
			return nil
		},
		func(githubCheckoutOptions *GithubCheckoutOptions) error {
			var externalSecurityOptions *ExternalSecurityOptions
			var err error
			if githubCheckoutOptions.SecurityOptions != nil {
				externalSecurityOptions, err = ConvertSecurityOptions(githubCheckoutOptions.SecurityOptions)
				if err != nil {
					return err
				}
			}
			externalCheckoutOptions = &ExternalCheckoutOptions{
				Type:            "github",
				User:            githubCheckoutOptions.User,
				UserName:        githubCheckoutOptions.UserName,
				Email:           githubCheckoutOptions.Email,
				Repository:      githubCheckoutOptions.Repository,
				Branch:          githubCheckoutOptions.Branch,
				CommitID:        githubCheckoutOptions.CommitID,
				CommitMessage:   githubCheckoutOptions.CommitMessage,
				SecurityOptions: externalSecurityOptions,
			}
			return nil
		},
		func(hgCheckoutOptions *HgCheckoutOptions) error {
			var externalSecurityOptions *ExternalSecurityOptions
			var err error
			if hgCheckoutOptions.SecurityOptions != nil {
				externalSecurityOptions, err = ConvertSecurityOptions(hgCheckoutOptions.SecurityOptions)
				if err != nil {
					return err
				}
			}
			externalCheckoutOptions = &ExternalCheckoutOptions{
				Type:            "hg",
				User:            hgCheckoutOptions.User,
				UserName:        hgCheckoutOptions.UserName,
				Email:           hgCheckoutOptions.Email,
				Host:            hgCheckoutOptions.Host,
				Path:            hgCheckoutOptions.Path,
				ChangesetID:     hgCheckoutOptions.ChangesetID,
				CommitMessage:   hgCheckoutOptions.CommitMessage,
				SecurityOptions: externalSecurityOptions,
			}
			return nil
		},
		func(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error {
			var externalSecurityOptions *ExternalSecurityOptions
			var err error
			if bitbucketGitCheckoutOptions.SecurityOptions != nil {
				externalSecurityOptions, err = ConvertSecurityOptions(bitbucketGitCheckoutOptions.SecurityOptions)
				if err != nil {
					return err
				}
			}
			externalCheckoutOptions = &ExternalCheckoutOptions{
				Type:            "bitbucketGit",
				User:            bitbucketGitCheckoutOptions.User,
				UserName:        bitbucketGitCheckoutOptions.UserName,
				Email:           bitbucketGitCheckoutOptions.Email,
				Repository:      bitbucketGitCheckoutOptions.Repository,
				Branch:          bitbucketGitCheckoutOptions.Branch,
				CommitID:        bitbucketGitCheckoutOptions.CommitID,
				CommitMessage:   bitbucketGitCheckoutOptions.CommitMessage,
				SecurityOptions: externalSecurityOptions,
			}
			return nil
		},
		func(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error {
			var externalSecurityOptions *ExternalSecurityOptions
			var err error
			if bitbucketHgCheckoutOptions.SecurityOptions != nil {
				externalSecurityOptions, err = ConvertSecurityOptions(bitbucketHgCheckoutOptions.SecurityOptions)
				if err != nil {
					return err
				}
			}
			externalCheckoutOptions = &ExternalCheckoutOptions{
				Type:            "bitbucketHg",
				User:            bitbucketHgCheckoutOptions.User,
				UserName:        bitbucketHgCheckoutOptions.UserName,
				Email:           bitbucketHgCheckoutOptions.Email,
				Repository:      bitbucketHgCheckoutOptions.Repository,
				ChangesetID:     bitbucketHgCheckoutOptions.ChangesetID,
				CommitMessage:   bitbucketHgCheckoutOptions.CommitMessage,
				SecurityOptions: externalSecurityOptions,
			}
			return nil
		},
	); switchErr != nil {
		return nil, switchErr
	}
	return externalCheckoutOptions, nil
}

func convertSecurityOptions(securityOptions SecurityOptions) (*ExternalSecurityOptions, error) {
	var externalSecurityOptions *ExternalSecurityOptions
	if switchErr := SecurityOptionsSwitch(
		securityOptions,
		func(sshSecurityOptions *SSHSecurityOptions) error {
			var privateKeyString string
			if sshSecurityOptions.PrivateKey != nil {
				data, err := ioutil.ReadAll(sshSecurityOptions.PrivateKey)
				if err != nil {
					return err
				}
				privateKeyString = string(data)
			}
			externalSecurityOptions = &ExternalSecurityOptions{
				Type: "ssh",
				StrictHostKeyChecking: sshSecurityOptions.StrictHostKeyChecking,
				PrivateKey:            privateKeyString,
			}
			return nil
		},
		func(accessTokenSecurityOptions *AccessTokenSecurityOptions) error {
			externalSecurityOptions = &ExternalSecurityOptions{
				Type:        "accessToken",
				AccessToken: accessTokenSecurityOptions.AccessToken,
			}
			return nil
		},
	); switchErr != nil {
		return nil, switchErr
	}
	return externalSecurityOptions, nil
}

func convertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	var securityOptions SecurityOptions
	if externalCheckoutOptions.SecurityOptions != nil {
		securityOptionsType, err := SecurityOptionsTypeOf(externalCheckoutOptions.SecurityOptions.Type)
		if err != nil {
			return nil, err
		}
		securityOptions, err = securityOptionsType.NewSecurityOptions(
			func() (*SSHSecurityOptions, error) {
				var privateKey bytes.Buffer
				if _, err := privateKey.WriteString(externalCheckoutOptions.SecurityOptions.PrivateKey); err != nil {
					return nil, err
				}
				return &SSHSecurityOptions{
					StrictHostKeyChecking: externalCheckoutOptions.SecurityOptions.StrictHostKeyChecking,
					PrivateKey:            &privateKey,
				}, nil
			},
			func() (*AccessTokenSecurityOptions, error) {
				return &AccessTokenSecurityOptions{
					AccessToken: externalCheckoutOptions.SecurityOptions.AccessToken,
				}, nil
			},
		)
		if err != nil {
			return nil, err
		}
	}
	checkoutOptionsType, err := CheckoutOptionsTypeOf(externalCheckoutOptions.Type)
	if err != nil {
		return nil, err
	}
	return checkoutOptionsType.NewCheckoutOptions(
		func() (*GitCheckoutOptions, error) {
			return &GitCheckoutOptions{
				User:            externalCheckoutOptions.User,
				UserName:        externalCheckoutOptions.UserName,
				Email:           externalCheckoutOptions.Email,
				Host:            externalCheckoutOptions.Host,
				Path:            externalCheckoutOptions.Path,
				Branch:          externalCheckoutOptions.Branch,
				CommitID:        externalCheckoutOptions.CommitID,
				CommitMessage:   externalCheckoutOptions.CommitMessage,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*GithubCheckoutOptions, error) {
			return &GithubCheckoutOptions{
				User:            externalCheckoutOptions.User,
				UserName:        externalCheckoutOptions.UserName,
				Email:           externalCheckoutOptions.Email,
				Repository:      externalCheckoutOptions.Repository,
				Branch:          externalCheckoutOptions.Branch,
				CommitID:        externalCheckoutOptions.CommitID,
				CommitMessage:   externalCheckoutOptions.CommitMessage,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*HgCheckoutOptions, error) {
			return &HgCheckoutOptions{
				User:            externalCheckoutOptions.User,
				UserName:        externalCheckoutOptions.UserName,
				Email:           externalCheckoutOptions.Email,
				Host:            externalCheckoutOptions.Host,
				Path:            externalCheckoutOptions.Path,
				ChangesetID:     externalCheckoutOptions.ChangesetID,
				CommitMessage:   externalCheckoutOptions.CommitMessage,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*BitbucketGitCheckoutOptions, error) {
			return &BitbucketGitCheckoutOptions{
				User:            externalCheckoutOptions.User,
				UserName:        externalCheckoutOptions.UserName,
				Email:           externalCheckoutOptions.Email,
				Repository:      externalCheckoutOptions.Repository,
				Branch:          externalCheckoutOptions.Branch,
				CommitID:        externalCheckoutOptions.CommitID,
				CommitMessage:   externalCheckoutOptions.CommitMessage,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*BitbucketHgCheckoutOptions, error) {
			return &BitbucketHgCheckoutOptions{
				User:            externalCheckoutOptions.User,
				UserName:        externalCheckoutOptions.UserName,
				Email:           externalCheckoutOptions.Email,
				Repository:      externalCheckoutOptions.Repository,
				ChangesetID:     externalCheckoutOptions.ChangesetID,
				CommitMessage:   externalCheckoutOptions.CommitMessage,
				SecurityOptions: securityOptions,
			}, nil
		},
	)
}

func checkout(
	checkoutOptions CheckoutOptions,
	absolutePath string,
) error {
	if err := validateCheckoutOptions(checkoutOptions); err != nil {
		return err
	}
	execClientProvider, err := exec.NewClientProvider(&exec.OsExecOptions{})
	if err != nil {
		return err
	}
	baseDir, path := filepath.Split(absolutePath)
	executor, err := exec.NewOsExecutor(baseDir)
	if err != nil {
		return err
	}
	return CheckoutOptionsSwitch(
		checkoutOptions,
		func(gitCheckoutOptions *GitCheckoutOptions) error {
			return checkoutGit(execClientProvider, gitCheckoutOptions, executor, path)
		},
		func(githubCheckoutOptions *GithubCheckoutOptions) error {
			return checkoutGithub(execClientProvider, githubCheckoutOptions, executor, path)
		},
		func(hgCheckoutOptions *HgCheckoutOptions) error {
			return checkoutHg(execClientProvider, hgCheckoutOptions, executor, path)
		},
		func(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error {
			return checkoutBitbucketGit(execClientProvider, bitbucketGitCheckoutOptions, executor, path)
		},
		func(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error {
			return checkoutBitbucketHg(execClientProvider, bitbucketHgCheckoutOptions, executor, path)
		},
	)
}

func checkoutGit(
	execClientProvider exec.ClientProvider,
	gitCheckoutOptions *GitCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string
	var client exec.Client
	var err error
	if gitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSSHCommand(execClientProvider, gitCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
		if client != nil {
			defer func() {
				if err := client.Destroy(); err != nil && retErr == nil {
					retErr = err
				}
			}()
		}
	}
	url, err := getGitURL(gitCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitWithExecutor(executor, sshCommand, url, gitCheckoutOptions.Branch, gitCheckoutOptions.CommitID, path)
}

func checkoutGithub(
	execClientProvider exec.ClientProvider,
	githubCheckoutOptions *GithubCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string
	var client exec.Client
	var err error
	if githubCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSSHCommand(execClientProvider, githubCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
		if client != nil {
			defer func() {
				if err := client.Destroy(); err != nil && retErr == nil {
					retErr = err
				}
			}()
		}
	}
	url, err := getGithubURL(githubCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitWithExecutor(executor, sshCommand, url, githubCheckoutOptions.Branch, githubCheckoutOptions.CommitID, path)
}

func checkoutHg(
	execClientProvider exec.ClientProvider,
	hgCheckoutOptions *HgCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string
	var client exec.Client
	var err error
	if hgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSSHCommand(execClientProvider, hgCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
		if client != nil {
			defer func() {
				if err := client.Destroy(); err != nil && retErr == nil {
					retErr = err
				}
			}()
		}
	}
	url, err := getHgURL(hgCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutHgWithExecutor(executor, sshCommand, url, hgCheckoutOptions.ChangesetID, path)
}

func checkoutBitbucketGit(
	execClientProvider exec.ClientProvider,
	bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string
	var client exec.Client
	var err error
	if bitbucketGitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSSHCommand(execClientProvider, bitbucketGitCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
		if client != nil {
			defer func() {
				if err := client.Destroy(); err != nil && retErr == nil {
					retErr = err
				}
			}()
		}
	}
	url, err := getBitbucketGitURL(bitbucketGitCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitWithExecutor(executor, sshCommand, url, bitbucketGitCheckoutOptions.Branch, bitbucketGitCheckoutOptions.CommitID, path)
}

func checkoutBitbucketHg(
	execClientProvider exec.ClientProvider,
	bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string
	var client exec.Client
	var err error
	if bitbucketHgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSSHCommand(execClientProvider, bitbucketHgCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
		if client != nil {
			defer func() {
				if err := client.Destroy(); err != nil && retErr == nil {
					retErr = err
				}
			}()
		}
	}
	url, err := getBitbucketHgURL(bitbucketHgCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutHgWithExecutor(executor, sshCommand, url, bitbucketHgCheckoutOptions.ChangesetID, path)
}

func getSSHCommand(execClientProvider exec.ClientProvider, securityOptions SecurityOptions) (string, exec.Client, error) {
	var sshCommand string
	var client exec.Client
	var err error
	if err = SecurityOptionsSwitch(
		securityOptions,
		func(sshSecurityOptions *SSHSecurityOptions) error {
			sshCommandArgs := []string{"ssh", "-o"}
			if sshSecurityOptions.StrictHostKeyChecking {
				sshCommandArgs = append(sshCommandArgs, "StrictHostKeyChecking=yes")
			} else {
				sshCommandArgs = append(sshCommandArgs, "StrictHostKeyChecking=no")
			}
			if sshSecurityOptions.PrivateKey != nil {
				client, err = execClientProvider.NewTempDirClient()
				if err != nil {
					return err
				}
				writeFile, err := client.Create("id_rsa")
				if err != nil {
					return err
				}
				data, err := ioutil.ReadAll(sshSecurityOptions.PrivateKey)
				if err != nil {
					return err
				}
				_, err = writeFile.Write(data)
				if err != nil {
					return err
				}
				err = writeFile.Chmod(0400)
				if err != nil {
					return err
				}
				sshCommandArgs = append(sshCommandArgs, "-i", client.Join(client.DirPath(), "id_rsa"))
				sshCommand = strings.Join(sshCommandArgs, " ")
			}
			return nil
		},
		func(accessTokenSecurityOptions *AccessTokenSecurityOptions) error {
			return nil
		},
	); err != nil {
		return "", nil, err
	}
	return sshCommand, client, nil
}

func getGitURL(gitCheckoutOptions *GitCheckoutOptions) (string, error) {
	if gitCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyURL(
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	if gitCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSSHURL(
			"",
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getGithubURL(githubCheckoutOptions *GithubCheckoutOptions) (string, error) {
	if githubCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyURL(
			"git",
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSSHURL(
			"",
			"git",
			"github.com",
			joinStrings(":", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeAccessToken {
		return getAccessTokenURL(
			(githubCheckoutOptions.SecurityOptions.(*AccessTokenSecurityOptions)).AccessToken,
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getHgURL(hgCheckoutOptions *HgCheckoutOptions) (string, error) {
	if hgCheckoutOptions.SecurityOptions == nil || hgCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSSHURL(
			"ssh://",
			hgCheckoutOptions.User,
			hgCheckoutOptions.Host,
			hgCheckoutOptions.Path,
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getBitbucketGitURL(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) (string, error) {
	if bitbucketGitCheckoutOptions.SecurityOptions == nil {
		return getSSHURL(
			"ssh://",
			"git",
			"bitbucket.org",
			joinStrings(":", bitbucketGitCheckoutOptions.User, "/", bitbucketGitCheckoutOptions.Repository, ".git"),
		), nil
	}
	if bitbucketGitCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSSHURL(
			"",
			"git",
			"bitbucket.org",
			joinStrings(":", bitbucketGitCheckoutOptions.User, "/", bitbucketGitCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getBitbucketHgURL(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) (string, error) {
	if bitbucketHgCheckoutOptions.SecurityOptions == nil || bitbucketHgCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSSHURL(
			"ssh://",
			"hg",
			"bitbucket.org",
			joinStrings("/", bitbucketHgCheckoutOptions.User, "/", bitbucketHgCheckoutOptions.Repository),
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

// TODO(pedge): user?
func getGitReadOnlyURL(user string, host string, path string) string {
	return joinStrings("git://", host, path)
}

func getSSHURL(base string, user string, host string, path string) string {
	return joinStrings(base, user, "@", host, path)
}

func getAccessTokenURL(accessToken string, host string, path string) string {
	return joinStrings("https://", accessToken, ":x-oauth-basic@", host, path)
}

func checkoutGitWithExecutor(
	executor exec.Executor,
	gitSSHCommand string,
	url string,
	branch string,
	commitID string,
	path string,
) error {
	var cloneStderr bytes.Buffer
	cmd := exec.Cmd{
		// TODO(peter): if the commit id is more than 50 back, the checkout will fail
		Args:   []string{"git", "clone", "--branch", branch, "--depth", "50", "--recursive", url, path},
		Stderr: &cloneStderr,
	}
	if gitSSHCommand != "" {
		cmd.Env = []string{"GIT_SSH_COMMAND=" + gitSSHCommand}
	}
	if err := executor.Execute(&cmd)(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotClone: %v %v", err.Error(), cloneStderr.String())
	}
	var checkoutStderr bytes.Buffer
	if err := executor.Execute(
		&exec.Cmd{
			Args:   []string{"git", "checkout", "-f", commitID},
			SubDir: path,
			Stderr: &checkoutStderr,
		},
	)(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotCheckout: %v %v", err.Error(), checkoutStderr.String())
	}
	return nil
}

func checkoutHgWithExecutor(
	executor exec.Executor,
	sshCommand string,
	url string,
	changesetID string,
	path string,
) error {
	args := []string{"hg", "clone", url, path}
	if sshCommand != "" {
		args = []string{"hg", "clone", "--ssh", sshCommand, url, path}
	}
	var cloneStderr bytes.Buffer
	if err := executor.Execute(
		&exec.Cmd{
			Args:   args,
			Stderr: &cloneStderr,
		},
	)(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotClone: %v %v", err.Error(), cloneStderr.String())
	}
	var updateStderr bytes.Buffer
	if err := executor.Execute(
		&exec.Cmd{
			Args:   []string{"hg", "update", "--cwd", path, changesetID},
			Stderr: &updateStderr,
		},
	)(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotUpdate: %v %v", err.Error(), updateStderr.String())
	}
	return nil
}

func joinStrings(elems ...string) string {
	return strings.Join(elems, "")
}

type validationError struct {
	errorType ValidationErrorType
	tags      map[string]string
}

func newValidationError(errorType ValidationErrorType, tags map[string]string) *validationError {
	if tags == nil {
		tags = make(map[string]string)
	}
	return &validationError{errorType, tags}
}

func (v *validationError) Error() string {
	return fmt.Sprintf("%v %v", v.errorType, v.tags)
}

func (v *validationError) Type() ValidationErrorType {
	return v.errorType
}

func validateCheckoutOptions(checkoutOptions CheckoutOptions) error {
	return CheckoutOptionsSwitch(
		checkoutOptions,
		validateGitCheckoutOptions,
		validateGithubCheckoutOptions,
		validateHgCheckoutOptions,
		validateBitbucketGitCheckoutOptions,
		validateBitbucketHgCheckoutOptions,
	)
}

func validateGitCheckoutOptions(gitCheckoutOptions *GitCheckoutOptions) error {
	if gitCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "User")
	}
	if gitCheckoutOptions.Host == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "Host")
	}
	if gitCheckoutOptions.Path == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "Path")
	}
	if gitCheckoutOptions.Branch == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "Branch")
	}
	if gitCheckoutOptions.CommitID == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "CommitID")
	}
	if gitCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(gitCheckoutOptions.SecurityOptions, CheckoutOptionsTypeGit, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateGithubCheckoutOptions(githubCheckoutOptions *GithubCheckoutOptions) error {
	if githubCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "User")
	}
	if githubCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "Repository")
	}
	if githubCheckoutOptions.Branch == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "Branch")
	}
	if githubCheckoutOptions.CommitID == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "CommitID")
	}
	if githubCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(githubCheckoutOptions.SecurityOptions, CheckoutOptionsTypeGithub, SecurityOptionsTypeSsh, SecurityOptionsTypeAccessToken); err != nil {
			return nil
		}
	}
	return nil
}

func validateHgCheckoutOptions(hgCheckoutOptions *HgCheckoutOptions) error {
	if hgCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "User")
	}
	if hgCheckoutOptions.Host == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "Host")
	}
	if hgCheckoutOptions.Path == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "Path")
	}
	if hgCheckoutOptions.ChangesetID == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "ChangesetID")
	}
	if hgCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(hgCheckoutOptions.SecurityOptions, CheckoutOptionsTypeHg, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateBitbucketGitCheckoutOptions(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error {
	if bitbucketGitCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "User")
	}
	if bitbucketGitCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "Repository")
	}
	if bitbucketGitCheckoutOptions.Branch == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "Branch")
	}
	if bitbucketGitCheckoutOptions.CommitID == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "CommitID")
	}
	if bitbucketGitCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(bitbucketGitCheckoutOptions.SecurityOptions, CheckoutOptionsTypeBitbucketGit, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateBitbucketHgCheckoutOptions(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error {
	if bitbucketHgCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "User")
	}
	if bitbucketHgCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "Repository")
	}
	if bitbucketHgCheckoutOptions.ChangesetID == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "ChangesetID")
	}
	if bitbucketHgCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(bitbucketHgCheckoutOptions.SecurityOptions, CheckoutOptionsTypeBitbucketHg, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateSecurityOptions(securityOptions SecurityOptions, checkoutType CheckoutOptionsType, allowedTypes ...SecurityOptionsType) error {
	if !isAllowedSecurityOptionsType(securityOptions.Type(), allowedTypes) {
		return newValidationErrorSecurityNotImplementedForCheckoutOptionsType(securityOptions.Type().String(), checkoutType.String())
	}
	return SecurityOptionsSwitch(
		securityOptions,
		validateSSHSecurityOptions,
		validateAccessTokenSecurityOptions,
	)
}

func validateSSHSecurityOptions(sshSecurityOptions *SSHSecurityOptions) error {
	return nil
}

func validateAccessTokenSecurityOptions(accessTokenSecurityOptions *AccessTokenSecurityOptions) error {
	if accessTokenSecurityOptions.AccessToken == "" {
		return newValidationErrorRequiredFieldMissing("AccessTokenSecurityOptions", "AccessToken")
	}
	return nil
}

func isAllowedSecurityOptionsType(securityType SecurityOptionsType, allowedTypes []SecurityOptionsType) bool {
	for _, allowedType := range allowedTypes {
		if securityType == allowedType {
			return true
		}
	}
	return false
}

func newValidationErrorRequiredFieldMissing(objectType string, fieldPath ...string) ValidationError {
	return newValidationError(ValidationErrorTypeRequiredFieldMissing, map[string]string{"type": objectType, "fieldPath": strings.Join(fieldPath, ".")})
}

func newValidationErrorFieldShouldNotBeSet(objectType string, fieldPath ...string) ValidationError {
	return newValidationError(ValidationErrorTypeFieldShouldNotBeSet, map[string]string{"type": objectType, "fieldPath": strings.Join(fieldPath, ".")})
}

func newValidationErrorSecurityNotImplementedForCheckoutOptionsType(securityType string, checkoutType string) ValidationError {
	return newValidationError(ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType, map[string]string{"securityType": securityType, "checkoutType": checkoutType})
}
