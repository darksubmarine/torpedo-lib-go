package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDependencyGraph_TopologicalSortDFS(t *testing.T) {
	g := NewDependencyGraph()
	g.AddVertex("CONTROLLER")
	g.AddVertex("LOGGER")
	g.AddVertex("HTTP")

	g.AddEdge("CONTROLLER", "HTTP")
	g.AddEdge("HTTP", "LOGGER")
	g.AddEdge("CONTROLLER", "LOGGER")

	sortList, err := g.TopologicalSortDFS()
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"LOGGER", "HTTP", "CONTROLLER"}, sortList)

	assert.False(t, g.IsCyclic())
}

func TestDependencyGraph_TopologicalSortKahn(t *testing.T) {
	g := NewDependencyGraph()
	g.AddVertex("CONTROLLER")
	g.AddVertex("LOGGER")
	g.AddVertex("HTTP")

	g.AddEdge("CONTROLLER", "HTTP")
	g.AddEdge("HTTP", "LOGGER")
	g.AddEdge("CONTROLLER", "LOGGER")

	sortList, err := g.TopologicalSortKahn()
	assert.Nil(t, err)
	assert.EqualValues(t, []string{"LOGGER", "HTTP", "CONTROLLER"}, sortList)

	assert.False(t, g.IsCyclic())
}

func TestDependencyGraph_TopologicalSort_DFS_Kahn(t *testing.T) {
	g := NewDependencyGraph()
	g.AddVertex("CONTROLLER")
	g.AddVertex("LOGGER")
	g.AddVertex("HTTP")

	g.AddEdge("CONTROLLER", "HTTP")
	g.AddEdge("HTTP", "LOGGER")
	g.AddEdge("CONTROLLER", "LOGGER")

	// first run DFS because Kahn run over original graph modifying it.
	sortListDFS, errDFS := g.TopologicalSortDFS()
	assert.Nil(t, errDFS)

	sortListKahn, errKahn := g.TopologicalSortKahn()
	assert.Nil(t, errKahn)

	assert.EqualValues(t, sortListKahn, sortListDFS)
}

func TestDependencyGraph_IsCyclic(t *testing.T) {
	g := NewDependencyGraph()
	g.AddVertex("A")
	g.AddVertex("B")
	g.AddVertex("C")
	g.AddVertex("D")

	g.AddEdge("D", "C")
	g.AddEdge("C", "B")
	g.AddEdge("C", "A")
	g.AddEdge("B", "A")
	g.AddEdge("A", "D")

	assert.True(t, g.IsCyclic())
}
