package scm

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/codeship/go-exec"
)

const (
	testSmartystreetsCommitID = "a40e854c17df0b1a98c90c250dc20e6cb2474dfa"
	testHgGitChangesetID      = "4538981d2c3f3fcb594ad7f2ae7622380929e226"
)

func TestGit(t *testing.T) {
	t.Parallel()
	tempDir := getTempDir(t)
	if err := Checkout(
		&GitCheckoutOptions{
			User:     "git",
			Host:     "github.com",
			Path:     "/peter-edge/smartystreets_ruby.git",
			Branch:   "master",
			CommitID: testSmartystreetsCommitID,
		},
		tempDir,
	); err != nil {
		t.Fatal(err)
	}
	testSmartystreetsCheckoutTarball(t, tempDir)
}

func TestGithub(t *testing.T) {
	t.Parallel()
	tempDir := getTempDir(t)
	if err := Checkout(
		&GithubCheckoutOptions{
			User:       "peter-edge",
			Repository: "smartystreets_ruby",
			Branch:     "master",
			CommitID:   testSmartystreetsCommitID,
			//SecurityOptions: NewGithubSecurityOptionsSSH(getSSHOptions()),
		},
		tempDir,
	); err != nil {
		t.Fatal(err)
	}
	testSmartystreetsCheckoutTarball(t, tempDir)
}

func testSmartystreetsCheckoutTarball(t *testing.T, tempDir string) {
	if _, err := os.Stat(filepath.Join(tempDir, "smartystreets.gemspec")); err != nil {
		t.Error(err)
	}
	file, err := os.Open(filepath.Join(tempDir, ".git/HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Error(err)
		}
	}()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	commitID := strings.TrimSpace(string(data))
	if testSmartystreetsCommitID != commitID {
		t.Errorf("expected %s, got %s", testSmartystreetsCommitID, commitID)
	}
}

func TestHg(t *testing.T) {
	t.Parallel()
	tempDir := getTempDir(t)
	if err := Checkout(
		&HgCheckoutOptions{
			User:        "hg",
			Host:        "bitbucket.org",
			Path:        "/durin42/hg-git",
			ChangesetID: testHgGitChangesetID,
			//SecurityOptions: NewHgSecurityOptionsSSH(getSSHOptions()),
		},
		tempDir,
	); err != nil {
		t.Fatal(err)
	}
	testHgGitCheckoutTarball(t, tempDir)
}

func TestBitbucketHg(t *testing.T) {
	t.Parallel()
	tempDir := getTempDir(t)
	if err := Checkout(
		&BitbucketHgCheckoutOptions{
			User:        "durin42",
			Repository:  "hg-git",
			ChangesetID: testHgGitChangesetID,
			//SecurityOptions: NewHgSecurityOptionsSSH(getSSHOptions()),
		},
		tempDir,
	); err != nil {
		t.Fatal(err)
	}
	testHgGitCheckoutTarball(t, tempDir)
}

func testHgGitCheckoutTarball(t *testing.T, tempDir string) {
	if _, err := os.Stat(filepath.Join(tempDir, "hggit/overlay.py")); err != nil {
		t.Error(err)
	}
	executor, err := exec.NewOsExecutor(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	var buffer bytes.Buffer
	if err := executor.Execute(
		&exec.Cmd{
			Args:   []string{"hg", "parent"},
			Stdout: &buffer,
		},
	)(); err != nil {
		t.Fatal(err)
	}
	output := buffer.String()
	if !strings.Contains(output, testHgGitChangesetID[0:12]) {
		t.Errorf("expected %v, got %v", testHgGitChangesetID, output)
	}
}

func getSSHOptions(t *testing.T) *SSHSecurityOptions {
	privateKeyReader, err := os.Open(os.Getenv("HOME") + "/.ssh/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := privateKeyReader.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	data, err := ioutil.ReadAll(privateKeyReader)
	if err != nil {
		t.Fatal(err)
	}
	var buffer bytes.Buffer
	if _, err := buffer.Write(data); err != nil {
		t.Fatal(err)
	}
	return &SSHSecurityOptions{
		StrictHostKeyChecking: false,
		PrivateKey:            &buffer,
	}
}

func getTempDir(t *testing.T) string {
	tempDir, err := ioutil.TempDir("", "go-scm-test")
	if err != nil {
		t.Error(err)
	}
	return tempDir
}
