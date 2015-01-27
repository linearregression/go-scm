package scm

import "github.com/peter-edge/exec"

type baseGitCheckoutClient struct {
	exec.ExecutorReadFileManagerProvider
}

func newBaseGitCheckoutClient(executorReadFileManagerProvider exec.ExecutorReadFileManagerProvider) *baseGitCheckoutClient {
	return &baseGitCheckoutClient{executorReadFileManagerProvider}
}

// return a reader for a tarball for this checkout
func (this *baseGitCheckoutClient) checkoutWithExecutor(executor exec.Executor, baseCheckoutOptions *baseCheckoutOptions, path string) error {
	if err := executor.Execute(
		&exec.Cmd{
			// TODO(peter): if the commit id is more than 50 back, the checkout will fail
			Args: []string{"git", "clone", "--branch", baseCheckoutOptions.branch, "--depth", "50", "--recursive", baseCheckoutOptions.url, path},
		},
	)(); err != nil {
		return err
	}
	return executor.Execute(
		&exec.Cmd{
			Args:   []string{"git", "checkout", "-f", baseCheckoutOptions.commitId},
			SubDir: path,
		},
	)()
}

func (this *baseGitCheckoutClient) ignoreCheckoutFilePatterns(readFileManager exec.ReadFileManager) []string {
	return []string{
		".git",
		".gitignore",
	}
}
