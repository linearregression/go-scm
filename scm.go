package scm

import (
	"io"

	"github.com/peter-edge/exec"
)

type GitCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	Branch          string
	CommitId        string
	SecurityOptions *GitSecurityOptions
}

type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions *GithubSecurityOptions
}

type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions *HgSecurityOptions
}

type BitbucketCheckoutOptions struct {
	Type            BitbucketType
	User            string
	Repository      string
	Branch          string // only set if BitbucketType == BitbucketTypeGit
	CommitId        string // only set if BitbucketType == BitbucketTypeGit
	ChangesetId     string // only set if BitbucketType == BitbucketTypeHg
	SecurityOptions *BitbucketSecurityOptions
}

type ClientOptions struct {
	IgnoreCheckoutFiles bool
}

type Client interface {
	CheckoutGitTarball(*GitCheckoutOptions) (io.Reader, error)
	CheckoutGithubTarball(*GithubCheckoutOptions) (io.Reader, error)
	CheckoutHgTarball(*HgCheckoutOptions) (io.Reader, error)
	CheckoutBitbucketTarball(*BitbucketCheckoutOptions) (io.Reader, error)
}

func NewClient(execClientProvider exec.ClientProvider, clientOptions *ClientOptions) Client {
	return newClient(execClientProvider, clientOptions)
}
