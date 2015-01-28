package scm

import (
	"errors"
	"fmt"
)

const (
	securityTypeSsh = iota
	securityTypeAccessToken
)

var (
	securityTypeToString = map[securityType]string{
		securityTypeSsh:         "ssh",
		securityTypeAccessToken: "accessToken",
	}
	lenSecurityTypeToString = len(securityTypeToString)
	stringToSecurityType    = map[string]securityType{
		"ssh":         securityTypeSsh,
		"accessToken": securityTypeAccessToken,
	}
)

type securityType uint

func allSecurityTypes() []securityType {
	return []securityType{
		securityTypeSsh,
		securityTypeAccessToken,
	}
}

func securityTypeOf(s string) (securityType, error) {
	securityType, ok := stringToSecurityType[s]
	if !ok {
		return 0, errors.New(unknownSecurityType(s))
	}
	return securityType, nil
}

func (this securityType) String() string {
	if int(this) < lenSecurityTypeToString {
		return securityTypeToString[this]
	}
	panic(unknownSecurityType(this))
}

func unknownSecurityType(unknownSecurityType interface{}) string {
	return fmt.Sprintf("Unknown securityType: %v", unknownSecurityType)
}
