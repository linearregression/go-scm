package scm

func validateCheckoutOptions(checkoutOptions CheckoutOptions) ValidationError {
	switch checkoutOptions.Type() {
	case CheckoutOptionsTypeGit:
		return validateGitCheckoutOptions(checkoutOptions.(*GitCheckoutOptions))
	case CheckoutOptionsTypeGithub:
		return validateGithubCheckoutOptions(checkoutOptions.(*GithubCheckoutOptions))
	case CheckoutOptionsTypeHg:
		return validateHgCheckoutOptions(checkoutOptions.(*HgCheckoutOptions))
	case CheckoutOptionsTypeBitbucket:
		return validateBitbucketCheckoutOptions(checkoutOptions.(*BitbucketCheckoutOptions))
	default:
		return newValidationErrorUnknownCheckoutOptionsType(checkoutOptions.Type().String())
	}
}

func validateGitCheckoutOptions(gitCheckoutOptions *GitCheckoutOptions) ValidationError {
	if gitCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "User")
	}
	if gitCheckoutOptions.Host == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "Host")
	}
	if gitCheckoutOptions.Path == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "Path")
	}
	if gitCheckoutOptions.Branch == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "Branch")
	}
	if gitCheckoutOptions.CommitId == "" {
		return newValidationErrorRequiredFieldMissing("*GitCheckoutOptions", "CommitId")
	}
	if gitCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(gitCheckoutOptions.SecurityOptions, CheckoutOptionsTypeGit, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateGithubCheckoutOptions(githubCheckoutOptions *GithubCheckoutOptions) ValidationError {
	if githubCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "User")
	}
	if githubCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "Repository")
	}
	if githubCheckoutOptions.Branch == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "Branch")
	}
	if githubCheckoutOptions.CommitId == "" {
		return newValidationErrorRequiredFieldMissing("*GithubCheckoutOptions", "CommitId")
	}
	if githubCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(githubCheckoutOptions.SecurityOptions, CheckoutOptionsTypeGithub, SecurityOptionsTypeSsh, SecurityOptionsTypeAccessToken); err != nil {
			return nil
		}
	}
	return nil
}

func validateHgCheckoutOptions(hgCheckoutOptions *HgCheckoutOptions) ValidationError {
	if hgCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "User")
	}
	if hgCheckoutOptions.Host == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "Host")
	}
	if hgCheckoutOptions.Path == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "Path")
	}
	if hgCheckoutOptions.ChangesetId == "" {
		return newValidationErrorRequiredFieldMissing("*HgCheckoutOptions", "ChangesetId")
	}
	if hgCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(hgCheckoutOptions.SecurityOptions, CheckoutOptionsTypeHg, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateBitbucketCheckoutOptions(bitbucketCheckoutOptions *BitbucketCheckoutOptions) ValidationError {
	if bitbucketCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketCheckoutOptions", "User")
	}
	if bitbucketCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketCheckoutOptions", "Repository")
	}
	switch bitbucketCheckoutOptions.BitbucketType {
	case BitbucketTypeGit:
		if bitbucketCheckoutOptions.Branch == "" {
			return newValidationErrorRequiredFieldMissing("*BitbucketCheckoutOptions", "Branch")
		}
		if bitbucketCheckoutOptions.CommitId == "" {
			return newValidationErrorRequiredFieldMissing("*BitbucketCheckoutOptions", "CommitId")
		}
		if bitbucketCheckoutOptions.ChangesetId != "" {
			return newValidationErrorFieldShouldNotBeSet("*BitbucketCheckoutOptions", "ChangesetId")
		}
	case BitbucketTypeHg:
		if bitbucketCheckoutOptions.Branch != "" {
			return newValidationErrorFieldShouldNotBeSet("*BitbucketCheckoutOptions", "Branch")
		}
		if bitbucketCheckoutOptions.CommitId != "" {
			return newValidationErrorFieldShouldNotBeSet("*BitbucketCheckoutOptions", "CommitId")
		}
		if bitbucketCheckoutOptions.ChangesetId == "" {
			return newValidationErrorRequiredFieldMissing("*BitbucketCheckoutOptions", "ChangesetId")
		}
	default:
		return newValidationErrorUnknownBitbucketType(bitbucketCheckoutOptions.BitbucketType.String())
	}
	if bitbucketCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(bitbucketCheckoutOptions.SecurityOptions, CheckoutOptionsTypeBitbucket, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateSecurityOptions(securityOptions SecurityOptions, checkoutType CheckoutOptionsType, allowedTypes ...SecurityOptionsType) ValidationError {
	if !isAllowedSecurityOptionsType(securityOptions.Type(), allowedTypes) {
		return newValidationErrorSecurityNotImplementedForCheckoutOptionsType(securityOptions.Type().String(), checkoutType.String())
	}
	switch securityOptions.Type() {
	case SecurityOptionsTypeSsh:
		return validateSshSecurityOptions(securityOptions.(*SshSecurityOptions))
	case SecurityOptionsTypeAccessToken:
		return validateAccessTokenSecurityOptions(securityOptions.(*AccessTokenSecurityOptions))
	default:
		return newValidationErrorUnknownSecurityOptionsType(securityOptions.Type().String())
	}
	return nil
}

func validateSshSecurityOptions(sshSecurityOptions *SshSecurityOptions) ValidationError {
	//if sshSecurityOptions.PrivateKey == nil {
	//return newValidationErrorRequiredFieldMissing("SshSecurityOptions", "PrivateKey")
	//}
	return nil
}

func validateAccessTokenSecurityOptions(accessTokenSecurityOptions *AccessTokenSecurityOptions) ValidationError {
	if accessTokenSecurityOptions.AccessToken == "" {
		return newValidationErrorRequiredFieldMissing("AccessTokenSecurityOptions", "AccessToken")
	}
	return nil
}

func isAllowedSecurityOptionsType(securityType SecurityOptionsType, allowedTypes []SecurityOptionsType) bool {
	for _, allowedType := range allowedTypes {
		if securityType == allowedType {
			return true
		}
	}
	return false
}
