package dag

import (
	"fmt"
	"testing"

	"github.com/yu31/gostructs/container"
)

func buildDAG() *DAG {
	g := New()
	vertexes := []container.Int64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// Add vertex
	for i := 0; i < len(vertexes); i++ {
		_, _ = g.Insert(vertexes[i], int64(vertexes[i]*2+1))
	}

	// Attach edges.
	_ = g.Attach(vertexes[0], vertexes[1])

	_ = g.Attach(vertexes[1], vertexes[2])
	_ = g.Attach(vertexes[1], vertexes[3])
	_ = g.Attach(vertexes[1], vertexes[6])

	_ = g.Attach(vertexes[2], vertexes[4])
	_ = g.Attach(vertexes[2], vertexes[5])
	_ = g.Attach(vertexes[2], vertexes[6])

	_ = g.Attach(vertexes[3], vertexes[2])
	_ = g.Attach(vertexes[3], vertexes[6])
	_ = g.Attach(vertexes[3], vertexes[7])

	_ = g.Attach(vertexes[4], vertexes[8])

	_ = g.Attach(vertexes[5], vertexes[8])

	_ = g.Attach(vertexes[6], vertexes[8])

	_ = g.Attach(vertexes[7], vertexes[8])

	return g
}

func DFSRecursive(g *DAG) []*Vertex {
	result := make([]*Vertex, 0, g.Len())
	visited := make(map[*Vertex]bool, g.Len())

	var visit func(vex *Vertex)

	visit = func(vex *Vertex) {
		if !visited[vex] {
			visited[vex] = true
			result = append(result, vex)
		}

		it := vex.out.Iter(nil, nil)
		for it.Valid() {
			v1 := it.Next().Value().(*Vertex)
			if !visited[v1] {
				visit(v1)
			}
		}
	}

	it := g.vertexes.Iter(nil, nil)
	for it.Valid() {
		vex := it.Next().Value().(*Vertex)
		visit(vex)
	}

	return result
}

func TestDAG_BFS(t *testing.T) {
	g := buildDAG()

	vexes := make([]*Vertex, 0)
	g.BFS(func(vex *Vertex) bool {
		vexes = append(vexes, vex)
		return true
	})
	fmt.Println(vexes)
}

func TestDAG_DFS(t *testing.T) {
	g := buildDAG()

	vexes := make([]*Vertex, 0)
	g.DFS(func(vex *Vertex) bool {
		vexes = append(vexes, vex)
		return true
	})

	fmt.Println(vexes)
	fmt.Println(DFSRecursive(g))
}

func TestDAG_Topological(t *testing.T) {
	g := buildDAG()

	vexes := make([]*Vertex, 0)
	g.Topological(func(vex *Vertex) bool {
		vexes = append(vexes, vex)
		return true
	})
	fmt.Println(vexes)
}
