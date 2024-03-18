package examples

import (
	"fmt"
	"math"

	"github.com/zefrenchwan/nodz.git/graphs"
	"github.com/zefrenchwan/nodz.git/internal"
	"github.com/zefrenchwan/nodz.git/internal/local"
)

// GNPCritialPointAppearance uses GNP undirected models to make appear critical point and connected regime.
// For theory about it, see Barabasi, Network Science, chapter 3, and on the 2016 edition, page 84.
func GNPCriticalPointAppearance() {
	generator := local.RandomGenerator[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]{}

	const N = 50
	fmt.Printf("Critical point at p = %0.5f and expected size %0.5f\n", 1.0/(N-1), math.Pow(float64(N), 0.66))
	fmt.Printf("Connected regime when average degree is way larger than %0.5f and expected size %d \n\n", math.Log(float64(N)), N)
	fmt.Println("Starting simulation\n\n\n")

	setBuilder := func(f graphs.SetEqualsFunction[internal.IdNode]) (graphs.AbstractSet[internal.IdNode], error) {
		result := local.NewSlicesSet(f)
		return &result, nil
	}

	itBuilder := func() (graphs.DynamicIterator[internal.IdNode], error) {
		result := local.NewDynamicSlicesIterator[internal.IdNode]()
		return &result, nil
	}

	undirectedCounter := func(n graphs.Neighborhood[internal.IdNode, internal.UndirectedSimpleLink[internal.IdNode]]) int64 {
		return n.UndirectedDegree()
	}

	// details about what we print
	fmt.Println("PROBA,AVERAGE DEGREE,PREDICTED AVERAGE DEGREE, CONNECTED COMPONENTS MAX SIZE")

	for p := 0.001; p <= 1.0; p += 0.001 {
		graph, errGraph := generator.UndirectedGNP(N, p, internal.NewRandomIdNode, internal.NewUndirectedSimpleLink)
		if errGraph != nil {
			panic(errGraph)
		}

		stats, errStats := graphs.CalculateNetworkStatistics(graph, undirectedCounter)
		if errStats != nil {
			panic(errStats)
		}

		var maxSize int64
		if count, err := graphs.ConnectedComponentsSize(graph, setBuilder, itBuilder); err != nil {
			panic(err)
		} else {
			// find size of greater component
			for _, v := range count {
				if v > maxSize {
					maxSize = v
				}
			}
		}

		// print stats:
		// In order: probability, real average degree, predicted average degree, max connected component size
		fmt.Printf("%0.5f,%0.5f,%0.5f,%d\n", p, stats.AverageUndirectedDegree(), p*(N-1), maxSize)
	}
}
