package scm

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/peter-edge/exec"
	execos "github.com/peter-edge/exec/os"
	tarexec "github.com/peter-edge/tar/exec"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testSmartystreetsCommitId = "a40e854c17df0b1a98c90c250dc20e6cb2474dfa"
)

type Suite struct {
	suite.Suite

	clientProvider exec.ClientProvider
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (this *Suite) SetupSuite() {
}

func (this *Suite) SetupTest() {
	clientProvider, err := execos.NewClientProvider()
	require.NoError(this.T(), err)
	this.clientProvider = clientProvider
}

func (this *Suite) TearDownTest() {
	require.NoError(this.T(), this.clientProvider.Destroy())
}

func (this *Suite) TearDownSuite() {
}

func (this *Suite) TestGit() {
	this.testGit(false)
}

func (this *Suite) TestGitIgnore() {
	this.testGit(true)
}

func (this *Suite) testGit(ignoreCheckoutFiles bool) {
	checkoutTarball, err := NewGitCheckoutClient(this.clientProvider).CheckoutTarball(
		&GitCheckoutOptions{
			Url:                 "https://github.com/peter-edge/smartystreets.git",
			Branch:              "master",
			CommitId:            testSmartystreetsCommitId,
			IgnoreCheckoutFiles: ignoreCheckoutFiles,
		},
	)
	require.NoError(this.T(), err)
	this.testSmartystreetsCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) TestGithub() {
	this.testGithub(false)
}

func (this *Suite) TestGithubIgnore() {
	this.testGithub(true)
}

func (this *Suite) testGithub(ignoreCheckoutFiles bool) {
	checkoutTarball, err := NewGithubCheckoutClient(this.clientProvider).CheckoutTarball(
		&GithubCheckoutOptions{
			User:                "peter-edge",
			Repository:          "smartystreets_ruby",
			Branch:              "master",
			CommitId:            testSmartystreetsCommitId,
			IgnoreCheckoutFiles: ignoreCheckoutFiles,
		},
	)
	require.NoError(this.T(), err)
	this.testSmartystreetsCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) testSmartystreetsCheckoutTarball(checkoutTarball CheckoutTarball, ignoreCheckoutFiles bool) {
	clientProvider, err := execos.NewClientProvider()
	require.NoError(this.T(), err)
	client, err := clientProvider.NewTempDirClient()
	require.NoError(this.T(), err)
	err = tarexec.NewUntarClient(client, nil).Untar(checkoutTarball, ".")
	require.NoError(this.T(), err)
	_, err = os.Stat(client.Join(client.DirPath(), "smartystreets.gemspec"))
	require.NoError(this.T(), err)
	if !ignoreCheckoutFiles {
		file, err := os.Open(client.Join(client.DirPath(), ".git/HEAD"))
		require.NoError(this.T(), err)
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		require.NoError(this.T(), err)
		var buffer bytes.Buffer
		buffer.Write(data)
		require.Equal(this.T(), testSmartystreetsCommitId, strings.TrimSpace(buffer.String()))
	} else {
		_, err := os.Open(client.Join(client.DirPath(), ".git"))
		require.Error(this.T(), err)
		require.True(this.T(), os.IsNotExist(err))
		_, err = os.Open(client.Join(client.DirPath(), ".gitignore"))
		require.Error(this.T(), err)
		require.True(this.T(), os.IsNotExist(err))
	}
	require.NoError(this.T(), clientProvider.Destroy())
}
