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

type NormalComponent struct{ value int }
type SingletonComponent struct{ value int }

func main() {
	world := ento.NewWorldBuilder().
		WithSparseComponents(NormalComponent{}).
		// Use concrete value for singleton component as it is shared across all entities
		WithSingletonComponents(SingletonComponent{1}).
		Build(8)

	// Use zero-value for singleton component as the value is ignored
	// Only value used at registration is used
	entity1 := world.AddEntity(NormalComponent{1}, SingletonComponent{})
	entity2 := world.AddEntity(NormalComponent{2}, SingletonComponent{})

	// Singleton components must still be explicitly added to entities
	// Get(c **SingletonComponent) on below entity will set c to nil
	world.AddEntity(NormalComponent{3})

	var singleton1, singleton2 *SingletonComponent
	entity1.Get(&singleton1)
	entity2.Get(&singleton2)

	println(singleton1 == singleton2) // Prints: true

	singleton1.value = 2
	println(singleton1.value == singleton2.value) // Prints: true
}
