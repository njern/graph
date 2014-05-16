package main

type Graph interface {
	VertexCount() int
	EdgeCount() int
	String() string
}

type Vertex struct {
	id string
}

func (v *Vertex) Equals(v2 *Vertex) bool {
	return v.id == v2.id
}

type Edge struct {
	start  Vertex
	end    Vertex
	weight int64
	id     string
}

func (e *Edge) Equals(e2 *Edge) bool {
	return e.id == e2.id && e.start.Equals(e2.start) && e.end.Equals(e2.end) && e.weight == e2.weight
}

type Vertices []Vertex
type Edges []Edge

func (v *Vertices) contains(vertex Vertex) bool {
	for _, existingVertex := range *v {
		if existingVertex == vertex {
			return true
		}
	}
	return false
}

func (e *Edges) contains(edge Edge) bool {
	for _, existingEdge := range *e {
		if existingEdge == edge {
			return true
		}
	}
	return false
}
