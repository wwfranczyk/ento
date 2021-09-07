# About
Ento is an Entity Component System written in Go. 

[![Go](https://github.com/wfranczyk/ento/actions/workflows/go.yml/badge.svg)](https://github.com/wfranczyk/ento/actions/workflows/go.yml)

# Getting Started

See [examples](./examples) folder.

From [hello.go](./examples/hello/hello.go):

```go
import (
	"github.com/wfranczyk/ento"
)

type Component1 struct{ value int }
type Component2 struct{ value int }
type Component3 struct{ value int }

type System struct {
	// Always use pointer to base component type (!)
	// Add `ento` tag to automatically bind components
	C1 *Component1 `ento:"required"`
	C2 *Component2 `ento:"required"`
	C3 *Component2 `ento:"optional"`

	// Systems can have non-component-bound fields
	calls int
}

// Update implements ento.System
// Based on `ento` tag, entities will be selected for update by the system.
// Entities that do not contain all the `required` components will be skipped.
func (s *System) Update(entity *ento.Entity) {
	// When Update is called, all the tagged components
	// will be already replaced with values from entity
	s.C1.value += s.C2.value

	// Optional field will be set to null if entity does not contain them
	if s.C3 != nil {
		s.C1.value += s.C3.value
	}
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
	entity := world.NewEntity(Component1{1}, Component2{1})

	// Use Set to add or change their components
	entity.Set(Component3{1})

	// Update the world
	world.Update()

	// Use Get to receive component value (or nil if not present)
	var c1 *Component1
	entity.Get(&c1)
	println(c1.value == 3) // true

	// Use Rem to remove component from entity
	// As Component2 is `required` in the System
	// it will no longer be updated by it
	entity.Rem(Component2{})

	world.Update()

	entity.Get(&c1)
	println(c1.value == 3) // true - the entity is no longer updated by the system
}
```

# License

Ento is distributed under MIT license.

```
MIT License

Copyright (c) 2021 Wojciech Franczyk

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
