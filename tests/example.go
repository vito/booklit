package tests

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vito/booklit"
	"github.com/vito/booklit/baselit"
	"github.com/vito/booklit/dangeval"
	"github.com/vito/booklit/load"
	"github.com/vito/booklit/render"
)

type Example struct {
	Input       string
	Inputs      Files
	Outputs     Files
	SearchIndex string
	LoadErr     string
	RenderErr   string
	Ext         string // file extension, defaults to ".md"
}

type Files map[string]string

func (example Example) Run(t *testing.T) {
	t.Helper()

	dir := t.TempDir()

	dang, err := dangeval.New(context.Background(), dir)
	require.NoError(t, err)
	t.Cleanup(dang.Close)

	processor := &load.Processor{
		Dang: dang,
	}

	pluginFactories := []booklit.PluginFactory{
		baselit.NewPlugin,
	}

	ext := example.Ext
	if ext == "" {
		ext = ".md"
	}

	// Use only the leaf subtest name for the filename, matching the old
	// ginkgo behavior of CurrentSpecReport().LeafNodeText.
	name := t.Name()
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	sectionPath := filepath.Join(dir, name+ext)

	err = os.WriteFile(sectionPath, []byte(example.Input), 0644)
	require.NoError(t, err)

	for file, contents := range example.Inputs {
		err := os.MkdirAll(filepath.Join(dir, filepath.Dir(file)), 0755)
		require.NoError(t, err)

		err = os.WriteFile(filepath.Join(dir, file), []byte(contents), 0644)
		require.NoError(t, err)
	}

	section, err := processor.LoadFile(sectionPath, pluginFactories)
	if example.LoadErr != "" {
		require.Error(t, err)
		assert.ErrorContains(t, err, example.LoadErr)
		return
	}

	require.NoError(t, err)

	engine := render.NewHTMLEngine()

	err = engine.LoadTemplates("fixtures")
	require.NoError(t, err)

	writer := render.Writer{
		Engine:      engine,
		Destination: dir,
	}

	err = writer.WriteSection(section)
	if example.RenderErr != "" {
		require.Error(t, err)
		assert.ErrorContains(t, err, example.RenderErr)
		return
	}

	require.NoError(t, err)

	for file, contents := range example.Outputs {
		fileContents, err := os.ReadFile(filepath.Join(dir, file))
		require.NoError(t, err)
		assertXMLEqual(t, contents, string(fileContents), "file %s", file)
	}

	if example.SearchIndex != "" {
		err := writer.WriteSearchIndex(section, "search_index.json")
		require.NoError(t, err)

		fileContents, err := os.ReadFile(filepath.Join(dir, "search_index.json"))
		require.NoError(t, err)
		assert.JSONEq(t, example.SearchIndex, string(fileContents))
	}

	assert.NotEmpty(t, stringifyEverything(section))
}

// assertXMLEqual compares two XML/XHTML strings structurally, replicating
// gomega's MatchXML behavior: the first element acts as a pseudo-root,
// text before it is discarded, whitespace in parent-node content is trimmed,
// and attributes are sorted.
func assertXMLEqual(t *testing.T, expected, actual string, msgAndArgs ...any) {
	t.Helper()
	expectedNode, err := parseXMLContent(expected)
	require.NoError(t, err, "failed to parse expected XML")
	actualNode, err := parseXMLContent(actual)
	require.NoError(t, err, "failed to parse actual XML")
	assert.Equal(t, expectedNode, actualNode, msgAndArgs...)
}

// xmlNode mirrors gomega's internal XML representation.
type xmlNode struct {
	XMLName   xml.Name
	XMLAttr   []xml.Attr
	Content   []byte
	Nodes     []*xmlNode
}

// parseXMLContent replicates gomega's parseXmlContent: the first element
// becomes the root node and all subsequent sibling elements are appended
// as its children. Text before the first element is discarded. For nodes
// with children, Content is trimmed of leading/trailing whitespace.
func parseXMLContent(content string) (*xmlNode, error) {
	allNodes := []*xmlNode{}

	dec := xml.NewDecoder(strings.NewReader(content))
	dec.Strict = false
	dec.AutoClose = xml.HTMLAutoClose
	dec.Entity = xml.HTMLEntity

	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		lastNodeIndex := len(allNodes) - 1
		var lastNode *xmlNode
		if len(allNodes) > 0 {
			lastNode = allNodes[lastNodeIndex]
		} else {
			lastNode = &xmlNode{}
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			attrs := make([]xml.Attr, len(tok.Attr))
			copy(attrs, tok.Attr)
			sort.Slice(attrs, func(i, j int) bool {
				if attrs[i].Name.Space != attrs[j].Name.Space {
					return attrs[i].Name.Space < attrs[j].Name.Space
				}
				return attrs[i].Name.Local < attrs[j].Name.Local
			})
			allNodes = append(allNodes, &xmlNode{XMLName: tok.Name, XMLAttr: attrs})
		case xml.EndElement:
			if len(allNodes) > 1 {
				allNodes[lastNodeIndex-1].Nodes = append(allNodes[lastNodeIndex-1].Nodes, lastNode)
				allNodes = allNodes[:lastNodeIndex]
			}
		case xml.CharData:
			lastNode.Content = append(lastNode.Content, tok.Copy()...)
		}
	}

	if len(allNodes) == 0 {
		return nil, nil
	}

	firstNode := allNodes[0]
	trimParentContent(firstNode)

	return firstNode, nil
}

// trimParentContent replicates gomega's trimParentNodesContentSpaces:
// for any node with child elements, trim leading/trailing whitespace
// from its Content.
func trimParentContent(node *xmlNode) {
	if len(node.Nodes) > 0 {
		node.Content = bytes.TrimSpace(node.Content)
		for _, child := range node.Nodes {
			trimParentContent(child)
		}
	}
}

// NB: this is really just to cut down on "missing" non-critical test
// coverage. this should recursively stringify all the content.
func stringifyEverything(section *booklit.Section) string {
	var str strings.Builder
	str.WriteString(section.String() + " " + section.Body.String())

	for _, sub := range section.Children {
		str.WriteString(stringifyEverything(sub))
	}

	return str.String()
}
