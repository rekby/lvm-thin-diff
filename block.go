package lvm_thin_diff

import (
	"fmt"
	"sort"
)

type Block struct {
	OriginOffset int64
	DataOffset   int64
	Length       int64
}

func (this *Block) IsEmpty() bool {
	if this == nil {
		return true
	}
	return *this == Block{}
}

// Offset for next then last byte of origin data
func (this *Block) OriginLast() int64 {
	return this.OriginOffset + this.Length
}

func (this *Block) Split(length int64) (left, right Block) {
	if length >= this.Length {
		left = *this
	} else {
		left.DataOffset = this.DataOffset
		left.OriginOffset = this.OriginOffset
		left.Length = length
	}

	right.DataOffset = left.DataOffset + left.Length
	right.OriginOffset = left.OriginOffset + left.Length
	right.Length = this.Length - left.Length
	return left, right
}

type Device struct {
	Id     int
	Blocks BlockArr
	Size   int64
}

// Operation
const (
	NONE = iota
	WRITE
	DELETE
	SET_SIZE
)

type Diff struct {
	Operation int
	Offset    int64
	Length    int64
}

type BlockArr []Block

var _ sort.Interface = BlockArr{}

func (arr BlockArr) Len() int {
	return len(arr)
}

func (arr BlockArr) Less(i, j int) bool {
	return arr[i].DataOffset < arr[j].DataOffset
}

func (arr BlockArr) Swap(i, j int) {
	tmp := arr[i]
	arr[i] = arr[j]
	arr[j] = tmp
}

/*
Основная задача структуры - предоставление интерфейс "откусывания" равных кусков из from и to массивов для удобства
дальнейшей работы.
После каждого Cut на выходе получаются два блока данных, каждый из которых может быть пустым или описывающим расположение данных.
Если оба блока данных не пустые, то их логические смещения (OriginOffset) и длины - совпадают.
*/
type DataBlockCutter struct {
	From, To BlockArr // Рабочие массивы, отсортированы по Originffset. ВАЖНО - портятся в процессе работы.
}

/*
ВАЖНО - массивы from, to ПОРТЯТСЯ в процессе работы. После вызова функции их можно использвоаться только для продолжения
работы этой же функции.

from, to - отсортированных по OriginOffset.
bFrom и bTo - два блока данных. Если данные присутствую в обоих блоках, то их начало в OriginOffset и длина совпадают.

Например (0 - пропуск)
from: 0ABCDEFGH0000
to:   000IGKLMNOPQR
      0112222223333

функция будет последовательно возвращать:
bFrom: AB, bTo: 00
bFrom: CDEFGH, bTo: IJKLMN
bFrom: 0000, bTo: OPQR

newFrom, newTo - оставшиеся блоки для продолжения работы.
*/
func (this *DataBlockCutter) Cut() (ok bool, bFrom, bTo Block) {
	switch {
	case len(this.From) == 0 && len(this.To) == 0:
		return // возвращаем пустые данные
	case len(this.From) == 0 && len(this.To) != 0:
		// EmptyFrom
		ok = true
		bTo = this.To[0]
		this.To = this.To[1:]
		return
	case len(this.From) != 0 && len(this.To) == 0:
		// Empty To
		ok = true
		bFrom = this.From[0]
		this.From = this.From[1:]
		return
	default:
		// pass
	}

	firstFrom := &this.From[0]
	firstTo := &this.To[0]

	switch {
	case firstFrom.Length == 0:
		// firstFrom empty
		this.From = this.From[1:]
		return this.Cut()
	case firstTo.Length == 0:
		// firstTo empty
		this.To = this.To[1:]
		return this.Cut()

	case firstFrom.OriginLast() <= firstTo.OriginOffset:
		// firstFrom end before firstTo start
		ok = true
		bFrom = *firstFrom
		this.From = this.From[1:]
		return
	case firstTo.OriginLast() <= firstFrom.OriginOffset:
		// firstTo end before firstFromStart
		ok = true
		bTo = *firstTo
		this.To = this.To[1:]
		return

	case firstFrom.OriginOffset < firstTo.OriginOffset:
		// firstFrom start before firstTo. Overlap.
		ok = true
		length := firstTo.OriginOffset - firstFrom.OriginOffset
		bFrom, *firstFrom = firstFrom.Split(length)
		return

	case firstTo.OriginOffset < firstFrom.OriginOffset:
		// firstTo start before firstFrom. Overlap
		ok = true
		length := firstFrom.OriginOffset - firstTo.OriginOffset
		bTo, *firstTo = firstTo.Split(length)
		return

	case firstFrom.OriginOffset == firstTo.OriginOffset:
		// Equal start
		switch {
		case firstFrom.Length < firstTo.Length:
			// firstFrom shorter then firstTo
			ok = true
			bFrom = *firstFrom
			this.From = this.From[1:]
			bTo, *firstTo = firstTo.Split(firstFrom.Length)
			return
		case firstTo.Length < firstFrom.Length:
			// firstTo shorter then firstFrom
			ok = true
			bTo = *firstTo
			this.To = this.To[1:]
			bFrom, *firstFrom = firstFrom.Split(firstTo.Length)
			return
		case firstTo.Length == firstFrom.Length:
			// firstTo equal length to firstFrom
			ok = true
			bFrom = *firstFrom
			bTo = *firstTo
			this.From = this.From[1:]
			this.To = this.To[1:]
			return
		default:
			panic(fmt.Errorf("Unhandled variant in cutHeader 2 %#v %#v:", *firstFrom, *firstTo))
		}
	default:
		panic(fmt.Errorf("Unhandled variant in cutHead: %#v %#v", *firstFrom, *firstTo))
	}
}
