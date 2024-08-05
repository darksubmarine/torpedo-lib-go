package app

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
	"strings"
)

// gNode is the DependencyGraph node representation
type gNode struct {
	label         string
	incomingEdges map[string]*gNode
}

// DependencyGraph graph to validate injected dependencies
type DependencyGraph struct {
	vertices map[string]*gNode
	edges    map[string]map[string]struct{}
}

// NewDependencyGraph constructor of DependencyGraph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		vertices: map[string]*gNode{},
		edges:    map[string]map[string]struct{}{}}
}

// AddVertex adds a vertex to the graph
func (g *DependencyGraph) AddVertex(v string) {
	if _, exists := g.vertices[v]; !exists {
		g.vertices[v] = &gNode{label: v, incomingEdges: map[string]*gNode{}}
	}
}

// AddEdge adds edge from v to w
func (g *DependencyGraph) AddEdge(v, w string) error {
	if g.edges[v] == nil {
		g.edges[v] = map[string]struct{}{}
	}
	// edge e from v to w
	g.edges[v][w] = struct{}{}

	// add incoming edge e in w from v
	g.vertices[w].incomingEdges[v] = g.vertices[v]

	return nil
}

// TopologicalSort by default calls a DFS based algorithm
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	return g.TopologicalSortDFS()
}

// nodesWithNoIncomingEdge used as helper function in Kahn's algorithm
func (g *DependencyGraph) nodesWithNoIncomingEdge() []string {
	toRet := []string{}
	for _, n := range g.vertices {
		if len(n.incomingEdges) == 0 {
			toRet = append(toRet, n.label)
		}
	}

	return toRet
}

// hasIncomingEdges used as helper function in Kahn's algorithm
func (g *DependencyGraph) hasIncomingEdges(n string) bool {
	return len(g.vertices[n].incomingEdges) > 0
}

// TopologicalSortKahn based on Kahn's algorithm
//
// IMPORTANT:
//
//	This implementation removes edges from original graph, so, once that has been executed the original
//	graph will be modified.
//
// TODO: clone graph before execution.
func (g *DependencyGraph) TopologicalSortKahn() ([]string, error) {
	sortedElements := []string{}                        // Empty list that will contain the sorted elements
	nodesNoIncomingEdges := g.nodesWithNoIncomingEdge() // Set of all nodes with no incoming edge

	for len(nodesNoIncomingEdges) > 0 { // while nodesNoIncomingEdges is not empty do
		n := nodesNoIncomingEdges[0] // remove a node n from nodesNoIncomingEdges
		nodesNoIncomingEdges = nodesNoIncomingEdges[1:]
		sortedElements = append(sortedElements, n) // add n to sortedElements

		for m, _ := range g.edges[n] { //for each node m with an edge e from n to m do
			delete(g.edges[n], m) // remove edge e from the graph
			delete(g.vertices[m].incomingEdges, n)
			if len(g.edges[n]) == 0 {
				delete(g.edges, n)
			}
			if !g.hasIncomingEdges(m) { // if m has no other incoming edges then
				nodesNoIncomingEdges = append(nodesNoIncomingEdges, m)
			}
		}
	}

	if len(g.edges) > 0 {
		return nil, errors.New("is cyclic")
	}

	// Reverse sorted list to match with the DFS implementation
	for i, j := 0, len(sortedElements)-1; i < j; i, j = i+1, j-1 {
		sortedElements[i], sortedElements[j] = sortedElements[j], sortedElements[i]
	}
	return sortedElements, nil
}

// TopologicalSortDFS based on Depth-first search
func (g *DependencyGraph) TopologicalSortDFS() ([]string, error) {
	sortedList := []string{}
	permanent := map[string]struct{}{}
	temporal := map[string]struct{}{}

	for len(permanent) < len(g.vertices) {
		for v, _ := range g.vertices {
			if _, exists := permanent[v]; !exists {
				if aux, err := g.visit(v, permanent, temporal); err == nil {
					sortedList = append(sortedList, aux...)
				} else {
					return nil, err
				}
			}
		}
	}

	return sortedList, nil
}

// visit is the DFS recursive function
func (g *DependencyGraph) visit(n string, permanent, temporal map[string]struct{}) ([]string, error) {
	if _, exists := permanent[n]; exists {
		return nil, nil
	}

	if _, exists := temporal[n]; exists {
		return nil, errors.New(fmt.Sprintf("graph has at least one cycle: %s",
			strings.Join(append(maps.Keys(temporal), n), "->")))
	}

	temporal[n] = struct{}{}

	var toRet = make([]string, 0)
	for m, _ := range g.edges[n] {
		if aux, err := g.visit(m, permanent, temporal); err == nil {
			toRet = append(toRet, aux...)
		} else {
			return nil, err
		}
	}

	delete(temporal, n)
	permanent[n] = struct{}{}
	return append(toRet, n), nil
}

// IsCyclic validates if the constructed graph has at least one cycle built on top of DFS topological sort.
func (g *DependencyGraph) IsCyclic() bool {
	if _, err := g.TopologicalSortDFS(); err != nil {
		return true
	}

	return false
}
