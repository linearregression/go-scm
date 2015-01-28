package scm

import "github.com/peter-edge/exec"

type baseCheckoutOptions struct {
	baseCloneArgs []string
	url           string
	branch        string
	commitId      string
}

type baseClient interface {
	checkoutWithExecutor(exec.Executor, *baseCheckoutOptions, string) error
	ignoreCheckoutFilePatterns(exec.ReadFileManager) []string
}

type baseGitClient struct {
}

func newBaseGitClient() *baseGitClient {
	return &baseGitClient{}
}

func (this *baseGitClient) checkoutWithExecutor(executor exec.Executor, baseCheckoutOptions *baseCheckoutOptions, path string) error {
	if err := executor.Execute(
		&exec.Cmd{
			// TODO(peter): if the commit id is more than 50 back, the checkout will fail
			Args: append(baseCheckoutOptions.baseCloneArgs, []string{"--branch", baseCheckoutOptions.branch, "--depth", "50", "--recursive", baseCheckoutOptions.url, path}...),
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

func (this *baseGitClient) ignoreCheckoutFilePatterns(readFileManager exec.ReadFileManager) []string {
	return []string{
		".git",
		".gitignore",
	}
}

type baseHgClient struct {
}

func newBaseHgClient() *baseHgClient {
	return &baseHgClient{}
}

// return a reader for a tarball for this checkout
func (this *baseHgClient) checkoutWithExecutor(executor exec.Executor, baseCheckoutOptions *baseCheckoutOptions, path string) error {
	if err := executor.Execute(
		&exec.Cmd{
			Args: append(baseCheckoutOptions.baseCloneArgs, []string{"hg", "clone", baseCheckoutOptions.url, path}...),
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

func (this *baseHgClient) ignoreCheckoutFilePatterns(readFileManager exec.ReadFileManager) []string {
	return []string{
		".hg",
		".hgignore",
		".hgsigs",
		".hgtags",
	}
}
