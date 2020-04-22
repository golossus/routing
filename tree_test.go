package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func nodePrefixHandler(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, prefix)
	}
}

func parseAndInsertSchema(tree *tree, schema, prefixHandler string) {
	parser := newParser(schema)
	_, _ = parser.parse()
	tree.insert(parser.chunks, nodePrefixHandler(prefixHandler))
}

func assertNodeValid(t *testing.T, node *node, nodeType int, prefix string, hasHandler bool) {
	if nodeType != node.t {
		t.Errorf("node prefix %s is not static", prefix)
	}

	if prefix != node.prefix {
		t.Errorf("node prefix %s not equals to prefix %s ", node.prefix, prefix)
	}

	if hasHandler == true && node.handler == nil {
		t.Errorf("node prefix %s don't have handler", node.prefix)
	}

	if hasHandler == false && node.handler != nil {
		t.Errorf("node prefix %s has handler", node.prefix)
	}

	if node.handler != nil {
		w := httptest.NewRecorder()
		node.handler(w, nil)
		if w.Body.String() != prefix {
			t.Errorf("invalid handler in node prefix %s", node.prefix)
		}
	}
}

func assertNodeStatic(t *testing.T, node *node, prefix string, hasHandler bool) {
	assertNodeValid(t, node, nodeTypeStatic, prefix, hasHandler)
}

func assertNodeDynamic(t *testing.T, node *node, prefix string, pattern string, hasHandler bool) {
	assertNodeValid(t, node, nodeTypeDynamic, prefix, hasHandler)

	if pattern != node.regexpToString() {
		t.Errorf("node regExp %s not equals to regExp %s ", node.regexpToString(), pattern)
	}

	if node.handler != nil {
		w := httptest.NewRecorder()
		node.handler(w, nil)
		if w.Body.String() != prefix {
			t.Errorf("invalid handler in node prefix %s", node.prefix)
		}
	}
}

func TestInsertChild(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "/path1")
	parseAndInsertSchema(tree, "/path1/path2", "/path2")

	assertNodeStatic(t, tree.root, "/path1", true)
	assertNodeStatic(t, tree.root.child, "/path2", true)
}

func TestInsertDynamicChild(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id}", "id")

	assertNodeStatic(t, tree.root, "/path1/", true)
	assertNodeDynamic(t, tree.root.child, "id", "", true)
}

func TestInsertDynamicChildWithRegularExpression(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id:[0-9]+}", "id")

	assertNodeStatic(t, tree.root, "/path1/", true)
	assertNodeDynamic(t, tree.root.child, "id", "^[0-9]+$", true)
}

func TestInsertDynamicChildHasNoHandler(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id:[0-9]{4}-[0-9]{2}-[0-9]{2}}/", "/")

	assertNodeStatic(t, tree.root, "/path1/", true)
	assertNodeDynamic(t, tree.root.child, "id", "^[0-9]{4}-[0-9]{2}-[0-9]{2}$", false)
	assertNodeStatic(t, tree.root.child.stops['/'], "/", true)
}

func TestInsertDynamicChildHasNoHandlerWithSiblings(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id}/", "/")
	parseAndInsertSchema(tree, "/path1/{id}-", "-")

	assertNodeStatic(t, tree.root, "/path1/", true)
	assertNodeDynamic(t, tree.root.child, "id", "", false)

	if len(tree.root.child.stops) != 2 {
		t.Errorf("")
	}
	assertNodeStatic(t, tree.root.child.stops['/'], "/", true)
	assertNodeStatic(t, tree.root.child.stops['-'], "-", true)
}

func TestInsertHandlerIsOnlyOnLeaf(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "/path1")
	parseAndInsertSchema(tree, "/path1/path2", "/path2")
	parseAndInsertSchema(tree, "/path1/path2/path3", "3")
	parseAndInsertSchema(tree, "/path1/path2/path4", "4")

	assertNodeStatic(t, tree.root, "/path1", true)
	assertNodeStatic(t, tree.root.child, "/path2", true)
	assertNodeStatic(t, tree.root.child.child, "/path", false)
	assertNodeStatic(t, tree.root.child.child.child, "3", true)
	assertNodeStatic(t, tree.root.child.child.child.sibling, "4", true)
}

func TestInsertHandlerNotRemovePreviousHandler(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/{id}", "id")
	parseAndInsertSchema(tree, "/path1/{id}/path2", "2")
	parseAndInsertSchema(tree, "/path1/{id}/path3", "3")
	parseAndInsertSchema(tree, "/path1/{id}/path2/path4", "/path4")

	assertNodeStatic(t, tree.root, "/path1/", false)
	assertNodeDynamic(t, tree.root.child, "id", "", true)
	assertNodeStatic(t, tree.root.child.stops['/'], "/path", false)
	assertNodeStatic(t, tree.root.child.stops['/'].child, "2", true)
	assertNodeStatic(t, tree.root.child.stops['/'].child.sibling, "3", true)
	assertNodeStatic(t, tree.root.child.stops['/'].child.child, "/path4", true)
}

func TestInsertChildOnSibling(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")
	parseAndInsertSchema(tree, "/path2/path3", "/path3")

	assertNodeStatic(t, tree.root, "/path", false)
	assertNodeStatic(t, tree.root.child, "1", true)
	assertNodeStatic(t, tree.root.child.sibling, "2", true)
	assertNodeStatic(t, tree.root.child.sibling.child, "/path3", true)
}

func TestInsertSiblingOnSibling(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")
	parseAndInsertSchema(tree, "/path3", "3")

	assertNodeStatic(t, tree.root, "/path", false)
	assertNodeStatic(t, tree.root.child, "1", true)
	assertNodeStatic(t, tree.root.child.sibling, "2", true)
	assertNodeStatic(t, tree.root.child.sibling.sibling, "3", true)
}

func TestInsertPrioritisesStaticPaths(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/{id}", "id")
	parseAndInsertSchema(tree, "/{name}", "name")
	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")

	assertNodeStatic(t, tree.root, "/", false)
	assertNodeStatic(t, tree.root.child, "path", false)
	assertNodeStatic(t, tree.root.child.child, "1", true)
	assertNodeStatic(t, tree.root.child.child.sibling, "2", true)
	assertNodeDynamic(t, tree.root.child.sibling, "id", "", true)
	assertNodeDynamic(t, tree.root.child.sibling.sibling, "name", "", true)
}

func TestCreateTreeFromChunksWorks(t *testing.T) {

	chunks := []chunk{
		{t: tChunkStatic, v: "/"},
		{t: tChunkDynamic, v: "id"},
		{t: tChunkStatic, v: "/abc"},
	}

	root := createTreeFromChunks(chunks, nodePrefixHandler("/abc"))

	assertNodeStatic(t, root, "/", false)
	assertNodeDynamic(t, root.child, "id", "", false)
	assertNodeStatic(t, root.child.stops['/'], "/abc", true)
}
