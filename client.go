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

func newClient(execClientProvider exec.ClientProvider, clientOptions *ClientOptions) *client {
	return &client{execClientProvider, clientOptions}
}

type client struct {
	exec.ClientProvider
	clientOptions *ClientOptions
}

func (this *client) CheckoutTarball(checkoutOptions CheckoutOptions) (io.Reader, error) {
	if err := validateCheckoutOptions(checkoutOptions); err != nil {
		return nil, err
	}
	switch checkoutOptions.Type() {
	case CheckoutTypeGit:
		return this.checkoutGitTarball(checkoutOptions.(*GitCheckoutOptions))
	case CheckoutTypeGithub:
		return this.checkoutGithubTarball(checkoutOptions.(*GithubCheckoutOptions))
	case CheckoutTypeHg:
		return this.checkoutHgTarball(checkoutOptions.(*HgCheckoutOptions))
	case CheckoutTypeBitbucket:
		return this.checkoutBitbucketTarball(checkoutOptions.(*BitbucketCheckoutOptions))
	default:
		return nil, newInternalError(newValidationErrorUnknownCheckoutType(checkoutOptions.Type().String()))
	}
}

func (this *client) checkoutGitTarball(gitCheckoutOptions *GitCheckoutOptions) (reader io.Reader, retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if gitCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(gitCheckoutOptions.SecurityOptions)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	return this.checkoutGit(sshCommand, url, gitCheckoutOptions.Branch, gitCheckoutOptions.CommitId)
}

func (this *client) checkoutGithubTarball(githubCheckoutOptions *GithubCheckoutOptions) (reader io.Reader, retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if githubCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(githubCheckoutOptions.SecurityOptions)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	return this.checkoutGit(sshCommand, url, githubCheckoutOptions.Branch, githubCheckoutOptions.CommitId)
}

func (this *client) checkoutHgTarball(hgCheckoutOptions *HgCheckoutOptions) (reader io.Reader, retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if hgCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(hgCheckoutOptions.SecurityOptions)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	return this.checkoutHg(sshCommand, url, hgCheckoutOptions.ChangesetId)
}

func (this *client) checkoutBitbucketTarball(bitbucketCheckoutOptions *BitbucketCheckoutOptions) (reader io.Reader, retErr error) {
	var sshCommand string = ""
	var client exec.Client
	var err error
	if bitbucketCheckoutOptions.SecurityOptions != nil {
		sshCommand, client, err = this.getSshCommand(bitbucketCheckoutOptions.SecurityOptions)
		if err != nil {
			return nil, err
		}
		if client != nil {
			defer func() {
				if err := client.Destroy(); err != nil && retErr == nil {
					retErr = err
				}
			}()
		}
	}
	url, err := getBitbucketUrl(bitbucketCheckoutOptions)
	if err != nil {
		return nil, err
	}
	switch bitbucketCheckoutOptions.BitbucketType {
	case BitbucketTypeGit:
		return this.checkoutGit(sshCommand, url, bitbucketCheckoutOptions.Branch, bitbucketCheckoutOptions.CommitId)
	case BitbucketTypeHg:
		return this.checkoutHg(sshCommand, url, bitbucketCheckoutOptions.ChangesetId)
	default:
		return nil, newInternalError(newValidationErrorUnknownBitbucketType(bitbucketCheckoutOptions.BitbucketType.String()))
	}
}

func (this *client) checkoutGit(gitSshCommand string, url string, branch string, commitId string) (io.Reader, error) {
	client, err := this.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	err = checkoutGitWithExecutor(client, gitSshCommand, url, branch, commitId, clonePath)
	if err != nil {
		return nil, err
	}
	return this.tarAndDestroy(client, ignoreGitCheckoutFilePatterns(), clonePath)
}

func (this *client) checkoutHg(sshCommand string, url string, changesetId string) (io.Reader, error) {
	client, err := this.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	err = checkoutHgWithExecutor(client, sshCommand, url, changesetId, clonePath)
	if err != nil {
		return nil, err
	}
	return this.tarAndDestroy(client, ignoreHgCheckoutFilePatterns(), clonePath)
}

func (this *client) tarAndDestroy(executorReadFileManager exec.ExecutorReadFileManager, ignoreCheckoutFilePatterns []string, path string) (io.Reader, error) {
	var reader io.Reader
	var err error
	if this.clientOptions.IgnoreCheckoutFiles {
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

func (this *client) getSshCommand(securityOptions SecurityOptions) (string, exec.Client, error) {
	if securityOptions.Type() != SecurityTypeSsh {
		return "", nil, nil
	}
	sshSecurityOptions := securityOptions.(*SshSecurityOptions)

	sshCommand := []string{"ssh", "-o"}
	if sshSecurityOptions.StrictHostKeyChecking {
		sshCommand = append(sshCommand, "StrictHostKeyChecking=yes")
	} else {
		sshCommand = append(sshCommand, "StrictHostKeyChecking=no")
	}
	var client exec.Client
	if sshSecurityOptions.PrivateKey != nil {
		client, err := this.NewTempDirClient()
		if err != nil {
			return "", nil, err
		}
		writeFile, err := client.Create("id_rsa")
		if err != nil {
			return "", nil, err
		}
		data, err := ioutil.ReadAll(sshSecurityOptions.PrivateKey)
		if err != nil {
			return "", nil, err
		}
		_, err = writeFile.Write(data)
		if err != nil {
			return "", nil, err
		}
		err = writeFile.Chmod(0400)
		if err != nil {
			return "", nil, err
		}
		sshCommand = append(sshCommand, "-i", client.Join(client.DirPath(), "id_rsa"))
	}
	return strings.Join(sshCommand, " "), client, nil
}

func getGitUrl(gitCheckoutOptions *GitCheckoutOptions) (string, error) {
	if gitCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyUrl(
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	if gitCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		return getSshUrl(
			"",
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
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
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		return getSshUrl(
			"",
			"git",
			"github.com",
			joinStrings(":", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityTypeAccessToken {
		return getAccessTokenUrl(
			(githubCheckoutOptions.SecurityOptions.(*AccessTokenSecurityOptions)).AccessToken,
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			githubCheckoutOptions.Type().String(),
			githubCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
}

func getHgUrl(hgCheckoutOptions *HgCheckoutOptions) (string, error) {
	if hgCheckoutOptions.SecurityOptions == nil || hgCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		return getSshUrl(
			"ssh://",
			hgCheckoutOptions.User,
			hgCheckoutOptions.Host,
			hgCheckoutOptions.Path,
		), nil
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			hgCheckoutOptions.Type().String(),
			hgCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
}

func getBitbucketUrl(bitbucketCheckoutOptions *BitbucketCheckoutOptions) (string, error) {
	if bitbucketCheckoutOptions.SecurityOptions == nil || bitbucketCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		switch bitbucketCheckoutOptions.BitbucketType {
		case BitbucketTypeGit:
			return getSshUrl(
				"ssh://",
				"git",
				"bitbucket.org",
				joinStrings(":", bitbucketCheckoutOptions.User, "/", bitbucketCheckoutOptions.Repository, ".git"),
			), nil
		case BitbucketTypeHg:
			return getSshUrl(
				"ssh://",
				"hg",
				"bitbucket.org",
				joinStrings("/", bitbucketCheckoutOptions.User, "/", bitbucketCheckoutOptions.Repository),
			), nil
		default:
			return "", newInternalError(newValidationErrorUnknownBitbucketType(bitbucketCheckoutOptions.Type().String()))
		}
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			bitbucketCheckoutOptions.Type().String(),
			bitbucketCheckoutOptions.SecurityOptions.Type().String(),
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

func ignoreGitCheckoutFilePatterns() []string {
	return []string{
		".git",
		".gitignore",
	}
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

func ignoreHgCheckoutFilePatterns() []string {
	return []string{
		".hg",
		".hgignore",
		".hgsigs",
		".hgtags",
	}
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
