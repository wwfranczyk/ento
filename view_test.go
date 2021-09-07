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

func TestView(t *testing.T) {
	const N = 8

	type C1 int
	type C2 int

	world := NewWorldBuilder().WithSparseComponents(C1(0), C2(0)).Build(N)

	for i := 0; i < N; i++ {
		entity := world.AddEntity(C1(0))
		if i%2 == 0 {
			entity.Set(C2(0))
		}
	}

	view := NewView(world, C1(0), C2(0))

	visited := make(map[*Entity]bool, N)
	view.Each(func(entity *Entity) {
		visited[entity] = true
	})

	assert.Equal(t, (N+1)/2, len(visited))
}
