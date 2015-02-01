package scm

import "io"

func newExternalClient(client Client) *externalClient {
	return &externalClient{client}
}

type externalClient struct {
	client Client
}

func (this *externalClient) CheckoutTarball(externalCheckoutOptions ExternalCheckoutOptions) (io.Reader, error) {
	return nil, nil
}
