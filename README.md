[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/scm)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/scm/blob/master/LICENSE)

Scm utilities for Go.

## Installation
```bash
go get -u github.com/peter-edge/scm
```

## Import
```go
import (
    "github.com/peter-edge/scm"
)
```

Git SSH requires Git 2.3.0.

## Usage

```go
const (
	BitbucketTypeGit = iota
	BitbucketTypeHg
)
```

```go
var (
	ErrNil                    = errors.New("scm: nil")
	ErrWrongSecurityType      = errors.New("scm: wrong security type")
	ErrRequiredFieldMissing   = errors.New("scm: required field missing")
	ErrFieldShouldNotBeSet    = errors.New("scm: field should not be set")
	ErrSecurityNotImplemented = errors.New("scm: security not implemented")
	ErrUnknownBitbucketType   = errors.New("scm: unknown BitbucketType")
)
```

#### func  AllBitbucketTypes

```go
func AllBitbucketTypes() []BitbucketType
```

#### type AccessTokenOptions

```go
type AccessTokenOptions struct {
	AccessToken string
}
```


#### type BitbucketCheckoutOptions

```go
type BitbucketCheckoutOptions struct {
	Type            BitbucketType
	User            string
	Repository      string
	Branch          string // only set if BitbucketType == BitbucketTypeGit
	CommitId        string // only set if BitbucketType == BitbucketTypeGit
	ChangesetId     string // only set if BitbucketType == BitbucketTypeHg
	SecurityOptions *BitbucketSecurityOptions
}
```


#### type BitbucketSecurityOptions

```go
type BitbucketSecurityOptions struct {
}
```


#### func  NewBitbucketSecurityOptionsSsh

```go
func NewBitbucketSecurityOptionsSsh(sshOptions *SshOptions) *BitbucketSecurityOptions
```

#### type BitbucketType

```go
type BitbucketType uint
```


#### func  BitbucketTypeOf

```go
func BitbucketTypeOf(s string) (BitbucketType, error)
```

#### func (BitbucketType) String

```go
func (this BitbucketType) String() string
```

#### type Client

```go
type Client interface {
	CheckoutGitTarball(*GitCheckoutOptions) (io.Reader, error)
	CheckoutGithubTarball(*GithubCheckoutOptions) (io.Reader, error)
	CheckoutHgTarball(*HgCheckoutOptions) (io.Reader, error)
	CheckoutBitbucketTarball(*BitbucketCheckoutOptions) (io.Reader, error)
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
	SecurityOptions *GitSecurityOptions
}
```


#### type GitSecurityOptions

```go
type GitSecurityOptions struct {
}
```


#### func  NewGitSecurityOptionsSsh

```go
func NewGitSecurityOptionsSsh(sshOptions *SshOptions) *GitSecurityOptions
```

#### type GithubCheckoutOptions

```go
type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions *GithubSecurityOptions
}
```


#### type GithubSecurityOptions

```go
type GithubSecurityOptions struct {
}
```


#### func  NewGithubSecurityOptionsAccessToken

```go
func NewGithubSecurityOptionsAccessToken(accessTokenOptions *AccessTokenOptions) *GithubSecurityOptions
```

#### func  NewGithubSecurityOptionsSsh

```go
func NewGithubSecurityOptionsSsh(sshOptions *SshOptions) *GithubSecurityOptions
```

#### type HgCheckoutOptions

```go
type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions *HgSecurityOptions
}
```


#### type HgSecurityOptions

```go
type HgSecurityOptions struct {
}
```


#### func  NewHgSecurityOptionsSsh

```go
func NewHgSecurityOptionsSsh(sshOptions *SshOptions) *HgSecurityOptions
```

#### type SecurityOptions

```go
type SecurityOptions interface {
	// contains filtered or unexported methods
}
```


#### type SshOptions

```go
type SshOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}
```
