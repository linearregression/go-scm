package scm

import (
	"fmt"
	"strings"
)

var (
	ValidationErrorTypeRequiredFieldMissing                         ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                          ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType ValidationErrorType = "SecurityNotImplementedForCheckoutOptionsType"
)

type ValidationErrorType string

type ValidationError interface {
	error
	Type() ValidationErrorType
}
type validationError struct {
	errorType ValidationErrorType
	tags      map[string]string
}

func newValidationError(errorType ValidationErrorType, tags map[string]string) *validationError {
	if tags == nil {
		tags = make(map[string]string)
	}
	return &validationError{errorType, tags}
}

func (this *validationError) Error() string {
	return fmt.Sprintf("%v %v", this.errorType, this.tags)
}

func (this *validationError) Type() ValidationErrorType {
	return this.errorType
}

func ValidateCheckoutOptions(checkoutOptions CheckoutOptions) error {
	return validateCheckoutOptions(checkoutOptions)
}

func validateCheckoutOptions(checkoutOptions CheckoutOptions) error {
	return CheckoutOptionsSwitch(
		checkoutOptions,
		validateGitCheckoutOptions,
		validateGithubCheckoutOptions,
		validateHgCheckoutOptions,
		validateBitbucketGitCheckoutOptions,
		validateBitbucketHgCheckoutOptions,
	)
}

func validateGitCheckoutOptions(gitCheckoutOptions *GitCheckoutOptions) error {
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

func validateGithubCheckoutOptions(githubCheckoutOptions *GithubCheckoutOptions) error {
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

func validateHgCheckoutOptions(hgCheckoutOptions *HgCheckoutOptions) error {
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

func validateBitbucketGitCheckoutOptions(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error {
	if bitbucketGitCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "User")
	}
	if bitbucketGitCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "Repository")
	}
	if bitbucketGitCheckoutOptions.Branch == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "Branch")
	}
	if bitbucketGitCheckoutOptions.CommitId == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketGitCheckoutOptions", "CommitId")
	}
	if bitbucketGitCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(bitbucketGitCheckoutOptions.SecurityOptions, CheckoutOptionsTypeBitbucketGit, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateBitbucketHgCheckoutOptions(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error {
	if bitbucketHgCheckoutOptions.User == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "User")
	}
	if bitbucketHgCheckoutOptions.Repository == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "Repository")
	}
	if bitbucketHgCheckoutOptions.ChangesetId == "" {
		return newValidationErrorRequiredFieldMissing("*BitbucketHgCheckoutOptions", "ChangesetId")
	}
	if bitbucketHgCheckoutOptions.SecurityOptions != nil {
		if err := validateSecurityOptions(bitbucketHgCheckoutOptions.SecurityOptions, CheckoutOptionsTypeBitbucketHg, SecurityOptionsTypeSsh); err != nil {
			return nil
		}
	}
	return nil
}

func validateSecurityOptions(securityOptions SecurityOptions, checkoutType CheckoutOptionsType, allowedTypes ...SecurityOptionsType) error {
	if !isAllowedSecurityOptionsType(securityOptions.Type(), allowedTypes) {
		return newValidationErrorSecurityNotImplementedForCheckoutOptionsType(securityOptions.Type().String(), checkoutType.String())
	}
	return SecurityOptionsSwitch(
		securityOptions,
		validateSshSecurityOptions,
		validateAccessTokenSecurityOptions,
	)
}

func validateSshSecurityOptions(sshSecurityOptions *SshSecurityOptions) error {
	return nil
}

func validateAccessTokenSecurityOptions(accessTokenSecurityOptions *AccessTokenSecurityOptions) error {
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

func newValidationErrorRequiredFieldMissing(objectType string, fieldPath ...string) ValidationError {
	return newValidationError(ValidationErrorTypeRequiredFieldMissing, map[string]string{"type": objectType, "fieldPath": strings.Join(fieldPath, ".")})
}

func newValidationErrorFieldShouldNotBeSet(objectType string, fieldPath ...string) ValidationError {
	return newValidationError(ValidationErrorTypeFieldShouldNotBeSet, map[string]string{"type": objectType, "fieldPath": strings.Join(fieldPath, ".")})
}

func newValidationErrorSecurityNotImplementedForCheckoutOptionsType(securityType string, checkoutType string) ValidationError {
	return newValidationError(ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType, map[string]string{"securityType": securityType, "checkoutType": checkoutType})
}
