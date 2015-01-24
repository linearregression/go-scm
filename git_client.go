package scm

import "github.com/peter-edge/exec"

type GitCheckoutOptions struct {
	Url      string
	Branch   string
	CommitId string
}

type GitCheckoutClient interface {
	CheckoutTarball(*GitCheckoutOptions) (CheckoutTarball, error)
}

func NewGitCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) GitCheckoutClient {
	return &gitCheckoutClient{newBaseGitCheckoutClient(executorReadFileManagerProvider)}
}

type gitCheckoutClient struct {
	*baseGitCheckoutClient
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
	tarballReader, err := this.checkout(gitCheckoutOptions.Url, gitCheckoutOptions.Branch, gitCheckoutOptions.CommitId)
	if err != nil {
		return nil, err
	}
	return newCheckoutTarball(tarballReader, gitCheckoutOptions.Branch, gitCheckoutOptions.CommitId), nil
}
