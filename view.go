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

import "reflect"

type View struct {
	world *World
	mask  mask
}

func NewView(world *World, components ...interface{}) *View {
	mask := makeMask(len(world.componentStores))

	for _, component := range components {
		componentType := reflect.TypeOf(component)
		componentId := world.componentIds[componentType]
		mask.Set(componentId)
	}

	return &View{world: world, mask: mask}
}

func (v *View) Each(consumer func(entity *Entity)) {
	for element := v.world.entities.Front(); element != nil; element = element.Next() {
		entity := element.Value.(*Entity)
		if entity.mask.Contains(v.mask) {
			consumer(entity)
		}
	}
}
