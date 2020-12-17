package dag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yu31/gostructs/container"
)

// InDegreeToString return strings by in-degree.
func inDegreeToString(g *DAG) string {
	s := "{ "
	it1 := g.vertexes.Iter(nil, nil)
	for it1.Valid() {
		e1 := it1.Next()
		v1 := e1.Value().(*Vertex)
		s += fmt.Sprintf("%v: [ ", e1.Key())

		it2 := v1.in.Iter(nil, nil)
		for it2.Valid() {
			e2 := it2.Next()
			v2 := e2.Value().(*Vertex)
			_ = v2
			s += fmt.Sprintf("%v ", e2.Key())
		}
		s += "] "
	}
	s += "}"
	return s
}

// OutDegreeToString return strings by out-degree.
func outDegreeToString(g *DAG) string {
	s := "{ "
	it1 := g.vertexes.Iter(nil, nil)
	for it1.Valid() {
		e1 := it1.Next()
		v1 := e1.Value().(*Vertex)
		s += fmt.Sprintf("%v: [ ", e1.Key())

		it2 := v1.out.Iter(nil, nil)
		for it2.Valid() {
			e2 := it2.Next()
			v2 := e2.Value().(*Vertex)
			_ = v2
			s += fmt.Sprintf("%v ", e2.Key())
		}
		s += "] "
	}
	s += "}"
	return s
}

func TestNew(t *testing.T) {
	g := New()
	require.NotNil(t, g)
	require.NotNil(t, g.vertexes)
}

func TestDAG_Insert(t *testing.T) {
	g := New()
	vertexes := []container.Int64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// Add vertex
	for i := 0; i < len(vertexes); i++ {
		vex, ok := g.Insert(vertexes[i], int64(vertexes[i]*2+1))
		require.True(t, ok)
		require.Equal(t, vex.Key(), vertexes[i])
		require.Equal(t, vex.Value(), int64(vertexes[i]*2+1))

		vex, ok = g.Insert(vertexes[i], int64(vertexes[i]*2+1)*2)
		require.False(t, ok)
		require.Equal(t, vex.Key(), vertexes[i])
		require.Equal(t, vex.Value(), int64(vertexes[i]*2+1))
	}

	// test internal element
	it1 := g.vertexes.Iter(nil, nil)
	for it1.Valid() {
		element := it1.Next()
		require.NotNil(t, element)
		require.NotNil(t, element.Key())
		require.NotNil(t, element.Value())
		n, ok := element.Value().(*Vertex)
		require.True(t, ok)
		require.NotNil(t, n.value)
		require.NotNil(t, n.in)
		require.NotNil(t, n.out)
	}
}

func TestDAG_Attach(t *testing.T) {
	g := New()
	vertexes := []container.Int64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// Test Add vertex
	for i := 0; i < len(vertexes); i++ {
		g.Insert(vertexes[i], int64(vertexes[i]*2+1))
	}

	// Test Add edge positive
	require.True(t, g.Attach(vertexes[0], vertexes[1]))

	require.True(t, g.Attach(vertexes[1], vertexes[2]))
	require.True(t, g.Attach(vertexes[1], vertexes[3]))
	require.True(t, g.Attach(vertexes[1], vertexes[6]))

	require.True(t, g.Attach(vertexes[2], vertexes[4]))
	require.True(t, g.Attach(vertexes[2], vertexes[5]))
	require.True(t, g.Attach(vertexes[2], vertexes[6]))

	require.True(t, g.Attach(vertexes[3], vertexes[2]))
	require.True(t, g.Attach(vertexes[3], vertexes[6]))
	require.True(t, g.Attach(vertexes[3], vertexes[7]))

	require.True(t, g.Attach(vertexes[4], vertexes[8]))

	require.True(t, g.Attach(vertexes[5], vertexes[8]))

	require.True(t, g.Attach(vertexes[6], vertexes[8]))

	require.True(t, g.Attach(vertexes[7], vertexes[8]))

	fmt.Println("outDegreeToString:", outDegreeToString(g))
	fmt.Println("inDegreeToString:", inDegreeToString(g))

	// Test Add edge negative
	require.Panics(t, func() {
		g.Attach(vertexes[1], vertexes[1])
	})
	//require.False(t, g.Attach(vertexes[1], vertexes[1]))
	require.False(t, g.Attach(vertexes[6], vertexes[1]))
	require.False(t, g.Attach(vertexes[2], vertexes[1]))
	require.False(t, g.Attach(vertexes[8], vertexes[3]))
}

func TestDAG_Delete(t *testing.T) {
	g := New()
	vertexes := []container.Int64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// Test Add vertex
	for i := 0; i < len(vertexes); i++ {
		g.Insert(vertexes[i], int64(vertexes[i]*2+1))
	}

	// Test Add edge positive
	require.True(t, g.Attach(vertexes[0], vertexes[1]))

	require.True(t, g.Attach(vertexes[1], vertexes[2]))
	require.True(t, g.Attach(vertexes[1], vertexes[3]))
	require.True(t, g.Attach(vertexes[1], vertexes[6]))

	require.True(t, g.Attach(vertexes[2], vertexes[4]))
	require.True(t, g.Attach(vertexes[2], vertexes[5]))
	require.True(t, g.Attach(vertexes[2], vertexes[6]))

	require.True(t, g.Attach(vertexes[3], vertexes[2]))
	require.True(t, g.Attach(vertexes[3], vertexes[6]))
	require.True(t, g.Attach(vertexes[3], vertexes[7]))

	require.True(t, g.Attach(vertexes[4], vertexes[8]))

	require.True(t, g.Attach(vertexes[5], vertexes[8]))

	require.True(t, g.Attach(vertexes[6], vertexes[8]))

	require.True(t, g.Attach(vertexes[7], vertexes[8]))

	// Test Delete vertex
	require.NotNil(t, g.Delete(vertexes[3]))
	require.Nil(t, g.Delete(vertexes[3]))

	it1 := g.vertexes.Iter(nil, nil)
	for it1.Valid() {
		e1 := it1.Next()
		n1 := e1.Value().(*Vertex)
		require.NotEqual(t, e1.Key(), vertexes[3])

		it2 := n1.in.Iter(nil, nil)
		for it2.Valid() {
			e2 := it2.Next()
			require.NotEqual(t, e2.Key(), vertexes[3])
		}
		it3 := n1.out.Iter(nil, nil)
		for it3.Valid() {
			e3 := it3.Next()
			require.NotEqual(t, e3.Key(), vertexes[3])
		}
	}
}

func TestDAG_Detach(t *testing.T) {
	g := New()
	vertexes := []container.Int64{0, 1, 2, 3, 4, 5, 6, 7, 8}

	// Test Add vertex
	for i := 0; i < len(vertexes); i++ {
		g.Insert(vertexes[i], int64(vertexes[i]*2+1))
	}

	// Test Add edge positive
	require.True(t, g.Attach(vertexes[0], vertexes[1]))

	require.True(t, g.Attach(vertexes[1], vertexes[2]))
	require.True(t, g.Attach(vertexes[1], vertexes[3]))
	require.True(t, g.Attach(vertexes[1], vertexes[6]))

	require.True(t, g.Attach(vertexes[2], vertexes[4]))
	require.True(t, g.Attach(vertexes[2], vertexes[5]))
	require.True(t, g.Attach(vertexes[2], vertexes[6]))

	require.True(t, g.Attach(vertexes[3], vertexes[2]))
	require.True(t, g.Attach(vertexes[3], vertexes[6]))
	require.True(t, g.Attach(vertexes[3], vertexes[7]))

	require.True(t, g.Attach(vertexes[4], vertexes[8]))

	require.True(t, g.Attach(vertexes[5], vertexes[8]))

	require.True(t, g.Attach(vertexes[6], vertexes[8]))

	require.True(t, g.Attach(vertexes[7], vertexes[8]))

	// Test Delete edge
	require.True(t, g.Detach(vertexes[7], vertexes[8]))
	require.False(t, g.Detach(vertexes[7], vertexes[8]))
}
