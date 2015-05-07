package scm

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/peter-edge/go-exec"
	tarexec "github.com/peter-edge/go-tar/exec"
)

const (
	clonePath = "clone"
)

var (
	ValidationErrorTypeRequiredFieldMissing                         ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                          ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType ValidationErrorType = "SecurityNotImplementedForCheckoutOptionsType"

	errorSecurityNotImplementedForCheckoutOptionsType = errors.New("SecurityNotImplementedForCheckoutOptionsType")
	ignoreGitCheckoutFilePatterns                     = []string{
		".git",
		".gitignore",
	}
	ignoreHgCheckoutFilePatterns = []string{
		".hg",
		".hgignore",
		".hgsigs",
		".hgtags",
	}
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
	Host            string
	Path            string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions github 1
type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions hg 2
type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions bitbucketGit 3
type BitbucketGitCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions bitbucketHg 4
type BitbucketHgCheckoutOptions struct {
	User            string
	Repository      string
	ChangesetId     string
	SecurityOptions SecurityOptions
}

// @gen-enumtype SecurityOptions ssh 0
type SshSecurityOptions struct {
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
	Host            string                   `json:"host,omitempty" yaml:"host,omitempty"`
	Path            string                   `json:"path,omitempty" yaml:"path,omitempty"`
	Repository      string                   `json:"repository,omitempty" yaml:"repository,omitempty"`
	Branch          string                   `json:"branch,omitempty" yaml:"branch,omitempty"`
	CommitId        string                   `json:"commit_id,omitempty" yaml:"commit_id,omitempty"`
	ChangesetId     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
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
	execClientProvider exec.ClientProvider,
	checkoutOptions CheckoutOptions,
	executor exec.Executor,
	path string,
) error {
	return checkout(
		execClientProvider,
		checkoutOptions,
		executor,
		path,
	)
}

func CheckoutTarball(
	execClientProvider exec.ClientProvider,
	checkoutOptions CheckoutOptions,
	ignoreCheckoutFiles bool,
) (io.Reader, error) {
	return checkoutTarball(
		execClientProvider,
		checkoutOptions,
		ignoreCheckoutFiles,
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
				Host:            gitCheckoutOptions.Host,
				Path:            gitCheckoutOptions.Path,
				Branch:          gitCheckoutOptions.Branch,
				CommitId:        gitCheckoutOptions.CommitId,
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
				Repository:      githubCheckoutOptions.Repository,
				Branch:          githubCheckoutOptions.Branch,
				CommitId:        githubCheckoutOptions.CommitId,
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
				Host:            hgCheckoutOptions.Host,
				Path:            hgCheckoutOptions.Path,
				ChangesetId:     hgCheckoutOptions.ChangesetId,
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
				Repository:      bitbucketGitCheckoutOptions.Repository,
				Branch:          bitbucketGitCheckoutOptions.Branch,
				CommitId:        bitbucketGitCheckoutOptions.CommitId,
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
				Repository:      bitbucketHgCheckoutOptions.Repository,
				ChangesetId:     bitbucketHgCheckoutOptions.ChangesetId,
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
		func(sshSecurityOptions *SshSecurityOptions) error {
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
			func() (*SshSecurityOptions, error) {
				var privateKey bytes.Buffer
				if _, err := privateKey.WriteString(externalCheckoutOptions.SecurityOptions.PrivateKey); err != nil {
					return nil, err
				}
				return &SshSecurityOptions{
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
				Host:            externalCheckoutOptions.Host,
				Path:            externalCheckoutOptions.Path,
				Branch:          externalCheckoutOptions.Branch,
				CommitId:        externalCheckoutOptions.CommitId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*GithubCheckoutOptions, error) {
			return &GithubCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Repository:      externalCheckoutOptions.Repository,
				Branch:          externalCheckoutOptions.Branch,
				CommitId:        externalCheckoutOptions.CommitId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*HgCheckoutOptions, error) {
			return &HgCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Host:            externalCheckoutOptions.Host,
				Path:            externalCheckoutOptions.Path,
				ChangesetId:     externalCheckoutOptions.ChangesetId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*BitbucketGitCheckoutOptions, error) {
			return &BitbucketGitCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Repository:      externalCheckoutOptions.Repository,
				Branch:          externalCheckoutOptions.Branch,
				CommitId:        externalCheckoutOptions.CommitId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*BitbucketHgCheckoutOptions, error) {
			return &BitbucketHgCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Repository:      externalCheckoutOptions.Repository,
				ChangesetId:     externalCheckoutOptions.ChangesetId,
				SecurityOptions: securityOptions,
			}, nil
		},
	)
}

func checkout(
	execClientProvider exec.ClientProvider,
	checkoutOptions CheckoutOptions,
	executor exec.Executor,
	path string,
) error {
	if err := validateCheckoutOptions(checkoutOptions); err != nil {
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

func checkoutTarball(
	execClientProvider exec.ClientProvider,
	checkoutOptions CheckoutOptions,
	ignoreCheckoutFiles bool,
) (io.Reader, error) {
	executorReadFileManager, err := execClientProvider.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	if err := checkout(execClientProvider, checkoutOptions, executorReadFileManager, clonePath); err != nil {
		return nil, err
	}
	return tarAndDestroy(executorReadFileManager, checkoutOptions, clonePath, ignoreCheckoutFiles)
}

func checkoutGit(
	execClientProvider exec.ClientProvider,
	gitCheckoutOptions *GitCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if gitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSshCommand(execClientProvider, gitCheckoutOptions.SecurityOptions)
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
	url, err := getGitUrl(gitCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitWithExecutor(executor, sshCommand, url, gitCheckoutOptions.Branch, gitCheckoutOptions.CommitId, path)
}

func checkoutGithub(
	execClientProvider exec.ClientProvider,
	githubCheckoutOptions *GithubCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if githubCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSshCommand(execClientProvider, githubCheckoutOptions.SecurityOptions)
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
	url, err := getGithubUrl(githubCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitWithExecutor(executor, sshCommand, url, githubCheckoutOptions.Branch, githubCheckoutOptions.CommitId, path)
}

func checkoutHg(
	execClientProvider exec.ClientProvider,
	hgCheckoutOptions *HgCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if hgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSshCommand(execClientProvider, hgCheckoutOptions.SecurityOptions)
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
	url, err := getHgUrl(hgCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutHgWithExecutor(executor, sshCommand, url, hgCheckoutOptions.ChangesetId, path)
}

func checkoutBitbucketGit(
	execClientProvider exec.ClientProvider,
	bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if bitbucketGitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSshCommand(execClientProvider, bitbucketGitCheckoutOptions.SecurityOptions)
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
	url, err := getBitbucketGitUrl(bitbucketGitCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitWithExecutor(executor, sshCommand, url, bitbucketGitCheckoutOptions.Branch, bitbucketGitCheckoutOptions.CommitId, path)
}

func checkoutBitbucketHg(
	execClientProvider exec.ClientProvider,
	bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions,
	executor exec.Executor,
	path string,
) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if bitbucketHgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = getSshCommand(execClientProvider, bitbucketHgCheckoutOptions.SecurityOptions)
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
	url, err := getBitbucketHgUrl(bitbucketHgCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutHgWithExecutor(executor, sshCommand, url, bitbucketHgCheckoutOptions.ChangesetId, path)
}

func getSshCommand(execClientProvider exec.ClientProvider, securityOptions SecurityOptions) (string, exec.Client, error) {
	var sshCommand string
	var client exec.Client
	var err error
	if err = SecurityOptionsSwitch(
		securityOptions,
		func(sshSecurityOptions *SshSecurityOptions) error {
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

func getGitUrl(gitCheckoutOptions *GitCheckoutOptions) (string, error) {
	if gitCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyUrl(
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	if gitCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSshUrl(
			"",
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getGithubUrl(githubCheckoutOptions *GithubCheckoutOptions) (string, error) {
	if githubCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyUrl(
			"git",
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSshUrl(
			"",
			"git",
			"github.com",
			joinStrings(":", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeAccessToken {
		return getAccessTokenUrl(
			(githubCheckoutOptions.SecurityOptions.(*AccessTokenSecurityOptions)).AccessToken,
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getHgUrl(hgCheckoutOptions *HgCheckoutOptions) (string, error) {
	if hgCheckoutOptions.SecurityOptions == nil || hgCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSshUrl(
			"ssh://",
			hgCheckoutOptions.User,
			hgCheckoutOptions.Host,
			hgCheckoutOptions.Path,
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getBitbucketGitUrl(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) (string, error) {
	if bitbucketGitCheckoutOptions.SecurityOptions == nil || bitbucketGitCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSshUrl(
			"ssh://",
			"git",
			"bitbucket.org",
			joinStrings(":", bitbucketGitCheckoutOptions.User, "/", bitbucketGitCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

func getBitbucketHgUrl(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) (string, error) {
	if bitbucketHgCheckoutOptions.SecurityOptions == nil || bitbucketHgCheckoutOptions.SecurityOptions.Type() == SecurityOptionsTypeSsh {
		return getSshUrl(
			"ssh://",
			"hg",
			"bitbucket.org",
			joinStrings("/", bitbucketHgCheckoutOptions.User, "/", bitbucketHgCheckoutOptions.Repository),
		), nil
	}
	return "", errorSecurityNotImplementedForCheckoutOptionsType
}

// TODO(pedge): user?
func getGitReadOnlyUrl(user string, host string, path string) string {
	return joinStrings("git://", host, path)
}

func getSshUrl(base string, user string, host string, path string) string {
	return joinStrings(base, user, "@", host, path)
}

func getAccessTokenUrl(accessToken string, host string, path string) string {
	return joinStrings("https://", accessToken, ":x-oauth-basic@", host, path)
}

func checkoutGitWithExecutor(
	executor exec.Executor,
	gitSshCommand string,
	url string,
	branch string,
	commitId string,
	path string,
) error {
	var cloneStderr bytes.Buffer
	cmd := exec.Cmd{
		// TODO(peter): if the commit id is more than 50 back, the checkout will fail
		Args:   []string{"git", "clone", "--branch", branch, "--depth", "50", "--recursive", url, path},
		Stderr: &cloneStderr,
	}
	if gitSshCommand != "" {
		cmd.Env = []string{"GIT_SSH_COMMAND=" + gitSshCommand}
	}
	if err := executor.Execute(&cmd)(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotClone: %v %v", err.Error(), cloneStderr.String())
	}
	var checkoutStderr bytes.Buffer
	if err := executor.Execute(
		&exec.Cmd{
			Args:   []string{"git", "checkout", "-f", commitId},
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
	changesetId string,
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
			Args:   []string{"hg", "update", "--cwd", path, changesetId},
			Stderr: &updateStderr,
		},
	)(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotUpdate: %v %v", err.Error(), updateStderr.String())
	}
	return nil
}

func tarAndDestroy(
	executorReadFileManager exec.ExecutorReadFileManager,
	checkoutOptions CheckoutOptions,
	path string,
	ignoreCheckoutFiles bool,
) (io.Reader, error) {
	var reader io.Reader
	var err error
	if ignoreCheckoutFiles {
		ignoreCheckoutFilePatterns, err := ignoreCheckoutFilePatterns(checkoutOptions)
		if err != nil {
			return nil, err
		}
		reader, err = tarFiles(executorReadFileManager, ignoreCheckoutFilePatterns, path)
	} else {
		reader, err = tarFiles(executorReadFileManager, nil, path)
	}
	if err != nil {
		return nil, err
	}
	if err := executorReadFileManager.Destroy(); err != nil {
		return nil, err
	}
	return reader, nil
}

func tarFiles(
	readFileManager exec.ReadFileManager,
	ignoreCheckoutFilePatterns []string,
	path string,
) (io.Reader, error) {
	fileList, err := readFileManager.ListRegularFiles(path)
	if err != nil {
		return nil, err
	}
	if ignoreCheckoutFilePatterns != nil && len(ignoreCheckoutFilePatterns) > 0 {
		filterFileList := make([]string, 0)
		for _, file := range fileList {
			matches, err := fileMatches(readFileManager, ignoreCheckoutFilePatterns, file, path)
			if err != nil {
				return nil, err
			}
			if !matches {
				filterFileList = append(filterFileList, file)
			}
		}
		fileList = filterFileList
	}
	var buffer bytes.Buffer
	if err := tarexec.NewTarClient(readFileManager, nil).Tar(fileList, path, &buffer); err != nil {
		return nil, err
	}
	return &buffer, nil
}

func ignoreCheckoutFilePatterns(checkoutOptions CheckoutOptions) ([]string, error) {
	var ignoreCheckoutFilePatterns []string
	if err := CheckoutOptionsSwitch(
		checkoutOptions,
		func(gitCheckoutOptions *GitCheckoutOptions) error {
			ignoreCheckoutFilePatterns = ignoreGitCheckoutFilePatterns
			return nil
		},
		func(githubCheckoutOptions *GithubCheckoutOptions) error {
			ignoreCheckoutFilePatterns = ignoreGitCheckoutFilePatterns
			return nil
		},
		func(hgCheckoutOptions *HgCheckoutOptions) error {
			ignoreCheckoutFilePatterns = ignoreHgCheckoutFilePatterns
			return nil
		},
		func(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error {
			ignoreCheckoutFilePatterns = ignoreGitCheckoutFilePatterns
			return nil
		},
		func(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error {
			ignoreCheckoutFilePatterns = ignoreHgCheckoutFilePatterns
			return nil
		},
	); err != nil {
		return nil, err
	}
	return ignoreCheckoutFilePatterns, nil
}

func fileMatches(
	readFileManager exec.ReadFileManager,
	patterns []string,
	path string,
	basePath string,
) (bool, error) {
	for _, pattern := range patterns {
		if strings.HasPrefix(path, readFileManager.Join(basePath, pattern)) {
			return true, nil
		}
	}
	return false, nil
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
	if gitCheckoutOptions.CommitId == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "CommitId")
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
	if githubCheckoutOptions.CommitId == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "CommitId")
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
	if hgCheckoutOptions.ChangesetId == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "ChangesetId")
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
	if bitbucketGitCheckoutOptions.CommitId == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "CommitId")
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
	if bitbucketHgCheckoutOptions.ChangesetId == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "ChangesetId")
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
		validateSshSecurityOptions,
		validateAccessTokenSecurityOptions,
	)
}

func validateSshSecurityOptions(sshSecurityOptions *SshSecurityOptions) error {
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
