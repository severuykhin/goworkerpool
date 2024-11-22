package goconcurrentpool

import "errors"

var (
	ErrPoolNotActive = errors.New("pool is not active")
)
