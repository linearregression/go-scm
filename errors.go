package scm

import "errors"

var (
	ErrNil                    = errors.New("scm: nil")
	ErrWrongSecurityType      = errors.New("scm: wrong security type")
	ErrRequiredFieldMissing   = errors.New("scm: required field missing")
	ErrFieldShouldNotBeSet    = errors.New("scm: field should not be set")
	ErrSecurityNotImplemented = errors.New("scm: security not implemented")
	ErrUnknownBitbucketType   = errors.New("scm: unknown BitbucketType")
)
