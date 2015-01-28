package scm

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/peter-edge/exec"
	tarexec "github.com/peter-edge/tar/exec"
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

func (this *client) CheckoutGitTarball(gitCheckoutOptions *GitCheckoutOptions) (reader io.Reader, retErr error) {
	if err := validateGitCheckoutOptions(gitCheckoutOptions); err != nil {
		return nil, err
	}
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

func (this *client) CheckoutGithubTarball(githubCheckoutOptions *GithubCheckoutOptions) (reader io.Reader, retErr error) {
	if err := validateGithubCheckoutOptions(githubCheckoutOptions); err != nil {
		return nil, err
	}
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

func (this *client) CheckoutHgTarball(hgCheckoutOptions *HgCheckoutOptions) (reader io.Reader, retErr error) {
	if err := validateHgCheckoutOptions(hgCheckoutOptions); err != nil {
		return nil, err
	}
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
	if securityOptions.securityType() != securityTypeSsh {
		return "", nil, nil
	}
	sshOptions := securityOptions.sshOptions()

	sshCommand := []string{"ssh", "-o"}
	if sshOptions.StrictHostKeyChecking {
		sshCommand = append(sshCommand, "StrictHostKeyChecking=yes")
	} else {
		sshCommand = append(sshCommand, "StrictHostKeyChecking=no")
	}
	var client exec.Client
	if sshOptions.PrivateKey != nil {
		client, err := this.NewTempDirClient()
		if err != nil {
			return "", nil, err
		}
		writeFile, err := client.Create("id_rsa")
		if err != nil {
			return "", nil, err
		}
		data, err := ioutil.ReadAll(sshOptions.PrivateKey)
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
	if gitCheckoutOptions.SecurityOptions == nil || gitCheckoutOptions.SecurityOptions.securityType() == securityTypeSsh {
		return getSshUrl(
			"",
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	return "", ErrSecurityNotImplemented
}

func getGithubUrl(githubCheckoutOptions *GithubCheckoutOptions) (string, error) {
	if githubCheckoutOptions.SecurityOptions == nil || githubCheckoutOptions.SecurityOptions.securityType() == securityTypeSsh {
		return getSshUrl(
			"",
			"git",
			"github.com",
			joinStrings(":", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.securityType() == securityTypeAccessToken {
		return getAccessTokenUrl(
			githubCheckoutOptions.SecurityOptions.accessTokenOptions().AccessToken,
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", ErrSecurityNotImplemented
}

func getHgUrl(hgCheckoutOptions *HgCheckoutOptions) (string, error) {
	if hgCheckoutOptions.SecurityOptions == nil || hgCheckoutOptions.SecurityOptions.securityType() == securityTypeSsh {
		return getSshUrl(
			"ssh://",
			hgCheckoutOptions.User,
			hgCheckoutOptions.Host,
			hgCheckoutOptions.Path,
		), nil
	}
	return "", ErrSecurityNotImplemented
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
	var stderr bytes.Buffer
	cmd := exec.Cmd{
		// TODO(peter): if the commit id is more than 50 back, the checkout will fail
		Args:   []string{"git", "clone", "--branch", branch, "--depth", "50", "--recursive", url, path},
		Stderr: &stderr,
	}
	if gitSshCommand != "" {
		cmd.Env = []string{"GIT_SSH_COMMAND=" + gitSshCommand}
	}
	if err := executor.Execute(&cmd)(); err != nil {
		return err
	}
	return executor.Execute(
		&exec.Cmd{
			Args:   []string{"git", "checkout", "-f", commitId},
			SubDir: path,
		},
	)()
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
	if err := executor.Execute(&exec.Cmd{Args: args})(); err != nil {
		return err
	}
	return executor.Execute(
		&exec.Cmd{
			Args: []string{"hg", "update", "--cwd", path, changesetId},
		},
	)()
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

func validateGitCheckoutOptions(gitCheckoutOptions *GitCheckoutOptions) error {
	if gitCheckoutOptions.User == "" {
		return ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Host == "" {
		return ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Path == "" {
		return ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Branch == "" {
		return ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.CommitId == "" {
		return ErrRequiredFieldMissing
	}
	return nil
}

func validateGithubCheckoutOptions(githubCheckoutOptions *GithubCheckoutOptions) error {
	if githubCheckoutOptions.User == "" {
		return ErrRequiredFieldMissing
	}
	if githubCheckoutOptions.Repository == "" {
		return ErrRequiredFieldMissing
	}
	if githubCheckoutOptions.Branch == "" {
		return ErrRequiredFieldMissing
	}
	if githubCheckoutOptions.CommitId == "" {
		return ErrRequiredFieldMissing
	}
	return nil
}

func validateHgCheckoutOptions(hgCheckoutOptions *HgCheckoutOptions) error {
	if hgCheckoutOptions.User == "" {
		return ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.Host == "" {
		return ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.Path == "" {
		return ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.ChangesetId == "" {
		return ErrRequiredFieldMissing
	}
	return nil
}

func joinStrings(elems ...string) string {
	return strings.Join(elems, "")
}
