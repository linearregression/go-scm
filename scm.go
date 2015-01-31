package scm

import (
	"io"

	"github.com/peter-edge/go-exec"
)

type CheckoutOptions interface {
	Type() CheckoutType
}

type GitCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

func (this *GitCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeGit
}

type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

func (this *GithubCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeGithub
}

type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions SecurityOptions
}

func (this *HgCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeHg
}

type BitbucketCheckoutOptions struct {
	BitbucketType   BitbucketType
	User            string
	Repository      string
	Branch          string // only set if BitbucketType == BitbucketTypeGit
	CommitId        string // only set if BitbucketType == BitbucketTypeGit
	ChangesetId     string // only set if BitbucketType == BitbucketTypeHg
	SecurityOptions SecurityOptions
}

func (this *BitbucketCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeBitbucket
}

type ClientOptions struct {
	IgnoreCheckoutFiles bool
}

type Client interface {
	CheckoutTarball(CheckoutOptions) (io.Reader, error)
}

func NewClient(execClientProvider exec.ClientProvider, clientOptions *ClientOptions) Client {
	return newClient(execClientProvider, clientOptions)
}
