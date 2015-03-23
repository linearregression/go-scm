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
cat cmd/go-scm/_testdata/external_checkout_options.json | go-scm
```

To run with Docker:

```bash
make container
cat cmd/go-scm/_testdata/external_checkout_options.json | docker run -i pedge/goscm
```

To run with Docker using etcd:

```bash
make container
docker run pedge/goscm --etcd_url=http://host:port --etcd_input_key=key_for_json_external_checkout_options --etcd_output_key=key_for_path_within_docker_volume_of_result
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
func (this *AccessTokenSecurityOptions) Type() SecurityOptionsType
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
func (this *BitbucketCheckoutOptions) Type() CheckoutOptionsType
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

#### func (BitbucketType) String

```go
func (this BitbucketType) String() string
```

#### type CheckoutOptions

```go
type CheckoutOptions interface {
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
var (
	CheckoutOptionsTypeGit       CheckoutOptionsType = 0
	CheckoutOptionsTypeGithub    CheckoutOptionsType = 1
	CheckoutOptionsTypeHg        CheckoutOptionsType = 2
	CheckoutOptionsTypeBitbucket CheckoutOptionsType = 3
)
```

#### func (CheckoutOptionsType) String

```go
func (this CheckoutOptionsType) String() string
```

#### type Client

```go
type Client interface {
	CheckoutTarball(checkoutOptions CheckoutOptions) (io.Reader, error)
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


#### type DirectClient

```go
type DirectClient interface {
	Checkout(checkoutOptions CheckoutOptions, executor exec.Executor, path string) error
}
```


#### func  NewDirectClient

```go
func NewDirectClient(execClientProvider exec.ClientProvider) DirectClient
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
	CheckoutTarball(externalCheckoutOptions *ExternalCheckoutOptions) (io.Reader, error)
}
```


#### func  NewExternalClient

```go
func NewExternalClient(client Client) ExternalClient
```

#### type ExternalDirectClient

```go
type ExternalDirectClient interface {
	Checkout(externalCheckoutOptions *ExternalCheckoutOptions, executor exec.Executor, path string) error
}
```


#### func  NewExternalDirectClient

```go
func NewExternalDirectClient(directClient DirectClient) ExternalDirectClient
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


#### func (*HgCheckoutOptions) Type

```go
func (this *HgCheckoutOptions) Type() CheckoutOptionsType
```

#### type SecurityOptions

```go
type SecurityOptions interface {
	Type() SecurityOptionsType
}
```


#### type SecurityOptionsType

```go
type SecurityOptionsType uint
```


```go
var (
	SecurityOptionsTypeSsh         SecurityOptionsType = 0
	SecurityOptionsTypeAccessToken SecurityOptionsType = 1
)
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


#### func (*SshSecurityOptions) Type

```go
func (this *SshSecurityOptions) Type() SecurityOptionsType
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
	ValidationErrorTypeSecurityNotImplementedForCheckoutOptionsType ValidationErrorType = "SecurityNotImplementedForCheckoutOptionsType"
	ValidationErrorTypeUnknownCheckoutOptionsType                   ValidationErrorType = "UnknownCheckoutOptionsType"
	ValidationErrorTypeUnknownSecurityOptionsType                   ValidationErrorType = "UnknownSecurityOptionsType"
	ValidationErrorTypeUnknownBitbucketType                  ValidationErrorType = "UnknownBitbucketType"
)
```
