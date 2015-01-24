package scm

import "io"

type CheckoutTarball interface {
	io.Reader
	Branch() string
	CommitId() string
}
