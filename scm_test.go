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
)

const (
	testSmartystreetsCommitId = "a40e854c17df0b1a98c90c250dc20e6cb2474dfa"
	testHgGitChangesetId      = "4538981d2c3f3fcb594ad7f2ae7622380929e226"
)

func TestGit(t *testing.T) {
	t.Parallel()
	testGit(t, false)
}

func TestGitIgnore(t *testing.T) {
	t.Parallel()
	testGit(t, true)
}

func testGit(t *testing.T, ignoreCheckoutFiles bool) {
	clientProvider := getClientProvider(t)
	checkoutTarball, err := CheckoutTarball(
		clientProvider,
		&GitCheckoutOptions{
			User:     "git",
			Host:     "github.com",
			Path:     "/peter-edge/smartystreets_ruby.git",
			Branch:   "master",
			CommitId: testSmartystreetsCommitId,
		},
		ignoreCheckoutFiles,
	)
	require.NoError(t, err)
	testSmartystreetsCheckoutTarball(t, clientProvider, checkoutTarball, ignoreCheckoutFiles)
}

func TestGithub(t *testing.T) {
	t.Parallel()
	testGithub(t, false)
}

func TestGithubIgnore(t *testing.T) {
	t.Parallel()
	testGithub(t, true)
}

func testGithub(t *testing.T, ignoreCheckoutFiles bool) {
	clientProvider := getClientProvider(t)
	checkoutTarball, err := CheckoutTarball(
		clientProvider,
		&GithubCheckoutOptions{
			User:       "peter-edge",
			Repository: "smartystreets_ruby",
			Branch:     "master",
			CommitId:   testSmartystreetsCommitId,
			//SecurityOptions: NewGithubSecurityOptionsSsh(getSshOptions()),
		},
		ignoreCheckoutFiles,
	)
	require.NoError(t, err)
	testSmartystreetsCheckoutTarball(t, clientProvider, checkoutTarball, ignoreCheckoutFiles)
}

func testSmartystreetsCheckoutTarball(t *testing.T, clientProvider exec.ClientProvider, checkoutTarball io.Reader, ignoreCheckoutFiles bool) {
	client, err := clientProvider.NewTempDirClient()
	require.NoError(t, err)
	err = tarexec.NewUntarClient(client, nil).Untar(checkoutTarball, ".")
	require.NoError(t, err)
	_, err = os.Stat(client.Join(client.DirPath(), "smartystreets.gemspec"))
	require.NoError(t, err)
	if !ignoreCheckoutFiles {
		file, err := os.Open(client.Join(client.DirPath(), ".git/HEAD"))
		require.NoError(t, err)
		defer func() {
			if err := file.Close(); err != nil {
				t.Fatal(err)
			}
		}()
		data, err := ioutil.ReadAll(file)
		require.NoError(t, err)
		var buffer bytes.Buffer
		if _, err = buffer.Write(data); err != nil {
			t.Fatal(err)
		}
		require.Equal(t, testSmartystreetsCommitId, strings.TrimSpace(buffer.String()))
	} else {
		_, err := os.Open(client.Join(client.DirPath(), ".git"))
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
		_, err = os.Open(client.Join(client.DirPath(), ".gitignore"))
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
	}
	require.NoError(t, clientProvider.Destroy())
}

func TestHg(t *testing.T) {
	t.Parallel()
	testHg(t, false)
}

func TestHgIgnore(t *testing.T) {
	t.Parallel()
	testHg(t, true)
}

func testHg(t *testing.T, ignoreCheckoutFiles bool) {
	clientProvider := getClientProvider(t)
	checkoutTarball, err := CheckoutTarball(
		clientProvider,
		&HgCheckoutOptions{
			User:        "hg",
			Host:        "bitbucket.org",
			Path:        "/durin42/hg-git",
			ChangesetId: testHgGitChangesetId,
			//SecurityOptions: NewHgSecurityOptionsSsh(getSshOptions()),
		},
		ignoreCheckoutFiles,
	)
	require.NoError(t, err)
	testHgGitCheckoutTarball(t, clientProvider, checkoutTarball, ignoreCheckoutFiles)
}

func TestBitbucketHg(t *testing.T) {
	t.Parallel()
	testBitbucketHg(t, false)
}

func TestBitbucketHgIgnore(t *testing.T) {
	t.Parallel()
	testBitbucketHg(t, true)
}

func testBitbucketHg(t *testing.T, ignoreCheckoutFiles bool) {
	clientProvider := getClientProvider(t)
	checkoutTarball, err := CheckoutTarball(
		clientProvider,
		&BitbucketHgCheckoutOptions{
			User:        "durin42",
			Repository:  "hg-git",
			ChangesetId: testHgGitChangesetId,
			//SecurityOptions: NewHgSecurityOptionsSsh(getSshOptions()),
		},
		ignoreCheckoutFiles,
	)
	require.NoError(t, err)
	testHgGitCheckoutTarball(t, clientProvider, checkoutTarball, ignoreCheckoutFiles)
}

func testHgGitCheckoutTarball(t *testing.T, clientProvider exec.ClientProvider, checkoutTarball io.Reader, ignoreCheckoutFiles bool) {
	client, err := clientProvider.NewTempDirClient()
	require.NoError(t, err)
	err = tarexec.NewUntarClient(client, nil).Untar(checkoutTarball, ".")
	require.NoError(t, err)
	_, err = os.Stat(client.Join(client.DirPath(), "hggit/overlay.py"))
	require.NoError(t, err)
	if !ignoreCheckoutFiles {
		var buffer bytes.Buffer
		err = client.Execute(
			&exec.Cmd{
				Args:   []string{"hg", "parent"},
				Stdout: &buffer,
			},
		)()
		require.NoError(t, err)
		require.True(t, strings.Contains(buffer.String(), testHgGitChangesetId[0:12]))
	} else {
		_, err := os.Open(client.Join(client.DirPath(), ".hg"))
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
		_, err = os.Open(client.Join(client.DirPath(), ".hgignore"))
		require.Error(t, err)
		require.True(t, os.IsNotExist(err))
	}
	require.NoError(t, clientProvider.Destroy())
}

func getClientProvider(t *testing.T) exec.ClientProvider {
	clientProvider, err := exec.NewClientProvider(&exec.OsExecOptions{})
	require.NoError(t, err)
	return clientProvider
}

func getSshOptions(t *testing.T) *SshSecurityOptions {
	privateKeyReader, err := os.Open(os.Getenv("HOME") + "/.ssh/id_rsa")
	require.NoError(t, err)
	defer func() {
		if err := privateKeyReader.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	data, err := ioutil.ReadAll(privateKeyReader)
	require.NoError(t, err)
	var buffer bytes.Buffer
	_, err = buffer.Write(data)
	require.NoError(t, err)
	return &SshSecurityOptions{
		StrictHostKeyChecking: false,
		PrivateKey:            &buffer,
	}
}
