package scm

import "bytes"

func convertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error) {
	var securityOptions SecurityOptions
	if externalCheckoutOptions.SecurityOptions != nil {
		securityType, err := SecurityOptionsTypeOf(externalCheckoutOptions.SecurityOptions.Type)
		if err != nil {
			return nil, err
		}
		switch securityType {
		case SecurityOptionsTypeSsh:
			var privateKey bytes.Buffer
			privateKey.WriteString(externalCheckoutOptions.SecurityOptions.PrivateKey)
			securityOptions = &SshSecurityOptions{
				StrictHostKeyChecking: externalCheckoutOptions.SecurityOptions.StrictHostKeyChecking,
				PrivateKey:            &privateKey,
			}
		case SecurityOptionsTypeAccessToken:
			securityOptions = &AccessTokenSecurityOptions{
				AccessToken: externalCheckoutOptions.SecurityOptions.AccessToken,
			}
		default:
			return nil, UnknownSecurityOptionsType(securityType)
		}
	}
	checkoutType, err := CheckoutOptionsTypeOf(externalCheckoutOptions.Type)
	if err != nil {
		return nil, err
	}
	switch checkoutType {
	case CheckoutOptionsTypeGit:
		return &GitCheckoutOptions{
			User:            externalCheckoutOptions.User,
			Host:            externalCheckoutOptions.Host,
			Path:            externalCheckoutOptions.Path,
			Branch:          externalCheckoutOptions.Branch,
			CommitId:        externalCheckoutOptions.CommitId,
			SecurityOptions: securityOptions,
		}, nil
	case CheckoutOptionsTypeGithub:
		return &GithubCheckoutOptions{
			User:            externalCheckoutOptions.User,
			Repository:      externalCheckoutOptions.Repository,
			Branch:          externalCheckoutOptions.Branch,
			CommitId:        externalCheckoutOptions.CommitId,
			SecurityOptions: securityOptions,
		}, nil
	case CheckoutOptionsTypeHg:
		return &HgCheckoutOptions{
			User:            externalCheckoutOptions.User,
			Host:            externalCheckoutOptions.Host,
			Path:            externalCheckoutOptions.Path,
			ChangesetId:     externalCheckoutOptions.ChangesetId,
			SecurityOptions: securityOptions,
		}, nil
	case CheckoutOptionsTypeBitbucket:
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
		return nil, UnknownCheckoutOptionsType(checkoutType)
	}
}
