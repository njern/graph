package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type UndirectedGraph struct {
	edges    map[Vertex][]Edge // [Start]Edges
	vertices []Vertex
}

func (g *UndirectedGraph) VertexCount() int {
	return len(g.vertices)
}

func (g *UndirectedGraph) EdgeCount() int {
	return len(g.edges)
}

func (g *UndirectedGraph) String() string {
	s := "\n"
	for k, v := range g.edges {
		s += fmt.Sprintf("%s\n\t%v\n", k.id, v)
	}

	return s
}

func (g *UndirectedGraph) AddEdge(e Edge) {
	if g.edges == nil {
		// Lazily initialize
		g.edges = make(map[Vertex][]Edge)
	}

	reverseEdge := Edge{start: e.end, end: e.start, weight: e.weight, id: e.id}
	edgeFound := false
	reverseEdgeFound := false

	if edges, ok := g.edges[e.start]; ok {
		// Check that edge is not already in list
		for _, existingEdge := range edges {
			if existingEdge == e {
				edgeFound = true
			}
		}
	}

	if edges, ok := g.edges[e.end]; ok {
		// Check that reverse edge is not already in list
		for _, existingEdge := range edges {
			if existingEdge == reverseEdge {
				reverseEdgeFound = true
			}
		}
	}

	if edgeFound == false {
		// Add the edge
		g.edges[e.start] = append(g.edges[e.start], e)
	}

	if reverseEdgeFound == false {
		// Add the reverse edge
		g.edges[e.end] = append(g.edges[e.end], reverseEdge)
	}

	// Check whether the two vertices already existed. If not, add them too
	startExists := false
	endExists := false
	for _, existingVertex := range g.vertices {
		if existingVertex == e.start {
			startExists = true
		} else if existingVertex == e.end {
			endExists = true
		}
	}

	if startExists == false {
		g.vertices = append(g.vertices, e.start)
	}
	if endExists == false {
		g.vertices = append(g.vertices, e.end)
	}
}

// NewUndirectedGraphFromFile reads in a graph from the given path to a CSV.
// It expects values in the form [startNodeID, endNodeID, weight, edgeID]
// for every row. It will skip rows where the length is not 4 or the third
// value can not be parsed as an integer.
// It returns an UndirectedGraph, built from the values in the csv file
// It will throw an error if the file can not be read for some reason.
func NewUndirectedGraphFromFile(filePath string, valueSeparator rune) (*UndirectedGraph, error) {
	g := UndirectedGraph{}

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

func (g *UndirectedGraph) depthFirstSearch(v Vertex, visited Vertices, depth int, depthFirstIndex map[Vertex]int) (Vertices, int, map[Vertex]int) {
	visited = append(visited, v)
	log.Printf("Visiting vertex '%s' with DFI: %d\n", v.id, depth)

	// Note the
	depthFirstIndex[v] = depth
	depth += 1
	//log.Printf("Already visited %s\n", visited)

	if edges, ok := g.edges[v]; ok {
		// There are some "children" here to explore
		for _, edge := range edges {
			if visited.contains(edge.end) {
				// If already visited, continue
				continue
			}
			visited, depth, depthFirstIndex = g.depthFirstSearch(edge.end, visited, depth, depthFirstIndex)
		}
	}

	return visited, depth, depthFirstIndex
}
