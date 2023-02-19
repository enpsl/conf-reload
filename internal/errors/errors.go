// Copyright 2023 enpsl. All rights reserved.

// Package errors defines the error type and functions used by
// conf-reload and its internal packages.

package errors

import (
	"errors"
	"fmt"
)

/******************************************
    Domain Specific Error Types & Values
*******************************************/

type ErrType error

func ErrFormat(errType ErrType, err error) error {
	if errors.Unwrap(err) == nil {
		return fmt.Errorf("%w", errType)
	}
	return fmt.Errorf("%s :%w", errType, errors.Unwrap(err))
}

var (
	// ErrInvalidFilePath indicates that we can't get valid files
	ErrInvalidFilePath ErrType = errors.New("invalid file path")
	// ErrInvalidFileExt indicates that we can't get valid file ext type
	ErrInvalidFileExt ErrType = errors.New("invalid file ext type")
	// ErrUnmarshaller indicates that we can't unmarshal content
	ErrUnmarshaller ErrType = errors.New("unmarshal error")
	// ErrInvalidKey indicates that we can't get valid key
	ErrInvalidKey ErrType = errors.New("key is invalid")
	// ErrBrokerDecode indicates that broker can't decode content
	ErrBrokerDecode ErrType = errors.New("broker can not decode")
)

/***************************************************************
	To Replace Go Error Package And Used As Internal Method
*****************************************************************/

// New Create internal errors
func New(text string) error { return errors.New(text) }

// Is strict check assert target err is internal errors
func Is(err, target error) bool { return errors.Is(err, target) }

// As not strict check assert target err type is internal errors
func As(err error, target interface{}) bool { return errors.As(err, target) }

// Unwrap you can recursively unpack to get the innermost error:
func Unwrap(err error) error { return errors.Unwrap(err) }
