package scm

import "fmt"

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

func SecurityTypeOf(s string) (SecurityType, error) {
	SecurityType, ok := stringToSecurityType[s]
	if !ok {
		return 0, UnknownSecurityType(s)
	}
	return SecurityType, nil
}

func (this SecurityType) String() string {
	if int(this) < lenSecurityTypeToString {
		return securityTypeToString[this]
	}
	panic(UnknownSecurityType(this).Error())
}

func UnknownSecurityType(unknownSecurityType interface{}) error {
	return fmt.Errorf("scm: unknown SecurityType: %v", unknownSecurityType)
}
