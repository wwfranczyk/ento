// MIT License
//
// Copyright (c) 2021 Wojciech Franczyk
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ento

import (
	"container/list"
	"reflect"
)

type Entity struct {
	world   *World
	element *list.Element
	index   int
	mask    mask
}

func (e *Entity) Get(components ...interface{}) bool {
	all := true

	for _, component := range components {
		componentValue := reflect.ValueOf(component).Elem()
		componentType := componentValue.Type().Elem()
		componentId := e.world.componentIds[componentType]

		if e.mask.Get(componentId) {
			e.world.componentStores[componentId].Get(e.index, componentValue)
		} else {
			componentValue.Set(reflect.Zero(reflect.PtrTo(componentType)))
			all = false
		}
	}

	return all
}

func (e *Entity) Set(components ...interface{}) {
	for _, component := range components {
		componentValue := reflect.ValueOf(component)
		componentType := componentValue.Type()
		componentId := e.world.componentIds[componentType]

		if e.mask.Get(componentId) {
			e.world.componentStores[componentId].Set(e.index, componentValue)
			return
		}

		e.world.componentStores[componentId].Add(e.index, componentValue)
		e.mask.Set(componentId)
	}
}

func (e *Entity) Rem(components ...interface{}) {
	for _, component := range components {
		componentValue := reflect.ValueOf(component)
		componentType := componentValue.Type()
		componentId := e.world.componentIds[componentType]

		if e.mask.Get(componentId) {
			e.world.componentStores[componentId].Rem(e.index)
			e.mask.Clear(componentId)
		}
	}
}
