[![Codeship Status](http://img.shields.io/codeship/34b974b0-6dfa-0132-51b4-66f2bf861e14/master.svg?style=flat-square)](https://codeship.com/projects/59076)
[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/scm)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/scm/blob/master/LICENSE)

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

## Usage

```go
var (
	ErrRequiredFieldMissing = errors.New("scm: required field missing")
)
```

#### type Client

```go
type Client interface {
	CheckoutGitTarball(*GitCheckoutOptions) (io.Reader, error)
	CheckoutGithubTarball(*GithubCheckoutOptions) (io.Reader, error)
	CheckoutHgTarball(*HgCheckoutOptions) (io.Reader, error)
}
```


#### func  NewClient

```go
func NewClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider, clientOptions *ClientOptions) Client
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
	Url      string
	Branch   string
	CommitId string
}
```


#### type GithubCheckoutOptions

```go
type GithubCheckoutOptions struct {
	User        string
	Repository  string
	Branch      string
	CommitId    string
	AccessToken string
}
```


#### type HgCheckoutOptions

```go
type HgCheckoutOptions struct {
	Url                 string
	ChangesetId         string
	IgnoreCheckoutFiles bool
}
```
