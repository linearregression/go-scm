package scm

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/peter-edge/go-osutils"
	taros "github.com/peter-edge/go-tar/os"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testSmartystreetsCommitId = "a40e854c17df0b1a98c90c250dc20e6cb2474dfa"
	testHgGitChangesetId      = "4538981d2c3f3fcb594ad7f2ae7622380929e226"
)

type Suite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (this *Suite) SetupSuite() {
}

func (this *Suite) SetupTest() {
}

func (this *Suite) TearDownTest() {
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
	checkoutTarball, err := CheckoutTarball(
		&GitCheckoutOptions{
			User:     "git",
			Host:     "github.com",
			Path:     "/peter-edge/smartystreets_ruby.git",
			Branch:   "master",
			CommitId: testSmartystreetsCommitId,
		},
		&Options{
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
	checkoutTarball, err := CheckoutTarball(
		&GithubCheckoutOptions{
			User:       "peter-edge",
			Repository: "smartystreets_ruby",
			Branch:     "master",
			CommitId:   testSmartystreetsCommitId,
			//SecurityOptions: NewGithubSecurityOptionsSsh(this.getSshOptions()),
		},
		&Options{
			IgnoreCheckoutFiles: ignoreCheckoutFiles,
		},
	)
	require.NoError(this.T(), err)
	this.testSmartystreetsCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) testSmartystreetsCheckoutTarball(checkoutTarball io.Reader, ignoreCheckoutFiles bool) {
	tempDir, err := osutils.NewTempDir()
	require.NoError(this.T(), err)
	err = taros.NewUntarClient(nil).Untar(checkoutTarball, tempDir)
	require.NoError(this.T(), err)
	_, err = os.Stat(filepath.Join(tempDir, "smartystreets.gemspec"))
	require.NoError(this.T(), err)
	if !ignoreCheckoutFiles {
		file, err := os.Open(filepath.Join(tempDir, ".git/HEAD"))
		require.NoError(this.T(), err)
		defer file.Close()
		data, err := ioutil.ReadAll(file)
		require.NoError(this.T(), err)
		var buffer bytes.Buffer
		buffer.Write(data)
		require.Equal(this.T(), testSmartystreetsCommitId, strings.TrimSpace(buffer.String()))
	} else {
		_, err := os.Open(filepath.Join(tempDir, ".git"))
		require.Error(this.T(), err)
		require.True(this.T(), os.IsNotExist(err))
		_, err = os.Open(filepath.Join(tempDir, ".gitignore"))
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
	checkoutTarball, err := CheckoutTarball(
		&HgCheckoutOptions{
			User:        "hg",
			Host:        "bitbucket.org",
			Path:        "/durin42/hg-git",
			ChangesetId: testHgGitChangesetId,
			//SecurityOptions: NewHgSecurityOptionsSsh(this.getSshOptions()),
		},
		&Options{
			IgnoreCheckoutFiles: ignoreCheckoutFiles,
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
	checkoutTarball, err := CheckoutTarball(
		&BitbucketCheckoutOptions{
			BitbucketType: BitbucketTypeHg,
			User:          "durin42",
			Repository:    "hg-git",
			ChangesetId:   testHgGitChangesetId,
			//SecurityOptions: NewHgSecurityOptionsSsh(this.getSshOptions()),
		},
		&Options{
			IgnoreCheckoutFiles: ignoreCheckoutFiles,
		},
	)
	require.NoError(this.T(), err)
	this.testHgGitCheckoutTarball(checkoutTarball, ignoreCheckoutFiles)
}

func (this *Suite) testHgGitCheckoutTarball(checkoutTarball io.Reader, ignoreCheckoutFiles bool) {
	tempDir, err := osutils.NewTempDir()
	require.NoError(this.T(), err)
	err = taros.NewUntarClient(nil).Untar(checkoutTarball, tempDir)
	require.NoError(this.T(), err)
	_, err = os.Stat(filepath.Join(tempDir, "hggit/overlay.py"))
	require.NoError(this.T(), err)
	if !ignoreCheckoutFiles {
		var buffer bytes.Buffer
		wait, err := osutils.Execute(
			&osutils.Cmd{
				Args:        []string{"hg", "parent"},
				AbsoluteDir: tempDir,
				Stdout:      &buffer,
			},
		)
		require.NoError(this.T(), err)
		require.NoError(this.T(), wait())
		require.True(this.T(), strings.Contains(buffer.String(), testHgGitChangesetId[0:12]))
	} else {
		_, err := os.Open(filepath.Join(tempDir, ".hg"))
		require.Error(this.T(), err)
		require.True(this.T(), os.IsNotExist(err))
		_, err = os.Open(filepath.Join(tempDir, ".hgignore"))
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
