package scm

import "io"

type SecurityOptions interface {
	Type() SecurityType
}

type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}

func (this *SshSecurityOptions) Type() SecurityType {
	return SecurityTypeSsh
}

type AccessTokenSecurityOptions struct {
	AccessToken string
}

func (this *AccessTokenSecurityOptions) Type() SecurityType {
	return SecurityTypeAccessToken
}
