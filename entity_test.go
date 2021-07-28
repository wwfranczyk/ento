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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntity(t *testing.T) {
	type C1 struct{ v int }
	type C2 struct{ v int }
	type C3 struct{ v int }

	w := NewWorldBuilder().
		WithSparseComponents(C1{}, C2{}).
		WithSingletonComponents(C3{3}).
		Build(2)

	e1 := w.NewEntity()
	e2 := w.NewEntity()

	e1.Set(C1{1}, C2{2}, C3{1})
	e2.Set(C1{3}, C3{2})

	var c1 *C1
	var c2 *C2
	var c3e1 *C3
	var c3e2 *C3

	all := e1.Get(&c1, &c2, &c3e1)
	assert.True(t, all)
	assert.Equal(t, 1, c1.v)
	assert.Equal(t, 2, c2.v)
	assert.Equal(t, 3, c3e1.v)

	all = e2.Get(&c1, &c2, &c3e2)
	assert.False(t, all)
	assert.Equal(t, 3, c1.v)
	assert.Nil(t, c2)
	assert.Equal(t, 3, c3e2.v)

	assert.Equal(t, c3e1, c3e2)

	e1.Rem(C2{}, C3{})
	all = e1.Get(&c1, &c2, &c3e1)
	assert.False(t, all)
	assert.Equal(t, 1, c1.v)
	assert.Nil(t, c2)
	assert.Nil(t, c3e1)
}
