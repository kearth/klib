//
package types

import (
	"hash"
	"io"
)

type Array struct {
}

func (a *Array) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (a *Array) Sum(b []byte) []byte {
	return nil
}

func (a *Array) Reset() {

}

func (a *Array) Size() int {
	return 0
}

func (a *Array) BlockSize() int {
	return 0
}

func Make() {

}

var _ hash.Hash = &Array{}
var _ io.Writer = &Array{}
