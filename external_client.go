package scm

import (
	"bytes"
	"io"
)

func newExternalClient(client Client) *externalClient {
	return &externalClient{client}
}

type externalClient struct {
	client Client
}

func (this *externalClient) CheckoutTarball(externalCheckoutOptions *ExternalCheckoutOptions) (io.Reader, error) {
	checkoutOptions, err := ConvertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return this.client.CheckoutTarball(checkoutOptions)
}

func ConvertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	var securityOptions SecurityOptions
	if externalCheckoutOptions.SecurityOptions != nil {
		if !validSecurityType(externalCheckoutOptions.SecurityOptions.Type) {
			return nil, newValidationErrorUnknownSecurityType(externalCheckoutOptions.Type)
		}
		securityType, err := securityTypeOf(externalCheckoutOptions.SecurityOptions.Type)
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
			return nil, newInternalError(newValidationErrorUnknownSecurityType(securityType.String()))
		}
	}
	if !validCheckoutType(externalCheckoutOptions.Type) {
		return nil, newValidationErrorUnknownCheckoutType(externalCheckoutOptions.Type)
	}
	checkoutType, err := checkoutTypeOf(externalCheckoutOptions.Type)
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
		if !validBitbucketType(externalCheckoutOptions.BitbucketType) {
			return nil, newValidationErrorUnknownBitbucketType(externalCheckoutOptions.BitbucketType)
		}
		bitbucketType, err := bitbucketTypeOf(externalCheckoutOptions.BitbucketType)
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
		return nil, newInternalError(newValidationErrorUnknownCheckoutType(checkoutType.String()))
	}
}
