package scm

import (
	"io"

	"github.com/peter-edge/exec"
	tarexec "github.com/peter-edge/tar/exec"
)

type baseGitCheckoutClient struct {
	exec.ExecutorReadFileManagerProvider
}

func newBaseGitCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) *baseGitCheckoutClient {
	return &baseGitCheckoutClient{executorReadFileManagerProvider}
}

// return a reader for a tarball for this checkout
func (this *baseGitCheckoutClient) checkout(url string, branch string, commitId string) (io.Reader, error) {
	client, err := this.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	clonePath := client.Join(client.DirPath(), "clone")
	err = client.Execute(
		&exec.Cmd{
			// TODO(peter): if the commit id is more than 50 back, the checkout will fail
			Args: []string{"git", "clone", "--branch", branch, "--depth", "50", "--recursive", url, clonePath},
		},
	)()
	if err != nil {
		return nil, err
	}
	fileList, err := client.ListRegularFiles("clone")
	if err != nil {
		return nil, err
	}
	reader, err := tarexec.NewTarClient(client, nil).Tar(fileList, "clone")
	if err != nil {
		return nil, err
	}
	err = client.Destroy()
	if err != nil {
		return nil, err
	}
	return reader, nil
}
