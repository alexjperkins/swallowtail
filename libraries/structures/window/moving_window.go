package window

import (
	"fmt"
	"swallowtail/libraries/structures/queues"
)

// MovingWindow is a moving window for floats.
type MovingWindow struct {
	q *queues.FIFOQueue
}

func NewMovingWindow(maxSize int) *MovingWindow {
	return &MovingWindow{
		q: queues.NewFIFOQueue(maxSize),
	}
}

func (mw *MovingWindow) Push(item float32) error {
	if mw.q.AtCapacity() {
		mw.q.Pop()
	}
	return mw.q.Push(item)
}

func (mw *MovingWindow) Mean() (float32, error) {
	var (
		total   float32
		counter float32
	)
	for _, v := range mw.q.GetAsArray() {
		counter++
		switch i := v.(type) {
		case float32:
			total += i
		default:
			return total / counter, fmt.Errorf("Error converting queue value to float32: %v", v)
		}
	}
	if counter == 0.0 {
		return 0.0, nil
	}
	return total / counter, nil
}

func (mw *MovingWindow) StdDev(nMeanAwayFromNormal float32) float32 {
	// TODO
	return 0.0
}

func (mw *MovingWindow) Len() int {
	return mw.q.Len()
}
