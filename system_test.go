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

type c1 struct{ v int }
type c2 struct{ v int }

type testSystem struct {
	C1 *c1 `ento:"required"`
	C2 *c2 `ento:"optional"`

	sums []int
}

func (t *testSystem) Update(entity *Entity) {
	sum := t.C1.v

	if t.C2 != nil {
		sum += t.C2.v
	}

	t.sums[entity.index] = sum
}

func TestSystem(t *testing.T) {
	const N = 8

	w := NewWorldBuilder().WithSparseComponents(c1{}, c2{}).Build(N)

	ts := &testSystem{sums: make([]int, N)}
	w.AddSystems(ts)

	for i := 0; i < N; i++ {
		e := w.NewEntity()
		e.Set(c1{i})
		if i%2 == 0 {
			e.Set(c2{i})
		}
	}

	w.Update()

	for i := 0; i < N; i++ {
		sum := ts.sums[i]
		expected := i
		if i%2 == 0 {
			expected += i
		}

		assert.Equal(t, expected, sum)
	}
}
