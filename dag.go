// Copyright (c) 2020, Yu Wu <yu.771991@gmail.com> All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package dag

import (
	"github.com/yu31/structs-go/container"
	"github.com/yu31/structs-go/queue"
	"github.com/yu31/structs-go/rb"
	"github.com/yu31/structs-go/stack"
)

// DAG implements data struct of Directed Acyclic Graph.
type DAG struct {
	vertexes container.Container
}

// New creates an DAG.
func New() *DAG {
	g := new(DAG)
	g.vertexes = rb.New()
	return g
}

// Len returns number of vertexes in DAG.
func (g *DAG) Len() int {
	return g.vertexes.Len()
}

// Insert inserts and returns a Vertex with given key and value if key doesn't exist.
// Or else, returns the existing Vertex for the key if present.
// The bool result is true if a Vertex was created, false if searched.
func (g *DAG) Insert(k container.Key, v container.Value) (*Vertex, bool) {
	vex := &Vertex{key: k, value: v}
	if ele, ok := g.vertexes.Insert(k, vex); !ok {
		return ele.Value().(*Vertex), false
	}
	vex.in = rb.New()
	vex.out = rb.New()
	return vex, true
}

// Delete deletes and returns the vertex by giving key.
// Returns nil if vertex not exists.
func (g *DAG) Delete(k container.Key) *Vertex {
	// The vertex not exists.
	ele := g.vertexes.Delete(k)
	if ele == nil {
		return nil
	}

	vex := ele.Value().(*Vertex)
	vex.in = nil
	vex.out = nil

	// Delete edges form other vertices that attach to this vertex.
	it := g.vertexes.Iter(nil, nil)
	for it.Valid() {
		vt := it.Next().Value().(*Vertex)
		_ = vt.in.Delete(k)
		_ = vt.out.Delete(k)
	}

	return vex
}

// Search get the vertex of a given key.
func (g *DAG) Search(k container.Key) *Vertex {
	ele := g.vertexes.Search(k)
	if ele == nil {
		return nil
	}
	return ele.Value().(*Vertex)
}

// Attach attaches an edge from srcKey to destKey.
// Returns false if their a ring between srcKey and destKey after attaching.
//
// And will be crashing in following cases:
//   - srcKey equal to destKey.
//   - srcKey or destKey does not exist.
func (g *DAG) Attach(srcKey, destKey container.Key) bool {
	if srcKey.Compare(destKey) == 0 {
		panic("dag: srcKey is same as the destKey")
	}

	e1 := g.vertexes.Search(srcKey)
	if e1 == nil {
		panic("dag: vertex of srcKey not exists")
	}
	e2 := g.vertexes.Search(destKey)
	if e2 == nil {
		panic("dag: vertex of destKey not exists")
	}

	srcVex := e1.Value().(*Vertex)
	destVex := e2.Value().(*Vertex)

	// Check whether there is a ring after attaching.
	s := stack.Default()
	s.Push(destVex)
	for !s.Empty() {
		vex := s.Pop().(*Vertex)
		if vex == srcVex {
			// Has ring, returns false.
			return false
		}
		it := vex.out.Iter(nil, nil)
		for it.Valid() {
			s.Push(it.Next().Value())
		}
	}

	// Attaching edges
	_, _ = srcVex.out.Insert(destKey, destVex)
	_, _ = destVex.in.Insert(srcKey, srcVex)
	return true
}

// Detach detaches edges from srcKey to destKey.
// Returns false if no edges between srcKey and destKey.
//
// And will be crashing in following cases:
//   - srcKey equal to destKey.
//   - srcKey or destKey does not exist.
func (g *DAG) Detach(srcKey, destKey container.Key) bool {
	if srcKey.Compare(destKey) == 0 {
		panic("dag: srcKey is same as the destKey")
	}

	e1 := g.vertexes.Search(srcKey)
	if e1 == nil {
		panic("dag: vertex of srcKey not exists")
	}
	e2 := g.vertexes.Search(destKey)
	if e2 == nil {
		panic("dag: vertex of destKey not exists")
	}

	if ele := e1.Value().(*Vertex).out.Delete(destKey); ele == nil {
		return false
	}
	if ele := e2.Value().(*Vertex).in.Delete(srcKey); ele == nil {
		return false
	}
	return true
}

// BFS iteration.
// If f returns false, range stops the iteration.
func (g *DAG) BFS(f func(vex *Vertex) bool) {
	q := queue.Default()
	visited := make(map[*Vertex]bool, g.Len())

	it := g.vertexes.Iter(nil, nil)
	for it.Valid() {
		v1 := it.Next().Value().(*Vertex)
		if !visited[v1] {
			q.Push(v1)
			visited[v1] = true
		}

		for !q.Empty() {
			vex := q.Pop().(*Vertex)
			// Adjacency
			adj := vex.out.Iter(nil, nil)
			for adj.Valid() {
				v0 := adj.Next().Value().(*Vertex)
				if !visited[v0] {
					q.Push(v0)
					visited[v0] = true
				}
			}

			// callback.
			if !f(vex) {
				return
			}
		}
	}
}

// DFS iteration.
// If f returns false, range stops the iteration.
func (g *DAG) DFS(f func(vex *Vertex) bool) {
	s := stack.Default()
	visited := make(map[*Vertex]bool, g.Len())

	it := g.vertexes.Iter(nil, nil)
	for it.Valid() {
		s.Push(it.Next().Value().(*Vertex))

		for !s.Empty() {
			vex := s.Pop().(*Vertex)
			adj := vex.out.IterReverse(nil, nil)
			for adj.Valid() {
				v0 := adj.Next().Value().(*Vertex)
				if !visited[v0] {
					s.Push(v0)
				}
			}

			if !visited[vex] {
				visited[vex] = true
				// callback.
				if !f(vex) {
					return
				}
			}
		}
	}
}

// Topological iteration.
// If f returns false, range stops the iteration.
func (g *DAG) Topological(f func(vex *Vertex) bool) {
	q := queue.Default()
	inDegrees := make(map[*Vertex]int, g.Len())

	it := g.vertexes.Iter(nil, nil)
	for it.Valid() {
		vex := it.Next().Value().(*Vertex)
		degrees := vex.in.Len()
		if degrees == 0 {
			q.Push(vex)
		} else {
			inDegrees[vex] = degrees
		}
	}

	for !q.Empty() {
		vex := q.Pop().(*Vertex)
		it = vex.out.Iter(nil, nil)
		for it.Valid() {
			v1 := it.Next().Value().(*Vertex)
			inDegrees[v1]--
			if inDegrees[v1] == 0 {
				delete(inDegrees, v1)
				q.Push(v1)
			}
		}

		// Callback
		if !f(vex) {
			return
		}
	}
}

// Iter return an Iterator.
func (g *DAG) Iter() *Iterator {
	return newIterator(g)
}
