package scm

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/peter-edge/go-exec"
	tarexec "github.com/peter-edge/go-tar/exec"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testSmartystreetsCommitId = "a40e854c17df0b1a98c90c250dc20e6cb2474dfa"
	testHgGitChangesetId      = "4538981d2c3f3fcb594ad7f2ae7622380929e226"
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
	clientProvider, err := exec.NewClientProvider(&exec.OsExecOptions{})
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
	client := NewClient(this.clientProvider, &ClientOptions{ignoreCheckoutFiles})
	checkoutTarball, err := client.CheckoutTarball(
		&GitCheckoutOptions{
			User:     "git",
			Host:     "github.com",
			Path:     ":peter-edge/smartystreets_ruby.git",
			Branch:   "master",
			CommitId: testSmartystreetsCommitId,
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
	client := NewClient(this.clientProvider, &ClientOptions{ignoreCheckoutFiles})
	checkoutTarball, err := client.CheckoutTarball(
		&GithubCheckoutOptions{
			User:       "peter-edge",
			Repository: "smartystreets_ruby",
			Branch:     "master",
			CommitId:   testSmartystreetsCommitId,
			//SecurityOptions: NewGithubSecurityOptionsSsh(this.getSshOptions()),
		},
	)
	require.NoError(this.T(), err)
	this.testSmartystreetsCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) testSmartystreetsCheckoutTarball(checkoutTarball io.Reader, ignoreCheckoutFiles bool) {
	client, err := this.clientProvider.NewTempDirClient()
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
}

func (this *Suite) TestHg() {
	this.testHg(false)
}

func (this *Suite) TestHgIgnore() {
	this.testHg(true)
}

func (this *Suite) testHg(ignoreCheckoutFiles bool) {
	client := NewClient(this.clientProvider, &ClientOptions{ignoreCheckoutFiles})
	checkoutTarball, err := client.CheckoutTarball(
		&HgCheckoutOptions{
			User:        "hg",
			Host:        "bitbucket.org",
			Path:        "/durin42/hg-git",
			ChangesetId: testHgGitChangesetId,
			//SecurityOptions: NewHgSecurityOptionsSsh(this.getSshOptions()),
		},
	)
	require.NoError(this.T(), err)
	this.testHgGitCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) TestBitbucketHg() {
	this.testBitbucketHg(false)
}

func (this *Suite) TestBitbucketHgIgnore() {
	this.testBitbucketHg(true)
}

func (this *Suite) testBitbucketHg(ignoreCheckoutFiles bool) {
	client := NewClient(this.clientProvider, &ClientOptions{ignoreCheckoutFiles})
	checkoutTarball, err := client.CheckoutTarball(
		&BitbucketCheckoutOptions{
			BitbucketType: BitbucketTypeHg,
			User:          "durin42",
			Repository:    "hg-git",
			ChangesetId:   testHgGitChangesetId,
			//SecurityOptions: NewHgSecurityOptionsSsh(this.getSshOptions()),
		},
	)
	require.NoError(this.T(), err)
	this.testHgGitCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) testHgGitCheckoutTarball(checkoutTarball io.Reader, ignoreCheckoutFiles bool) {
	client, err := this.clientProvider.NewTempDirClient()
	require.NoError(this.T(), err)
	err = tarexec.NewUntarClient(client, nil).Untar(checkoutTarball, ".")
	require.NoError(this.T(), err)
	_, err = os.Stat(client.Join(client.DirPath(), "hggit/overlay.py"))
	require.NoError(this.T(), err)
	if !ignoreCheckoutFiles {
		var buffer bytes.Buffer
		err = client.Execute(
			&exec.Cmd{
				Args:   []string{"hg", "parent"},
				Stdout: &buffer,
			},
		)()
		require.NoError(this.T(), err)
		require.True(this.T(), strings.Contains(buffer.String(), testHgGitChangesetId[0:12]))
	} else {
		_, err := os.Open(client.Join(client.DirPath(), ".hg"))
		require.Error(this.T(), err)
		require.True(this.T(), os.IsNotExist(err))
		_, err = os.Open(client.Join(client.DirPath(), ".hgignore"))
		require.Error(this.T(), err)
		require.True(this.T(), os.IsNotExist(err))
	}
}

func (this *Suite) getSshOptions() *SshSecurityOptions {
	privateKeyReader, err := os.Open(os.Getenv("HOME") + "/.ssh/id_rsa")
	require.NoError(this.T(), err)
	defer privateKeyReader.Close()
	data, err := ioutil.ReadAll(privateKeyReader)
	require.NoError(this.T(), err)
	var buffer bytes.Buffer
	_, err = buffer.Write(data)
	require.NoError(this.T(), err)
	return &SshSecurityOptions{
		StrictHostKeyChecking: false,
		PrivateKey:            &buffer,
	}
}
