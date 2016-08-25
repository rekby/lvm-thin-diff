package lvm_thin_diff

import (
	"flag"
	"strings"
	"os"
	"io"
	"encoding/gob"
)

const (
	BUF_SIZE = 4*1024*1024 // bytes
)

var (
	MetadataDumpFile = flag.String("metadata-dump-file", "", "Path to xml metadata file from thin-dump")
	Output = flag.String("output", "-", "Path to output file. '-' mean stdout")
	Operation = flag.String("operation", "", "makediff. makediff - create diff of snapshots")
	FromDevId = flag.Int("from-dev-id", 0, "DevID of old snapshot")
	ToDevId = flag.Int("to-dev-id", 0, "DevID of new snapshot")
	DataFile = flag.String("data-file", "", "path to device or file with underliing data")
)

func Main(){
	flag.Parse()

	switch strings.ToLower(*Operation) {
	case "makediff":
		makeDiff()
	}
}


func makeDiff(){
	var err error
	var writer io.WriteCloser
	if *Output == "-"{
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(*Output, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
	}
	defer writer.Close()
	enc := gob.NewEncoder(writer)

	reader, err := os.OpenFile(*DataFile, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(*MetadataDumpFile)

	if err != nil {
		panic(err)
	}

	devices, err := parseMetaDataXML(f)
	if err != nil {
		panic(err)
	}

	var from, to dataDevice
	for _, dev := range devices {
		if dev.Id == *FromDevId {
			from = dev
		}
		if dev.Id == *ToDevId {
			to = dev
		}
	}

	cutter := newDataBlockArrCutter(from.Blocks, to.Blocks)
	buf := make([]byte, BUF_SIZE)
	for {
		ok, bFrom, bTo := cutter.Cut()
		if !ok {
			return
		}
		diff := calcDiff(bFrom, bTo)
		enc.Encode(diff)
		if diff.Operation == WRITE {
			_, err = reader.Seek(bTo.DataOffset, 0)
			if err != nil {
				panic(err)
			}
			var writedBytes int64
			for writedBytes = 0;writedBytes < diff.Length; {
				bytesForRead := diff.Length - writedBytes
				localBuf := buf[:minInt64(BUF_SIZE, bytesForRead)]
				_, err := reader.Read(localBuf)
				if err != nil {
					panic(err)
				}
				err = enc.Encode(localBuf)
				if err != nil {
					panic(err)
				}
				writedBytes += int64(len(localBuf))
			}
			if writedBytes != diff.Length {
				panic("writedBytes != diff.Length")
			}
		}
	}
}

/*
Создать команду для патча данных from так чтобы получились данные to.
bFrom и bTo - два блока данных. Если оба блока не пустые - то они должны начинаться с одного логического смещения и
быть равной длины.

Пустой блок означает что в месте, указанном вторым блоком данных нет.
*/
func calcDiff(bFrom, bTo dataBlock) dataPatch {
	if bFrom.IsEmpty() && bTo.IsEmpty() {
		return dataPatch{Operation: NONE}
	}

	if bFrom.IsEmpty() {
		return dataPatch{Offset: bTo.OriginOffset, Operation: WRITE, Length: bTo.Length}
	}

	if bTo.IsEmpty() {
		return dataPatch{Offset: bFrom.OriginOffset, Operation: DELETE, Length: bFrom.Length}
	}

	if bFrom.OriginOffset != bTo.OriginOffset || bFrom.Length != bTo.Length {
		panic("bFrom and bTo must have same start and length")
	}

	if bFrom.DataOffset == bTo.DataOffset {
		return dataPatch{Operation: NONE} // Data is equal. Do nothing.
	}

	return dataPatch{Offset: bTo.OriginOffset, Operation: WRITE, Length: bTo.Length}
}

func minInt64(a,b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}