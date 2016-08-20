package lvm_thin_diff

import "testing"

func TestCutHead(t *testing.T){
	var data dataBlockArrCutter
	var bFrom, bTo, expectedBFrom, expectedBTo dataBlock
	var expectedFromArr, expectedToArr blockArr
	var ok, expectedOk bool

	equals := func(a,b blockArr)bool{
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
	// FROM:
	// TO:   DDDDDDDD
	data = dataBlockArrCutter{
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = dataBlock{}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedFromArr = blockArr{}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// EmptyTo
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = dataBlock{}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}


	// firstFrom empty
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:10,DataOffset:20,Length:0},
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:250,Length:300},
			dataBlock{OriginOffset:400,DataOffset:550,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:250,Length:300}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:550,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// firstTo empty
	// FROM: DDDDDDDD
	// TO:
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:10,DataOffset:20,Length:0},
			dataBlock{OriginOffset:100,DataOffset:250,Length:300},
			dataBlock{OriginOffset:400,DataOffset:550,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:250,Length:300}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:550,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// firstFrom end before firstTo start 1 (have space between blocks)
	// FROM: DDDDDDDD
	// TO:                 DDDDDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:800,DataOffset:550,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = dataBlock{}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:800,DataOffset:550,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// firstFrom end before firstTo start 2 (have no space between blocks)
	// FROM: DDDDDDDD
	// TO:           DDDDDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:400,DataOffset:550,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedBTo = dataBlock{}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:550,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// firstTo end before firstFrom start 1 (have space between blocks)
	// FROM:               DDDDDDDDDD
	// TO:   DDDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:800,DataOffset:550,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = dataBlock{}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:800,DataOffset:550,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// firstTo end before firstFrom start 2 (have no space between blocks)
	// FROM:         DDDDDDDDDD
	// TO:   DDDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:400,DataOffset:550,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = dataBlock{}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:200,Length:300}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:550,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}


	// firstFrom start before firstTo. Overlap.
	// FROM: DDDDDDD
	// TO:       DDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:250,DataOffset:550,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:150}
	expectedBTo = dataBlock{}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:250,DataOffset:350,Length:150},
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:250,DataOffset:550,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// firstTo start before firstFrom. Overlap.
	// FROM:     DDDDDDD
	// TO:   DDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:250,DataOffset:550,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = dataBlock{}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:200,Length:150}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:250,DataOffset:550,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:250,DataOffset:350,Length:150},
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// Equal start. firstFrom shorter then firstTo
	// FROM:  DDDDD
	// TO:    DDDDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:550,Length:200},
			dataBlock{OriginOffset:300,DataOffset:550,Length:250},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:550,Length:200}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:200,Length:200}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:300,DataOffset:550,Length:250},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:300,DataOffset:400,Length:100},
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}

	// Equal start. firstTo shorter then firstFrom
	// FROM:  DDDDD
	// TO:    DDDDDDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:300},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:550,Length:200},
			dataBlock{OriginOffset:300,DataOffset:550,Length:250},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:200}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:550,Length:200}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:300,DataOffset:400,Length:100},
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:300,DataOffset:550,Length:250},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}


	// Equal start. firstTo equal length to firstFrom
	// FROM:  DDDDD
	// TO:    DDDDD
	data = dataBlockArrCutter{
		From:blockArr{
			dataBlock{OriginOffset:100,DataOffset:200,Length:200},
			dataBlock{OriginOffset:400,DataOffset:500,Length:600},
		},
		To:blockArr{
			dataBlock{OriginOffset:100,DataOffset:550,Length:200},
			dataBlock{OriginOffset:300,DataOffset:550,Length:250},
		},
	}
	expectedBFrom = dataBlock{OriginOffset:100,DataOffset:200,Length:200}
	expectedBTo = dataBlock{OriginOffset:100,DataOffset:550,Length:200}
	expectedFromArr = blockArr{
		dataBlock{OriginOffset:400,DataOffset:500,Length:600},
	}
	expectedToArr = blockArr{
		dataBlock{OriginOffset:300,DataOffset:550,Length:250},
	}
	ok, bFrom, bTo = data.Cut()
	if !isOk(){
		t.Error()
	}
}

func TestMakePatch(t *testing.T){
	t.Error("TODO")
}

func TestSplit(t *testing.T){
	b := dataBlock{DataOffset:100, OriginOffset:200, Length:50}

	l,r := b.Split(0)
	lOK := dataBlock{DataOffset:100, OriginOffset:200, Length:0}
	rOK := dataBlock{DataOffset:100, OriginOffset:200, Length:50}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}


	l,r = b.Split(10)
	lOK = dataBlock{DataOffset:100, OriginOffset:200, Length:10}
	rOK = dataBlock{DataOffset:110, OriginOffset:210, Length:40}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}

	l,r = b.Split(50)
	lOK = dataBlock{DataOffset:100, OriginOffset:200, Length:50}
	rOK = dataBlock{DataOffset:150, OriginOffset:250, Length:0}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}

	l,r = b.Split(100)
	lOK = dataBlock{DataOffset:100, OriginOffset:200, Length:50}
	rOK = dataBlock{DataOffset:150, OriginOffset:250, Length:0}
	if l != lOK {
		t.Errorf("%#v != %#v", l, lOK)
	}
	if r != rOK{
		t.Errorf("%#v != %#v", r, rOK)
	}
}
