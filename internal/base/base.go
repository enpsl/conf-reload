// Copyright 2023 enpsl. All rights reserved.

// Broker interface

package base

import (
	"io"
)

type Broker interface {
	Parse(content []byte) (error, map[string]interface{})
	LoadContent() ([]byte, error)
	Watch()
	Decode(input interface{}, output interface{}, weaklyTypedInput bool) error
	Notify() <-chan struct{}
	io.Closer
}
