package lvm_thin_diff

import "testing"

func TestCalcDiff(t *testing.T){
	var diff, expectedDiff dataPatch

	diff = calcDiff(dataBlock{}, dataBlock{})
	if diff != expectedDiff {
		t.Errorf("%#v", diff)
	}

	expectedDiff = dataPatch{Offset:100, Length: 50, Operation:DELETE}
	diff = calcDiff(dataBlock{OriginOffset:100, DataOffset:200, Length:50}, dataBlock{})
	if diff != expectedDiff {
		t.Errorf("%#v", diff)
	}

	expectedDiff = dataPatch{Offset:100, Length: 50, Operation:WRITE}
	diff = calcDiff(dataBlock{}, dataBlock{OriginOffset:100, DataOffset:200, Length:50})
	if diff != expectedDiff {
		t.Errorf("%#v", diff)
	}

	expectedDiff = dataPatch{Operation:NONE}
	diff = calcDiff(dataBlock{OriginOffset:100, DataOffset:200, Length:50}, dataBlock{OriginOffset:100, DataOffset:200, Length:50})
	if diff != expectedDiff {
		t.Errorf("%#v", diff)
	}

	expectedDiff = dataPatch{Operation:WRITE, Offset:100, Length: 50}
	diff = calcDiff(dataBlock{OriginOffset:100, DataOffset:500, Length:50}, dataBlock{OriginOffset:100, DataOffset:200, Length:50})
	if diff != expectedDiff {
		t.Errorf("%#v", diff)
	}
}

