# Nodz

A project to play with graphs and test some ideas about networks science. 
Full code is under MIT License. 

## Features

So far, implementing graphs core definitions. 

### Next features (working on it)

* basic stats: degree distribution, average size, etc
* random graphs (preferential attachment, GNP)
* walkthroughs

### Features to implement one day

* gephi import and export
* neo4j import and export
* observability: observer over nodes to detect changes (node creation, deletion, or links changes. Even, for some nodes, changes of states)

And, still digging, but a must-have:
* graph data visualization 

### Features that sound like good ideas, but not sure yet

* DSL for non technical use (create a random graph, print basic stats, etc)
* Sort of jupyter interface to use that DSL
* Data viz via an embedded server

## Implementation details 

* `graphs/` is about main definitions.  
* `internal/` is about main implementations. 
* Use of iterators: Some graphs are HUGE, so  using **iterators** is more flexible and efficient than slices (lazy loading / pagination)
* `internal/local/`: split definitions and local implementations. **Local implementations** are "in memory" implementations of general definitions 

### Types of graph

1. **Value**BasedGraph: nodes holding content (NV type) are linked with values (LV type). Use it when you do not care about the graph structure. 
2. **Central**StructureGraph: linked nodes and you need global operations (all nodes iterations, adjacency matrix, etc)
3. **Peers**StructureGraph: linked nodes but you discover graph from a node. **No** global operations

So far: 

| Type | Implementation | Local | Directed |
|------|----------------|-------|-----------|
| Value | [DirectedValuesGraph](https://github.com/zefrenchwan/nodz/blob/main/internal/local/directed_value_graphs.go) | YES | YES |
| Central | [MapGraph](https://github.com/zefrenchwan/nodz/blob/main/internal/local/map_graphs.go) | YES | MIXED |

### Wait, what ? How do I start with your project ? 

1. Start with [general definition of a graph](https://github.com/zefrenchwan/nodz/blob/main/graphs/structures.go) 
2. Read `graphs/` interfaces if you need more details about `nodes`, `neighbors` or `links`. There should be no surprise, it is basic definition 
3. Have a look at `internal_test/local_test` and start with value based graphs. Go on with central graphs tests
4. Dig into implementation details if you like

## You like graphs or network science ? 

This part is about links or mentions about graphs / network science related stuff. 
It just is "hey, look at that if you want, I found it interesting". 
I don't make money by advertising, I am not in position of any conflict of interest, it is pure personal opinion. 

### frameworks / show me some code

* [Neo4j database](https://neo4j.com/), community version on premise is an excellent graph database. Neo4j is, to me, **very** pushy about its cloud solution (Aura). 
* [Apache GraphX](https://spark.apache.org/graphx/): Not a fan, but it exists and I wanted to mention it. Sounds more like an abandoned POC to me
* [NetworkX](https://networkx.org/): perfect for its purpose, easy to use, powerful. If language is not a question, I would recommand Python and NetworkX for sure ! 

### Tools 

* [Gephi](https://gephi.org/) is a classic and efficient vizualisation tool for graphs in general, excellent for big graphs

### Books 

* **Barabasi: Networks Science**: Author has a style, many ideas, not a lot of details about some key parts. Brilliant, really good to understand complex networks ideas. But... Not for a first read 
* **Mentzer, Fortunato, Davis: a first course in network science** : very good too, covers more topics and is easier, first read material for sure! 
* **Boullier: Propagations** (in french). WHAT A BOOK ! Very clever ideas about propagations and related use of data. A source of inspiration to go further than technical implementations

### Videos

* FASCINATING phenomenon : [percolation](https://www.youtube.com/watch?v=a-767WnbaCQ)

### Not about graphs, but related somehow

* **Russel, Norvig: Artificial intelligence**: THE book about AI. But get ready, it is a huge book that goes really in depth. Chapter 3 is excellent about exploration, and explains applications of graphs for walkthroughs
