package scm

import (
	"io"
	"strings"

	"github.com/peter-edge/exec"
	tarexec "github.com/peter-edge/tar/exec"
)

const (
	clonePath = "clone"
)

type baseCheckoutOptions struct {
	url                 string
	branch              string
	commitId            string
	ignoreCheckoutFiles bool
}

type baseCheckoutClient interface {
	exec.ExecutorReadFileManagerProvider
	checkoutWithExecutor(exec.Executor, *baseCheckoutOptions, string) error
	ignoreCheckoutFilePatterns(exec.ReadFileManager) []string
}

func checkout(baseCheckoutClient baseCheckoutClient, baseCheckoutOptions *baseCheckoutOptions) (CheckoutTarball, error) {
	client, err := baseCheckoutClient.NewTempDirExecutorReadFileManager()
	if err != nil {
		return nil, err
	}
	err = baseCheckoutClient.checkoutWithExecutor(client, baseCheckoutOptions, clonePath)
	if err != nil {
		return nil, err
	}
	var ignoreCheckoutFilePatterns []string = nil
	if baseCheckoutOptions.ignoreCheckoutFiles {
		ignoreCheckoutFilePatterns = baseCheckoutClient.ignoreCheckoutFilePatterns(client)
	}

	reader, err := tarFiles(client, ignoreCheckoutFilePatterns, clonePath)
	if err != nil {
		return nil, err
	}
	if err := client.Destroy(); err != nil {
		return nil, err
	}
	return newCheckoutTarball(reader, baseCheckoutOptions.branch, baseCheckoutOptions.commitId), nil
}

func tarFiles(readFileManager exec.ReadFileManager, ignoreCheckoutFilePatterns []string, path string) (io.Reader, error) {
	fileList, err := readFileManager.ListRegularFiles(path)
	if err != nil {
		return nil, err
	}
	if ignoreCheckoutFilePatterns != nil && len(ignoreCheckoutFilePatterns) > 0 {
		filterFileList := make([]string, 0)
		for _, file := range fileList {
			matches, err := fileMatches(readFileManager, ignoreCheckoutFilePatterns, file, path)
			if err != nil {
				return nil, err
			}
			if !matches {
				filterFileList = append(filterFileList, file)
			}
		}
		fileList = filterFileList
	}
	return tarexec.NewTarClient(readFileManager, nil).Tar(fileList, path)
}

func fileMatches(readFileManager exec.ReadFileManager, patterns []string, path string, basePath string) (bool, error) {
	for _, pattern := range patterns {
		if strings.HasPrefix(path, readFileManager.Join(basePath, pattern)) {
			return true, nil
		}
	}
	return false, nil
}
