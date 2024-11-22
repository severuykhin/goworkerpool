package goconcurrentpool

type Future struct {
	job job
	Err error
}
