package scm

import "fmt"

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

func CheckoutTypeOf(s string) (CheckoutType, error) {
	checkoutType, ok := stringToCheckoutType[s]
	if !ok {
		return 0, UnknownCheckoutType(s)
	}
	return checkoutType, nil
}

func (this CheckoutType) String() string {
	if int(this) < lenCheckoutTypeToString {
		return checkoutTypeToString[this]
	}
	panic(UnknownCheckoutType(this).Error())
}

func UnknownCheckoutType(unknownCheckoutType interface{}) error {
	return fmt.Errorf("scm: unknown CheckoutType: %v", unknownCheckoutType)
}
