package scm

import "io"

type ExternalCheckoutOptions struct {
	Type            string                   `json:"type,omitempty" yaml:"type,omitempty"`
	User            string                   `json:"user,omitempty" yaml:"user,omitempty"`
	Host            string                   `json:"host,omitempty" yaml:"host,omitempty"`
	Path            string                   `json:"path,omitempty" yaml:"path,omitempty"`
	Repository      string                   `json:"repository,omitempty" yaml:"repository,omitempty"`
	Branch          string                   `json:"branch,omitempty" yaml:"branch,omitempty"`
	CommitId        string                   `json:"commit_id,omitempty" yaml:"commit_id,omitempty"`
	BitbucketType   string                   `json:"bitbucket_type,omitempty" yaml:"bitbucket_type,omitempty"`
	ChangesetId     string                   `json:"changeset_id,omitempty" yaml:"changeset_id,omitempty"`
	SecurityOptions *ExternalSecurityOptions `json:"security_options,omitempty" yaml:"security_options,omitempty"`
}

type ExternalSecurityOptions struct {
	Type                  string `json:"type,omitempty" yaml:"type,omitempty"`
	StrictHostKeyChecking bool   `json:"strict_host_key_checking,omitempty" yaml:"strict_host_key_checking,omitempty"`
	PrivateKey            string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
	AccessToken           string `json:"access_token,omitempty" yaml:"access_token,omitempty"`
}

type ExternalClient interface {
	CheckoutTarball(*ExternalCheckoutOptions) (io.Reader, error)
}

func NewExternalClient(client Client) ExternalClient {
	return newExternalClient(client)
}
