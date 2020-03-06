package http_router

import (
	"net/http"
	"regexp"
	"testing"
)

func generateChunk(path string, typeof int, regExp *regexp.Regexp) []chunk {
	return []chunk{
		{v: path, t: typeof, exp: regExp},
	}
}

func generateStaticChunk(path string) []chunk {
	return generateChunk(path, TChunkStatic, nil)
}

func generateDynamicChunk(ident string, regExp *regexp.Regexp) []chunk {
	return generateChunk(ident, TChunkDynamic, regExp)
}

func TestInsertOnEmptyTree(t *testing.T) {
	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), nil)

	if "/path1" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
}

func TestInsertChild(t *testing.T) {
	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("/path1/path2"), nil)

	if "/path1" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}

	if "/path2" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
}

func TestInsertDynamicChild(t *testing.T) {

	h := func(res http.ResponseWriter, req *http.Request) {}

	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1/"), nil)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), generateDynamicChunk("id", nil)...), h)

	if "/path1/" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].t {
		t.Errorf("")
	}

	if "id" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root[http.MethodGet].child.t {
		t.Errorf("")
	}
}

func TestInsertDynamicChildWithRegularExpression(t *testing.T) {

	h := func(res http.ResponseWriter, req *http.Request) {}

	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1/"), nil)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), generateDynamicChunk("id", regexp.MustCompile("[0-9]+"))...), h)

	if "/path1/" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].t {
		t.Errorf("")
	}

	if "id" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root[http.MethodGet].child.t {
		t.Errorf("")
	}
	if "[0-9]+" != tree.root[http.MethodGet].child.regexp.String() {
		t.Errorf("")
	}
}

func TestInsertDynamicChildHasNoHandler(t *testing.T) {

	h := func(res http.ResponseWriter, req *http.Request) {}

	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1/"), nil)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id", regexp.MustCompile("[0-9]{4}-[0-9]{2}-[0-9]{2}")), generateStaticChunk("/")...)...), h)

	if "/path1/" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].t {
		t.Errorf("")
	}

	if "id" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root[http.MethodGet].child.t {
		t.Errorf("")
	}
	if "[0-9]{4}-[0-9]{2}-[0-9]{2}" != tree.root[http.MethodGet].child.regexp.String() {
		t.Errorf("")
	}

	child, ok := tree.root[http.MethodGet].child.stops["/"]

	if !ok {
		t.Errorf("")
	}

	if child.prefix != "/" {
		t.Errorf("")
	}

	if child.handler == nil {
		t.Errorf("")
	}

}

func TestInsertDynamicChildHasNoHandlerWithSiblings(t *testing.T) {

	h := func(res http.ResponseWriter, req *http.Request) {}

	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1/"), nil)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id", nil), generateStaticChunk("/")...)...), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id", nil), generateStaticChunk("-")...)...), h)

	if "/path1/" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].t {
		t.Errorf("")
	}

	if "id" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root[http.MethodGet].child.t {
		t.Errorf("")
	}

	_, ok := tree.root[http.MethodGet].child.stops[""]

	if ok {
		t.Errorf("")
	}

	var child, sibling *Node

	child, ok = tree.root[http.MethodGet].child.stops["/"]

	if !ok {
		t.Errorf("")
	}

	if child.prefix != "/" {
		t.Errorf("")
	}

	if child.handler == nil {
		t.Errorf("")
	}

	sibling, ok = tree.root[http.MethodGet].child.stops["-"]

	if !ok {
		t.Errorf("")
	}

	if sibling.prefix != "-" {
		t.Errorf("")
	}

	if sibling.handler == nil {
		t.Errorf("")
	}
}

func TestInsertHandlerIsOnlyOnLeaf(t *testing.T) {
	tree := Tree{}
	h := func(w http.ResponseWriter, r *http.Request) {}
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1")), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/path2")), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/path2/path3")), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/path2/path4")), h)

	if "/path1" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].handler {
		t.Errorf("")
	}

	if "/path2" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.handler {
		t.Errorf("")
	}

	if "/path" != tree.root[http.MethodGet].child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.t {
		t.Errorf("")
	}
	if nil != tree.root[http.MethodGet].child.child.handler {
		t.Errorf("")
	}

	if "3" != tree.root[http.MethodGet].child.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.child.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.child.child.handler {
		t.Errorf("")
	}

	if "4" != tree.root[http.MethodGet].child.child.child.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.child.sibling.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.child.child.sibling.handler {
		t.Errorf("")
	}
}

func TestInsertHandlerNotRemovePreviousHandler(t *testing.T) {
	tree := Tree{}
	h := func(w http.ResponseWriter, r *http.Request) {}
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), generateDynamicChunk("id", nil)...), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id", nil), generateStaticChunk("/path2")...)...), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id", nil), generateStaticChunk("/path3")...)...), h)
	tree.Insert(http.MethodGet, append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id", nil), generateStaticChunk("/path2/path4")...)...), h)

	if "/path1/" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].t {
		t.Errorf("")
	}
	if nil != tree.root[http.MethodGet].handler {
		t.Errorf("")
	}

	if "id" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root[http.MethodGet].child.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.handler {
		t.Errorf("")
	}

	if "/path" != tree.root[http.MethodGet].child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.t {
		t.Errorf("")
	}
	if nil != tree.root[http.MethodGet].child.child.handler {
		t.Errorf("")
	}

	if "2" != tree.root[http.MethodGet].child.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.child.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.child.child.handler {
		t.Errorf("")
	}

	if "3" != tree.root[http.MethodGet].child.child.child.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.child.sibling.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.child.child.sibling.handler {
		t.Errorf("")
	}

	if "/path4" != tree.root[http.MethodGet].child.child.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root[http.MethodGet].child.child.child.child.t {
		t.Errorf("")
	}
	if nil == tree.root[http.MethodGet].child.child.child.child.handler {
		t.Errorf("")
	}
}

func TestInsertSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("/path2"), nil)

	if "/path" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}

	if "1" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}

	if "2" != tree.root[http.MethodGet].child.sibling.prefix {
		t.Errorf("")
	}
}

func TestInsertSiblingNoCommon(t *testing.T) {
	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("path2"), nil)

	if "/path1" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}

	if "path2" != tree.root[http.MethodGet].sibling.prefix {
		t.Errorf("")
	}

}

func TestInsertChildOnSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("/path2"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("/path1/path3"), nil)

	if "/path" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}

	if "1" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}

	if "2" != tree.root[http.MethodGet].child.sibling.prefix {
		t.Errorf("")
	}

	if "/path3" != tree.root[http.MethodGet].child.child.prefix {
		t.Errorf("")
	}
}

func TestInsertSiblingOnSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("/path2"), nil)
	tree.Insert(http.MethodGet, generateStaticChunk("/path3"), nil)

	if "/path" != tree.root[http.MethodGet].prefix {
		t.Errorf("")
	}

	if "1" != tree.root[http.MethodGet].child.prefix {
		t.Errorf("")
	}

	if "2" != tree.root[http.MethodGet].child.sibling.prefix {
		t.Errorf("")
	}

	if "3" != tree.root[http.MethodGet].child.sibling.sibling.prefix {
		t.Errorf("")
	}
}

func TestInsertWithHandler(t *testing.T) {
	tree := Tree{}
	handler1, flag1 := generateHandler("/path1")
	handler2, _ := generateHandler("/path2")
	handler3, _ := generateHandler("/path3")
	handler4, flag4 := generateHandler("/path3/path4")
	handler5, _ := generateHandler("/path5/path4")
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), handler1)
	tree.Insert(http.MethodGet, generateStaticChunk("/path2"), handler2)
	tree.Insert(http.MethodGet, generateStaticChunk("/path3"), handler3)
	tree.Insert(http.MethodGet, generateStaticChunk("/path3/path4"), handler4)
	tree.Insert(http.MethodGet, generateStaticChunk("/path4/path5"), handler5)

	if nil != tree.root[http.MethodGet].handler {
		t.Errorf("")
	}

	handler := tree.root[http.MethodGet].child.handler
	if nil == handler {
		t.Errorf("")
	}
	handler(nil, nil)
	if "/path1" != *flag1 {
		t.Errorf("")
	}

	if nil == tree.root[http.MethodGet].child.sibling.handler {
		t.Errorf("")
	}

	if nil == tree.root[http.MethodGet].child.sibling.sibling.handler {
		t.Errorf("")
	}

	handler = tree.root[http.MethodGet].child.sibling.sibling.child.handler
	if nil == handler {
		t.Errorf("")
	}
	handler(nil, nil)
	if "/path3/path4" != *flag4 {
		t.Errorf("")
	}

	if nil == tree.root[http.MethodGet].child.sibling.sibling.sibling.handler {
		t.Errorf("")
	}

}

func TestFindHandler(t *testing.T) {
	tree := Tree{}
	handler1, _ := generateHandler("/path1")
	handler2, _ := generateHandler("/path2")
	handler3, _ := generateHandler("/path3")
	handler4, flag4 := generateHandler("/path3/path4")
	handler5, _ := generateHandler("/path5/path4")
	tree.Insert(http.MethodGet, generateStaticChunk("/path1"), handler1)
	tree.Insert(http.MethodGet, generateStaticChunk("/path2"), handler2)
	tree.Insert(http.MethodGet, generateStaticChunk("/path3"), handler3)
	tree.Insert(http.MethodGet, generateStaticChunk("/path3/path4"), handler4)
	tree.Insert(http.MethodGet, generateStaticChunk("/path4/path5"), handler5)

	handler, paramBag := tree.Find(http.MethodGet, "/path3/path4")
	handler(nil, nil)

	if *flag4 != "/path3/path4" {
		t.Errorf("")
	}

	if len(paramBag.params) != 0 {
		t.Errorf("")
	}
}

type findResult struct {
	path   string
	ok     bool
	f      *string
	schema string
}

type findResultWithParams struct {
	path   string
	ok     bool
	values []string
	keys   []string
}

func TestFindHandlerWithDynamic(t *testing.T) {
	tree := Tree{}
	handler1, flag1 := generateHandler("/path1/{id}")
	handler2, flag2 := generateHandler("/path1/{id}/path2")
	handler3, flag3 := generateHandler("/path1/{id}-path2")
	handler4, flag4 := generateHandler("/path1/{name}")
	handler5, flag5 := generateHandler("/{date}")
	handler6, flag6 := generateHandler("/path3/{slug}")
	handler7, flag7 := generateHandler("/path4/{id:[0-9]+}")
	handler8, flag8 := generateHandler("/path4/{id:[0-9]+}/{slug:[a-z]+}")

	parser := newParser("/path1/{id}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler1)

	parser = newParser("/path1/{id}/path2")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler2)

	parser = newParser("/path1/{id}-path2")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler3)

	parser = newParser("/path1/{name}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler4)

	parser = newParser("/{date:[0-9]{4}-[0-9]{2}-[0-9]{2}}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler5)

	parser = newParser("/path3/{slug}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler6)

	parser = newParser("/path4/{id:[0-9]+}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler7)

	parser = newParser("/path4/{id:[0-9]+}/{slug:[a-z]+}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler8)

	data := []findResult{
		{path: "/path4/1s2a3", ok: false, f: nil},
		{path: "/path1/123", ok: true, f: flag1, schema: "/path1/{id}"},
		{path: "/path1/123/", ok: true, f: flag4, schema: "/path1/{name}"},
		{path: "/path1/123/path2", ok: true, f: flag2, schema: "/path1/{id}/path2"},
		{path: "/path1/123-path2", ok: true, f: flag3, schema: "/path1/{id}-path2"},
		{path: "/path1/pepe", ok: true, f: flag1, schema: "/path1/{id}"},
		{path: "/path1/pepe_path2", ok: true, f: flag1, schema: "/path1/{id}"}, //two siblings dynamic not allowed
		{path: "/2019-20-11", ok: true, f: flag5, schema: "/{date}"},
		{path: "/path3/123", ok: true, f: flag6, schema: "/path3/{slug}"},
		{path: "/path3/123/asdf", ok: true, f: flag6, schema: "/path3/{slug}"},
		{path: "/", ok: false, f: nil},

		{path: "/path4/123", ok: true, f: flag7, schema: "/path4/{id:[0-9]+}"},
		{path: "/path4/123/123", ok: false, f: nil},
		{path: "/path4/123/john", ok: true, f: flag8, schema: "/path4/{id:[0-9]+}/{slug:[a-z]+}"},
	}

	for _, item := range data {
		handler, _ := tree.Find(http.MethodGet, item.path)

		if handler != nil {
			if item.ok == false {
				t.Errorf("")
			}

			handler(nil, nil)
			if *item.f != item.schema {
				t.Errorf("")
			}
		} else {
			if item.ok == true {
				t.Errorf("")
			}
		}
	}

}

func TestFindHandlerWithDynamicAndParameters(t *testing.T) {
	tree := Tree{}

	h := func(response http.ResponseWriter, request *http.Request) {}

	parser := newParser("/path1/{id}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, h)
	parser = newParser("/path1/{id}/path2/{slug}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, h)
	parser = newParser("/path1")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, h)
	parser = newParser("/data1/{id}/data2/{id}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, h)

	data := []findResultWithParams{
		{path: "/path1/123", ok: true, keys: []string{"id"}, values: []string{"123"}},
		{path: "/path1/123/path2/this-is-a-slug", ok: true, keys: []string{"slug", "id"}, values: []string{"this-is-a-slug", "123"}},
		{path: "/path1", ok: true, keys: []string{}, values: []string{}},
		{path: "/data1/123/data2/456", keys: []string{"id", "id"}, values: []string{"456", "123"}},
		{path: "/", ok: false, keys: []string{}, values: []string{}},
	}

	for _, item := range data {
		_, paramBag := tree.Find(http.MethodGet, item.path)

		if len(paramBag.params) != len(item.values) {
			t.Errorf("")
		}

		for index, _ := range item.keys {
			if paramBag.params[index].name != item.keys[index] {
				t.Errorf("")
			}

			if paramBag.params[index].value != item.values[index] {
				t.Errorf("")
			}
		}
	}
}

func TestCreateTreeFromChunksWorks(t *testing.T) {

	chunks := []chunk{
		{t: TChunkStatic, v: "/"},
		{t: TChunkDynamic, v: "id"},
		{t: TChunkStatic, v: "/abc"},
	}

	handler := func(http.ResponseWriter, *http.Request) {}

	tree, err := createTreeFromChunks(chunks, handler)

	if err != nil {
		t.Errorf("Unable to create tree for chunks %v", chunks)
	}

	if tree.prefix != "/" || tree.handler != nil || tree.t != NodeTypeStatic || tree.child == nil {
		t.Errorf("Invalid root node %v", tree)
	}

	tree = tree.child
	if tree.prefix != "id" || tree.handler != nil || tree.t != NodeTypeDynamic || tree.child == nil {
		t.Errorf("Invalid root node %v", tree)
	}
	if _, ok := tree.stops["/"]; !ok {
		t.Errorf("Invalid stops %v for node %v", tree.stops, tree)
	}

	tree = tree.child
	if tree.prefix != "/abc" || tree.handler == nil || tree.t != NodeTypeStatic || tree.child != nil {
		t.Errorf("Invalid leaf node %v", tree)
	}
}

func TestInsertPrioritisesStaticPaths(t *testing.T) {

	handler := func(http.ResponseWriter, *http.Request) {}

	tree := Tree{}

	parser := newParser("/{id}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler)
	parser = newParser("/{name}")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler)
	parser = newParser("/path1")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler)
	parser = newParser("/path2")
	_, _ = parser.parse()
	tree.Insert(http.MethodGet, parser.chunks, handler)

	node := tree.root[http.MethodGet]
	if nil != node.handler {
		t.Errorf("")
	}
	if "/" != node.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != node.t {
		t.Errorf("")
	}

	if nil != node.child.handler {
		t.Errorf("")
	}
	if "path" != node.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != node.child.t {
		t.Errorf("")
	}

	if nil == node.child.sibling.handler {
		t.Errorf("")
	}
	if "id" != node.child.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != node.child.sibling.t {
		t.Errorf("")
	}

	if nil == node.child.sibling.sibling.handler {
		t.Errorf("")
	}
	if "name" != node.child.sibling.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != node.child.sibling.sibling.t {
		t.Errorf("")
	}

	if nil == node.child.child.handler {
		t.Errorf("")
	}
	if "1" != node.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != node.child.child.t {
		t.Errorf("")
	}

	if nil == node.child.child.sibling.handler {
		t.Errorf("")
	}
	if "2" != node.child.child.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != node.child.child.sibling.t {
		t.Errorf("")
	}
}

func generateHandler(path string) (HandlerFunction, *string) {
	var flag string
	return func(response http.ResponseWriter, request *http.Request) {
		flag = path
	}, &flag
}
