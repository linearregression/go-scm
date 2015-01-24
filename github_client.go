package scm

import (
	"strings"

	"github.com/peter-edge/exec"
)

type GithubCheckoutOptions struct {
	User        string
	Repository  string
	Branch      string
	CommitId    string
	AccessToken string
}

type GithubCheckoutClient interface {
	CheckoutTarball(*GithubCheckoutOptions) (CheckoutTarball, error)
}

func NewGithubCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) GithubCheckoutClient {
	return &githubCheckoutClient{newBaseGitCheckoutClient(executorReadFileManagerProvider)}
}

type githubCheckoutClient struct {
	*baseGitCheckoutClient
}

func (this *githubCheckoutClient) CheckoutTarball(githubCheckoutOptions *GithubCheckoutOptions) (CheckoutTarball, error) {
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
	tarballReader, err := this.checkout(this.getGithubUrl(githubCheckoutOptions), githubCheckoutOptions.Branch, githubCheckoutOptions.CommitId)
	if err != nil {
		return nil, err
	}
	return newCheckoutTarball(tarballReader, githubCheckoutOptions.Branch, githubCheckoutOptions.CommitId), nil
}

func (this *githubCheckoutClient) getGithubUrl(githubCheckoutOptions *GithubCheckoutOptions) string {
	return strings.Join(
		[]string{
			this.getBaseUrl(githubCheckoutOptions.AccessToken),
			"/",
			githubCheckoutOptions.User,
			"/",
			githubCheckoutOptions.Repository,
			".git",
		},
		"",
	)
}

func (this *githubCheckoutClient) getBaseUrl(accessToken string) string {
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
