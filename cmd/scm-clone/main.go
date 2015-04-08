package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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
	var clonePath string
	flag.StringVar(&input, "input", "", "The JSON ExternalCheckoutOptions")
	flag.StringVar(&baseDirPath, "base_dir_path", "", "The directory to clone into (defaults to a temporary directory)")
	flag.StringVar(&hostBaseDirPath, "host_base_dir_path", "", "The equivalent directory within the host if base_dir_path is a linked volume (base_dir_path must be set)")
	flag.StringVar(&clonePath, "clone_path", "", "The name of the clone directory (defaults to clone)")
	flag.Parse()
	if input == "" {
		checkError(errors.New("must pass JSON ExternalCheckoutOptions as --input"))
	}
	if hostBaseDirPath != "" {
		checkTrue(baseDirPath != "", "--base_dir_path must be set if --host_base_dir_path is set")
	}
	if clonePath == "" {
		clonePath = "clone"
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

	executor, err := execClientProvider.NewTempDirExecutorReadFileManager()
	checkError(err)
	checkError(scm.Checkout(execClientProvider, checkoutOptions, executor, clonePath))
	path := filepath.Join(executor.DirPath(), clonePath)
	if hostBaseDirPath != "" {
		path = strings.NewReplacer(baseDirPath, hostBaseDirPath).Replace(path)
	}
	recordConverterHandler, err := scm.NewRecordConverterHandler()
	checkError(err)
	record.NewWriterMapMarshallerRecorder(
		os.Stdout,
		record.NewJSONMapMarshaller(),
		recordConverterHandler,
		record.SystemTimer,
		record.RecordLevelInfo,
		nil,
	).RecordUserInfo(
		&scm.CloneRecord{
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
