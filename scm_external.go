package scm

import (
	"bytes"
	"io"
)

type ExternalCheckoutOptions struct {
	Type            string                   `json:"type,omitempty" yaml:"type,omitempty"`
	User            string                   `json:"user,omitempty" yaml:"user,omitempty"`
	Host            string                   `json:"host,omitempty" yaml:"host,omitempty"`
	Path            string                   `json:"path,omitempty" yaml:"path,omitempty"`
	Repository      string                   `json:"repository,omitempty" yaml:"repository,omitempty"`
	Branch          string                   `json:"branch,omitempty" yaml:"branch,omitempty"`
	CommitId        string                   `json:"commit_id,omitempty" yaml:"commit_id,omitempty"`
	BitbucketType   string                   `json:"bitbucket_type,omitempty" yaml:"bitbucket_type,omitempty"`
	ChangesetId     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
	SecurityOptions *ExternalSecurityOptions `json:"security_options,omitempty" yaml:"security_options,omitempty"`
}

type ExternalSecurityOptions struct {
	Type                  string `json:"type,omitempty" yaml:"type,omitempty"`
	StrictHostKeyChecking bool   `json:"strict_host_key_checking,omitempty" yaml:"strict_host_key_checking,omitempty"`
	PrivateKey            string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	AccessToken           string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
}

func ExternalCheckout(absolutePath string, externalCheckoutOptions *ExternalCheckoutOptions, options *Options) error {
	checkoutOptions, err := convertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return err
	}
	return Checkout(absolutePath, checkoutOptions, options)
}

func ExternalCheckoutTarball(externalCheckoutOptions *ExternalCheckoutOptions, options *Options) (io.Reader, error) {
	checkoutOptions, err := convertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return CheckoutTarball(checkoutOptions, options)
}

// ***** PRIVATE *****

func convertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	var securityOptions SecurityOptions
	if externalCheckoutOptions.SecurityOptions != nil {
		if !ValidSecurityType(externalCheckoutOptions.SecurityOptions.Type) {
			return nil, newValidationErrorUnknownSecurityType(externalCheckoutOptions.Type)
		}
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
	if !ValidCheckoutType(externalCheckoutOptions.Type) {
		return nil, newValidationErrorUnknownCheckoutType(externalCheckoutOptions.Type)
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
		if !ValidBitbucketType(externalCheckoutOptions.BitbucketType) {
			return nil, newValidationErrorUnknownBitbucketType(externalCheckoutOptions.BitbucketType)
		}
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
		return nil, newInternalError(newValidationErrorUnknownCheckoutType(checkoutType.String()))
	}
}
