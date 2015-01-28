package scm

import "io"

type SecurityOptions interface {
	securityType() securityType
	sshOptions() *SshOptions
	accessTokenOptions() *AccessTokenOptions
}

type SshOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}

type AccessTokenOptions struct {
	AccessToken string
}

type GitSecurityOptions struct {
	securityType_ securityType
	sshOptions_   *SshOptions
}

func NewGitSecurityOptionsSsh(sshOptions *SshOptions) *GitSecurityOptions {
	return &GitSecurityOptions{securityTypeSsh, sshOptions}
}

func (this *GitSecurityOptions) securityType() securityType {
	return this.securityType_
}

func (this *GitSecurityOptions) sshOptions() *SshOptions {
	return this.sshOptions_
}

func (this *GitSecurityOptions) accessTokenOptions() *AccessTokenOptions {
	return nil
}

type GithubSecurityOptions struct {
	securityType_       securityType
	sshOptions_         *SshOptions
	accessTokenOptions_ *AccessTokenOptions
}

func NewGithubSecurityOptionsSsh(sshOptions *SshOptions) *GithubSecurityOptions {
	return &GithubSecurityOptions{securityTypeSsh, sshOptions, nil}
}

func NewGithubSecurityOptionsAccessToken(accessTokenOptions *AccessTokenOptions) *GithubSecurityOptions {
	return &GithubSecurityOptions{securityTypeAccessToken, nil, accessTokenOptions}
}

func (this *GithubSecurityOptions) securityType() securityType {
	return this.securityType_
}

func (this *GithubSecurityOptions) sshOptions() *SshOptions {
	return this.sshOptions_
}

func (this *GithubSecurityOptions) accessTokenOptions() *AccessTokenOptions {
	return this.accessTokenOptions_
}

type HgSecurityOptions struct {
	securityType_ securityType
	sshOptions_   *SshOptions
}

func NewHgSecurityOptionsSsh(sshOptions *SshOptions) *HgSecurityOptions {
	return &HgSecurityOptions{securityTypeSsh, sshOptions}
}

func (this *HgSecurityOptions) securityType() securityType {
	return this.securityType_
}

func (this *HgSecurityOptions) sshOptions() *SshOptions {
	return this.sshOptions_
}

func (this *HgSecurityOptions) accessTokenOptions() *AccessTokenOptions {
	return nil
}
