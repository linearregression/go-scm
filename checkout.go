package scm

import "io"

type checkoutTarball struct {
	io.Reader
	branch   string
	commitId string
}

func newCheckoutTarball(tarballReader io.Reader, branch string, commitId string) *checkoutTarball {
	return &checkoutTarball{tarballReader, branch, commitId}
}

func (this *checkoutTarball) Branch() string {
	return this.branch
}

func (this *checkoutTarball) CommitId() string {
	return this.commitId
}
