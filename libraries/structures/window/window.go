package window

// Window is the interface implementation of windows. They essentially behave as queues
// with nice properties
type Window interface {
	Push(interface{}) error
	Mean() float32
	StdDev(float32) float32
}
