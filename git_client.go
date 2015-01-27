package scm

import "github.com/peter-edge/exec"

type GitCheckoutOptions struct {
	Url                 string
	Branch              string
	CommitId            string
	IgnoreCheckoutFiles bool
}

type GitCheckoutClient interface {
	CheckoutTarball(*GitCheckoutOptions) (CheckoutTarball, error)
}

func NewGitCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) GitCheckoutClient {
	return &gitCheckoutClient{newBaseGitCheckoutClient(executorReadFileManagerProvider)}
}

type gitCheckoutClient struct {
	baseCheckoutClient
}

func (this *gitCheckoutClient) CheckoutTarball(gitCheckoutOptions *GitCheckoutOptions) (CheckoutTarball, error) {
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
		this,
		&baseCheckoutOptions{
			url:                 gitCheckoutOptions.Url,
			branch:              gitCheckoutOptions.Branch,
			commitId:            gitCheckoutOptions.CommitId,
			ignoreCheckoutFiles: gitCheckoutOptions.IgnoreCheckoutFiles,
		},
	)
}
