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

Git SSH requires Git 2.3.0.

## Usage

#### func  ValidBitbucketType

```go
func ValidBitbucketType(s string) bool
```

#### func  ValidCheckoutType

```go
func ValidCheckoutType(s string) bool
```

#### func  ValidSecurityType

```go
func ValidSecurityType(s string) bool
```

#### type AccessTokenSecurityOptions

```go
type AccessTokenSecurityOptions struct {
	AccessToken string
}
```


#### func (*AccessTokenSecurityOptions) SecurityType

```go
func (this *AccessTokenSecurityOptions) SecurityType() SecurityType
```

#### type BitbucketCheckoutOptions

```go
type BitbucketCheckoutOptions struct {
	BitbucketType   BitbucketType
	User            string
	Repository      string
	Branch          string // only set if BitbucketType == BitbucketTypeGit
	CommitId        string // only set if BitbucketType == BitbucketTypeGit
	ChangesetId     string // only set if BitbucketType == BitbucketTypeHg
	SecurityOptions SecurityOptions
}
```


#### func (*BitbucketCheckoutOptions) Type

```go
func (this *BitbucketCheckoutOptions) Type() CheckoutType
```

#### type BitbucketType

```go
type BitbucketType uint
```


```go
var (
	BitbucketTypeGit BitbucketType = 0
	BitbucketTypeHg  BitbucketType = 1
)
```

#### func  BitbucketTypeOf

```go
func BitbucketTypeOf(s string) (BitbucketType, error)
```

#### func (BitbucketType) String

```go
func (this BitbucketType) String() string
```

#### type CheckoutOptions

```go
type CheckoutOptions interface {
	Type() CheckoutType
}
```


#### type CheckoutType

```go
type CheckoutType uint
```


```go
var (
	CheckoutTypeGit       CheckoutType = 0
	CheckoutTypeGithub    CheckoutType = 1
	CheckoutTypeHg        CheckoutType = 2
	CheckoutTypeBitbucket CheckoutType = 3
)
```

#### func  CheckoutTypeOf

```go
func CheckoutTypeOf(s string) (CheckoutType, error)
```

#### func (CheckoutType) String

```go
func (this CheckoutType) String() string
```

#### type Client

```go
type Client interface {
	CheckoutTarball(CheckoutOptions) (io.Reader, error)
}
```


#### func  NewClient

```go
func NewClient(execClientProvider exec.ClientProvider, clientOptions *ClientOptions) Client
```

#### type ClientOptions

```go
type ClientOptions struct {
	IgnoreCheckoutFiles bool
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


#### func (*GitCheckoutOptions) Type

```go
func (this *GitCheckoutOptions) Type() CheckoutType
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


#### func (*GithubCheckoutOptions) Type

```go
func (this *GithubCheckoutOptions) Type() CheckoutType
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


#### func (*HgCheckoutOptions) Type

```go
func (this *HgCheckoutOptions) Type() CheckoutType
```

#### type SecurityOptions

```go
type SecurityOptions interface {
	SecurityType() SecurityType
}
```


#### type SecurityType

```go
type SecurityType uint
```


```go
var (
	SecurityTypeSsh         SecurityType = 0
	SecurityTypeAccessToken SecurityType = 1
)
```

#### func  SecurityTypeOf

```go
func SecurityTypeOf(s string) (SecurityType, error)
```

#### func (SecurityType) String

```go
func (this SecurityType) String() string
```

#### type SshSecurityOptions

```go
type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}
```


#### func (*SshSecurityOptions) SecurityType

```go
func (this *SshSecurityOptions) SecurityType() SecurityType
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
	ValidationErrorTypeRequiredFieldMissing                  ValidationErrorType = "RequiredFieldMissing"
	ValidationErrorTypeFieldShouldNotBeSet                   ValidationErrorType = "FieldShouldNotBeSet"
	ValidationErrorTypeSecurityNotImplementedForCheckoutType ValidationErrorType = "SecurityNotImplementedForCheckoutType"
	ValidationErrorTypeUnknownCheckoutType                   ValidationErrorType = "UnknownCheckoutType"
	ValidationErrorTypeUnknownSecurityType                   ValidationErrorType = "UnknownSecurityType"
	ValidationErrorTypeUnknownBitbucketType                  ValidationErrorType = "UnknownBitbucketType"
)
```
