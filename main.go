package main

import (
	"flag"
	"log"
)

var (
	csvFilePath = flag.String("file", "", "The CSV file from which to read the input graph.")
)

func parseFlags() {
	flag.Parse()
	if *csvFilePath == "" {
		log.Fatalln("Please specify a file (with --file) from which to read the input graph!")
	}
}

func init() {
	parseFlags() // Parse flags
}

func main() {
	/*g, err := NewUndirectedGraphFromFile(*csvFilePath, '\t')
	if err != nil {
		log.Fatalf("Parsing graph with failed error '%s'\n", err)
	}

	// Test DFS
	startVertex := Vertex{id: "n34"}
	// g.depthFirstSearch(startVertex, nil, 0, make(map[Vertex]int))*/

	d, err := NewDirectedGraphFromFile(*csvFilePath, '\t')
	if err != nil {
		log.Fatalf("Parsing graph failed with error: %s\n", err)
	}

	for _, v := range d.vertices {
		d.shortestPathsFrom(v)
		//fmt.Println("\n\n")
	}
}
