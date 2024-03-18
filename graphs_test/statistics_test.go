package graphs_test

import (
	"math"
	"testing"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

func TestNetworkStatsForCompleteGraphs(t *testing.T) {
	graph, errGraph := local.GenerateCompleteUndirectedGraph[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]](
		10,
		internal.NewRandomIdNode,
		internal.NewUndirectedSimpleLink,
	)

	if errGraph != nil {
		t.Fail()
	}

	counterForUndirected := func(n graphs.Neighborhood[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]) int64 {
		return n.UndirectedDegree()
	}

	stats, errStats := graphs.CalculateNetworkStatistics(&graph, counterForUndirected)
	if errStats != nil {
		t.Fail()
	}

	if stats.NodesSize != 10 {
		t.Fail()
	}

	if stats.UndirectedSize != 45 {
		t.Fail()
	}

	if len(stats.DegreeDistribution) != 1 {
		t.Fail()
	}

	if math.Abs(stats.DegreeDistribution[9]-1.0) > 0.001 {
		t.Fail()
	}
}

func TestDensities(t *testing.T) {
	stats := graphs.NetworkStatistics{
		DirectedSize:   10,
		UndirectedSize: 8,
		NodesSize:      5,
	}

	if math.Abs(stats.UndirectedDensity()-0.8) > 0.001 {
		t.Errorf("max is 10, current is 8, expected 0.8, got %f", stats.UndirectedDensity())
	}

	if math.Abs(stats.DirectedDensity()-0.5) > 0.001 {
		t.Errorf("max is 20, current is 10, expected 0.5, got %f", stats.DirectedDensity())
	}
}

func TestAverageDegrees(t *testing.T) {
	stats := graphs.NetworkStatistics{
		DirectedSize:   10,
		UndirectedSize: 8,
		NodesSize:      5,
	}

	if math.Abs(stats.AverageUndirectedDegree()-3.2) > 0.001 {
		t.Errorf("max is 10, current is 8, expected 0.8, got %f", stats.UndirectedDensity())
	}

	if math.Abs(stats.AverageDirectedDegree()-2.0) > 0.001 {
		t.Errorf("max is 20, current is 10, expected 0.5, got %f", stats.DirectedDensity())
	}
}
