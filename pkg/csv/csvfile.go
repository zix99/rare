package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSV interface {
	io.Closer
	Write(record []string) error
	WriteRow(data ...any) error // Convenience method
}

// Wrapper that represents a Closer on a CSV file
type CSVFile struct {
	*csv.Writer
	f io.WriteCloser
}

func (s *CSVFile) Close() error {
	s.Flush()
	return s.f.Close()
}

func OpenCSV(filename string) (*CSVFile, error) {
	f, err := openFile(filename)
	if err != nil {
		return nil, err
	}

	return NewCSV(f), nil
}

func NewCSV(w io.WriteCloser) *CSVFile {
	return &CSVFile{
		csv.NewWriter(w),
		w,
	}
}

func (s *CSVFile) WriteRow(data ...any) error {
	record := make([]string, len(data))
	for i, v := range data {
		record[i] = fmt.Sprintf("%v", v)
	}
	return s.Write(record)
}

func openFile(filename string) (io.WriteCloser, error) {
	if filename == "-" {
		return &nopWriteCloser{os.Stdout}, nil
	}
	return os.Create(filename)
}

type nopWriteCloser struct {
	io.Writer
}

func (s *nopWriteCloser) Close() error {
	return nil
}
