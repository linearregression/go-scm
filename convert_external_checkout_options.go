package scm

import "bytes"

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
		func() (*BitbucketCheckoutOptions, error) {
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
		},
	)
}
