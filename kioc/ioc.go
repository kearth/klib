package kioc

import (
	"sync"

	"github.com/kearth/tea/terrors"
)

// IContainer 容器
type IContainer interface {
	// 名称
	Name() string
	// 新容器
	New() IContainer
}

// ioc
type ioc struct{}

var (
	// iocMap
	iocMap map[string]IContainer = make(map[string]IContainer)

	// ioc instance
	iocInstance *ioc

	// once
	once sync.Once
)

// IOC
func IOC() *ioc {
	once.Do(func() {
		iocInstance = &ioc{}
	})
	return iocInstance
}

// Register
func (i *ioc) Register(ic IContainer) terrors.IError {
	if _, ok := iocMap[ic.Name()]; !ok {
		iocMap[ic.Name()] = ic
		return nil
	}
	return terrors.NameRegistered
}

// Get
func (i *ioc) Get(name string) (IContainer, terrors.IError) {
	if ic, ok := iocMap[name]; ok {
		return ic.New(), nil
	}
	return nil, terrors.NameNotRegistered
}
