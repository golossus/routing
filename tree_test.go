package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func nodePrefixHandler(prefix string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, prefix)
	}
}

func parseAndInsertSchema(tree *tree, schema, prefixHandler string) {
	parser := newParser(schema)
	_, _ = parser.parse()
	tree.insert(parser.chunks, nodePrefixHandler(prefixHandler))
}

func assertNodeRelative(t *testing.T, childNode *node, parentNode *node) {
	if !reflect.DeepEqual(childNode.parent, parentNode) {
		t.Errorf("parent node of child %s is not equal to node %s", childNode.prefix, parentNode.prefix)
	}
}

func assertNodeValid(t *testing.T, node *node, nodeType int, prefix string, hasHandler bool, parentNode *node) {
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

	assertNodeRelative(t, node, parentNode)
}

func assertNodeStatic(t *testing.T, node *node, prefix string, hasHandler bool, parentNode *node) {
	assertNodeValid(t, node, nodeTypeStatic, prefix, hasHandler, parentNode)
}

func assertNodeDynamic(t *testing.T, node *node, prefix string, pattern string, hasHandler bool, parentNode *node) {
	assertNodeValid(t, node, nodeTypeDynamic, prefix, hasHandler, parentNode)

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

func TestTree_Insert_Child(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "/path1")
	parseAndInsertSchema(tree, "/path1/path2", "/path2")

	assertNodeStatic(t, tree.root, "/path1", true, nil)
	assertNodeStatic(t, tree.root.child, "/path2", true, tree.root)
}

func TestTree_Insert_DynamicChild(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id}", "id")

	assertNodeStatic(t, tree.root, "/path1/", true, nil)
	assertNodeDynamic(t, tree.root.child, "id", "", true, tree.root)
}

func TestTree_Insert_DynamicChildWithRegularExpression(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id:[0-9]+}", "id")

	assertNodeStatic(t, tree.root, "/path1/", true, nil)
	assertNodeDynamic(t, tree.root.child, "id", "^[0-9]+$", true, tree.root)
}

func TestTree_Insert_DynamicChildHasNoHandler(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id:[0-9]{4}-[0-9]{2}-[0-9]{2}}/", "/")

	assertNodeStatic(t, tree.root, "/path1/", true, nil)
	assertNodeDynamic(t, tree.root.child, "id", "^[0-9]{4}-[0-9]{2}-[0-9]{2}$", false, tree.root)
	assertNodeStatic(t, tree.root.child.stops['/'], "/", true, tree.root.child)
}

func TestTree_Insert_DynamicChildHasNoHandlerWithSiblings(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/", "/path1/")
	parseAndInsertSchema(tree, "/path1/{id}/", "/")
	parseAndInsertSchema(tree, "/path1/{id}-", "-")

	assertNodeStatic(t, tree.root, "/path1/", true, nil)
	assertNodeDynamic(t, tree.root.child, "id", "", false, tree.root)

	if len(tree.root.child.stops) != 2 {
		t.Errorf("")
	}
	assertNodeStatic(t, tree.root.child.stops['/'], "/", true, tree.root.child)
	assertNodeStatic(t, tree.root.child.stops['-'], "-", true, tree.root.child)
}

func TestTree_Insert_HandlerIsOnlyOnLeaf(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "/path1")
	parseAndInsertSchema(tree, "/path1/path2", "/path2")
	parseAndInsertSchema(tree, "/path1/path2/path3", "3")
	parseAndInsertSchema(tree, "/path1/path2/path4", "4")

	assertNodeStatic(t, tree.root, "/path1", true, nil)
	assertNodeStatic(t, tree.root.child, "/path2", true, tree.root)
	assertNodeStatic(t, tree.root.child.child, "/path", false, tree.root.child)
	assertNodeStatic(t, tree.root.child.child.child, "3", true, tree.root.child.child)
	assertNodeStatic(t, tree.root.child.child.child.sibling, "4", true, tree.root.child.child)
}

func TestTree_Insert_HandlerNotRemovePreviousHandler(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/{id}", "id")
	parseAndInsertSchema(tree, "/path1/{id}/path2", "2")
	parseAndInsertSchema(tree, "/path1/{id}/path3", "3")
	parseAndInsertSchema(tree, "/path1/{id}/path2/path4", "/path4")

	assertNodeStatic(t, tree.root, "/path1/", false, nil)
	assertNodeDynamic(t, tree.root.child, "id", "", true, tree.root)
	assertNodeStatic(t, tree.root.child.stops['/'], "/path", false, tree.root.child)
	assertNodeStatic(t, tree.root.child.stops['/'].child, "2", true, tree.root.child.stops['/'])
	assertNodeStatic(t, tree.root.child.stops['/'].child.sibling, "3", true, tree.root.child.stops['/'])
	assertNodeStatic(t, tree.root.child.stops['/'].child.child, "/path4", true, tree.root.child.stops['/'].child)
}

func TestTree_Insert_ChildOnSibling(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")
	parseAndInsertSchema(tree, "/path2/path3", "/path3")

	assertNodeStatic(t, tree.root, "/path", false, nil)
	assertNodeStatic(t, tree.root.child, "1", true, tree.root)
	assertNodeStatic(t, tree.root.child.sibling, "2", true, tree.root)
	assertNodeStatic(t, tree.root.child.sibling.child, "/path3", true, tree.root.child.sibling)
}

func TestTree_Insert_SiblingOnSibling(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")
	parseAndInsertSchema(tree, "/path3", "3")

	assertNodeStatic(t, tree.root, "/path", false, nil)
	assertNodeStatic(t, tree.root.child, "1", true, tree.root)
	assertNodeStatic(t, tree.root.child.sibling, "2", true, tree.root)
	assertNodeStatic(t, tree.root.child.sibling.sibling, "3", true, tree.root)
}

func TestTree_Insert_PrioritisesStaticPaths(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/{id}", "id")
	parseAndInsertSchema(tree, "/{name}", "name")
	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")

	assertNodeStatic(t, tree.root, "/", false, nil)
	assertNodeStatic(t, tree.root.child, "path", false, tree.root)
	assertNodeStatic(t, tree.root.child.child, "1", true, tree.root.child)
	assertNodeStatic(t, tree.root.child.child.sibling, "2", true, tree.root.child)
	assertNodeDynamic(t, tree.root.child.sibling, "id", "", true, tree.root)
	assertNodeDynamic(t, tree.root.child.sibling.sibling, "name", "", true, tree.root)
}

func TestCreateTreeFromChunks_ReturnsNilIfEmptyChunks(t *testing.T) {

	chunks := []chunk{}

	root, leaf := createTreeFromChunks(chunks)

	assertNil(t, root)
	assertNil(t, leaf)
}

func TestCreateTreeFromChunks(t *testing.T) {

	chunks := []chunk{
		{t: tChunkStatic, v: "/"},
		{t: tChunkDynamic, v: "id"},
		{t: tChunkStatic, v: "/abc"},
	}

	root, leaf := createTreeFromChunks(chunks)

	assertNodeStatic(t, root, "/", false, nil)
	assertNodeDynamic(t, root.child, "id", "", false, root)
	assertNodeStatic(t, root.child.stops['/'], "/abc", false, root.child)
	assertNodeStatic(t, leaf, "/abc", false, root.child)
}

func TestTree_Insert_PrioritisesStaticPathsWithComplexPaths(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/path1/{id}/{name:[a-z]{1,5}}", "")
	parseAndInsertSchema(tree, "/path1/{name:.*}", "name")
	parseAndInsertSchema(tree, "/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", "date")
	parseAndInsertSchema(tree, "/posts/{id}/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}", "date")

	assertNodeStatic(t, tree.root, "/", false, nil)
	assertNodeStatic(t, tree.root.child, "p", false, tree.root)
	assertNodeStatic(t, tree.root.child.child, "ath1/", false, tree.root.child)
	assertNodeStatic(t, tree.root.child.child.sibling, "osts/", false, tree.root.child)
	assertNodeDynamic(t, tree.root.child.sibling, "date", "^[0-9]{4}-[0-9]{2}-[0-9]{2}$", true, tree.root)
	assertNodeDynamic(t, tree.root.child.child.sibling.child, "id", "", false, tree.root.child.child.sibling)
}

func TestTree_OptimizeByWeight_PrioritisesHeavierPathsAllStatic(t *testing.T) {
	tree := &tree{}

	parseAndInsertSchema(tree, "/data", "data")
	parseAndInsertSchema(tree, "/path1", "1")
	parseAndInsertSchema(tree, "/path2", "2")
	parseAndInsertSchema(tree, "/path2/id", "/id")
	parseAndInsertSchema(tree, "/path3", "3")
	parseAndInsertSchema(tree, "/path3/name", "name")
	parseAndInsertSchema(tree, "/path3/phone", "phone")
	parseAndInsertSchema(tree, "/path3/{name:[a-z]+}/phone", "phoneName")

	assertNodeStatic(t, tree.root, "/", false, nil)
	assertNodeStatic(t, tree.root.child, "data", true, tree.root)
	assertNodeStatic(t, tree.root.child.sibling, "path", false, tree.root)
	assertNodeStatic(t, tree.root.child.sibling.child, "1", true, tree.root.child.sibling)
	assertNodeStatic(t, tree.root.child.sibling.child.sibling, "2", true, tree.root.child.sibling)
	assertNodeStatic(t, tree.root.child.sibling.child.sibling.child, "/id", true, tree.root.child.sibling.child.sibling)
	assertNodeStatic(t, tree.root.child.sibling.child.sibling.sibling, "3", true, tree.root.child.sibling)
	assertNodeStatic(t, tree.root.child.sibling.child.sibling.sibling.child, "/", false, tree.root.child.sibling.child.sibling.sibling)
	assertNodeStatic(t, tree.root.child.sibling.child.sibling.sibling.child.child, "name", true, tree.root.child.sibling.child.sibling.sibling.child)
	assertNodeStatic(t, tree.root.child.sibling.child.sibling.sibling.child.child.sibling, "phone", true, tree.root.child.sibling.child.sibling.sibling.child)
	assertNodeDynamic(t, tree.root.child.sibling.child.sibling.sibling.child.child.sibling.sibling, "name", "^[a-z]+$", false, tree.root.child.sibling.child.sibling.sibling.child.child.sibling.parent)

	_ = calcWeight(tree.root)
	tree.root = sortByWeight(tree.root)

	assertNodeStatic(t, tree.root, "/", false, nil)
	assertNodeStatic(t, tree.root.child, "path", false, tree.root)
	assertNodeStatic(t, tree.root.child.sibling, "data", true, tree.root)
	assertNodeStatic(t, tree.root.child.child, "3", true, tree.root.child)
	assertNodeStatic(t, tree.root.child.child.sibling, "2", true, tree.root.child)
	assertNodeStatic(t, tree.root.child.child.sibling.child, "/id", true, tree.root.child.child.sibling)
	assertNodeStatic(t, tree.root.child.child.sibling.sibling, "1", true, tree.root.child)
	assertNodeStatic(t, tree.root.child.child.child, "/", false, tree.root.child.child)
	assertNodeStatic(t, tree.root.child.child.child.child, "name", true, tree.root.child.child.child)
	assertNodeStatic(t, tree.root.child.child.child.child.sibling, "phone", true, tree.root.child.child.child)
}
