package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/peter-edge/go-exec"
	"github.com/peter-edge/go-scm"
)

func main() {
	var baseDirPath string
	var clonePath string
	var tarballName string
	var ignoreCheckoutFiles bool
	flag.StringVar(&baseDirPath, "base_dir_path", "", "The directory to clone into (defaults to a temporary directory)")
	flag.StringVar(&clonePath, "clone_path", "", "The name of the clone directory (defaults to clone)")
	flag.StringVar(&tarballName, "tarball_name", "", "The name of the tarball to output (no tarball by default)")
	flag.BoolVar(&ignoreCheckoutFiles, "ignore_checkout_files", false, "Ignore checkout files if tarballing (false by default)")
	flag.Parse()
	checkTrue(!(clonePath != "" && tarballName != ""), "Cannot have both --clone_path and --tarball_name set")
	checkTrue(!(tarballName == "" && ignoreCheckoutFiles), "Cannot set --ignoreCheckoutFiles if --tarball_name is not set")

	data, err := ioutil.ReadAll(os.Stdin)
	checkError(err)
	var externalCheckoutOptions scm.ExternalCheckoutOptions
	checkError(json.Unmarshal(data, &externalCheckoutOptions))

	execClientProvider, err := exec.NewClientProvider(
		&exec.OsExecOptions{
			TmpDir: "",
		},
	)
	checkError(err)

	if tarballName != "" {
		client := scm.NewClient(execClientProvider, &scm.ClientOptions{IgnoreCheckoutFiles: ignoreCheckoutFiles})
		externalClient := scm.NewExternalClient(client)
		tarballReader, err := externalClient.CheckoutTarball(&externalCheckoutOptions)
		dirPath := baseDirPath
		if dirPath == "" {
			dirPath = os.TempDir()
		}
		file, err := os.Create(filepath.Join(dirPath, tarballName))
		checkError(err)
		_, err = io.Copy(file, tarballReader)
		checkError(err)
		checkError(file.Close())
		fmt.Printf("Created tarball as %s/%s\n", dirPath, tarballName)
	} else {
		if clonePath == "" {
			clonePath = "clone"
		}
		directClient := scm.NewDirectClient(execClientProvider)
		externalDirectClient := scm.NewExternalDirectClient(directClient)
		executor, err := execClientProvider.NewTempDirExecutorReadFileManager()
		checkError(err)
		checkError(externalDirectClient.Checkout(&externalCheckoutOptions, executor, clonePath))
		if baseDirPath != "" {
			checkError(os.Rename(filepath.Join(executor.DirPath(), clonePath), filepath.Join(baseDirPath, clonePath)))
			checkError(executor.Destroy())
			fmt.Printf("Checked out to %s/%s\n", baseDirPath, clonePath)
		} else {
			fmt.Printf("Checked out to %s/%s\n", executor.DirPath(), clonePath)
		}
	}
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func checkTrue(value bool, message string) {
	if !value {
		fmt.Println(message)
		os.Exit(1)
	}
}
