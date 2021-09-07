// +build examples

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

package main

import "github.com/wfranczyk/ento"

type Component1 struct{ value int }
type Component2 struct{ value int }

func main() {
	// Create the world and register components
	world := ento.NewWorldBuilder().
		// Use "zero-values" as values are ignored when registering
		WithSparseComponents(Component1{}, Component2{}).
		// Pre-allocate space for 256 entities (world can grow beyond that automatically)
		Build(256)

	// Add entities
	world.AddEntity(Component1{1})
	world.AddEntity(Component1{2}, Component2{2})

	// Create a view and define required components
	view := ento.NewView(world, Component1{}, Component2{})

	// Use Each to iterate all entities with required components
	view.Each(func(entity *ento.Entity) {
		var c1 *Component1
		var c2 *Component2
		entity.Get(&c1, &c2)

		// No need to check for nil
		println(c1.value, c2.value)
	})
}
