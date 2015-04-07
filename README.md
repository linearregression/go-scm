[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/go-scm)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/go-scm/blob/master/LICENSE)

Scm utilities for Go.

## Installation
```bash
go get -u github.com/peter-edge/go-scm
```

## Import
```go
import (
    "github.com/peter-edge/go-scm"
)
```

## Notes

To run:

```bash
make install
cat cmd/_testdata/external_checkout_options.json | scm-clone
```

To run with Docker:

```bash
make container
cat cmd/_testdata/external_checkout_options.json | docker run -i pedge/goscmclone
```
Git SSH requires Git 2.3.0.

## Usage

```go
var (
	AllRecordConverters = []record.RecordConverter{
		&CloneRecordRecordConverter{},
		&TarballRecordRecordConverter{},
	}
)
```

#### func  AllCheckoutOptionsTypes

```go
func AllCheckoutOptionsTypes() []CheckoutOptionsType
```

#### func  AllSecurityOptionsTypes

```go
func AllSecurityOptionsTypes() []SecurityOptionsType
```

#### func  Checkout

```go
func Checkout(
	execClientProvider exec.ClientProvider,
	checkoutOptions CheckoutOptions,
	executor exec.Executor,
	path string,
) error
```

#### func  CheckoutOptionsSwitch

```go
func CheckoutOptionsSwitch(
	checkoutOptions CheckoutOptions,
	gitCheckoutOptionsFunc func(gitCheckoutOptions *GitCheckoutOptions) error,
	githubCheckoutOptionsFunc func(githubCheckoutOptions *GithubCheckoutOptions) error,
	hgCheckoutOptionsFunc func(hgCheckoutOptions *HgCheckoutOptions) error,
	bitbucketGitCheckoutOptionsFunc func(bitbucketGitCheckoutOptions *BitbucketGitCheckoutOptions) error,
	bitbucketHgCheckoutOptionsFunc func(bitbucketHgCheckoutOptions *BitbucketHgCheckoutOptions) error,
) error
```

#### func  CheckoutTarball

```go
func CheckoutTarball(
	execClientProvider exec.ClientProvider,
	checkoutOptions CheckoutOptions,
	ignoreCheckoutFiles bool,
) (io.Reader, error)
```

#### func  SecurityOptionsSwitch

```go
func SecurityOptionsSwitch(
	securityOptions SecurityOptions,
	sshSecurityOptionsFunc func(sshSecurityOptions *SshSecurityOptions) error,
	accessTokenSecurityOptionsFunc func(accessTokenSecurityOptions *AccessTokenSecurityOptions) error,
) error
```

#### type AccessTokenSecurityOptions

```go
type AccessTokenSecurityOptions struct {
	AccessToken string
}
```

@gen-enumtype SecurityOptions accessToken 1

#### func (*AccessTokenSecurityOptions) String

```go
func (this *AccessTokenSecurityOptions) String() string
```

#### func (*AccessTokenSecurityOptions) Type

```go
func (this *AccessTokenSecurityOptions) Type() SecurityOptionsType
```

#### type BitbucketGitCheckoutOptions

```go
type BitbucketGitCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}
```

@gen-enumtype CheckoutOptions bitbucketGit 3

#### func (*BitbucketGitCheckoutOptions) String

```go
func (this *BitbucketGitCheckoutOptions) String() string
```

#### func (*BitbucketGitCheckoutOptions) Type

```go
func (this *BitbucketGitCheckoutOptions) Type() CheckoutOptionsType
```

#### type BitbucketHgCheckoutOptions

```go
type BitbucketHgCheckoutOptions struct {
	User            string
	Repository      string
	ChangesetId     string
	SecurityOptions SecurityOptions
}
```

@gen-enumtype CheckoutOptions bitbucketHg 4

#### func (*BitbucketHgCheckoutOptions) String

```go
func (this *BitbucketHgCheckoutOptions) String() string
```

#### func (*BitbucketHgCheckoutOptions) Type

```go
func (this *BitbucketHgCheckoutOptions) Type() CheckoutOptionsType
```

#### type CheckoutOptions

```go
type CheckoutOptions interface {
	fmt.Stringer
	Type() CheckoutOptionsType
}
```


#### func  ConvertExternalCheckoutOptions

```go
func ConvertExternalCheckoutOptions(externalCheckoutOptions *ExternalCheckoutOptions) (CheckoutOptions, error)
```

#### type CheckoutOptionsType

```go
type CheckoutOptionsType uint
```


```go
var CheckoutOptionsTypeBitbucketGit CheckoutOptionsType = 3
```

```go
var CheckoutOptionsTypeBitbucketHg CheckoutOptionsType = 4
```

```go
var CheckoutOptionsTypeGit CheckoutOptionsType = 0
```

```go
var CheckoutOptionsTypeGithub CheckoutOptionsType = 1
```

```go
var CheckoutOptionsTypeHg CheckoutOptionsType = 2
```

#### func  CheckoutOptionsTypeOf

```go
func CheckoutOptionsTypeOf(s string) (CheckoutOptionsType, error)
```

#### func (CheckoutOptionsType) Handle

```go
func (this CheckoutOptionsType) Handle(
	checkoutOptionsTypeGitFunc func() error,
	checkoutOptionsTypeGithubFunc func() error,
	checkoutOptionsTypeHgFunc func() error,
	checkoutOptionsTypeBitbucketGitFunc func() error,
	checkoutOptionsTypeBitbucketHgFunc func() error,
) error
```

#### func (CheckoutOptionsType) NewCheckoutOptions

```go
func (this CheckoutOptionsType) NewCheckoutOptions(
	gitCheckoutOptionsFunc func() (*GitCheckoutOptions, error),
	githubCheckoutOptionsFunc func() (*GithubCheckoutOptions, error),
	hgCheckoutOptionsFunc func() (*HgCheckoutOptions, error),
	bitbucketGitCheckoutOptionsFunc func() (*BitbucketGitCheckoutOptions, error),
	bitbucketHgCheckoutOptionsFunc func() (*BitbucketHgCheckoutOptions, error),
) (CheckoutOptions, error)
```

#### func (CheckoutOptionsType) Produce

```go
func (this CheckoutOptionsType) Produce(
	checkoutOptionsTypeGitFunc func() (interface{}, error),
	checkoutOptionsTypeGithubFunc func() (interface{}, error),
	checkoutOptionsTypeHgFunc func() (interface{}, error),
	checkoutOptionsTypeBitbucketGitFunc func() (interface{}, error),
	checkoutOptionsTypeBitbucketHgFunc func() (interface{}, error),
) (interface{}, error)
```

#### func (CheckoutOptionsType) String

```go
func (this CheckoutOptionsType) String() string
```

#### type CloneRecord

```go
type CloneRecord struct {
	Path string
}
```

@gen-record

#### func (*CloneRecord) ReadableName

```go
func (this *CloneRecord) ReadableName() string
```

#### type CloneRecordRecordConverter

```go
type CloneRecordRecordConverter struct{}
```


#### func (*CloneRecordRecordConverter) FromMap

```go
func (this *CloneRecordRecordConverter) FromMap(m map[string]string) (record.RecordObject, error)
```

#### func (*CloneRecordRecordConverter) ToMap

```go
func (this *CloneRecordRecordConverter) ToMap(object record.RecordObject) (map[string]string, error)
```

#### func (*CloneRecordRecordConverter) Type

```go
func (this *CloneRecordRecordConverter) Type() reflect.Type
```

#### type ExternalCheckoutOptions

```go
type ExternalCheckoutOptions struct {
	Type            string                   `json:"type,omitempty" yaml:"type,omitempty"`
	User            string                   `json:"user,omitempty" yaml:"user,omitempty"`
	Host            string                   `json:"host,omitempty" yaml:"host,omitempty"`
	Path            string                   `json:"path,omitempty" yaml:"path,omitempty"`
	Repository      string                   `json:"repository,omitempty" yaml:"repository,omitempty"`
	Branch          string                   `json:"branch,omitempty" yaml:"branch,omitempty"`
	CommitId        string                   `json:"commit_id,omitempty" yaml:"commit_id,omitempty"`
	ChangesetId     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
	SecurityOptions *ExternalSecurityOptions `json:"security_options,omitempty" yaml:"security_options,omitempty"`
}
```


#### type ExternalSecurityOptions

```go
type ExternalSecurityOptions struct {
	Type                  string `json:"type,omitempty" yaml:"type,omitempty"`
	StrictHostKeyChecking bool   `json:"strict_host_key_checking,omitempty" yaml:"strict_host_key_checking,omitempty"`
	PrivateKey            string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	AccessToken           string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
}
```


#### type GitCheckoutOptions

```go
type GitCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}
```

@gen-enumtype CheckoutOptions git 0

#### func (*GitCheckoutOptions) String

```go
func (this *GitCheckoutOptions) String() string
```

#### func (*GitCheckoutOptions) Type

```go
func (this *GitCheckoutOptions) Type() CheckoutOptionsType
```

#### type GithubCheckoutOptions

```go
type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}
```

@gen-enumtype CheckoutOptions github 1

#### func (*GithubCheckoutOptions) String

```go
func (this *GithubCheckoutOptions) String() string
```

#### func (*GithubCheckoutOptions) Type

```go
func (this *GithubCheckoutOptions) Type() CheckoutOptionsType
```

#### type HgCheckoutOptions

```go
type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions SecurityOptions
}
```

@gen-enumtype CheckoutOptions hg 2

#### func (*HgCheckoutOptions) String

```go
func (this *HgCheckoutOptions) String() string
```

#### func (*HgCheckoutOptions) Type

```go
func (this *HgCheckoutOptions) Type() CheckoutOptionsType
```

#### type SecurityOptions

```go
type SecurityOptions interface {
	fmt.Stringer
	Type() SecurityOptionsType
}
```


#### type SecurityOptionsType

```go
type SecurityOptionsType uint
```


```go
var SecurityOptionsTypeAccessToken SecurityOptionsType = 1
```

```go
var SecurityOptionsTypeSsh SecurityOptionsType = 0
```

#### func  SecurityOptionsTypeOf

```go
func SecurityOptionsTypeOf(s string) (SecurityOptionsType, error)
```

#### func (SecurityOptionsType) Handle

```go
func (this SecurityOptionsType) Handle(
	securityOptionsTypeSshFunc func() error,
	securityOptionsTypeAccessTokenFunc func() error,
) error
```

#### func (SecurityOptionsType) NewSecurityOptions

```go
func (this SecurityOptionsType) NewSecurityOptions(
	sshSecurityOptionsFunc func() (*SshSecurityOptions, error),
	accessTokenSecurityOptionsFunc func() (*AccessTokenSecurityOptions, error),
) (SecurityOptions, error)
```

#### func (SecurityOptionsType) Produce

```go
func (this SecurityOptionsType) Produce(
	securityOptionsTypeSshFunc func() (interface{}, error),
	securityOptionsTypeAccessTokenFunc func() (interface{}, error),
) (interface{}, error)
```

#### func (SecurityOptionsType) String

```go
func (this SecurityOptionsType) String() string
```

#### type SshSecurityOptions

```go
type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}
```

@gen-enumtype SecurityOptions ssh 0

#### func (*SshSecurityOptions) String

```go
func (this *SshSecurityOptions) String() string
```

#### func (*SshSecurityOptions) Type

```go
func (this *SshSecurityOptions) Type() SecurityOptionsType
```

#### type TarballRecord

```go
type TarballRecord struct {
	Path string
}
```

@gen-record

#### func (*TarballRecord) ReadableName

```go
func (this *TarballRecord) ReadableName() string
```

#### type TarballRecordRecordConverter

```go
type TarballRecordRecordConverter struct{}
```


#### func (*TarballRecordRecordConverter) FromMap

```go
func (this *TarballRecordRecordConverter) FromMap(m map[string]string) (record.RecordObject, error)
```

#### func (*TarballRecordRecordConverter) ToMap

```go
func (this *TarballRecordRecordConverter) ToMap(object record.RecordObject) (map[string]string, error)
```

#### func (*TarballRecordRecordConverter) Type

```go
func (this *TarballRecordRecordConverter) Type() reflect.Type
```

#### type ValidationError

```go
type ValidationError interface {
	error
	Type() ValidationErrorType
}
```


#### type ValidationErrorType

```go
type ValidationErrorType string
```


```go
var (
	ValidationErrorTypeRequiredFieldMissing                         ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                          ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType ValidationErrorType = "SecurityNotImplementedForCheckoutOptionsType"
)
```
