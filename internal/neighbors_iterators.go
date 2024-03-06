package internal

import "github.com/zefrenchwan/nodz.git/graphs"

type NeighborsIterator[N graphs.Node, L graphs.Link[N]] struct {
	IncomingCounter   int64
	OutgoingCounter   int64
	UndirectedCounter int64
	IteratorsFactory  func() graphs.LinksIterator[N, L]
}

func (it NeighborsIterator[N, L]) IncomingDegree() int64 {
	return it.IncomingCounter
}

func (it NeighborsIterator[N, L]) OutgoingDegree() int64 {
	return it.OutgoingCounter
}

func (it NeighborsIterator[N, L]) UndirectedDegree() int64 {
	return it.UndirectedCounter
}

func (it NeighborsIterator[N, L]) Links() (graphs.LinksIterator[N, L], error) {
	return it.IteratorsFactory(), nil
}
