package infrastructure

import "github.com/creasty/defaults"

type DefaultSetterImpl struct{}

func NewDefaultSetterImpl() *DefaultSetterImpl {
	return &DefaultSetterImpl{}
}

func (d *DefaultSetterImpl) Set(ptr interface{}) error {
	return defaults.Set(ptr)
}
