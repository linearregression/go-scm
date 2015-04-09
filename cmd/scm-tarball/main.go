package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/peter-edge/go-exec"
	"github.com/peter-edge/go-record"
	"github.com/peter-edge/go-scm"
)

func main() {
	var input string
	var baseDirPath string
	var hostBaseDirPath string
	var tarballName string
	var ignoreCheckoutFiles bool
	flag.StringVar(&input, "input", "", "The JSON ExternalCheckoutOptions")
	flag.StringVar(&baseDirPath, "base_dir_path", "", "The directory to clone into (defaults to a temporary directory)")
	flag.StringVar(&hostBaseDirPath, "host_base_dir_path", "", "The equivalent directory within the host if base_dir_path is a linked volume (base_dir_path must be set)")
	flag.StringVar(&tarballName, "tarball_name", "", "The name of the tarball to output (no tarball by default)")
	flag.BoolVar(&ignoreCheckoutFiles, "ignore_checkout_files", false, "Ignore checkout files if tarballing (false by default)")
	flag.Parse()
	if input == "" {
		checkError(errors.New("must pass JSON ExternalCheckoutOptions as --input"))
	}
	checkTrue(tarballName != "", "--tarball_name must be set")
	if hostBaseDirPath != "" {
		checkTrue(baseDirPath != "", "--base_dir_path must be set if --host_base_dir_path is set")
	}

	var externalCheckoutOptions scm.ExternalCheckoutOptions
	checkError(json.Unmarshal([]byte(input), &externalCheckoutOptions))
	checkoutOptions, err := scm.ConvertExternalCheckoutOptions(&externalCheckoutOptions)
	checkError(err)

	execClientProvider, err := exec.NewClientProvider(
		&exec.OsExecOptions{
			TmpDir: baseDirPath,
		},
	)
	checkError(err)

	tarballReader, err := scm.CheckoutTarball(execClientProvider, checkoutOptions, ignoreCheckoutFiles)
	dirPath := baseDirPath
	if dirPath == "" {
		dirPath = os.TempDir()
	}
	path := filepath.Join(dirPath, tarballName)
	file, err := os.Create(path)
	checkError(err)
	_, err = io.Copy(file, tarballReader)
	checkError(err)
	checkError(file.Close())
	if hostBaseDirPath != "" {
		path = strings.NewReplacer(baseDirPath, hostBaseDirPath).Replace(path)
	}
	recordConverterHandler, err := scm.NewRecordConverterHandler()
	checkError(err)
	record.NewWriterMapEncoderRecorder(
		os.Stdout,
		record.NewProtoMapEncoder(),
		recordConverterHandler,
		record.SystemTimer,
		record.RecordLevelInfo,
		nil,
	).RecordUserInfo(
		&scm.TarballRecord{
			Path: path,
		},
	)
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
