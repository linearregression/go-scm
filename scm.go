package scm

import (
	"io"

	"github.com/peter-edge/go-exec"
)

//go:generate gen-enumtype

// @gen-enumtype CheckoutOptions git 0
type GitCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions github 1
type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions hg 2
type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions bitbucketGit 3
type BitbucketGitCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

// @gen-enumtype CheckoutOptions bitbucketHg 4
type BitbucketHgCheckoutOptions struct {
	User            string
	Repository      string
	ChangesetId     string
	SecurityOptions SecurityOptions
}

// @gen-enumtype SecurityOptions ssh 0
type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}

// @gen-enumtype SecurityOptions accessToken 1
type AccessTokenSecurityOptions struct {
	AccessToken string
}

type ClientOptions struct {
	IgnoreCheckoutFiles bool
}

type Client interface {
	CheckoutTarball(checkoutOptions CheckoutOptions) (io.Reader, error)
}

func NewClient(execClientProvider exec.ClientProvider, clientOptions *ClientOptions) Client {
	return newClient(execClientProvider, clientOptions)
}

type DirectClient interface {
	Checkout(checkoutOptions CheckoutOptions, executor exec.Executor, path string) error
}

func NewDirectClient(execClientProvider exec.ClientProvider) DirectClient {
	return newClient(execClientProvider, &ClientOptions{IgnoreCheckoutFiles: false})
}
