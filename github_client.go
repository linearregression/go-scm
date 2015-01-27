package scm

import (
	"strings"

	"github.com/peter-edge/exec"
)

type GithubCheckoutOptions struct {
	User                string
	Repository          string
	Branch              string
	CommitId            string
	IgnoreCheckoutFiles bool
	AccessToken         string
}

type GithubCheckoutClient interface {
	CheckoutTarball(*GithubCheckoutOptions) (CheckoutTarball, error)
}

func NewGithubCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) GithubCheckoutClient {
	return &githubCheckoutClient{newBaseGitCheckoutClient(executorReadFileManagerProvider)}
}

type githubCheckoutClient struct {
	baseCheckoutClient
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
	return checkout(
		this,
		&baseCheckoutOptions{
			url:                 this.getGithubUrl(githubCheckoutOptions),
			branch:              githubCheckoutOptions.Branch,
			commitId:            githubCheckoutOptions.CommitId,
			ignoreCheckoutFiles: githubCheckoutOptions.IgnoreCheckoutFiles,
		},
	)
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
