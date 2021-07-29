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

type World struct {
	componentIds    map[reflect.Type]int
	componentStores []Store
	systems         []systemBinder

	entities          *list.List
	freeEntityIndexes []int
	nextEntityIndex   int
}

func (w *World) AddSystems(systems ...System) {
	for _, system := range systems {
		w.systems = append(w.systems, newSystemBinder(w, system))
	}
}

func (w *World) NewEntity(components ...interface{}) *Entity {
	entity := &Entity{world: w, index: w.getNewEntityIndex(), mask: makeMask(len(w.componentStores))}
	entity.element = w.entities.PushBack(entity)
	entity.Set(components...)
	return entity
}

func (w *World) RemoveEntity(entity *Entity) {
	w.freeEntityIndexes = append(w.freeEntityIndexes, entity.index)
	w.entities.Remove(entity.element)

	for i := 0; i < len(w.componentStores); i++ {
		if entity.mask.Get(i) {
			w.componentStores[i].Rem(entity.index)
		}
	}

	entity.world, entity.element, entity.index, entity.mask = nil, nil, -1, nil
}

func (w *World) Update() {
	for i := range w.systems {
		for element := w.entities.Front(); element != nil; element = element.Next() {
			entity := element.Value.(*Entity)
			w.systems[i].update(entity)
		}
	}
}

func (w *World) getNewEntityIndex() int {
	var index int

	if len(w.freeEntityIndexes) == 0 {
		index = w.nextEntityIndex
		w.nextEntityIndex++
	} else {
		index = w.freeEntityIndexes[len(w.freeEntityIndexes)-1]
		w.freeEntityIndexes = w.freeEntityIndexes[:len(w.freeEntityIndexes)-1]
	}

	return index
}
