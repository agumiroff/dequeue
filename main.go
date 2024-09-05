package main

import (
	"fmt"
	"log"
)

type Operation int64

const (
	Append Operation = iota
	Insert
)

func main() {
	deq := NewDequeue(3, 5)

	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Insert(9)
	deq.Append(1)
	deq.Append(2)
	deq.Append(3)
	deq.Append(4)
	deq.Append(5)
	deq.Append(6)
	deq.Append(1)
	deq.Append(2)
	deq.Append(3)
	deq.Append(4)
	deq.Append(1)
	deq.Append(2)
	deq.Append(3)
	deq.Append(4)
	deq.Append(5)
	fmt.Println(deq.storage)
}

type Dequeue struct {
	storage   [][]int
	fPointer  pointer
	bPointer  pointer
	bSize     int
	IsFirstOp bool
}

type pointer struct {
	bNumber int
	index   int
}

func NewDequeue(cap int, bSize int) *Dequeue {
	// Setting default values for case when buckets count is even
	fMiddle, bMiddle := cap/2, cap/2-1
	fIndex, bIndex := 0, 0
	storage := make([][]int, cap)

	// Calculate target bucket and index if buckets count is odd
	if cap%2 != 0 {
		bMiddle += 1
		if bSize%2 == 0 {
			fIndex = bSize / 2
			bIndex = bSize/2 - 1
		}
		if bSize%2 != 0 {
			fIndex = bSize / 2
			bIndex = bSize / 2
		}
	}

	fPointer := pointer{
		bNumber: fMiddle,
		index:   fIndex,
	}
	bPointer := pointer{
		bNumber: bMiddle,
		index:   bIndex,
	}

	deq := &Dequeue{
		storage:   storage,
		fPointer:  fPointer,
		bPointer:  bPointer,
		bSize:     bSize,
		IsFirstOp: true,
	}
	return deq
}

func (deq *Dequeue) Append(val int) {
	if deq.IsFirstOp {
		deq.firstOperation(Append, deq.fPointer.bNumber)
	}
	err := deq.addNewBucketIfNeeded(Append)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	deq.storage[deq.fPointer.bNumber][deq.fPointer.index] = val
	deq.fPointer.index++
	deq.calculatePointer(deq.fPointer.index, Append)
}

func (deq *Dequeue) Insert(val int) {
	if deq.IsFirstOp {
		deq.firstOperation(Insert, deq.bPointer.bNumber)
	}
	err := deq.addNewBucketIfNeeded(Insert)
	if err != nil {
		log.Printf(err.Error())
		return
	}
	deq.storage[deq.bPointer.bNumber][deq.bPointer.index] = val
	if deq.bPointer.bNumber == cap(deq.storage)/2 {
		deq.bPointer.index--
	} else {
		deq.bPointer.index++
	}
	deq.calculatePointer(deq.bPointer.index, Insert)
}

// Helpers
func (deq *Dequeue) calculatePointer(index int, op Operation) {
	// check if index out of range, change bucket and reset pointerIndex
	switch op {
	case Append:
		if index > len(deq.storage[deq.fPointer.bNumber])-1 {
			deq.fPointer.bNumber++
			deq.fPointer.index = 0
		}
	case Insert:
		if index > len(deq.storage[deq.bPointer.bNumber])-1 || index < 0 {
			deq.bPointer.bNumber--
			deq.bPointer.index = 0
		}
	}
}

func (deq *Dequeue) addNewBucketIfNeeded(op Operation) error {
	bucketNumber := 0
	switch op {
	case Append:
		bucketNumber = deq.fPointer.bNumber
	case Insert:
		bucketNumber = deq.bPointer.bNumber
	}
	fmt.Println(bucketNumber)
	if cap(deq.storage)-1 >= bucketNumber && bucketNumber >= 0 {
		if len(deq.storage[bucketNumber]) == 0 {
			deq.storage[bucketNumber] = make([]int, deq.bSize)
		}
		return nil
	} else {
		return fmt.Errorf("there is no bucket capacity to add")
	}
}

func (deq *Dequeue) firstOperation(op Operation, bNumber int) {
	if bNumber == len(deq.storage)/2 {
		switch op {
		case Append:
			deq.bPointer.index--
		case Insert:
			deq.fPointer.index++
		}
	}
	deq.IsFirstOp = false
}
