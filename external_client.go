package scm

import (
	"io"

	"github.com/peter-edge/go-exec"
)

type externalClient struct {
	client Client
}

func newExternalClient(client Client) *externalClient {
	return &externalClient{client}
}

func (this *externalClient) CheckoutTarball(externalCheckoutOptions *ExternalCheckoutOptions) (io.Reader, error) {
	checkoutOptions, err := ConvertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return nil, err
	}
	return this.client.CheckoutTarball(checkoutOptions)
}

type externalDirectClient struct {
	directClient DirectClient
}

func newExternalDirectClient(directClient DirectClient) *externalDirectClient {
	return &externalDirectClient{directClient}
}

func (this *externalDirectClient) Checkout(externalCheckoutOptions *ExternalCheckoutOptions, executor exec.Executor, path string) error {
	checkoutOptions, err := ConvertExternalCheckoutOptions(externalCheckoutOptions)
	if err != nil {
		return err
	}
	return this.directClient.Checkout(checkoutOptions, executor, path)
}
