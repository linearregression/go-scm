package scm

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ValidationErrorTypeRequiredFieldMissing                         ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                          ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType ValidationErrorType = "SecurityNotImplementedForCheckoutOptionsType"
	ValidationErrorTypeUnknownCheckoutOptionsType                   ValidationErrorType = "UnknownCheckoutOptionsType"
	ValidationErrorTypeUnknownSecurityOptionsType                   ValidationErrorType = "UnknownSecurityOptionsType"
	ValidationErrorTypeUnknownBitbucketType                         ValidationErrorType = "UnknownBitbucketType"
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

func newValidationErrorSecurityNotImplementedForCheckoutOptionsType(securityType string, checkoutType string) ValidationError {
	return newValidationError(ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType, map[string]string{"securityType": securityType, "checkoutType": checkoutType})
}

func newValidationErrorUnknownCheckoutOptionsType(checkoutType string) ValidationError {
	return newValidationError(ValidationErrorTypeUnknownCheckoutOptionsType, map[string]string{"checkoutType": checkoutType})
}

func newValidationErrorUnknownSecurityOptionsType(securityType string) ValidationError {
	return newValidationError(ValidationErrorTypeUnknownSecurityOptionsType, map[string]string{"securityType": securityType})
}

func newValidationErrorUnknownBitbucketType(bitbucketType string) ValidationError {
	return newValidationError(ValidationErrorTypeUnknownBitbucketType, map[string]string{"bitbucketType": bitbucketType})
}

func newInternalError(validationError ValidationError) error {
	return errors.New(validationError.Error())
}
