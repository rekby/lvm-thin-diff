package lvm_thin_diff

import (
	"testing"
	"bytes"
	"encoding/xml"
)

func TestGetAttr(t *testing.T){
	if "" != getAttr(nil, "test"){
		t.Error()
	}

	if "val2" != getAttr([]xml.Attr{
		{xml.Name{"","key1"},"val1"},
		{xml.Name{"","key2"}, "val2"},
		{xml.Name{"","key3"},"val3"}},
		"key2") {
		t.Error()
	}
}

func TestParser(t *testing.T){
	data:=`<superblock uuid="" time="5" transaction="17" data_block_size="128" nr_data_blocks="0">
  <device dev_id="1" mapped_blocks="4096000" transaction="0" creation_time="0" snap_time="0">
    <single_mapping origin_block="0" data_block="0" time="0"/>
    <range_mapping origin_begin="8190976" data_begin="8461326" length="80" time="5"/>
    <range_mapping origin_begin="8191056" data_begin="8461529" length="9" time="5"/>
    <single_mapping origin_block="8191999" data_block="4146006" time="5"/>
  </device>
  <device dev_id="2" mapped_blocks="40960001" transaction="1" creation_time="1" snap_time="1">
    <single_mapping origin_block="1" data_block="1" time="1"/>
    <range_mapping origin_begin="81909761" data_begin="84613261" length="801" time="53"/>
    <range_mapping origin_begin="81910561" data_begin="84615291" length="91" time="52"/>
  </device>
</superblock>
`
	const blockSize = 128*512

	res, err := parseMetaDataXML(bytes.NewBufferString(data))
	if err != nil {
		t.Error(err)
	}

	if len(res) != 2 {
		t.Error()
	}
	dev := res[0]
	if dev.Id != 1 {
		t.Error()
	}

	if len(dev.Blocks) != 4 {
		t.Error()
	}

	b := dataBlock{OriginOffset:0, DataOffset:0, Length: blockSize}
	if dev.Blocks[0] != b {
		t.Error()
	}

	b = dataBlock{OriginOffset:8190976*blockSize, DataOffset:8461326*blockSize, Length: 80*blockSize}
	if dev.Blocks[1] != b {
		t.Error()
	}

	b = dataBlock{OriginOffset:8191056*blockSize, DataOffset:8461529*blockSize, Length:9*blockSize}
	if dev.Blocks[2] != b {
		t.Error()
	}

	b = dataBlock{OriginOffset:8191999*blockSize, DataOffset:4146006*blockSize, Length:blockSize}
	if dev.Blocks[3] != b {
		t.Error()
	}

	// DEV2
	dev = res[1]
	if dev.Id != 2 {
		t.Error()
	}

	if len(dev.Blocks) != 3 {
		t.Error()
	}

	b = dataBlock{OriginOffset:1*blockSize, DataOffset:1*blockSize, Length: blockSize}
	if dev.Blocks[0] != b {
		t.Error()
	}

	b = dataBlock{OriginOffset:81909761*blockSize, DataOffset:84613261*blockSize, Length: 801*blockSize}
	if dev.Blocks[1] != b {
		t.Error()
	}

	b = dataBlock{OriginOffset:81910561*blockSize, DataOffset:84615291*blockSize, Length:91*blockSize}
	if dev.Blocks[2] != b {
		t.Error()
	}

}