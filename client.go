package scm

import (
	"io"
	"strings"

	"github.com/peter-edge/exec"
	tarexec "github.com/peter-edge/tar/exec"
)

const (
	clonePath = "clone"
)

func newClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider, clientOptions *ClientOptions) *client {
	return &client{executorReadFileManagerProvider, newBaseGitClient(), newBaseHgClient(), clientOptions}
}

type client struct {
	exec.ExecutorReadFileManagerProvider
	gitClient     *baseGitClient
	hgClient      *baseHgClient
	clientOptions *ClientOptions
}

func (this *client) CheckoutGitTarball(gitCheckoutOptions *GitCheckoutOptions) (io.Reader, error) {
	if gitCheckoutOptions.User == "" {
		return nil, ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Host == "" {
		return nil, ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Path == "" {
		return nil, ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Branch == "" {
		return nil, ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.CommitId == "" {
		return nil, ErrRequiredFieldMissing
	}
	baseCloneArgs := []string{"git", "clone"}
	url, err := getGitUrl(gitCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return this.checkout(
		this.gitClient,
		&baseCheckoutOptions{
			baseCloneArgs: baseCloneArgs,
			url:           url,
			branch:        gitCheckoutOptions.Branch,
			commitId:      gitCheckoutOptions.CommitId,
		},
	)
}

func (this *client) CheckoutGithubTarball(githubCheckoutOptions *GithubCheckoutOptions) (io.Reader, error) {
	if githubCheckoutOptions.User == "" {
		return nil, ErrRequiredFieldMissing
	}
	if githubCheckoutOptions.Repository == "" {
		return nil, ErrRequiredFieldMissing
	}
	if githubCheckoutOptions.Branch == "" {
		return nil, ErrRequiredFieldMissing
	}
	if githubCheckoutOptions.CommitId == "" {
		return nil, ErrRequiredFieldMissing
	}
	baseCloneArgs := []string{"git", "clone"}
	url, err := getGithubUrl(githubCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return this.checkout(
		this.gitClient,
		&baseCheckoutOptions{
			baseCloneArgs: baseCloneArgs,
			url:           url,
			branch:        githubCheckoutOptions.Branch,
			commitId:      githubCheckoutOptions.CommitId,
		},
	)
}

func (this *client) CheckoutHgTarball(hgCheckoutOptions *HgCheckoutOptions) (io.Reader, error) {
	if hgCheckoutOptions.User == "" {
		return nil, ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.Host == "" {
		return nil, ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.Path == "" {
		return nil, ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.ChangesetId == "" {
		return nil, ErrRequiredFieldMissing
	}
	baseCloneArgs := []string{"hg", "clone"}
	url, err := getHgUrl(hgCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return this.checkout(
		this.hgClient,
		&baseCheckoutOptions{
			baseCloneArgs: baseCloneArgs,
			url:           url,
			commitId:      hgCheckoutOptions.ChangesetId,
		},
	)
}

func (this *client) checkout(baseClient baseClient, baseCheckoutOptions *baseCheckoutOptions) (io.Reader, error) {
	client, err := this.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	err = baseClient.checkoutWithExecutor(client, baseCheckoutOptions, clonePath)
	if err != nil {
		return nil, err
	}
	var ignoreCheckoutFilePatterns []string = nil
	if this.clientOptions.IgnoreCheckoutFiles {
		ignoreCheckoutFilePatterns = baseClient.ignoreCheckoutFilePatterns(client)
	}

	reader, err := tarFiles(client, ignoreCheckoutFilePatterns, clonePath)
	if err != nil {
		return nil, err
	}
	if err := client.Destroy(); err != nil {
		return nil, err
	}
	return reader, nil
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
			getGithubPath(githubCheckoutOptions),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.securityType() == securityTypeAccessToken {
		return getAccessTokenUrl(
			githubCheckoutOptions.SecurityOptions.accessTokenOptions().AccessToken,
			"github.com",
			getGithubPath(githubCheckoutOptions),
		), nil
	}
	return "", ErrSecurityNotImplemented
}

func getGithubPath(githubCheckoutOptions *GithubCheckoutOptions) string {
	return joinStrings(":", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git")
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
	return joinStrings("https://", accessToken, ":x-oauth-basic@", host, "/", path)
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
