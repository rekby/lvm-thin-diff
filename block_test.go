package lvm_thin_diff

import "testing"

func TestCutHead(t *testing.T){
	var data DataBlockCutter
	var bFrom, bTo, expectedBFrom, expectedBTo Block
	var expectedFromArr, expectedToArr BlockArr
	var ok, expectedOk bool

	equals := func(a,b BlockArr)bool{
		if len(a) != len(b){
			return false
		}
		for i := range a{
			if a[i] != b[i]{
				return false
			}
		}
		return true
	}

	isOk := func()(res bool){
		res = true
		if ok != expectedOk {
			t.Errorf("ok: %#v != %#v", ok, expectedOk)
			res = false
		}
		if bFrom != expectedBFrom {
			t.Errorf("bFrom: %#v != %#v", bFrom, expectedBFrom)
			res = false
		}
		if bTo != expectedBTo {
			t.Errorf("bTo: %#v != %#v", bTo, expectedBTo)
			res = false
		}
		if !equals(data.From, expectedFromArr) {
			t.Errorf("newFrom: %#v != %#v", data.From, expectedFromArr)
			res = false
		}
		if !equals(data.To, expectedToArr) {
			t.Errorf("newTo: %#v != %#v", data.To, expectedToArr)
			res = false
		}
		return res
	}

	if !isOk() {
		t.Error()
	}

	expectedOk = true

	// EmptyFrom
	data = DataBlockCutter{
		To:BlockArr{
			Block{OriginOffset:100,DataOffset:200,Length:300},
			Block{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = Block{}
	expectedBTo = Block{OriginOffset:100,DataOffset:200,Length:300}
	expectedFromArr =BlockArr{}
	expectedToArr = BlockArr{
		Block{OriginOffset:400,DataOffset:500,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// EmptyTo
	data = DataBlockCutter{
		From:BlockArr{
			Block{OriginOffset:100,DataOffset:200,Length:300},
			Block{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = Block{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = Block{}
	expectedFromArr = BlockArr{
		Block{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = BlockArr{}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}


	// firstFrom empty
	data = DataBlockCutter{
		From:BlockArr{
			Block{OriginOffset:10,DataOffset:20,Length:0},
			Block{OriginOffset:100,DataOffset:200,Length:300},
			Block{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:BlockArr{
			Block{OriginOffset:100,DataOffset:250,Length:300},
			Block{OriginOffset:400,DataOffset:550,Length:600},
		},
	}
	expectedBFrom = Block{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = Block{OriginOffset:100,DataOffset:250,Length:300}
	expectedFromArr = BlockArr{
		Block{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = BlockArr{
		Block{OriginOffset:400,DataOffset:550,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}
}

func TestSplit(t *testing.T){
	b := Block{DataOffset:100, OriginOffset:200, Length:50}

	l,r := b.Split(0)
	lOK := Block{DataOffset:100, OriginOffset:200, Length:0}
	rOK := Block{DataOffset:100, OriginOffset:200, Length:50}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}


	l,r = b.Split(10)
	lOK = Block{DataOffset:100, OriginOffset:200, Length:10}
	rOK = Block{DataOffset:110, OriginOffset:210, Length:40}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}

	l,r = b.Split(50)
	lOK = Block{DataOffset:100, OriginOffset:200, Length:50}
	rOK = Block{DataOffset:150, OriginOffset:250, Length:0}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}

	l,r = b.Split(100)
	lOK = Block{DataOffset:100, OriginOffset:200, Length:50}
	rOK = Block{DataOffset:150, OriginOffset:250, Length:0}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}
}
