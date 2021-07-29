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

import (
	"github.com/wfranczyk/ento"
)

type Component1 struct{ value int }
type Component2 struct{ value int }
type Component3 struct{ value int }

type System struct {
	// Always use pointer to base component type (!)
	// Add `ento:"bind"` tag to automatically bind components
	Component1 *Component1 `ento:"bind"`

	// Components are bound by type, not name
	Renamed *Component2 `ento:"bind"`

	// Systems can have non-component-bound fields
	calls int
}

// Update implements ento.System
// It will be called only for entities that have both tagged components
func (s *System) Update(entity *ento.Entity) {
	// Access automatically bound components from currently passed entity
	s.Component1.value += s.Renamed.value

	// Get components from entity manually (somewhat slower than automatic binding)
	// Can be used for "optional" components
	var c3 *Component3
	entity.Get(&c3)

	// If entity does not have the component, the pointer value will be set to nil
	if c3 != nil {
		s.Component1.value += c3.value
	}

	s.calls++
}

func (s *System) reportCalls() {
	println(s.calls)
	s.calls = 0
}

func main() {
	// Create the world and register components
	world := ento.NewWorldBuilder().
		// Use "zero-values" as values are ignored when registering
		WithSparseComponents(Component1{}, Component2{}, Component3{}).
		// Pre-allocate space for 256 entities (world can grow beyond that automatically)
		Build(256)

	// Add systems
	system := &System{}
	world.AddSystems(system)

	// Create entities (they are added to the world immediately)
	world.NewEntity(Component1{1}, Component2{2}, Component3{3})

	// Use Set to add or update their components
	entity := world.NewEntity()
	entity.Set(Component1{0}) // Will not be handled by System

	// Update the world
	world.Update()
	system.reportCalls() // Prints: 1

	// Update the entity (add Component2)
	entity.Set(Component2{5})

	world.Update()
	system.reportCalls() // Prints: 2

	// Update the entity (remove Component2)
	entity.Rem(Component2{})

	world.Update()
	system.reportCalls() // Prints: 1
}
