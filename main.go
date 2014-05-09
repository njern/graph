package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
)

var (
	shortest_path = flag.String("shortest_path", "", "The CSV file from which to read the input graph.")
	prim          = flag.String("prim", "", "The CSV file from which to read the input graph for calculating Minimum Spanning Trees (exercise 3).")
	vertex_colors = flag.String("vertex_colors", "", "The CSV file from which to read the input graph for calculating minimum vertex coloring (exercise 4).")
	edge_colors   = flag.String("edge_colors", "", "The CSV file from which to read the input graph for calculating minimum edge coloring (exercise 4).")
)

func parseFlags() {
	flag.Parse()
}

func init() {
	parseFlags() // Parse flags
}

func main() {
	if *shortest_path != "" {
		d, err := NewDirectedGraphFromFile(*shortest_path, '\t')
		if err != nil {
			log.Fatalf("Parsing graph failed with error: %s\n", err)
		}

		for _, v := range d.vertices {
			d.shortestPathsFrom(v)
			//fmt.Println("\n\n")
		}
	} else if *prim != "" {
		d, err := NewUndirectedGraphFromFile(*prim, '\t')
		if err != nil {
			log.Fatalf("Parsing graph failed with error: %s\n", err)
		}

		edges := d.PrimMST(d.vertices[0])
		var edgeLabels []string
		//var totalWeight int64
		for _, edge := range edges {
			edgeLabels = append(edgeLabels, edge.id)
			//totalWeight += edge.weight
		}
		sort.Strings(edgeLabels)

		fmt.Printf("%s\n", strings.Join(edgeLabels, ","))
		// fmt.Printf("total weight: %d\n", totalWeight)
	} else if *vertex_colors != "" {
		d, err := NewUndirectedGraphFromFile(*vertex_colors, '\t')
		if err != nil {
			log.Fatalf("Parsing graph failed with error: %s\n", err)
		}

		vertexColors := d.VertexColors()
		for vertex, color := range vertexColors {
			fmt.Printf("%s: %d\n", vertex.id, color)
		}
	} else if *edge_colors != "" {
		d, err := NewUndirectedGraphFromFile(*edge_colors, '\t')
		if err != nil {
			log.Fatalf("Parsing graph failed with error: %s\n", err)
		}

		vertexColors := d.EdgeColors()
		for vertex, color := range vertexColors {
			fmt.Printf("%s: %d\n", vertex.id, color)
		}
	}
}
