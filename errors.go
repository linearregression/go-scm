package scm

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ValidationErrorTypeRequiredFieldMissing                  ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                   ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutType ValidationErrorType = "SecurityNotImplementedForCheckoutType"
	ValidationErrorTypeUnknownCheckoutType                   ValidationErrorType = "UnknownCheckoutType"
	ValidationErrorTypeUnknownSecurityType                   ValidationErrorType = "UnknownSecurityType"
	ValidationErrorTypeUnknownBitbucketType                  ValidationErrorType = "UnknownBitbucketType"
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

func newValidationErrorRequiredFieldMissing(objectType string, fieldPath ...string) ValidationError {
	return newValidationError(ValidationErrorTypeRequiredFieldMissing, map[string]string{"type": objectType, "fieldPath": strings.Join(fieldPath, ".")})
}

func newValidationErrorFieldShouldNotBeSet(objectType string, fieldPath ...string) ValidationError {
	return newValidationError(ValidationErrorTypeFieldShouldNotBeSet, map[string]string{"type": objectType, "fieldPath": strings.Join(fieldPath, ".")})
}

func newValidationErrorSecurityNotImplementedForCheckoutType(securityType string, checkoutType string) ValidationError {
	return newValidationError(ValidationErrorTypeSecurityNotImplementedForCheckoutType, map[string]string{"securityType": securityType, "checkoutType": checkoutType})
}

func newValidationErrorUnknownCheckoutType(checkoutType string) ValidationError {
	return newValidationError(ValidationErrorTypeUnknownCheckoutType, map[string]string{"checkoutType": checkoutType})
}

func newValidationErrorUnknownSecurityType(securityType string) ValidationError {
	return newValidationError(ValidationErrorTypeUnknownSecurityType, map[string]string{"securityType": securityType})
}

func newValidationErrorUnknownBitbucketType(bitbucketType string) ValidationError {
	return newValidationError(ValidationErrorTypeUnknownBitbucketType, map[string]string{"bitbucketType": bitbucketType})
}

func newInternalError(validationError ValidationError) error {
	return errors.New(validationError.Error())
}
