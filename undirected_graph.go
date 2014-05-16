package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type UndirectedGraph struct {
	edges    map[Vertex][]Edge // [Start]Edges
	edgeList []Edge
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

// vertexNeighbours returns a list of v's
// neighbour vertices to which it is directly
// connected with an edge.
func (g *UndirectedGraph) vertexNeighbours(v Vertex) Vertices {
	var vertices Vertices
	for _, edge := range g.edges[v] {
		vertices = append(vertices, edge.end)
	}

	return vertices
}

// edgeNeighbours returns all the edges neighbouring the given edge.
// It will only return one "version" of each Edge, not both the
// forwards and backwards version.
func (g *UndirectedGraph) edgeNeighbours(e Edge) []Edge {
	edges := make(map[Edge]bool)

	for _, edge := range g.edgeList {
		if e.start == edge.start || e.start == edge.end ||
			e.end == edge.start || e.end == edge.end {
			edges[edge] = true

		}
	}

	var result []Edge
	for edge := range edges {
		result = append(result, edge)
	}

	return result
}

func (g *UndirectedGraph) AddEdge(e Edge) {
	if g.edges == nil {
		// Lazily initialize
		g.edges = make(map[Vertex][]Edge)
	}

	g.edgeList = append(g.edgeList, e)

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

		// Some graphs leave out the weight
		if len(record) == 3 {
			// Dig out the values, skip the row if invalid values
			start := Vertex{id: record[0]}
			end := Vertex{id: record[1]}
			id := record[2]
			// Create & add Edge
			e := Edge{start: start, end: end, weight: -1, id: id}
			g.AddEdge(e)

			continue
		}

		// We're expecting three or four values
		if len(record) != 3 && len(record) != 4 {
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

// PrimMST implements Prim's algorithm. Shamelessly implemented
// as per it's Wikipedia description: http://en.wikipedia.org/wiki/Prim's_algorithm
func (g *UndirectedGraph) PrimMST(start Vertex) []Edge {
	vNew := Vertices{start}
	var eNew []Edge

	var count = 0

	for len(vNew) != len(g.vertices) {
		// Dirty hack :) - Break after we have run more than O(n^2) times - it's a disconnected graph
		count++
		if count > (len(g.vertices) * len(g.vertices) * 2) {
			break
		}

		var minWeightCandidate int64 = math.MaxInt64
		var vertexCandidate Vertex
		var edgeCandidate Edge

		for _, v := range vNew {
			for _, edge := range g.edges[v] {
				if vNew.contains(edge.end) == false {
					if edge.weight < minWeightCandidate {
						vertexCandidate = edge.end
						edgeCandidate = edge
						minWeightCandidate = edge.weight
					}
				}
			}
		}
		if vertexCandidate.id != "" {
			vNew = append(vNew, vertexCandidate)
			eNew = append(eNew, edgeCandidate)
		}
	}

	return eNew
}

// VertexColors iteratively applies a greedy algorithm of
// "smallest possible color" to find the minimum
// vertex coloring for a given undirected graph.
// It will loop over the list of vertices and set the color
// to the smallest integer not used by one of the vertex's
// neighbours until no more optimisations can be made.
func (g *UndirectedGraph) VertexColors() map[Vertex]int {
	// Track vertex colors & start off at math.Maxint32
	vertexColors := make(map[Vertex]int)
	for _, v := range g.vertices {
		vertexColors[v] = math.MaxInt32
	}

	for {
		graphChangedDuringCurrentPass := false
		for vertex := range vertexColors {
			// Get list of neighbour vertices
			neighbours := g.vertexNeighbours(vertex)

			// Find the smallest color(integer)
			// not already assigned to one of
			// the vertex's neighbours.
			color := 0
			for {
				shouldContinue := false
				for _, neighbour := range neighbours {
					if vertexColors[neighbour] == color {
						// If one of the neighbours are already
						// painted with this color, continue
						// on to the next.
						shouldContinue = true
					}
				}
				if shouldContinue == false {
					break
				}
				// Else, we try with the next color.
				color++
			}
			// If the old color was larger than the new
			// color, replace it.
			if vertexColors[vertex] > color {
				graphChangedDuringCurrentPass = true
				vertexColors[vertex] = color
			}
		}

		// Once we can no longer improve the graph...
		if graphChangedDuringCurrentPass == false {
			break
		}
	}

	return vertexColors
}

// EdgeColors iteratively applies a greedy algorithm of
// "smallest possible color" to find the minimum edge
// coloring for a given undirected graph in a way very
// similar to VertexColors().
// It will loop over the list of edges and set the color
// to the smallest integer not used by one of the edge's
// neighbours until no more optimisations can be made.
func (g *UndirectedGraph) EdgeColors() map[Edge]int {
	// Track edge colors & start off at math.Maxint32
	edgeColors := make(map[Edge]int)
	for _, e := range g.edgeList {
		edgeColors[e] = math.MaxInt32
	}

	for {
		graphChangedDuringCurrentPass := false
		for edge := range edgeColors {
			// Get list of neighbour vertices
			neighbours := g.edgeNeighbours(edge)

			// Find the smallest color(integer)
			// not already assigned to one of
			// the edge's neighbours.
			color := 0
			for {
				shouldContinue := false
				for _, neighbour := range neighbours {
					if edgeColors[neighbour] == color {
						// If one of the neighbours are already
						// painted with this color, continue
						// on to the next.
						shouldContinue = true
					}
				}
				if shouldContinue == false {
					break
				}
				// Else, we try with the next color.
				color++
			}
			// If the old color was larger than the new
			// color, replace it.
			if edgeColors[edge] > color {
				graphChangedDuringCurrentPass = true
				edgeColors[edge] = color
			}
		}

		// Once we can no longer improve the graph...
		if graphChangedDuringCurrentPass == false {
			break
		}
	}

	return edgeColors
}

// maxCardMatching finds a maximum-cardinality edge
// matching in a connected undirected graph.
func (g *UndirectedGraph) maxCardMatching([]Edge) {
	for {
		// Generate maximal matchings until we find a maximum matching.
	}
}
