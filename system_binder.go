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
	mask     mask
	ids      []int
	stores   []Store
	values   []reflect.Value
	zeros    []reflect.Value
	required []bool
}

func newSystemBinder(world *World, system System) systemBinder {
	const (
		componentTag         = "ento"
		componentTagRequired = "required"
		componentTagOptional = "optional"
	)

	systemValue := reflect.ValueOf(system).Elem()
	systemType := systemValue.Type()
	systemFieldsNum := systemType.NumField()

	mask := makeMask(len(world.componentStores))
	ids := make([]int, 0, systemFieldsNum)
	stores := make([]Store, 0, systemFieldsNum)
	values := make([]reflect.Value, 0, systemFieldsNum)
	zeros := make([]reflect.Value, 0, systemFieldsNum)
	required := make([]bool, 0, systemFieldsNum)

	for i := 0; i < systemFieldsNum; i++ {
		field := systemType.Field(i)
		tag, ok := field.Tag.Lookup(componentTag)
		if !ok {
			continue
		}

		componentValue := systemValue.Field(i)
		componentType := componentValue.Type().Elem()
		componentId := world.componentIds[componentType]

		ids = append(ids, componentId)
		stores = append(stores, world.componentStores[componentId])
		values = append(values, componentValue)
		zeros = append(zeros, reflect.Zero(componentValue.Type()))

		switch tag {
		case componentTagRequired:
			mask.Set(componentId)
			required = append(required, true)
		case componentTagOptional:
			required = append(required, false)
		default:
			panic(fmt.Sprintf("unrecognised tag: %s", tag))
		}
	}

	return systemBinder{
		System:   system,
		mask:     mask,
		ids:      ids,
		stores:   stores,
		values:   values,
		zeros:    zeros,
		required: required,
	}
}

func (b *systemBinder) update(entity *Entity) {
	if entity.mask.Contains(b.mask) {
		for i := range b.stores {
			if b.required[i] {
				b.stores[i].Get(entity.index, b.values[i])
			} else if entity.mask.Get(b.ids[i]) {
				b.stores[i].Get(entity.index, b.values[i])
			} else {
				b.values[i].Set(b.zeros[i])
			}
		}
		b.Update(entity)
	}
}
