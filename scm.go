package scm

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/peter-edge/go-osutils"
	taros "github.com/peter-edge/go-tar/os"
)

const (
	clonePath = "clone"
)

var (
	IgnoreGitCheckoutFilePatterns = []string{
		".git",
		".gitignore",
	}
	IgnoreHgCheckoutFilePatterns = []string{
		".hg",
		".hgignore",
		".hgsigs",
		".hgtags",
	}
)

type CheckoutOptions interface {
	Type() CheckoutType
}

type GitCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

func (this *GitCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeGit
}

type GithubCheckoutOptions struct {
	User            string
	Repository      string
	Branch          string
	CommitId        string
	SecurityOptions SecurityOptions
}

func (this *GithubCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeGithub
}

type HgCheckoutOptions struct {
	User            string
	Host            string
	Path            string
	ChangesetId     string
	SecurityOptions SecurityOptions
}

func (this *HgCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeHg
}

type BitbucketCheckoutOptions struct {
	BitbucketType   BitbucketType
	User            string
	Repository      string
	Branch          string // only set if BitbucketType == BitbucketTypeGit
	CommitId        string // only set if BitbucketType == BitbucketTypeGit
	ChangesetId     string // only set if BitbucketType == BitbucketTypeHg
	SecurityOptions SecurityOptions
}

func (this *BitbucketCheckoutOptions) Type() CheckoutType {
	return CheckoutTypeBitbucket
}

type SecurityOptions interface {
	Type() SecurityType
}

type SshSecurityOptions struct {
	StrictHostKeyChecking bool
	PrivateKey            io.Reader
}

func (this *SshSecurityOptions) Type() SecurityType {
	return SecurityTypeSsh
}

type AccessTokenSecurityOptions struct {
	AccessToken string
}

func (this *AccessTokenSecurityOptions) Type() SecurityType {
	return SecurityTypeAccessToken
}

type Options struct {
	IgnoreCheckoutFiles bool
}

func Checkout(absolutePath string, checkoutOptions CheckoutOptions, options *Options) error {
	return checkout(absolutePath, checkoutOptions, options)
}

func CheckoutTarball(checkoutOptions CheckoutOptions, options *Options) (io.Reader, error) {
	return checkoutTarball(checkoutOptions, options)
}

// ***** PRIVATE *****

func checkoutTarball(checkoutOptions CheckoutOptions, options *Options) (retValue io.Reader, retErr error) {
	tempDir, err := osutils.NewTempDir()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil && retErr == nil {
			retErr = err
		}
	}()
	if err := checkout(filepath.Join(tempDir, clonePath), checkoutOptions, options); err != nil {
		return nil, err
	}
	reader, err := tar(filepath.Join(tempDir, clonePath), checkoutOptions, options)
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func checkout(absolutePath string, checkoutOptions CheckoutOptions, options *Options) error {
	if !filepath.IsAbs(absolutePath) {
		return newInternalError(newValidationErrorNotAbsolutePath(absolutePath))
	}
	if err := validateCheckoutOptions(checkoutOptions); err != nil {
		return err
	}
	switch checkoutOptions.Type() {
	case CheckoutTypeGit:
		return checkoutGit(absolutePath, checkoutOptions.(*GitCheckoutOptions))
	case CheckoutTypeGithub:
		return checkoutGithub(absolutePath, checkoutOptions.(*GithubCheckoutOptions))
	case CheckoutTypeHg:
		return checkoutHg(absolutePath, checkoutOptions.(*HgCheckoutOptions))
	case CheckoutTypeBitbucket:
		return checkoutBitbucket(absolutePath, checkoutOptions.(*BitbucketCheckoutOptions))
	default:
		return UnknownCheckoutType(checkoutOptions.Type())
	}
}

func checkoutGit(path string, gitCheckoutOptions *GitCheckoutOptions) error {
	var sshCommand string = ""
	var err error
	if gitCheckoutOptions.SecurityOptions != nil {
		sshCommand, err = getSshCommand(gitCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
	}
	url, err := getGitUrl(gitCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitExecute(path, sshCommand, url, gitCheckoutOptions.Branch, gitCheckoutOptions.CommitId)
}

func checkoutGithub(path string, githubCheckoutOptions *GithubCheckoutOptions) error {
	var sshCommand string = ""
	var err error
	if githubCheckoutOptions.SecurityOptions != nil {
		sshCommand, err = getSshCommand(githubCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
	}
	url, err := getGithubUrl(githubCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutGitExecute(path, sshCommand, url, githubCheckoutOptions.Branch, githubCheckoutOptions.CommitId)
}

func checkoutHg(path string, hgCheckoutOptions *HgCheckoutOptions) error {
	var sshCommand string = ""
	var err error
	if hgCheckoutOptions.SecurityOptions != nil {
		sshCommand, err = getSshCommand(hgCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
	}
	url, err := getHgUrl(hgCheckoutOptions)
	if err != nil {
		return err
	}
	return checkoutHgExecute(path, sshCommand, url, hgCheckoutOptions.ChangesetId)
}

func checkoutBitbucket(path string, bitbucketCheckoutOptions *BitbucketCheckoutOptions) (retErr error) {
	var sshCommand string = ""
	var err error
	if bitbucketCheckoutOptions.SecurityOptions != nil {
		sshCommand, err = getSshCommand(bitbucketCheckoutOptions.SecurityOptions)
		if err != nil {
			return err
		}
	}
	url, err := getBitbucketUrl(bitbucketCheckoutOptions)
	if err != nil {
		return err
	}
	switch bitbucketCheckoutOptions.BitbucketType {
	case BitbucketTypeGit:
		return checkoutGitExecute(path, sshCommand, url, bitbucketCheckoutOptions.Branch, bitbucketCheckoutOptions.CommitId)
	case BitbucketTypeHg:
		return checkoutHgExecute(path, sshCommand, url, bitbucketCheckoutOptions.ChangesetId)
	default:
		return UnknownBitbucketType(bitbucketCheckoutOptions.BitbucketType)
	}
}

func getSshCommand(securityOptions SecurityOptions) (retValue string, retErr error) {
	if securityOptions.Type() != SecurityTypeSsh {
		return "", nil
	}
	sshSecurityOptions := securityOptions.(*SshSecurityOptions)

	sshCommand := []string{"ssh", "-o"}
	if sshSecurityOptions.StrictHostKeyChecking {
		sshCommand = append(sshCommand, "StrictHostKeyChecking=yes")
	} else {
		sshCommand = append(sshCommand, "StrictHostKeyChecking=no")
	}
	if sshSecurityOptions.PrivateKey != nil {
		tempDir, err := osutils.NewTempDir()
		if err != nil {
			return "", err
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil && retErr == nil {
				retErr = err
			}
		}()
		writeFile, err := os.Create(filepath.Join(tempDir, "id_rsa"))
		if err != nil {
			return "", err
		}
		defer func() {
			if err := writeFile.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		data, err := ioutil.ReadAll(sshSecurityOptions.PrivateKey)
		if err != nil {
			return "", err
		}
		_, err = writeFile.Write(data)
		if err != nil {
			return "", err
		}
		err = writeFile.Chmod(0400)
		if err != nil {
			return "", err
		}
		sshCommand = append(sshCommand, "-i", filepath.Join(tempDir, "id_rsa"))
	}
	return strings.Join(sshCommand, " "), nil
}

func getGitUrl(gitCheckoutOptions *GitCheckoutOptions) (string, error) {
	if gitCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyUrl(
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	if gitCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		return getSshUrl(
			"",
			gitCheckoutOptions.User,
			gitCheckoutOptions.Host,
			gitCheckoutOptions.Path,
		), nil
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			gitCheckoutOptions.Type().String(),
			gitCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
}

func getGithubUrl(githubCheckoutOptions *GithubCheckoutOptions) (string, error) {
	if githubCheckoutOptions.SecurityOptions == nil {
		return getGitReadOnlyUrl(
			"git",
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		return getSshUrl(
			"",
			"git",
			"github.com",
			joinStrings(":", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	if githubCheckoutOptions.SecurityOptions.Type() == SecurityTypeAccessToken {
		return getAccessTokenUrl(
			(githubCheckoutOptions.SecurityOptions.(*AccessTokenSecurityOptions)).AccessToken,
			"github.com",
			joinStrings("/", githubCheckoutOptions.User, "/", githubCheckoutOptions.Repository, ".git"),
		), nil
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			githubCheckoutOptions.Type().String(),
			githubCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
}

func getHgUrl(hgCheckoutOptions *HgCheckoutOptions) (string, error) {
	if hgCheckoutOptions.SecurityOptions == nil || hgCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		return getSshUrl(
			"ssh://",
			hgCheckoutOptions.User,
			hgCheckoutOptions.Host,
			hgCheckoutOptions.Path,
		), nil
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			hgCheckoutOptions.Type().String(),
			hgCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
}

func getBitbucketUrl(bitbucketCheckoutOptions *BitbucketCheckoutOptions) (string, error) {
	if bitbucketCheckoutOptions.SecurityOptions == nil || bitbucketCheckoutOptions.SecurityOptions.Type() == SecurityTypeSsh {
		switch bitbucketCheckoutOptions.BitbucketType {
		case BitbucketTypeGit:
			return getSshUrl(
				"ssh://",
				"git",
				"bitbucket.org",
				joinStrings(":", bitbucketCheckoutOptions.User, "/", bitbucketCheckoutOptions.Repository, ".git"),
			), nil
		case BitbucketTypeHg:
			return getSshUrl(
				"ssh://",
				"hg",
				"bitbucket.org",
				joinStrings("/", bitbucketCheckoutOptions.User, "/", bitbucketCheckoutOptions.Repository),
			), nil
		default:
			return "", UnknownBitbucketType(bitbucketCheckoutOptions.BitbucketType)
		}
	}
	return "", newInternalError(
		newValidationErrorSecurityNotImplementedForCheckoutType(
			bitbucketCheckoutOptions.Type().String(),
			bitbucketCheckoutOptions.SecurityOptions.Type().String(),
		),
	)
}

// TODO(pedge): user?
func getGitReadOnlyUrl(user string, host string, path string) string {
	return joinStrings("git://", host, path)
}

func getSshUrl(base string, user string, host string, path string) string {
	return joinStrings(base, user, "@", host, path)
}

func getAccessTokenUrl(accessToken string, host string, path string) string {
	return joinStrings("https://", accessToken, ":x-oauth-basic@", host, path)
}

func checkoutGitExecute(
	path string,
	gitSshCommand string,
	url string,
	branch string,
	commitId string,
) error {
	var cloneStderr bytes.Buffer
	cmd := osutils.Cmd{
		// TODO(peter): if the commit id is more than 50 back, the checkout will fail
		Args:        []string{"git", "clone", "--branch", branch, "--depth", "50", "--recursive", url, path},
		AbsoluteDir: filepath.Dir(path),
		Stderr:      &cloneStderr,
	}
	if gitSshCommand != "" {
		cmd.Env = []string{"GIT_SSH_COMMAND=" + gitSshCommand}
	}
	wait, err := osutils.Execute(&cmd)
	if err != nil {
		return err
	}
	if err := wait(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotClone: %v %v", err.Error(), cloneStderr.String())
	}
	var checkoutStderr bytes.Buffer
	wait, err = osutils.Execute(
		&osutils.Cmd{
			Args:        []string{"git", "checkout", "-f", commitId},
			AbsoluteDir: path,
			Stderr:      &checkoutStderr,
		},
	)
	if err != nil {
		return err
	}
	if err := wait(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotCheckout: %v %v", err.Error(), checkoutStderr.String())
	}
	return nil
}

func checkoutHgExecute(
	path string,
	sshCommand string,
	url string,
	changesetId string,
) error {
	args := []string{"hg", "clone", url, path}
	if sshCommand != "" {
		args = []string{"hg", "clone", "--ssh", sshCommand, url, path}
	}
	var cloneStderr bytes.Buffer
	wait, err := osutils.Execute(
		&osutils.Cmd{
			Args:        args,
			AbsoluteDir: filepath.Dir(path),
			Stderr:      &cloneStderr,
		},
	)
	if err != nil {
		return err
	}
	if err := wait(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotClone: %v %v", err.Error(), cloneStderr.String())
	}
	var updateStderr bytes.Buffer
	wait, err = osutils.Execute(
		&osutils.Cmd{
			Args:        []string{"hg", "update", "--cwd", path, changesetId},
			AbsoluteDir: path,
			Stderr:      &updateStderr,
		},
	)
	if err != nil {
		return err
	}
	if err := wait(); err != nil {
		// TODO(pedge)
		return fmt.Errorf("CouldNotUpdate: %v %v", err.Error(), updateStderr.String())
	}
	return nil
}

func tar(path string, checkoutOptions CheckoutOptions, options *Options) (io.Reader, error) {
	var reader io.Reader
	var err error
	if options.IgnoreCheckoutFiles {
		ignoreCheckoutFilePatterns, err := ignoreCheckoutFilePatterns(checkoutOptions)
		if err != nil {
			return nil, err
		}
		reader, err = tarFiles(path, ignoreCheckoutFilePatterns)
	} else {
		reader, err = tarFiles(path, nil)
	}
	if err != nil {
		return nil, err
	}
	return reader, nil
}

func ignoreCheckoutFilePatterns(checkoutOptions CheckoutOptions) ([]string, error) {
	switch checkoutOptions.Type() {
	case CheckoutTypeGit:
		return IgnoreGitCheckoutFilePatterns, nil
	case CheckoutTypeGithub:
		return IgnoreGitCheckoutFilePatterns, nil
	case CheckoutTypeHg:
		return IgnoreHgCheckoutFilePatterns, nil
	case CheckoutTypeBitbucket:
		bitbucketCheckoutOptions := checkoutOptions.(*BitbucketCheckoutOptions)
		switch bitbucketCheckoutOptions.BitbucketType {
		case BitbucketTypeGit:
			return IgnoreGitCheckoutFilePatterns, nil
		case BitbucketTypeHg:
			return IgnoreHgCheckoutFilePatterns, nil
		default:
			return nil, UnknownBitbucketType(bitbucketCheckoutOptions.BitbucketType)
		}
	default:
		return nil, UnknownCheckoutType(checkoutOptions.Type())
	}
}

func tarFiles(path string, ignoreCheckoutFilePatterns []string) (io.Reader, error) {
	fileList, err := osutils.ListRegularFiles(path)
	if err != nil {
		return nil, err
	}
	if ignoreCheckoutFilePatterns != nil && len(ignoreCheckoutFilePatterns) > 0 {
		filterFileList := make([]string, 0)
		for _, file := range fileList {
			rel, err := filepath.Rel(path, file)
			if err != nil {
				return nil, err
			}
			matches, err := fileMatches(rel, ignoreCheckoutFilePatterns)
			if err != nil {
				return nil, err
			}
			if !matches {
				filterFileList = append(filterFileList, file)
			}
		}
		fileList = filterFileList
	}
	return taros.NewTarClient(nil).Tar(fileList, path)
}

func fileMatches(file string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		if strings.HasPrefix(file, pattern) {
			return true, nil
		}
	}
	return false, nil
}

func joinStrings(elems ...string) string {
	return strings.Join(elems, "")
}
