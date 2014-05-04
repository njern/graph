package main

type Graph interface {
	VertexCount() int
	EdgeCount() int
	String() string
}

type Vertices []Vertex

func (v *Vertices) contains(vertex Vertex) bool {
	for _, existingVertex := range *v {
		if existingVertex == vertex {
			return true
		}
	}
	return false
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
