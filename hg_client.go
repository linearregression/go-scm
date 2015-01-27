package scm

import "github.com/peter-edge/exec"

type HgCheckoutOptions struct {
	Url                 string
	ChangesetId         string
	IgnoreCheckoutFiles bool
}

type HgCheckoutClient interface {
	CheckoutTarball(*HgCheckoutOptions) (CheckoutTarball, error)
}

func NewHgCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) HgCheckoutClient {
	return &hgCheckoutClient{newBaseHgCheckoutClient(executorReadFileManagerProvider)}
}

type hgCheckoutClient struct {
	baseCheckoutClient
}

func (this *hgCheckoutClient) CheckoutTarball(hgCheckoutOptions *HgCheckoutOptions) (CheckoutTarball, error) {
	if hgCheckoutOptions.Url == "" {
		return nil, ErrRequiredFieldMissing
	}
	if hgCheckoutOptions.ChangesetId == "" {
		return nil, ErrRequiredFieldMissing
	}
	return checkout(
		this,
		&baseCheckoutOptions{
			url:                 hgCheckoutOptions.Url,
			commitId:            hgCheckoutOptions.ChangesetId,
			ignoreCheckoutFiles: hgCheckoutOptions.IgnoreCheckoutFiles,
		},
	)
}
