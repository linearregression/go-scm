package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/peter-edge/go-exec"
	"github.com/peter-edge/go-record"
	"github.com/peter-edge/go-scm"
)

func main() {
	var baseDirPath string
	var hostBaseDirPath string
	var clonePath string
	flag.StringVar(&baseDirPath, "base_dir_path", "", "The directory to clone into (defaults to a temporary directory)")
	flag.StringVar(&hostBaseDirPath, "host_base_dir_path", "", "The equivalent directory within the host if base_dir_path is a linked volume (base_dir_path must be set)")
	flag.StringVar(&clonePath, "clone_path", "", "The name of the clone directory (defaults to clone)")
	flag.Parse()
	if hostBaseDirPath != "" {
		checkTrue(baseDirPath != "", "--base_dir_path must be set if --host_base_dir_path is set")
	}
	if clonePath == "" {
		clonePath = "clone"
	}

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

	executor, err := execClientProvider.NewTempDirExecutorReadFileManager()
	checkError(err)
	checkError(scm.Checkout(execClientProvider, checkoutOptions, executor, clonePath))
	path := filepath.Join(executor.DirPath(), clonePath)
	if hostBaseDirPath != "" {
		path = strings.NewReplacer(baseDirPath, hostBaseDirPath).Replace(path)
	}
	recordConverterRegistry, err := record.NewRecordConverterRegistry(
		&record.RecordConverterReservedKeys{
			Id:           "ID_",
			Type:         "TYPE",
			TimeUnixNsec: "TIME_UNIX_NSEC",
			Category:     "CATEGORY",
			RecordLevel:  "RECORD_LEVEL",
			WriterOutput: "WRITER_OUTPUT",
		},
	)
	checkError(err)
	for _, recordConverter := range scm.AllRecordConverters {
		checkError(recordConverterRegistry.Register(recordConverter))
	}
	record.NewWriterMapMarshallerRecorder(
		os.Stdout,
		record.NewJSONMapMarshaller(),
		recordConverterRegistry.Handler(),
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
