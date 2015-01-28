package scm

import "errors"

var (
	ErrNil                    = errors.New("scm: nil")
	ErrWrongSecurityType      = errors.New("scm: wrong security type")
	ErrRequiredFieldMissing   = errors.New("scm: required field missing")
	ErrSecurityNotImplemented = errors.New("scm: security not implemented")
)
