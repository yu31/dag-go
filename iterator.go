// Copyright (c) 2020, Yu Wu <yu.771991@gmail.com> All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package dag

import (
	"github.com/yu31/structs-go/queue"
)

// Iterator implements an iterator with Topological Sorting(Kahn Algorithm).
type Iterator struct {
	queue     *queue.Queue
	inDegrees map[*Vertex]int
}

// newIterator is an interval func helps creates an Iterator.
func newIterator(g *DAG) *Iterator {
	it := &Iterator{
		queue:     queue.Default(),
		inDegrees: make(map[*Vertex]int, g.Len()),
	}
	// Init the iterator.
	itVex := g.vertexes.Iter(nil, nil)
	for itVex.Valid() {
		vex := itVex.Next().Value().(*Vertex)

		degree := vex.in.Len()
		if degree == 0 {
			// Push the vertex of zero in-degree to queue.
			it.queue.Push(vex)
		} else {
			it.inDegrees[vex] = degree
		}
	}
	return it
}

// Valid represents whether has more vertex in iterator.
func (it *Iterator) Valid() bool {
	return !it.queue.Empty()
}

// Next returns a vertex that in-degree is zero.
// Returns nil if no more.
func (it *Iterator) Next() *Vertex {
	if !it.Valid() {
		return nil
	}
	vex := it.queue.Pop().(*Vertex)
	it.fillQueue(vex)
	return vex
}

// Batch returns all vertexes that in-degree is zero at once.
// Returns nil if no more.
func (it *Iterator) Batch() []*Vertex {
	if !it.Valid() {
		return nil
	}

	vexes := make([]*Vertex, 0, it.queue.Len())
	for !it.queue.Empty() {
		vex := it.queue.Pop().(*Vertex)
		vexes = append(vexes, vex)
	}

	for i := range vexes {
		it.fillQueue(vexes[i])
	}
	return vexes
}

func (it *Iterator) fillQueue(vex *Vertex) {
	itOut := vex.out.Iter(nil, nil)
	for itOut.Valid() {
		vex := itOut.Next().Value().(*Vertex)

		it.inDegrees[vex]--
		if it.inDegrees[vex] == 0 {
			delete(it.inDegrees, vex)
			// Push the vertex of zero in-degree to queue.
			it.queue.Push(vex)
		}
	}
}
