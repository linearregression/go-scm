package scm

import (
	"fmt"

	"github.com/peter-edge/go-stringhelper"
)

type SecurityOptionsType uint

var SecurityOptionsTypeSsh SecurityOptionsType = 0
var SecurityOptionsTypeAccessToken SecurityOptionsType = 1

var securityOptionsTypeToString = map[SecurityOptionsType]string{
	SecurityOptionsTypeSsh: "ssh",
	SecurityOptionsTypeAccessToken: "accessToken",
}

var stringToSecurityOptionsType = map[string]SecurityOptionsType{
	"ssh": SecurityOptionsTypeSsh,
	"accessToken": SecurityOptionsTypeAccessToken,
}

func AllSecurityOptionsTypes() []SecurityOptionsType {
	return []SecurityOptionsType{
		SecurityOptionsTypeSsh,
		SecurityOptionsTypeAccessToken,
	}
}

func SecurityOptionsTypeOf(s string) (SecurityOptionsType, error) {
	securityOptionsType, ok := stringToSecurityOptionsType[s]
	if !ok {
		return 0, newErrorUnknownSecurityOptionsType(s)
	}
	return securityOptionsType, nil
}

func (this SecurityOptionsType) String() string {
	if int(this) < len(securityOptionsTypeToString) {
		 return securityOptionsTypeToString[this]
	}
	panic(newErrorUnknownSecurityOptionsType(this).Error())
}

type SecurityOptions interface {
	fmt.Stringer
	Type() SecurityOptionsType
}

func (this *SshSecurityOptions) Type() SecurityOptionsType {
	return SecurityOptionsTypeSsh
}

func (this *AccessTokenSecurityOptions) Type() SecurityOptionsType {
	return SecurityOptionsTypeAccessToken
}

func (this *SshSecurityOptions) String() string {
	return stringhelper.String(this)
}

func (this *AccessTokenSecurityOptions) String() string {
	return stringhelper.String(this)
}

func SecurityOptionsSwitch(
	securityOptions SecurityOptions,
	sshSecurityOptionsFunc func(sshSecurityOptions *SshSecurityOptions) error,
	accessTokenSecurityOptionsFunc func(accessTokenSecurityOptions *AccessTokenSecurityOptions) error,
) error {
	switch securityOptions.Type() {
	case SecurityOptionsTypeSsh:
		return sshSecurityOptionsFunc(securityOptions.(*SshSecurityOptions))
	case SecurityOptionsTypeAccessToken:
		return accessTokenSecurityOptionsFunc(securityOptions.(*AccessTokenSecurityOptions))
	default:
		return newErrorUnknownSecurityOptionsType(securityOptions.Type())
	}
}

func (this SecurityOptionsType) NewSecurityOptions(
	sshSecurityOptionsFunc func() (*SshSecurityOptions, error),
	accessTokenSecurityOptionsFunc func() (*AccessTokenSecurityOptions, error),
) (SecurityOptions, error) {
	switch this {
	case SecurityOptionsTypeSsh:
		return sshSecurityOptionsFunc()
	case SecurityOptionsTypeAccessToken:
		return accessTokenSecurityOptionsFunc()
	default:
		return nil, newErrorUnknownSecurityOptionsType(this)
	}
}

func (this SecurityOptionsType) Produce(
	securityOptionsTypeSshFunc func() (interface{}, error),
	securityOptionsTypeAccessTokenFunc func() (interface{}, error),
) (interface{}, error) {
	switch this {
	case SecurityOptionsTypeSsh:
		return securityOptionsTypeSshFunc()
	case SecurityOptionsTypeAccessToken:
		return securityOptionsTypeAccessTokenFunc()
	default:
		return nil, newErrorUnknownSecurityOptionsType(this)
	}
}

func (this SecurityOptionsType) Handle(
	securityOptionsTypeSshFunc func() error,
	securityOptionsTypeAccessTokenFunc func() error,
) error {
	switch this {
	case SecurityOptionsTypeSsh:
		return securityOptionsTypeSshFunc()
	case SecurityOptionsTypeAccessToken:
		return securityOptionsTypeAccessTokenFunc()
	default:
		return newErrorUnknownSecurityOptionsType(this)
	}
}

func newErrorUnknownSecurityOptionsType(value interface{}) error {
	return fmt.Errorf("scm: UnknownSecurityOptionsType: %v", value)
}

type CheckoutOptionsType uint

var CheckoutOptionsTypeGit CheckoutOptionsType = 0
var CheckoutOptionsTypeGithub CheckoutOptionsType = 1
var CheckoutOptionsTypeHg CheckoutOptionsType = 2
var CheckoutOptionsTypeBitbucket CheckoutOptionsType = 3

var checkoutOptionsTypeToString = map[CheckoutOptionsType]string{
	CheckoutOptionsTypeGit: "git",
	CheckoutOptionsTypeGithub: "github",
	CheckoutOptionsTypeHg: "hg",
	CheckoutOptionsTypeBitbucket: "bitbucket",
}

var stringToCheckoutOptionsType = map[string]CheckoutOptionsType{
	"git": CheckoutOptionsTypeGit,
	"github": CheckoutOptionsTypeGithub,
	"hg": CheckoutOptionsTypeHg,
	"bitbucket": CheckoutOptionsTypeBitbucket,
}

func AllCheckoutOptionsTypes() []CheckoutOptionsType {
	return []CheckoutOptionsType{
		CheckoutOptionsTypeGit,
		CheckoutOptionsTypeGithub,
		CheckoutOptionsTypeHg,
		CheckoutOptionsTypeBitbucket,
	}
}

func CheckoutOptionsTypeOf(s string) (CheckoutOptionsType, error) {
	checkoutOptionsType, ok := stringToCheckoutOptionsType[s]
	if !ok {
		return 0, newErrorUnknownCheckoutOptionsType(s)
	}
	return checkoutOptionsType, nil
}

func (this CheckoutOptionsType) String() string {
	if int(this) < len(checkoutOptionsTypeToString) {
		 return checkoutOptionsTypeToString[this]
	}
	panic(newErrorUnknownCheckoutOptionsType(this).Error())
}

type CheckoutOptions interface {
	fmt.Stringer
	Type() CheckoutOptionsType
}

func (this *GitCheckoutOptions) Type() CheckoutOptionsType {
	return CheckoutOptionsTypeGit
}

func (this *GithubCheckoutOptions) Type() CheckoutOptionsType {
	return CheckoutOptionsTypeGithub
}

func (this *HgCheckoutOptions) Type() CheckoutOptionsType {
	return CheckoutOptionsTypeHg
}

func (this *BitbucketCheckoutOptions) Type() CheckoutOptionsType {
	return CheckoutOptionsTypeBitbucket
}

func (this *GitCheckoutOptions) String() string {
	return stringhelper.String(this)
}

func (this *GithubCheckoutOptions) String() string {
	return stringhelper.String(this)
}

func (this *HgCheckoutOptions) String() string {
	return stringhelper.String(this)
}

func (this *BitbucketCheckoutOptions) String() string {
	return stringhelper.String(this)
}

func CheckoutOptionsSwitch(
	checkoutOptions CheckoutOptions,
	gitCheckoutOptionsFunc func(gitCheckoutOptions *GitCheckoutOptions) error,
	githubCheckoutOptionsFunc func(githubCheckoutOptions *GithubCheckoutOptions) error,
	hgCheckoutOptionsFunc func(hgCheckoutOptions *HgCheckoutOptions) error,
	bitbucketCheckoutOptionsFunc func(bitbucketCheckoutOptions *BitbucketCheckoutOptions) error,
) error {
	switch checkoutOptions.Type() {
	case CheckoutOptionsTypeGit:
		return gitCheckoutOptionsFunc(checkoutOptions.(*GitCheckoutOptions))
	case CheckoutOptionsTypeGithub:
		return githubCheckoutOptionsFunc(checkoutOptions.(*GithubCheckoutOptions))
	case CheckoutOptionsTypeHg:
		return hgCheckoutOptionsFunc(checkoutOptions.(*HgCheckoutOptions))
	case CheckoutOptionsTypeBitbucket:
		return bitbucketCheckoutOptionsFunc(checkoutOptions.(*BitbucketCheckoutOptions))
	default:
		return newErrorUnknownCheckoutOptionsType(checkoutOptions.Type())
	}
}

func (this CheckoutOptionsType) NewCheckoutOptions(
	gitCheckoutOptionsFunc func() (*GitCheckoutOptions, error),
	githubCheckoutOptionsFunc func() (*GithubCheckoutOptions, error),
	hgCheckoutOptionsFunc func() (*HgCheckoutOptions, error),
	bitbucketCheckoutOptionsFunc func() (*BitbucketCheckoutOptions, error),
) (CheckoutOptions, error) {
	switch this {
	case CheckoutOptionsTypeGit:
		return gitCheckoutOptionsFunc()
	case CheckoutOptionsTypeGithub:
		return githubCheckoutOptionsFunc()
	case CheckoutOptionsTypeHg:
		return hgCheckoutOptionsFunc()
	case CheckoutOptionsTypeBitbucket:
		return bitbucketCheckoutOptionsFunc()
	default:
		return nil, newErrorUnknownCheckoutOptionsType(this)
	}
}

func (this CheckoutOptionsType) Produce(
	checkoutOptionsTypeGitFunc func() (interface{}, error),
	checkoutOptionsTypeGithubFunc func() (interface{}, error),
	checkoutOptionsTypeHgFunc func() (interface{}, error),
	checkoutOptionsTypeBitbucketFunc func() (interface{}, error),
) (interface{}, error) {
	switch this {
	case CheckoutOptionsTypeGit:
		return checkoutOptionsTypeGitFunc()
	case CheckoutOptionsTypeGithub:
		return checkoutOptionsTypeGithubFunc()
	case CheckoutOptionsTypeHg:
		return checkoutOptionsTypeHgFunc()
	case CheckoutOptionsTypeBitbucket:
		return checkoutOptionsTypeBitbucketFunc()
	default:
		return nil, newErrorUnknownCheckoutOptionsType(this)
	}
}

func (this CheckoutOptionsType) Handle(
	checkoutOptionsTypeGitFunc func() error,
	checkoutOptionsTypeGithubFunc func() error,
	checkoutOptionsTypeHgFunc func() error,
	checkoutOptionsTypeBitbucketFunc func() error,
) error {
	switch this {
	case CheckoutOptionsTypeGit:
		return checkoutOptionsTypeGitFunc()
	case CheckoutOptionsTypeGithub:
		return checkoutOptionsTypeGithubFunc()
	case CheckoutOptionsTypeHg:
		return checkoutOptionsTypeHgFunc()
	case CheckoutOptionsTypeBitbucket:
		return checkoutOptionsTypeBitbucketFunc()
	default:
		return newErrorUnknownCheckoutOptionsType(this)
	}
}

func newErrorUnknownCheckoutOptionsType(value interface{}) error {
	return fmt.Errorf("scm: UnknownCheckoutOptionsType: %v", value)
}

