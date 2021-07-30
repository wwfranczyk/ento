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

type DenseStore struct {
	sparse  *SparseStore
	mapping map[int]int
	indexes IndexPool
}

func DenseStoreProvider(component interface{}) StoreProvider {
	return func(capacity int) Store { return NewDenseStore(reflect.TypeOf(component), capacity) }
}

func NewDenseStore(componentType reflect.Type, capacity int) *DenseStore {
	return &DenseStore{
		sparse:  NewSparseStore(componentType, capacity),
		mapping: make(map[int]int, capacity),
		indexes: NewIndexPool(capacity),
	}
}

func (d *DenseStore) Add(id int, value reflect.Value) {
	index := d.indexes.GetFree()
	d.mapping[id] = index
	d.sparse.Add(index, value)
}

func (d *DenseStore) Get(id int, value reflect.Value) {
	index := d.mapping[id]
	d.sparse.Get(index, value)
}

func (d *DenseStore) Set(id int, value reflect.Value) {
	index := d.mapping[id]
	d.sparse.Set(index, value)
}

func (d *DenseStore) Rem(id int) {
	index := d.mapping[id]
	d.sparse.Rem(index)

	delete(d.mapping, id)
	d.indexes.Release(index)
}
