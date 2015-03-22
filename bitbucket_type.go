package scm

import "fmt"

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

func BitbucketTypeOf(s string) (BitbucketType, error) {
	bitbucketType, ok := stringToBitbucketType[s]
	if !ok {
		return 0, UnknownBitbucketType(s)
	}
	return bitbucketType, nil
}

func (this BitbucketType) String() string {
	if int(this) < lenBitbucketTypeToString {
		return bitbucketTypeToString[this]
	}
	panic(UnknownBitbucketType(this).Error())
}

func UnknownBitbucketType(unknownBitbucketType interface{}) error {
	return fmt.Errorf("scm: unknown BitbucketType: %v", unknownBitbucketType)
}
