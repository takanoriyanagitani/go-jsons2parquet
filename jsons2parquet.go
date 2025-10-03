package jsons2parquet

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"iter"
	"os"

	pq "github.com/parquet-go/parquet-go"
)

func WriteRowsGeneric[T any](rows []T, wtr *pq.GenericWriter[T]) error {
	_, e := wtr.Write(rows)
	return e
}

func WriteRowGeneric[T any](row T, wtr *pq.GenericWriter[T]) error {
	return WriteRowsGeneric([]T{row}, wtr)
}

func WriteIterGeneric[T any](rows iter.Seq[T], wtr *pq.GenericWriter[T]) error {
	for row := range rows {
		e := WriteRowGeneric(row, wtr)
		if nil != e {
			return e
		}
	}
	return nil
}

type GenericRow[T any] struct{ Raw T }

type ReaderToStructs[T any] func(io.Reader) iter.Seq2[GenericRow[T], error]

func ReaderToStructsDefault[T any](rdr io.Reader) iter.Seq2[GenericRow[T], error] {
	return func(yield func(GenericRow[T], error) bool) {
		var row GenericRow[T]
		var br io.Reader = bufio.NewReader(rdr)
		var dec *json.Decoder = json.NewDecoder(br)
		for {
			e := dec.Decode(&row.Raw)
			if io.EOF == e {
				return
			}

			if !yield(row, e) {
				return
			}
		}
	}
}

func (j ReaderToStructs[T]) WriteAll(
	jsons io.Reader,
	wtr *pq.GenericWriter[T],
) error {
	var rows iter.Seq2[GenericRow[T], error] = j(jsons)
	for row, e := range rows {
		if nil != e {
			return e
		}
		e := WriteRowGeneric(row.Raw, wtr)
		if nil != e {
			return e
		}
	}
	return wtr.Flush()
}

func (j ReaderToStructs[T]) ToWriter(
	jsons io.Reader,
	wtr io.Writer,
	opts ...pq.WriterOption,
) error {
	var gwtr *pq.GenericWriter[T] = pq.NewGenericWriter[T](
		wtr,
		opts...,
	)
	return errors.Join(j.WriteAll(jsons, gwtr), gwtr.Close())
}

func (j ReaderToStructs[T]) ToStdout(
	jsons io.Reader,
	opts ...pq.WriterOption,
) error {
	return j.ToWriter(jsons, os.Stdout, opts...)
}

func (j ReaderToStructs[T]) StdinToJsonsToParquetToStdout(
	opts ...pq.WriterOption,
) error {
	return j.ToStdout(os.Stdin, opts...)
}

func (j ReaderToStructs[T]) StdinToJsonsToParquetToStdoutDefault() error {
	return j.StdinToJsonsToParquetToStdout()
}
