package scm

import (
	"errors"
	"fmt"
)

var (
	SecurityTypeSsh         SecurityType = 0
	SecurityTypeAccessToken SecurityType = 1

	securityTypeToString = map[SecurityType]string{
		SecurityTypeSsh:         "ssh",
		SecurityTypeAccessToken: "accessToken",
	}
	lenSecurityTypeToString = len(securityTypeToString)
	stringToSecurityType    = map[string]SecurityType{
		"ssh":         SecurityTypeSsh,
		"accessToken": SecurityTypeAccessToken,
	}
)

type SecurityType uint

func validSecurityType(s string) bool {
	_, ok := stringToSecurityType[s]
	return ok
}

func securityTypeOf(s string) (SecurityType, error) {
	SecurityType, ok := stringToSecurityType[s]
	if !ok {
		return 0, errors.New(unknownSecurityType(s))
	}
	return SecurityType, nil
}

func (this SecurityType) string() string {
	if int(this) < lenSecurityTypeToString {
		return securityTypeToString[this]
	}
	panic(unknownSecurityType(this))
}

func unknownSecurityType(unknownSecurityType interface{}) string {
	return fmt.Sprintf("scm: unknown SecurityType: %v", unknownSecurityType)
}
