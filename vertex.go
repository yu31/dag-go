// Copyright (c) 2020, Yu Wu <yu.771991@gmail.com> All rights reserved.
//
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package dag

import (
	"fmt"

	"github.com/yu31/structs-go/container"
)

// Vertex represents a vertex in DAG.
type Vertex struct {
	key   container.Key
	value container.Value
	out   container.Container // edges of out-degrees.
	in    container.Container // edges of in-degrees.
}

func (vex *Vertex) String() string {
	return fmt.Sprintf("%v", vex.key)
}

// Key returns the key of vertex.
func (vex *Vertex) Key() container.Key {
	return vex.key
}

// Value returns the value of vertex.
func (vex *Vertex) Value() container.Value {
	return vex.value
}
