package scm

import (
	"errors"
	"fmt"
)

const (
	BitbucketTypeGit = iota
	BitbucketTypeHg
)

var (
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

func AllBitbucketTypes() []BitbucketType {
	return []BitbucketType{
		BitbucketTypeGit,
		BitbucketTypeHg,
	}
}

func BitbucketTypeOf(s string) (BitbucketType, error) {
	bitbucketType, ok := stringToBitbucketType[s]
	if !ok {
		return 0, errors.New(unknownBitbucketType(s))
	}
	return bitbucketType, nil
}

func (this BitbucketType) String() string {
	if int(this) < lenBitbucketTypeToString {
		return bitbucketTypeToString[this]
	}
	panic(unknownBitbucketType(this))
}

func unknownBitbucketType(unknownBitbucketType interface{}) string {
	return fmt.Sprintf("Unknown BitbucketType: %v", unknownBitbucketType)
}
