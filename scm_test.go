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
	testSmartystreetsCommitId = "26edc5a1950af7102ab178cc87c7836ee748d003"
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
	checkoutTarball, err := NewGitCheckoutClient(this.clientProvider).CheckoutTarball(
		&GitCheckoutOptions{
			Url:      "https://github.com/peter-edge/smartystreets.git",
			Branch:   "master",
			CommitId: testSmartystreetsCommitId,
		},
	)
	require.NoError(this.T(), err)
	this.testSmartystreetsCheckoutTarball(checkoutTarball)
}

func (this *Suite) TestGithub() {
	checkoutTarball, err := NewGithubCheckoutClient(this.clientProvider).CheckoutTarball(
		&GithubCheckoutOptions{
			User:       "peter-edge",
			Repository: "smartystreets_ruby",
			Branch:     "master",
			CommitId:   testSmartystreetsCommitId,
		},
	)
	require.NoError(this.T(), err)
	this.testSmartystreetsCheckoutTarball(checkoutTarball)
}

func (this *Suite) testSmartystreetsCheckoutTarball(checkoutTarball CheckoutTarball) {
	clientProvider, err := execos.NewClientProvider()
	require.NoError(this.T(), err)
	client, err := clientProvider.NewTempDirClient()
	require.NoError(this.T(), err)
	err = tarexec.NewUntarClient(client, nil).Untar(checkoutTarball, ".")
	require.NoError(this.T(), err)
	_, err = os.Stat(client.Join(client.DirPath(), "smartystreets.gemspec"))
	require.NoError(this.T(), err)
	file, err := os.Open(client.Join(client.DirPath(), ".git/refs/heads/master"))
	require.NoError(this.T(), err)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	require.NoError(this.T(), err)
	var buffer bytes.Buffer
	buffer.Write(data)
	require.Equal(this.T(), testSmartystreetsCommitId, strings.TrimSpace(buffer.String()))
	require.NoError(this.T(), clientProvider.Destroy())
}
