package lvm_thin_diff

import (
	"fmt"
	"sort"
)

type dataBlock struct {
	OriginOffset int64
	DataOffset   int64
	Length       int64
}

func (this *dataBlock) IsEmpty() bool {
	if this == nil {
		return true
	}
	return *this == dataBlock{}
}

// Offset for next then last byte of origin data
func (this *dataBlock) OriginLast() int64 {
	return this.OriginOffset + this.Length
}

func (this *dataBlock) Split(length int64) (left, right dataBlock) {
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

type dataDevice struct {
	Id     int
	Blocks blockArr
	Size   int64
}

// Operation
const (
	NONE = iota
	WRITE
	DELETE
	SET_SIZE
)

type dataPatch struct {
	Operation int
	Offset    int64
	Length    int64
}

type blockArr []dataBlock

var _ sort.Interface = blockArr{}

func (arr blockArr) Len() int {
	return len(arr)
}

func (arr blockArr) Less(i, j int) bool {
	return arr[i].DataOffset < arr[j].DataOffset
}

func (arr blockArr) Swap(i, j int) {
	tmp := arr[i]
	arr[i] = arr[j]
	arr[j] = tmp
}

/*
Основная задача структуры - предоставление интерфейс "откусывания" равных кусков из from и to массивов для удобства
дальнейшей работы.
После каждого Cut на выходе получаются два блока данных, каждый из которых может быть пустым или описывающим расположение данных.
Если оба блока данных не пустые, то их логические смещения (OriginOffset) и длины - совпадают.
ВАЖНО - From и To портятся в процессе работы.
*/
type dataBlockArrCutter struct {
	From, To blockArr // Рабочие массивы, отсортированы по Originffset. ВАЖНО - портятся в процессе работы.
}

/*
ВАЖНО - массивы from, to ПОРТЯТСЯ в процессе работы. После вызова функции их можно использвоаться только для продолжения
работы этой же функции.

"Подравнивает" from и to.
Если from и to начинаются с одного смещения - возвращает блоки с этого смещения и с наименьшей длиной.
Если блок в from или to начинается и заканчивается раньше блока данных другого массива - возвращает этот блок (+ пустой блок вторым)
Если from и to начинаются в разных местах, но перекрываются - возвращает кусок данных. Который начинается раньше и длиной до начала
        блока данных второго массива. Чтобы при следующем вызове вернуться в ситуацию, когда массивы начинаются по одному смешению.
*/
func (this *dataBlockArrCutter) Cut() (ok bool, bFrom, bTo dataBlock) {
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

/*
Создать команду для патча данных from так чтобы получились данные to.
bFrom и bTo - два блока данных. Если оба блока не пустые - то они должны начинаться с одного логического смещения и
быть равной длины.

Пустой блок означает что в месте, указанном вторым блоком данных нет.
 */
func makePatch(bFrom, bTo dataBlock) dataPatch {
	if bFrom.IsEmpty() && bTo.IsEmpty() {
		return dataPatch{Operation:NONE}
	}

	if bFrom.IsEmpty() {
		return dataPatch{Offset:bTo.OriginOffset, Operation:WRITE, Length: bTo.Length}
	}

	if bTo.IsEmpty() {
		return dataPatch{Offset:bFrom.OriginOffset, Operation:DELETE, Length:bFrom.Length}
	}

	if bFrom.OriginOffset != bTo.OriginOffset || bFrom.Length != bTo.Length {
		panic("bFrom and bTo must have same start and length")
	}

	if bFrom.DataOffset == bTo.DataOffset {
		return dataPatch{Operation:NONE} // Data is equal. Do nothing.
	}

	return dataPatch{Offset:bTo.OriginOffset, Operation:WRITE, Length: bTo.Length}
}
