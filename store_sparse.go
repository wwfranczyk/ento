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
	"reflect"
)

type SparseStore struct {
	zero     reflect.Value
	slice    reflect.Value
	capacity int
}

func SparseStoreProvider(component interface{}) StoreProvider {
	return func(capacity int) Store { return NewSparseStore(reflect.TypeOf(component), capacity) }
}

func NewSparseStore(componentType reflect.Type, capacity int) *SparseStore {
	return &SparseStore{
		zero:     reflect.Zero(componentType),
		slice:    reflect.MakeSlice(reflect.SliceOf(componentType), capacity, capacity),
		capacity: capacity,
	}
}

func (s *SparseStore) Add(id int, value reflect.Value) {
	s.EnsureCapacity(id + 1)
	s.Set(id, value)
}

func (s *SparseStore) Get(id int, value reflect.Value) {
	value.Set(s.slice.Index(id).Addr())
}

func (s *SparseStore) Set(id int, value reflect.Value) {
	s.slice.Index(id).Set(value)
}

func (s *SparseStore) Rem(id int) {
	s.slice.Index(id).Set(s.zero)
}

func (s *SparseStore) EnsureCapacity(capacity int) {
	if capacity < s.capacity {
		return
	}

	capacity = nextPowerOf2(capacity)

	slice := reflect.MakeSlice(s.slice.Type(), capacity, capacity)
	reflect.Copy(slice, s.slice)

	s.slice, s.capacity = slice, capacity
}
