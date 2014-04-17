package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type DirectedGraph struct {
	edges    map[Vertex][]Edge // [Start]Edges
	vertices []Vertex
}

func (d *DirectedGraph) VertexCount() int {
	return len(d.vertices)
}

func (d *DirectedGraph) EdgeCount() int {
	return len(d.edges)
}

func (d *DirectedGraph) String() string {
	s := "\n"
	for k, v := range d.edges {
		s += fmt.Sprintf("%s\n\t%v\n", k.id, v)
	}

	return s
}

func (d *DirectedGraph) AddEdge(e Edge) {
	if d.edges == nil {
		// Lazily initialize
		d.edges = make(map[Vertex][]Edge)
	}

	edgeFound := false

	if edges, ok := d.edges[e.start]; ok {
		// Check that edge is not already in list
		for _, existingEdge := range edges {
			if existingEdge == e {
				edgeFound = true
			}
		}
	}

	if edgeFound == false {
		// Add the edge
		d.edges[e.start] = append(d.edges[e.start], e)
	}

	// Check whether the two vertices already existed. If not, add them too
	startExists := false
	endExists := false
	for _, existingVertex := range d.vertices {
		if existingVertex == e.start {
			startExists = true
		} else if existingVertex == e.end {
			endExists = true
		}
	}

	if startExists == false {
		d.vertices = append(d.vertices, e.start)
	}
	if endExists == false {
		d.vertices = append(d.vertices, e.end)
	}
}

// NewUndirectedGraphFromFile reads in a graph from the given path to a CSV.
// It expects values in the form [startNodeID, endNodeID, weight, edgeID]
// for every row. It will skip rows where the length is not 4 or the third
// value can not be parsed as an integer.
// It returns an UndirectedGraph, built from the values in the csv file
// It will throw an error if the file can not be read for some reason.
func NewDirectedGraphFromFile(filePath string, valueSeparator rune) (*DirectedGraph, error) {
	g := DirectedGraph{}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// automatically call Close() at the end of current method
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = valueSeparator

	// options are available at:
	// http://golang.org/src/pkg/encoding/csv/reader.go?s=3213:3671#L94
	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// EOF is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// We're expecting four values
		if len(record) != 4 {
			continue
		}

		// Dig out the values, skip the row if invalid values
		start := Vertex{id: record[0]}
		end := Vertex{id: record[1]}
		weight, err := strconv.ParseInt(record[2], 0, 64)
		if err != nil {
			// log.Printf("Skipping row '%v', not a valid edge!\n", record)
			continue
		}
		id := record[3]

		// Create & add Edge
		e := Edge{start: start, end: end, weight: weight, id: id}
		g.AddEdge(e)
	}

	return &g, nil
}

func smallestDistanceVertex(dists map[Vertex]int64, Q map[Vertex]bool) Vertex {
	var smallestDistVertex Vertex
	var smallestDist int64 = -1

	for v, _ := range Q {
		dist := dists[v]
		if smallestDist == -1 || dist < smallestDist {
			smallestDist = dist
			smallestDistVertex = v
		}
	}
	return smallestDistVertex
}

func (d *DirectedGraph) shortestPathsFrom(source Vertex) {
	dists := make(map[Vertex]int64)
	Q := make(map[Vertex]bool)
	previousOptimalPathNode := make(map[Vertex]Vertex)

	for _, vertex := range d.vertices {
		Q[vertex] = true
		dists[vertex] = math.MaxInt64
	}
	dists[source] = 0 // Dist to source is 0

	for _, edge := range d.edges[source] {
		// Dist to nearest neighbours is easy
		dists[edge.end] = edge.weight
	}

	for len(Q) > 0 {
		// Find vertex u in Q with smallest distance in dists[]
		u := smallestDistanceVertex(dists, Q)
		// Remove u from Q
		delete(Q, u)
		if dists[u] == math.MaxInt64 {
			break // all remaining vertices are inaccessible from source
		}

		for _, v := range d.edges[u] {
			// where v has not yet been removed from Q
			if _, ok := Q[v.end]; ok {
				alt := dists[u] + v.weight
				if alt < dists[v.end] {
					dists[v.end] = alt
					previousOptimalPathNode[v.end] = u
				}
			}
		}
	}

	for v, dist := range dists {
		if dist == math.MaxInt64 {
			fmt.Printf("%s\t%s\tNo path!\n", source, v)
		} else {
			fmt.Printf("%s\t%s\t%d\n", source, v, dist)
		}
	}

}
