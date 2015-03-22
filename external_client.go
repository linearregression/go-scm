package scm

import (
	"bytes"
	"io"

	"github.com/peter-edge/go-exec"
)

type externalClient struct {
	client Client
}

func newExternalClient(client Client) *externalClient {
	return &externalClient{client}
}

func (this *externalClient) CheckoutTarball(externalCheckoutOptions *ExternalCheckoutOptions) (io.Reader, error) {
	checkoutOptions, err := ConvertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return this.client.CheckoutTarball(checkoutOptions)
}

type externalDirectClient struct {
	directClient DirectClient
}

func newExternalDirectClient(directClient DirectClient) *externalDirectClient {
	return &externalDirectClient{directClient}
}

func (this *externalDirectClient) Checkout(externalCheckoutOptions *ExternalCheckoutOptions, executor exec.Executor, path string) error {
	checkoutOptions, err := ConvertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return err
	}
	return this.directClient.Checkout(checkoutOptions, executor, path)
}

func ConvertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	var securityOptions SecurityOptions
	if externalCheckoutOptions.SecurityOptions != nil {
		securityType, err := SecurityTypeOf(externalCheckoutOptions.SecurityOptions.Type)
		if err != nil {
			return nil, err
		}
		switch securityType {
		case SecurityTypeSsh:
			var privateKey bytes.Buffer
			privateKey.WriteString(externalCheckoutOptions.SecurityOptions.PrivateKey)
			securityOptions = &SshSecurityOptions{
				StrictHostKeyChecking: externalCheckoutOptions.SecurityOptions.StrictHostKeyChecking,
				PrivateKey:            &privateKey,
			}
		case SecurityTypeAccessToken:
			securityOptions = &AccessTokenSecurityOptions{
				AccessToken: externalCheckoutOptions.SecurityOptions.AccessToken,
			}
		default:
			return nil, UnknownSecurityType(securityType)
		}
	}
	checkoutType, err := CheckoutTypeOf(externalCheckoutOptions.Type)
	if err != nil {
		return nil, err
	}
	switch checkoutType {
	case CheckoutTypeGit:
		return &GitCheckoutOptions{
			User:            externalCheckoutOptions.User,
			Host:            externalCheckoutOptions.Host,
			Path:            externalCheckoutOptions.Path,
			Branch:          externalCheckoutOptions.Branch,
			CommitId:        externalCheckoutOptions.CommitId,
			SecurityOptions: securityOptions,
		}, nil
	case CheckoutTypeGithub:
		return &GithubCheckoutOptions{
			User:            externalCheckoutOptions.User,
			Repository:      externalCheckoutOptions.Repository,
			Branch:          externalCheckoutOptions.Branch,
			CommitId:        externalCheckoutOptions.CommitId,
			SecurityOptions: securityOptions,
		}, nil
	case CheckoutTypeHg:
		return &HgCheckoutOptions{
			User:            externalCheckoutOptions.User,
			Host:            externalCheckoutOptions.Host,
			Path:            externalCheckoutOptions.Path,
			ChangesetId:     externalCheckoutOptions.ChangesetId,
			SecurityOptions: securityOptions,
		}, nil
	case CheckoutTypeBitbucket:
		bitbucketType, err := BitbucketTypeOf(externalCheckoutOptions.BitbucketType)
		if err != nil {
			return nil, err
		}
		return &BitbucketCheckoutOptions{
			BitbucketType:   bitbucketType,
			User:            externalCheckoutOptions.User,
			Repository:      externalCheckoutOptions.Repository,
			Branch:          externalCheckoutOptions.Branch,
			CommitId:        externalCheckoutOptions.CommitId,
			ChangesetId:     externalCheckoutOptions.ChangesetId,
			SecurityOptions: securityOptions,
		}, nil
	default:
		return nil, UnknownCheckoutType(checkoutType)
	}
}
