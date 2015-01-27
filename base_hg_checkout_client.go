package scm

import "github.com/peter-edge/exec"

type baseHgCheckoutClient struct {
	exec.ExecutorReadFileManagerProvider
}

func newBaseHgCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) *baseHgCheckoutClient {
	return &baseHgCheckoutClient{executorReadFileManagerProvider}
}

// return a reader for a tarball for this checkout
func (this *baseHgCheckoutClient) checkoutWithExecutor(executor exec.Executor, baseCheckoutOptions *baseCheckoutOptions, path string) error {
	if err := executor.Execute(
		&exec.Cmd{
			// TODO(peter): if the commit id is more than 50 back, the checkout will fail
			Args: []string{"hg", "clone", baseCheckoutOptions.url, path},
		},
	)(); err != nil {
		return err
	}
	return executor.Execute(
		&exec.Cmd{
			Args: []string{"hg", "update", "--cwd", path, baseCheckoutOptions.commitId},
		},
	)()
}

func (this *baseHgCheckoutClient) ignoreCheckoutFilePatterns(readFileManager exec.ReadFileManager) []string {
	return []string{
		".hg",
		".hgignore",
		".hgsigs",
		".hgtags",
	}
}
