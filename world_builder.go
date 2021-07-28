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
	"fmt"
	"reflect"
)

type WorldBuilder struct {
	componentIds            map[reflect.Type]int
	componentStoreProviders []StoreProvider
}

func NewWorldBuilder() *WorldBuilder {
	return &WorldBuilder{
		componentIds:            make(map[reflect.Type]int, 256),
		componentStoreProviders: make([]StoreProvider, 0, 256),
	}
}

func (b *WorldBuilder) WithComponent(component interface{}, provider StoreProvider) *WorldBuilder {
	componentType := reflect.TypeOf(component)
	if _, ok := b.componentIds[componentType]; ok {
		panic(fmt.Sprintf("component of type '%t' already registered", component))
	}

	b.componentIds[componentType] = len(b.componentIds)
	b.componentStoreProviders = append(b.componentStoreProviders, provider)

	return b
}

func (b *WorldBuilder) WithSparseComponents(components ...interface{}) *WorldBuilder {
	for _, component := range components {
		b.WithComponent(component, SparseStoreProvider(component))
	}
	return b
}

func (b *WorldBuilder) WithSingletonComponents(components ...interface{}) *WorldBuilder {
	for _, component := range components {
		b.WithComponent(component, SingletonStoreProvider(component))
	}
	return b
}

func (b *WorldBuilder) Build(capacity int) *World {
	stores := make([]Store, len(b.componentStoreProviders))
	for i, provider := range b.componentStoreProviders {
		stores[i] = provider(capacity)
	}

	world := &World{
		componentIds:      b.componentIds,
		componentStores:   stores,
		entities:          list.New(),
		freeEntityIndexes: make([]int, 0, 1024),
	}

	b.componentIds = nil
	b.componentStoreProviders = nil

	return world
}
