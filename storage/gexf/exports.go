package gexf

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/template"

	"github.com/zefrenchwan/nodz.git/graphs"
)

// GexfNodeExporter exports a node to a label and a map of properties.
// Both label and properties are optional, just return default values
type GexfNodeExporter[N graphs.Node] func(N) (string, map[string]string)

// GexfBlankNodeExporter is a shortcut for returning default, value
func GexfBlankNodeExporter(N graphs.Node) (string, map[string]string) {
	return "", nil
}

// GexfLinkSerializer exports a link to its gexf link representation.
// Edges model in GEXF allows label, properties, types, weights.
// It would not make sense to
type GexfLinkSerializer[N graphs.Node, L graphs.Link[N]] func(source, destination int, link L) string

// GexfLinkBasicSerializer returns the xml for links with no more information
func GexfLinkBasicSerializer[N graphs.Node, L graphs.Link[N]](source, destination int, link L) string {
	if link.IsDirected() {
		return fmt.Sprintf(`<edge source="%d" target="%d" type="directed"/>`, source, destination)
	}

	return fmt.Sprintf(`<edge source="%d" target="%d" type="undirected" />`, source, destination)
}

//go:embed data_pattern.xml
var dataPattern embed.FS

// ExportDataGraph exports a graph to a gexf xml file with the "data" structure.
func ExportDataGraph[N graphs.Node, L graphs.Link[N]](
	path string, // output path
	g graphs.CentralStructureGraph[N, L], // graph to export
	nodesExporter GexfNodeExporter[N], // to export nodes to something gexf understands
	linksSerializer GexfLinkSerializer[N, L], // to serialize edges directly in GEXF format
) error {
	// load template
	var dataTemplate *template.Template
	if dp, errDp := dataPattern.ReadFile("data_pattern.xml"); errDp != nil {
		return errDp
	} else if tmpl, errTmpl := template.New("dataPattern").Parse(string(dp)); errTmpl != nil {
		return errTmpl
	} else {
		dataTemplate = tmpl
	}

	// graph iteration over nodes.
	// 1. Serialize each node
	// 2. For each node, find its neighbors. We may not know yet the indexes of the nodes
	// 3. Once all nodes index are known, then calculate edges
	it, errIt := g.AllNodes()
	if errIt != nil {
		return errIt
	}

	var globalErr error

	attributesIndex := 0
	attributeIndexMap := make(map[string]int)
	nodeValues := make([]string, 0)
	nodes := make([]N, 0)
	linkValues := make([]string, 0)
	attributeValues := make([]string, 0)

	for has, errHas := it.Next(); has; has, errHas = it.Next() {
		if errHas != nil {
			globalErr = errors.Join(globalErr, errHas)
			continue
		}

		node, errNode := it.Value()
		if errNode != nil {
			globalErr = errors.Join(globalErr, errNode)
			continue
		}

		// serialize the attributes of the node, and discover them for each node
		attributesForNode := make(map[int]string)

		label, properties := nodesExporter(node)
		// fill attributes index map
		for k, v := range properties {
			if index, found := attributeIndexMap[k]; !found {
				attributeIndexMap[k] = attributesIndex
				attributesForNode[attributesIndex] = v
				attributesIndex++
			} else {
				attributesForNode[index] = v
			}
		}

		nodeValue := serializeGexfNode(len(nodes), label, attributesForNode)
		nodeValues = append(nodeValues, nodeValue)
		nodes = append(nodes, node)
	}

	if globalErr != nil {
		return globalErr
	}

	// deal with attributes
	for value, index := range attributeIndexMap {
		attributeValues = append(attributeValues, serializeGexfAttribute(index, value))
	}

	// then, go for edges.
	it, errIt = g.AllNodes()
	if errIt != nil {
		return errIt
	}

	// For each edge, find source and destination index, and then, make the link
	for has, errHas := it.Next(); has; has, errHas = it.Next() {
		if errHas != nil {
			globalErr = errors.Join(globalErr, errHas)
			continue
		}

		node, errNode := it.Value()
		if errNode != nil {
			globalErr = errors.Join(globalErr, errNode)
			continue
		}

		// get its neighbors
		neighbors, errNeighbors := g.Neighbors(node)
		if errNeighbors != nil {
			globalErr = errors.Join(globalErr, errNeighbors)
			continue
		}

		itNeighbors, errItN := neighbors.Links()
		if errItN != nil {
			globalErr = errors.Join(globalErr, errItN)
			continue
		}

		for hasLink, errHasLink := itNeighbors.Next(); hasLink; hasLink, errHasLink = itNeighbors.Next() {
			if errHasLink != nil {
				globalErr = errors.Join(globalErr, errHasLink)
				continue
			}

			if link, errLink := itNeighbors.Value(); errLink != nil {
				globalErr = errors.Join(globalErr, errLink)
				continue
			} else {
				source := link.Source()
				destination := link.Destination()
				sourceIndex := slices.IndexFunc(nodes, func(a N) bool { return source.SameNode(a) })
				destIndex := slices.IndexFunc(nodes, func(a N) bool { return destination.SameNode(a) })

				linkValues = append(linkValues, linksSerializer(sourceIndex, destIndex, link))
			}
		}
	}

	if globalErr != nil {
		return globalErr
	}

	// generate content to write
	var localWriter bytes.Buffer
	var content gexfDataContent
	content.Attributes = strings.Join(attributeValues, "\n")
	content.Nodes = strings.Join(nodeValues, "\n")
	content.Edges = strings.Join(linkValues, "\n")

	// write it in buffer
	dataTemplate.Execute(&localWriter, content)

	// make output file
	if _, err := os.Stat(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if err == nil {
		if errRemove := os.Remove(path); errRemove != nil {
			return errRemove
		}
	}

	return os.WriteFile(path, localWriter.Bytes(), 0777)
}

// gexfDataContent goes with the pattern to form a complete gexf data file
type gexfDataContent struct {
	Attributes string // attributes are properties metadata
	Nodes      string // nodes are the nodes with their label, id and properties if any
	Edges      string // edges are the edges definition, based on the nodes index
}

// serializeGexfAttribute returns the gexf file for attributes
func serializeGexfAttribute(attrIndex int, name string) string {
	return fmt.Sprintf(`<attribute id="%d" title="%s" type="string"/>`, attrIndex, name)
}

// serializeGexfNode returns the gexf xml content for a node
func serializeGexfNode(nodeIndex int, label string, indexPropertyValues map[int]string) string {
	result := fmt.Sprintf(`<node id="%d" `, nodeIndex)
	if label != "" {
		result = result + fmt.Sprintf("label=%q", label)
	}

	if len(indexPropertyValues) == 0 {
		result = result + " />\n"
		return result
	}

	result = result + ">\n<attvalues>\n"
	for index, value := range indexPropertyValues {
		result = result + fmt.Sprintf(`<attvalue for="%d" value=%q/>`, index, value)
		result = result + "\n"
	}

	result = result + "</attvalues>\n</node>"
	return result
}
