package scm

func validateCheckoutOptions(checkoutOptions CheckoutOptions) ValidationError {
	switch checkoutOptions.Type() {
	case CheckoutTypeGit:
		return validateGitCheckoutOptions(checkoutOptions.(*GitCheckoutOptions))
	case CheckoutTypeGithub:
		return validateGithubCheckoutOptions(checkoutOptions.(*GithubCheckoutOptions))
	case CheckoutTypeHg:
		return validateHgCheckoutOptions(checkoutOptions.(*HgCheckoutOptions))
	case CheckoutTypeBitbucket:
		return validateBitbucketCheckoutOptions(checkoutOptions.(*BitbucketCheckoutOptions))
	default:
		return newValidationErrorUnknownCheckoutType(checkoutOptions.Type().string())
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
		if err := validateSecurityOptions(gitCheckoutOptions.SecurityOptions, CheckoutTypeGit, SecurityTypeSsh); err != nil {
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
		if err := validateSecurityOptions(githubCheckoutOptions.SecurityOptions, CheckoutTypeGithub, SecurityTypeSsh, SecurityTypeAccessToken); err != nil {
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
		if err := validateSecurityOptions(hgCheckoutOptions.SecurityOptions, CheckoutTypeHg, SecurityTypeSsh); err != nil {
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
		return newValidationErrorUnknownBitbucketType(bitbucketCheckoutOptions.BitbucketType.string())
	}
	if bitbucketCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(bitbucketCheckoutOptions.SecurityOptions, CheckoutTypeBitbucket, SecurityTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateSecurityOptions(securityOptions SecurityOptions, checkoutType CheckoutType, allowedTypes ...SecurityType) ValidationError {
	if !isAllowedSecurityType(securityOptions.Type(), allowedTypes) {
		return newValidationErrorSecurityNotImplementedForCheckoutType(securityOptions.Type().string(), checkoutType.string())
	}
	switch securityOptions.Type() {
	case SecurityTypeSsh:
		return validateSshSecurityOptions(securityOptions.(*SshSecurityOptions))
	case SecurityTypeAccessToken:
		return validateAccessTokenSecurityOptions(securityOptions.(*AccessTokenSecurityOptions))
	default:
		return newValidationErrorUnknownSecurityType(securityOptions.Type().string())
	}
	return nil
}

func validateSshSecurityOptions(sshSecurityOptions *SshSecurityOptions) ValidationError {
	if sshSecurityOptions.PrivateKey == nil {
		return newValidationErrorRequiredFieldMissing("SshSecurityOptions", "PrivateKey")
	}
	return nil
}

func validateAccessTokenSecurityOptions(accessTokenSecurityOptions *AccessTokenSecurityOptions) ValidationError {
	if accessTokenSecurityOptions.AccessToken == "" {
		return newValidationErrorRequiredFieldMissing("AccessTokenSecurityOptions", "AccessToken")
	}
	return nil
}

func isAllowedSecurityType(securityType SecurityType, allowedTypes []SecurityType) bool {
	for _, allowedType := range allowedTypes {
		if securityType == allowedType {
			return true
		}
	}
	return false
}
