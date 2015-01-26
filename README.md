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

#### type CheckoutTarball

```go
type CheckoutTarball interface {
	io.Reader
	Branch() string
	CommitId() string
}
```


#### type GitCheckoutClient

```go
type GitCheckoutClient interface {
	CheckoutTarball(*GitCheckoutOptions) (CheckoutTarball, error)
}
```


#### func  NewGitCheckoutClient

```go
func NewGitCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) GitCheckoutClient
```

#### type GitCheckoutOptions

```go
type GitCheckoutOptions struct {
	Url      string
	Branch   string
	CommitId string
}
```


#### type GithubCheckoutClient

```go
type GithubCheckoutClient interface {
	CheckoutTarball(*GithubCheckoutOptions) (CheckoutTarball, error)
}
```


#### func  NewGithubCheckoutClient

```go
func NewGithubCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) GithubCheckoutClient
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
