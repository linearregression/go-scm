package scm

import "io"

type SecurityOptions interface {
	SecurityType() SecurityType
}

type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}

func (this *SshSecurityOptions) SecurityType() SecurityType {
	return SecurityTypeSsh
}

type AccessTokenSecurityOptions struct {
	AccessToken string
}

func (this *AccessTokenSecurityOptions) SecurityType() SecurityType {
	return SecurityTypeAccessToken
}
