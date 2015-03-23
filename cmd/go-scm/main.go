package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/peter-edge/go-exec"
	"github.com/peter-edge/go-scm"
)

func main() {
	var baseDirPath string
	var hostBaseDirPath string
	var clonePath string
	var tarballName string
	var ignoreCheckoutFiles bool
	flag.StringVar(&baseDirPath, "base_dir_path", "", "The directory to clone into (defaults to a temporary directory)")
	flag.StringVar(&hostBaseDirPath, "host_base_dir_path", "", "The equivalent directory within the host if base_dir_path is a linked volume (base_dir_path must be set)")
	flag.StringVar(&clonePath, "clone_path", "", "The name of the clone directory (defaults to clone)")
	flag.StringVar(&tarballName, "tarball_name", "", "The name of the tarball to output (no tarball by default)")
	flag.BoolVar(&ignoreCheckoutFiles, "ignore_checkout_files", false, "Ignore checkout files if tarballing (false by default)")
	flag.Parse()
	checkTrue(!(clonePath != "" && tarballName != ""), "Cannot have both --clone_path and --tarball_name set")
	if hostBaseDirPath != "" {
		checkTrue(baseDirPath != "", "--base_dir_path must be set if --host_base_dir_path is set")
	}
	checkTrue(!(tarballName == "" && ignoreCheckoutFiles), "Cannot set --ignoreCheckoutFiles if --tarball_name is not set")

	data, err := ioutil.ReadAll(os.Stdin)
	checkError(err)
	var externalCheckoutOptions scm.ExternalCheckoutOptions
	checkError(json.Unmarshal(data, &externalCheckoutOptions))
	checkoutOptions, err := scm.ConvertExternalCheckoutOptions(&externalCheckoutOptions)
	checkError(err)

	execClientProvider, err := exec.NewClientProvider(
		&exec.OsExecOptions{
			TmpDir: baseDirPath,
		},
	)
	checkError(err)

	var path string
	if tarballName != "" {
		tarballReader, err := scm.CheckoutTarball(execClientProvider, checkoutOptions, ignoreCheckoutFiles)
		dirPath := baseDirPath
		if dirPath == "" {
			dirPath = os.TempDir()
		}
		path = filepath.Join(dirPath, tarballName)
		file, err := os.Create(path)
		checkError(err)
		_, err = io.Copy(file, tarballReader)
		checkError(err)
		checkError(file.Close())
	} else {
		if clonePath == "" {
			clonePath = "clone"
		}
		executor, err := execClientProvider.NewTempDirExecutorReadFileManager()
		checkError(err)
		path = filepath.Join(executor.DirPath(), clonePath)
		checkError(scm.Checkout(execClientProvider, checkoutOptions, executor, clonePath))
	}

	if hostBaseDirPath != "" {
		path = strings.NewReplacer(baseDirPath, hostBaseDirPath).Replace(path)
	}
	fmt.Println(path)
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}

func checkTrue(value bool, message string) {
	if !value {
		fmt.Fprintf(os.Stderr, "%s\n", message)
		os.Exit(1)
	}
}
