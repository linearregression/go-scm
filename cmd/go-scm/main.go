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

	"github.com/coreos/go-etcd/etcd"
	"github.com/peter-edge/go-etcdmarshal"
	"github.com/peter-edge/go-osutils"
	"github.com/peter-edge/go-scm"
)

func main() {
	var baseDirPath string
	var hostBaseDirPath string
	var clonePath string
	var tarballName string
	var ignoreCheckoutFiles bool
	var etcdUrl string
	var etcdInputKey string
	var etcdOutputKey string
	flag.StringVar(&baseDirPath, "base_dir_path", "", "The directory to clone into (defaults to a temporary directory)")
	flag.StringVar(&hostBaseDirPath, "host_base_dir_path", "", "The equivalent directory within the host if base_dir_path is a linked volume (base_dir_path must be set)")
	flag.StringVar(&clonePath, "clone_path", "", "The name of the clone directory (defaults to clone)")
	flag.StringVar(&tarballName, "tarball_name", "", "The name of the tarball to output (no tarball by default)")
	flag.BoolVar(&ignoreCheckoutFiles, "ignore_checkout_files", false, "Ignore checkout files if tarballing (false by default)")
	flag.StringVar(&etcdUrl, "etcd_url", "", "The etcd url")
	flag.StringVar(&etcdInputKey, "etcd_input_key", "", "The etcd input key")
	flag.StringVar(&etcdOutputKey, "etcd_output_key", "", "The etcd output key")
	flag.Parse()
	checkTrue(!(clonePath != "" && tarballName != ""), "Cannot have both --clone_path and --tarball_name set")
	if hostBaseDirPath != "" {
		checkTrue(baseDirPath != "", "--base_dir_path must be set if --host_base_dir_path is set")
	}
	checkTrue(!(tarballName == "" && ignoreCheckoutFiles), "Cannot set --ignoreCheckoutFiles if --tarball_name is not set")
	checkTrue(((etcdUrl == "") == (etcdInputKey == "")) && ((etcdInputKey == "") == (etcdOutputKey == "")), "All of --etcd_url, --etcd_input_key, --etcd_output_key must be set or not set")

	var externalCheckoutOptions scm.ExternalCheckoutOptions
	if etcdUrl == "" {
		data, err := ioutil.ReadAll(os.Stdin)
		checkError(err)
		checkError(json.Unmarshal(data, &externalCheckoutOptions))
	} else {
		etcdmarshalApi := etcdmarshal.NewJsonApi(
			etcd.NewClient(
				[]string{
					etcdUrl,
				},
			),
		)
		checkError(etcdmarshalApi.Read(etcdInputKey, &externalCheckoutOptions))
	}

	dirPath := baseDirPath
	var err error
	if dirPath == "" {
		dirPath, err = osutils.NewTempDir()
		checkError(err)
	} else {
		dirPath, err = osutils.NewSubDir(dirPath)
		checkError(err)
	}
	var path string
	if tarballName != "" {
		tarballReader, err := scm.ExternalCheckoutTarball(&externalCheckoutOptions, &scm.Options{IgnoreCheckoutFiles: ignoreCheckoutFiles})
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
		path = filepath.Join(dirPath, clonePath)
		checkError(scm.ExternalCheckout(path, &externalCheckoutOptions, &scm.Options{IgnoreCheckoutFiles: ignoreCheckoutFiles}))
	}

	if hostBaseDirPath != "" {
		path = strings.NewReplacer(baseDirPath, hostBaseDirPath).Replace(path)
	}
	if etcdUrl == "" {
		fmt.Println(path)
	} else {
		etcdmarshalApi := etcdmarshal.NewStringApi(
			etcd.NewClient(
				[]string{
					etcdUrl,
				},
			),
		)
		checkError(etcdmarshalApi.Write(etcdOutputKey, path))
	}
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
