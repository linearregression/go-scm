package scm

import "github.com/peter-edge/go-record"

//go:generate gen-record

// @gen-record
type CloneRecord struct {
	Path string
}

// @gen-record
type TarballRecord struct {
	Path string
}

func NewRecordConverterHandler() (record.RecordConverterHandler, error) {
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
	if err != nil {
		return nil, err
	}
	for _, recordConverter := range AllRecordConverters {
		if err := recordConverterRegistry.Register(recordConverter); err != nil {
			return nil, err
		}
	}
	return recordConverterRegistry.Handler(), nil
}
