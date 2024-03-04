# Nodz

A project to play with graphs and test some ideas about networks science. 

## Features

So far, almost nothing. But plan is:
* gephi export
* random graphs (preferential attachment, GNP and GNL)

## Implementation details 

* `graphs/` is about main definitions.  
* `internal/` is about main implementations. 
* Use of iterators: Some graphs are HUGE, so **iterators** is more flexible and efficient than slices (lazy loading / pagination)
* `internal/local/`: split definitions and local implementations. **Local implementations** means "in memory", but it is key to split general definitions and local implementations.

### Types of graph

1. **Value**BasedGraph: nodes holding content (NV type) are linked with values (LV type). Use it when you do not care about the graph structure. 
2. **Central**StructureGraph: linked nodes and you need global operations (all nodes iterations, adjacency matrix, etc)
3. **Peers**StructureGraph: linked nodes but you discover graph from a node. **No** global operations

## You like graphs or network science ? 

Anything provided here is my personal opinion.

### frameworks / show me some code

* [Neo4j database](https://neo4j.com/), community version on premise is an excellent graph database. Neo4j is, to me, **very** pushy about its cloud solution (Aura). 
* [Apache GraphX](https://spark.apache.org/graphx/): Played with it long ago, not the most active part of Spark, but something to dig  

### Tools 

* [Gephi](https://gephi.org/) is a classic and efficient vizualisation tool for graphs in general, excellent for big graphs

### Books 

* **Barabasi, Networks Science**: very interesting but it is more about ideas than a real course. I would not recommand it as a first read about network science.
