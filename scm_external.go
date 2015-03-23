package scm

import "bytes"

type ExternalCheckoutOptions struct {
	Type            string                   `json:"type,omitempty" yaml:"type,omitempty"`
	User            string                   `json:"user,omitempty" yaml:"user,omitempty"`
	Host            string                   `json:"host,omitempty" yaml:"host,omitempty"`
	Path            string                   `json:"path,omitempty" yaml:"path,omitempty"`
	Repository      string                   `json:"repository,omitempty" yaml:"repository,omitempty"`
	Branch          string                   `json:"branch,omitempty" yaml:"branch,omitempty"`
	CommitId        string                   `json:"commit_id,omitempty" yaml:"commit_id,omitempty"`
	ChangesetId     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
	SecurityOptions *ExternalSecurityOptions `json:"security_options,omitempty" yaml:"security_options,omitempty"`
}

type ExternalSecurityOptions struct {
	Type                  string `json:"type,omitempty" yaml:"type,omitempty"`
	StrictHostKeyChecking bool   `json:"strict_host_key_checking,omitempty" yaml:"strict_host_key_checking,omitempty"`
	PrivateKey            string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	AccessToken           string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
}

func ConvertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	return convertExternalCheckoutOptions(externalCheckoutOptions)
}

// ***** PRIVATE *****

func convertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	var securityOptions SecurityOptions
	if externalCheckoutOptions.SecurityOptions != nil {
		securityOptionsType, err := SecurityOptionsTypeOf(externalCheckoutOptions.SecurityOptions.Type)
		if err != nil {
			return nil, err
		}
		securityOptions, err = securityOptionsType.NewSecurityOptions(
			func() (*SshSecurityOptions, error) {
				var privateKey bytes.Buffer
				privateKey.WriteString(externalCheckoutOptions.SecurityOptions.PrivateKey)
				return &SshSecurityOptions{
					StrictHostKeyChecking: externalCheckoutOptions.SecurityOptions.StrictHostKeyChecking,
					PrivateKey:            &privateKey,
				}, nil
			},
			func() (*AccessTokenSecurityOptions, error) {
				return &AccessTokenSecurityOptions{
					AccessToken: externalCheckoutOptions.SecurityOptions.AccessToken,
				}, nil
			},
		)
		if err != nil {
			return nil, err
		}
	}
	checkoutOptionsType, err := CheckoutOptionsTypeOf(externalCheckoutOptions.Type)
	if err != nil {
		return nil, err
	}
	return checkoutOptionsType.NewCheckoutOptions(
		func() (*GitCheckoutOptions, error) {
			return &GitCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Host:            externalCheckoutOptions.Host,
				Path:            externalCheckoutOptions.Path,
				Branch:          externalCheckoutOptions.Branch,
				CommitId:        externalCheckoutOptions.CommitId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*GithubCheckoutOptions, error) {
			return &GithubCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Repository:      externalCheckoutOptions.Repository,
				Branch:          externalCheckoutOptions.Branch,
				CommitId:        externalCheckoutOptions.CommitId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*HgCheckoutOptions, error) {
			return &HgCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Host:            externalCheckoutOptions.Host,
				Path:            externalCheckoutOptions.Path,
				ChangesetId:     externalCheckoutOptions.ChangesetId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*BitbucketGitCheckoutOptions, error) {
			return &BitbucketGitCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Repository:      externalCheckoutOptions.Repository,
				Branch:          externalCheckoutOptions.Branch,
				CommitId:        externalCheckoutOptions.CommitId,
				SecurityOptions: securityOptions,
			}, nil
		},
		func() (*BitbucketHgCheckoutOptions, error) {
			return &BitbucketHgCheckoutOptions{
				User:            externalCheckoutOptions.User,
				Repository:      externalCheckoutOptions.Repository,
				ChangesetId:     externalCheckoutOptions.ChangesetId,
				SecurityOptions: securityOptions,
			}, nil
		},
	)
}
