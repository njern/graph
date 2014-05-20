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
	return e.id == e2.id
}

func (e *Edge) Reverse() Edge {
	return Edge{e.end, e.start, e.weight, e.id}
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

// VertexStack is a FIFO stack that holds Vertex structs.
type VertexStack struct {
	top  *Element
	size int
}

type Element struct {
	value Vertex
	next  *Element
}

// Return the stack's length
func (s *VertexStack) Len() int {
	return s.size
}

// Push a new element onto the stack
func (s *VertexStack) Push(value Vertex) {
	s.top = &Element{value, s.top}
	s.size++
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *VertexStack) Pop() Vertex {
	if s.size > 0 {
		vertex := s.top.value
		s.top = s.top.next
		s.size--
		return vertex
	}
	return Vertex{}
}

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []Vertex
	size  int
	head  int
	tail  int
	count int
}

// NewQueue returns a new queue with the given initial size.
func NewQueue(size int) *Queue {
	return &Queue{
		nodes: make([]Vertex, size),
		size:  size,
	}
}

// Push adds a node to the queue.
func (q *Queue) Push(n Vertex) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]Vertex, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() Vertex {
	if q.count == 0 {
		return Vertex{}
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node
}

func (q *Queue) Len() int {
	return q.count
}
