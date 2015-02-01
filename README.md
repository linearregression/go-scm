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

#### type AccessTokenSecurityOptions

```go
type AccessTokenSecurityOptions struct {
	AccessToken string
}
```


#### func (*AccessTokenSecurityOptions) Type

```go
func (this *AccessTokenSecurityOptions) Type() SecurityType
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
	BitbucketType   string                   `json:"bitbucket_type,omitempty" yaml:"bitbucket_type,omitempty"`
	ChangesetId     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
	SecurityOptions *ExternalSecurityOptions `json:"security_options,omitempty" yaml:"security_options,omitempty"`
}
```


#### type ExternalClient

```go
type ExternalClient interface {
	CheckoutTarball(*ExternalCheckoutOptions) (io.Reader, error)
}
```


#### func  NewExternalClient

```go
func NewExternalClient(client Client) ExternalClient
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
	Type() SecurityType
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

#### type SshSecurityOptions

```go
type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}
```


#### func (*SshSecurityOptions) Type

```go
func (this *SshSecurityOptions) Type() SecurityType
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
