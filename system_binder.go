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
	"fmt"
	"reflect"
)

type systemBinder struct {
	System
	mask   mask
	stores []Store
	values []reflect.Value
}

func newSystemBinder(world *World, system System) systemBinder {
	const (
		componentTag     = "ento"
		componentTagBind = "bind"
	)

	systemValue := reflect.ValueOf(system).Elem()
	systemType := systemValue.Type()
	systemFieldsNum := systemType.NumField()

	mask := makeMask(len(world.componentStores))
	stores := make([]Store, 0, systemFieldsNum)
	values := make([]reflect.Value, 0, systemFieldsNum)

	for i := 0; i < systemFieldsNum; i++ {
		field := systemType.Field(i)
		tag, ok := field.Tag.Lookup(componentTag)
		if !ok {
			continue
		}

		componentValue := systemValue.Field(i)
		componentType := componentValue.Type().Elem()
		componentId := world.componentIds[componentType]

		switch tag {
		case componentTagBind:
			mask.Set(componentId)
			stores = append(stores, world.componentStores[componentId])
			values = append(values, componentValue)
		default:
			panic(fmt.Sprintf("unrecognised tag: %s", tag))
		}
	}

	return systemBinder{System: system, mask: mask, stores: stores, values: values}
}

func (b *systemBinder) update(entity *Entity) {
	if entity.mask.Contains(b.mask) {
		for i, store := range b.stores {
			store.Get(entity.index, b.values[i])
		}
		b.Update(entity)
	}
}
