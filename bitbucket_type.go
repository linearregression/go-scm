package scm

import (
	"errors"
	"fmt"
)

var (
	BitbucketTypeGit BitbucketType = 0
	BitbucketTypeHg  BitbucketType = 1

	bitbucketTypeToString = map[BitbucketType]string{
		BitbucketTypeGit: "git",
		BitbucketTypeHg:  "hg",
	}
	lenBitbucketTypeToString = len(bitbucketTypeToString)
	stringToBitbucketType    = map[string]BitbucketType{
		"git": BitbucketTypeGit,
		"hg":  BitbucketTypeHg,
	}
)

type BitbucketType uint

func validBitbucketType(s string) bool {
	_, ok := stringToBitbucketType[s]
	return ok
}

func bitbucketTypeOf(s string) (BitbucketType, error) {
	bitbucketType, ok := stringToBitbucketType[s]
	if !ok {
		return 0, errors.New(unknownBitbucketType(s))
	}
	return bitbucketType, nil
}

func (this BitbucketType) string() string {
	if int(this) < lenBitbucketTypeToString {
		return bitbucketTypeToString[this]
	}
	panic(unknownBitbucketType(this))
}

func unknownBitbucketType(unknownBitbucketType interface{}) string {
	return fmt.Sprintf("scm: unknown BitbucketType: %v", unknownBitbucketType)
}
