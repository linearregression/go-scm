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
	return &client{newBaseGitClient(executorReadFileManagerProvider), newBaseHgClient(executorReadFileManagerProvider), clientOptions}
}

type client struct {
	gitClient     baseClient
	hgClient      baseClient
	clientOptions *ClientOptions
}

func (this *client) CheckoutGitTarball(gitCheckoutOptions *GitCheckoutOptions) (io.Reader, error) {
	if gitCheckoutOptions.Url == "" {
		return nil, ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.Branch == "" {
		return nil, ErrRequiredFieldMissing
	}
	if gitCheckoutOptions.CommitId == "" {
		return nil, ErrRequiredFieldMissing
	}
	return checkout(
		this.gitClient,
		&baseCheckoutOptions{
			url:                 gitCheckoutOptions.Url,
			branch:              gitCheckoutOptions.Branch,
			commitId:            gitCheckoutOptions.CommitId,
			ignoreCheckoutFiles: this.clientOptions.IgnoreCheckoutFiles,
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
	return checkout(
		this.gitClient,
		&baseCheckoutOptions{
			url:                 this.getGithubUrl(githubCheckoutOptions),
			branch:              githubCheckoutOptions.Branch,
			commitId:            githubCheckoutOptions.CommitId,
			ignoreCheckoutFiles: this.clientOptions.IgnoreCheckoutFiles,
		},
	)
}

func (this *client) CheckoutHgTarball(hgCheckoutOptions *HgCheckoutOptions) (io.Reader, error) {
	if hgCheckoutOptions.Url == "" {
		return nil, ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.ChangesetId == "" {
		return nil, ErrRequiredFieldMissing
	}
	return checkout(
		this.hgClient,
		&baseCheckoutOptions{
			url:                 hgCheckoutOptions.Url,
			commitId:            hgCheckoutOptions.ChangesetId,
			ignoreCheckoutFiles: this.clientOptions.IgnoreCheckoutFiles,
		},
	)
}

func (this *client) getGithubUrl(githubCheckoutOptions *GithubCheckoutOptions) string {
	return strings.Join(
		[]string{
			this.getGithubBaseUrl(githubCheckoutOptions.AccessToken),
			"/",
			githubCheckoutOptions.User,
			"/",
			githubCheckoutOptions.Repository,
			".git",
		},
		"",
	)
}

func (this *client) getGithubBaseUrl(accessToken string) string {
	if accessToken != "" {
		return strings.Join(
			[]string{
				"https://",
				accessToken,
				":x-oauth-basic@github.com",
			},
			"",
		)
	}
	return "https://github.com"
}

func checkout(baseClient baseClient, baseCheckoutOptions *baseCheckoutOptions) (io.Reader, error) {
	client, err := baseClient.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	err = baseClient.checkoutWithExecutor(client, baseCheckoutOptions, clonePath)
	if err != nil {
		return nil, err
	}
	var ignoreCheckoutFilePatterns []string = nil
	if baseCheckoutOptions.ignoreCheckoutFiles {
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
