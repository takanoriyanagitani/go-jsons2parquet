package main

import (
	"log"
	"time"

	j2p "github.com/takanoriyanagitani/go-jsons2parquet"
)

type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"`
	Status    int32     `json:"status"`
}

var ReaderToLogs j2p.ReaderToStructs[Log] = j2p.ReaderToStructsDefault

func main() {
	e := ReaderToLogs.StdinToJsonsToParquetToStdoutDefault()
	if nil != e {
		log.Fatal(e)
	}
}
