package scm

import (
	"io"

	"github.com/peter-edge/exec"
)

type GitCheckoutOptions struct {
	Url      string
	Branch   string
	CommitId string
}

type GithubCheckoutOptions struct {
	User        string
	Repository  string
	Branch      string
	CommitId    string
	AccessToken string
}

type HgCheckoutOptions struct {
	Url                 string
	ChangesetId         string
	IgnoreCheckoutFiles bool
}

type ClientOptions struct {
	IgnoreCheckoutFiles bool
}

type Client interface {
	CheckoutGitTarball(*GitCheckoutOptions) (io.Reader, error)
	CheckoutGithubTarball(*GithubCheckoutOptions) (io.Reader, error)
	CheckoutHgTarball(*HgCheckoutOptions) (io.Reader, error)
}

func NewClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider, clientOptions *ClientOptions) Client {
	return newClient(executorReadFileManagerProvider, clientOptions)
}
