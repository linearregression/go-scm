package scm

import (
	"errors"
	"fmt"
)

var (
	CheckoutTypeGit       CheckoutType = 0
	CheckoutTypeGithub    CheckoutType = 1
	CheckoutTypeHg        CheckoutType = 2
	CheckoutTypeBitbucket CheckoutType = 3

	checkoutTypeToString = map[CheckoutType]string{
		CheckoutTypeGit:       "git",
		CheckoutTypeGithub:    "github",
		CheckoutTypeHg:        "hg",
		CheckoutTypeBitbucket: "bitbucket",
	}
	lenCheckoutTypeToString = len(checkoutTypeToString)
	stringToCheckoutType    = map[string]CheckoutType{
		"git":       CheckoutTypeGit,
		"github":    CheckoutTypeGithub,
		"hg":        CheckoutTypeHg,
		"bitbucket": CheckoutTypeBitbucket,
	}
)

type CheckoutType uint

func validCheckoutType(s string) bool {
	_, ok := stringToCheckoutType[s]
	return ok
}

func checkoutTypeOf(s string) (CheckoutType, error) {
	checkoutType, ok := stringToCheckoutType[s]
	if !ok {
		return 0, errors.New(unknownCheckoutType(s))
	}
	return checkoutType, nil
}

func (this CheckoutType) String() string {
	if int(this) < lenCheckoutTypeToString {
		return checkoutTypeToString[this]
	}
	panic(unknownCheckoutType(this))
}

func unknownCheckoutType(unknownCheckoutType interface{}) string {
	return fmt.Sprintf("scm: unknown CheckoutType: %v", unknownCheckoutType)
}
