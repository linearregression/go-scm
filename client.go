package scm

import (
	"bytes"
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
	ignoreGitCheckoutFilePatterns = []string{
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

func newClient(execClientProvider exec.ClientProvider, clientOptions *ClientOptions) *client {
	return &client{execClientProvider, clientOptions}
}

type client struct {
	exec.ClientProvider
	clientOptions *ClientOptions
}

func (this *client) CheckoutTarball(checkoutOptions CheckoutOptions) (io.Reader, error) {
	executorReadFileManager, err := this.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	if err := this.Checkout(checkoutOptions, executorReadFileManager, clonePath); err != nil {
		return nil, err
	}
	return this.tarAndDestroy(executorReadFileManager, checkoutOptions, clonePath)
}

func (this *client) Checkout(checkoutOptions CheckoutOptions, executor exec.Executor, path string) error {
	if err := validateCheckoutOptions(checkoutOptions); err != nil {
		return err
	}
	return CheckoutOptionsSwitch(
		checkoutOptions,
		func(gitCheckoutOptions *GitCheckoutOptions) error {
			return this.checkoutGit(gitCheckoutOptions, executor, path)
		},
		func(githubCheckoutOptions *GithubCheckoutOptions) error {
			return this.checkoutGithub(githubCheckoutOptions, executor, path)
		},
		func(hgCheckoutOptions *HgCheckoutOptions) error {
			return this.checkoutHg(hgCheckoutOptions, executor, path)
		},
		func(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error {
			return this.checkoutBitbucketGit(bitbucketGitCheckoutOptions, executor, path)
		},
		func(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error {
			return this.checkoutBitbucketHg(bitbucketHgCheckoutOptions, executor, path)
		},
	)
}

func (this *client) checkoutGit(gitCheckoutOptions *GitCheckoutOptions, executor exec.Executor, path string) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if gitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(gitCheckoutOptions.SecurityOptions)
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

func (this *client) checkoutGithub(githubCheckoutOptions *GithubCheckoutOptions, executor exec.Executor, path string) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if githubCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(githubCheckoutOptions.SecurityOptions)
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

func (this *client) checkoutHg(hgCheckoutOptions *HgCheckoutOptions, executor exec.Executor, path string) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if hgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(hgCheckoutOptions.SecurityOptions)
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

func (this *client) checkoutBitbucketGit(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions, executor exec.Executor, path string) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if bitbucketGitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(bitbucketGitCheckoutOptions.SecurityOptions)
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

func (this *client) checkoutBitbucketHg(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions, executor exec.Executor, path string) (retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if bitbucketHgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(bitbucketHgCheckoutOptions.SecurityOptions)
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

func (this *client) tarAndDestroy(executorReadFileManager exec.ExecutorReadFileManager, checkoutOptions CheckoutOptions, path string) (io.Reader, error) {
	var reader io.Reader
	var err error
	if this.clientOptions.IgnoreCheckoutFiles {
		ignoreCheckoutFilePatterns, err := this.ignoreCheckoutFilePatterns(checkoutOptions)
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

func (this *client) ignoreCheckoutFilePatterns(checkoutOptions CheckoutOptions) ([]string, error) {
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

func (this *client) getSshCommand(securityOptions SecurityOptions) (string, exec.Client, error) {
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
				client, err = this.NewTempDirClient()
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
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutOptionsType(
			gitCheckoutOptions.Type().String(),
			gitCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
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
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutOptionsType(
			githubCheckoutOptions.Type().String(),
			githubCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
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
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutOptionsType(
			hgCheckoutOptions.Type().String(),
			hgCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
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
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutOptionsType(
			bitbucketGitCheckoutOptions.Type().String(),
			bitbucketGitCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
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
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutOptionsType(
			bitbucketHgCheckoutOptions.Type().String(),
			bitbucketHgCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
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

func tarFiles(readFileManager exec.ReadFileManager, ignoreCheckoutFilePatterns []string, path string) (io.Reader, error) {
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
	return tarexec.NewTarClient(readFileManager, nil).Tar(fileList, path)
}

func fileMatches(readFileManager exec.ReadFileManager, patterns []string, path string, basePath string) (bool, error) {
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
