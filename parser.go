package lvm_thin_diff

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

const sectorSize = 512 // bytes in sector

// Offset/length in bytes
type Block struct {
	OriginOffset int64
	DataOffset   int64
	Length       int64
}

type Device struct {
	Id     int
	Blocks []Block
}

func Parse(reader io.Reader) ([]Device, error) {
	arr := []Device{}
	xmlReader := xml.NewDecoder(reader)
	var dev *Device
	var blockSize int64 = 0
	for {
		token, err := xmlReader.Token()
		if err != nil {
			return arr, errors.New("Parse token error: " + err.Error())
		}
		switch t := token.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "superblock":
				blockSize, err = strconv.ParseInt(getAttr(t.Attr, "data_block_size"), 10, 64)
				if err != nil {
					return arr, errors.New("Can't parse blockSize: " + err.Error())
				}
				blockSize *= sectorSize
			case "device":
				arr = append(arr, Device{})
				dev = &arr[len(arr)-1]
				dev.Id, err = strconv.Atoi(getAttr(t.Attr, "dev_id"))
				if err != nil {
					return arr, errors.New("Can't parse device id: " + getAttr(t.Attr, "dev_id"))
				}
			case "single_mapping":
				dev.Blocks = append(dev.Blocks, Block{Length:blockSize})
				block := &dev.Blocks[len(dev.Blocks)-1]
				block.OriginOffset, err = strconv.ParseInt(getAttr(t.Attr, "origin_block"), 10, 64)
				block.OriginOffset *= blockSize
				if err != nil {
					return arr, errors.New("Can't parse single_mapping block origin offset '" + getAttr(t.Attr, "origin_block") + "' :" + err.Error())
				}
				block.DataOffset, err = strconv.ParseInt(getAttr(t.Attr, "data_block"), 10, 64)
				block.DataOffset *= blockSize
				if err != nil {
					return arr, errors.New("Can't parse single_mapping block data offset '" + getAttr(t.Attr, "data_block") + "' :" + err.Error())
				}

			case "range_mapping":
				dev.Blocks = append(dev.Blocks, Block{})
				block := &dev.Blocks[len(dev.Blocks)-1]
				block.OriginOffset, err = strconv.ParseInt(getAttr(t.Attr, "origin_begin"), 10, 64)
				block.OriginOffset *= blockSize
				if err != nil{
					return arr, errors.New("Can't parse range_mapping block origin offset '" + getAttr(t.Attr, "origin_begin") + "' :" + err.Error())
				}
				block.DataOffset, err = strconv.ParseInt(getAttr(t.Attr, "data_begin"), 10, 64)
				block.DataOffset *= blockSize
				if err != nil {
					return arr, errors.New("Can't parse range_mapping block data offset '" + getAttr(t.Attr, "data_begin") + "' :" + err.Error())
				}
				block.Length, err = strconv.ParseInt(getAttr(t.Attr, "length"), 10, 64)
				block.Length *= blockSize
				if err != nil {
					return arr, errors.New("Can't parse range_mapping length '" + getAttr(t.Attr, "length") + "': " + err.Error())
				}
			}
		case xml.EndElement:
			if t.Name.Local == "superblock" {
				return arr, nil
			}
		}
	}
	return arr, errors.New("Unexpected exit of funtion")
}

func getAttr(arr []xml.Attr, name string) string {
	for _, attr := range arr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}
