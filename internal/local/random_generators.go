package local

import (
	"errors"
	"math/rand"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// RandomGenerator generates random graphs.
// It is a struct in case you want to embed your own random generator.
type RandomGenerator[N graphs.Node, L graphs.Link[N]] struct {
	// empty struct so far, using default golang random numbers generator.
}

// DirectedGNP returns a directed GNP graph
func (rm RandomGenerator[N, L]) DirectedGNP(
	size int, // number of nodes
	probability float64, // linking probability
	nodeGenerator graphs.RandomNodeGenerator[N], // generates a new node at each call
	linkGenerator graphs.RandomLinkGenerator[N, L], // generates a new link at each call
) (
	graphs.CentralStructureGraph[N, L], // random graph
	error, // error during build
) {
	return rm.gnp(size, probability, true, nodeGenerator, linkGenerator)
}

// UndirectedGNP returns a undirected GNP graph
func (rm RandomGenerator[N, L]) UndirectedGNP(
	size int, // number of nodes
	probability float64, // linking probability
	nodeGenerator graphs.RandomNodeGenerator[N], // generates a new node at each call
	linkGenerator graphs.RandomLinkGenerator[N, L], // generates a new link at each call
) (
	graphs.CentralStructureGraph[N, L], // random graph
	error, // error during build
) {
	return rm.gnp(size, probability, false, nodeGenerator, linkGenerator)
}

// gnp returns a random graph of size n and with a probability of p to create links.
// For directed, we test each couple (source, destination).
// For undirected, we test EITHER (source, destination) OR (destination, source), never both
func (rm RandomGenerator[N, L]) gnp(
	size int, // number of nodes
	probability float64, // probability to create a link
	directedAlgorithm bool, // true for directed, false for undirected
	nodeGenerator graphs.RandomNodeGenerator[N], // generates a random node
	linkGenerator graphs.RandomLinkGenerator[N, L], // generates a random link
) (
	graphs.CentralStructureGraph[N, L], // local graph with size and random links
	error, // error for any parameter that makes no sense
) {
	if size < 0 {
		return nil, errors.New("invalid size")
	}

	if probability > 1.0 || probability < 0.0 {
		return nil, errors.New("invalid probability")
	}

	var result MapGraph[N, L]
	if size == 0 {
		return &result, nil
	}

	result = NewMapGraph[N, L]()

	nodes := make([]N, size)
	// generate nodes
	for index := 0; index < size; index++ {
		newNode := nodeGenerator()
		result.AddNode(newNode)
		nodes[index] = newNode
	}

	for i, source := range nodes {
		for j, dest := range nodes {
			// exclude links having source == destination
			if i == j || (!directedAlgorithm && i > j) {
				continue
			} else if rm.nextFloat() > probability {
				continue
			}

			link := linkGenerator(source, dest)
			if link.IsDirected() != directedAlgorithm {
				return &result, errors.New("inconsistent link type")
			} else if !result.HasLink(link) {
				result.AddLink(link)
			}
		}
	}

	return &result, nil
}

// UndirectedBarabasiAlbertGraph returns an undirected graph with maxSize nodes, using the preferential attachment.
// Algorithm is:
// * to generate a first complete graph with initialSize nodes
// * for each node, keep in mind its degree d(n) and the sum of all degrees so far (D)
// * add nodes (until maxSize is reached) from a new node to an existing one n with probability p = d(n) / D
func (rm RandomGenerator[N, L]) UndirectedBarabasiAlbertGraph(
	initialSize int, // initial number of nodes for complete base
	maxSize int, // total number of nodes
	nodeGenerator graphs.RandomNodeGenerator[N], // generates random nodes
	linkGenerator graphs.RandomLinkGenerator[N, L], // generate random undirected links
) (graphs.CentralStructureGraph[N, L], // result
	error, // error if parameters make no sense or linkGenerator makes directed links
) {
	if initialSize <= 0 || maxSize <= 0 || initialSize > maxSize {
		return nil, errors.New("invalid size")
	}

	var result MapGraph[N, L]
	// initial value, a complete graph of size initialSize
	result, errInit := GenerateCompleteUndirectedGraph[N, L](initialSize, nodeGenerator, linkGenerator)
	if errInit != nil {
		return &result, errInit
	}

	// From this point to the end of the function, there is an extensive use of technical details within map graphs.
	// Two reasons:
	// * no need to parse the whole graph again and again, we just need initial state of degrees and maintain degrees map
	// * using low level functions saves time by not going through maps to find indexes of the nodes again and again
	// There is a risk here for encapsulation (having to break it all for a change in mapgraphs).
	//
	// Preferential attachment means adding with a probability equals to current node degree / sum of degrees.
	// So we want to keep in mind sumDegrees (the sum of all degrees) and the map that contains each node degree
	var sumDegrees int64 = 0
	// degrees contains the degree of all the nodes within the graph to avoid recalculation.
	// Its keys are the index within map graph (that is result.content), values are undirected degree
	degrees := make(map[int]int64)
	// init degrees (node degree) and sumDegrees (sum of degrees)
	for nodeIndex := range result.nodes.toIncreasingIndexes() {
		currentDegree := result.content[nodeIndex].undirectedCounter
		degrees[nodeIndex] = currentDegree
		sumDegrees += currentDegree
	}

	// initial graph is complete, add links one by one
	for index := initialSize; index < maxSize; index++ {
		newNode := nodeGenerator()
		// not necessary, but doing it allows to get its index without a full rescan
		newNodeIndex := result.nodes.addValue(newNode)

		// find destination using preferential attachement
		var destNode N
		var destIndex int
		// To do so, generate a value, sum until the random value is reached.
		// When reached, the corresponding node is the destination node.
		// randomValue is between 0 included and sumDegrees excluded
		randomValue := rm.nextInt64(sumDegrees - 1)
		var sum int64
		for nodeIndex, degree := range degrees {
			sum += degree
			if randomValue < sum {
				destIndex = nodeIndex
				destNode = result.nodes.values[nodeIndex]
				break
			}
		}

		// add new link within the graph
		newLink := linkGenerator(newNode, destNode)
		result.AddLink(newLink)

		// ensure invariants
		sumDegrees += 2
		degrees[newNodeIndex] = 1
		degrees[destIndex] = degrees[destIndex] + 1
	}

	return &result, nil
}

// nextFloat returns a random float between 0.0 and 1.0.
// This is Golang default implementation.
func (rm RandomGenerator[N, L]) nextFloat() float64 {
	return rand.Float64()
}

// nextInt64 returns a new random positive int64 from 0 to max included
func (rm RandomGenerator[N, L]) nextInt64(max int64) int64 {
	return rand.Int63() % (max + 1)
}

// generateDistinctValues returns a slice of size size, with values from 0 to max (included), all different.
// Of course, it makes no sens if size > max, in case we return nil.
// About the algorithm, it is not that simple : worst idea is to make random values until we have size different ones.
// To generate different values in a "reasonable" time, solution came from the internet :
// https://stackoverflow.com/questions/3722430/most-efficient-way-of-randomly-choosing-a-set-of-distinct-integers/
func (rm RandomGenerator[N, L]) generateDistinctValues(size, max int) []int {
	if size > max {
		return nil
	}

	values := make(map[int]bool)
	count := max + 1
	for i := count - size; i < count; i++ {
		newValue := rand.Intn(i + 1)
		if values[newValue] {
			values[i] = true
		} else {
			values[newValue] = true
		}
	}

	result := make([]int, size)
	index := 0
	for k := range values {
		result[index] = k
		index++
	}

	return result
}
