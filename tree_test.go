package hw14_go

import (
	"net/http"
	"testing"
)

func generateChunk(path string, typeof int) []chunk {
	return []chunk{
		{v: path, t: typeof},
	}
}

func generateStaticChunk(path string) []chunk {
	return generateChunk(path, TChunkStatic)
}

func generateDynamicChunk(ident string) []chunk {
	return generateChunk(ident, TChunkDynamic)
}

func TestInsertOnEmptyTree(t *testing.T) {
	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1"), nil)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}
}

func TestInsertChild(t *testing.T) {
	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1"), nil)
	tree.Insert(generateStaticChunk("/path1/path2"), nil)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}

	if "/path2" != tree.root.child.prefix {
		t.Errorf("")
	}
}

func TestInsertDynamicChild(t *testing.T) {

	h := func(res http.ResponseWriter, req *http.Request) {}

	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1/"), nil)
	tree.Insert(append(generateStaticChunk("/path1/"), generateDynamicChunk("id")...), h)

	if "/path1/" != tree.root.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.t {
		t.Errorf("")
	}

	if "id" != tree.root.child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root.child.t {
		t.Errorf("")
	}
}

func TestInsertDynamicChildHasNoHandler(t *testing.T) {

	h := func(res http.ResponseWriter, req *http.Request) {}

	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1/"), nil)
	tree.Insert(append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id"), generateStaticChunk("/")...)...), h)

	if "/path1/" != tree.root.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.t {
		t.Errorf("")
	}

	if "id" != tree.root.child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root.child.t {
		t.Errorf("")
	}

	child, ok := tree.root.child.stops["/"]

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
	tree.Insert(generateStaticChunk("/path1/"), nil)
	tree.Insert(append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id"), generateStaticChunk("/")...)...), h)
	tree.Insert(append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id"), generateStaticChunk("-")...)...), h)

	if "/path1/" != tree.root.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.t {
		t.Errorf("")
	}

	if "id" != tree.root.child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root.child.t {
		t.Errorf("")
	}

	_, ok := tree.root.child.stops[""]

	if ok {
		t.Errorf("")
	}

	var child, sibling *Node

	child, ok = tree.root.child.stops["/"]

	if !ok {
		t.Errorf("")
	}

	if child.prefix != "/" {
		t.Errorf("")
	}

	if child.handler == nil {
		t.Errorf("")
	}

	sibling, ok = tree.root.child.stops["-"]

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
	tree.Insert(append(generateStaticChunk("/path1")), h)
	tree.Insert(append(generateStaticChunk("/path1/path2")), h)
	tree.Insert(append(generateStaticChunk("/path1/path2/path3")), h)
	tree.Insert(append(generateStaticChunk("/path1/path2/path4")), h)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.t {
		t.Errorf("")
	}
	if nil == tree.root.handler {
		t.Errorf("")
	}

	if "/path2" != tree.root.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.t {
		t.Errorf("")
	}
	if nil == tree.root.child.handler {
		t.Errorf("")
	}

	if "/path" != tree.root.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.t {
		t.Errorf("")
	}
	if nil != tree.root.child.child.handler {
		t.Errorf("")
	}

	if "3" != tree.root.child.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.child.t {
		t.Errorf("")
	}
	if nil == tree.root.child.child.child.handler {
		t.Errorf("")
	}

	if "4" != tree.root.child.child.child.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.child.sibling.t {
		t.Errorf("")
	}
	if nil == tree.root.child.child.child.sibling.handler {
		t.Errorf("")
	}
}

func TestInsertHandlerNotRemovePreviousHandler(t *testing.T) {
	tree := Tree{}
	h := func(w http.ResponseWriter, r *http.Request) {}
	tree.Insert(append(generateStaticChunk("/path1/"), generateDynamicChunk("id")...), h)
	tree.Insert(append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id"), generateStaticChunk("/path2")...)...), h)
	tree.Insert(append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id"), generateStaticChunk("/path3")...)...), h)
	tree.Insert(append(generateStaticChunk("/path1/"), append(generateDynamicChunk("id"), generateStaticChunk("/path2/path4")...)...), h)

	if "/path1/" != tree.root.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.t {
		t.Errorf("")
	}
	if nil != tree.root.handler {
		t.Errorf("")
	}

	if "id" != tree.root.child.prefix {
		t.Errorf("")
	}
	if NodeTypeDynamic != tree.root.child.t {
		t.Errorf("")
	}
	if nil == tree.root.child.handler {
		t.Errorf("")
	}

	if "/path" != tree.root.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.t {
		t.Errorf("")
	}
	if nil != tree.root.child.child.handler {
		t.Errorf("")
	}

	if "2" != tree.root.child.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.child.t {
		t.Errorf("")
	}
	if nil == tree.root.child.child.child.handler {
		t.Errorf("")
	}

	if "3" != tree.root.child.child.child.sibling.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.child.sibling.t {
		t.Errorf("")
	}
	if nil == tree.root.child.child.child.sibling.handler {
		t.Errorf("")
	}

	if "/path4" != tree.root.child.child.child.child.prefix {
		t.Errorf("")
	}
	if NodeTypeStatic != tree.root.child.child.child.child.t {
		t.Errorf("")
	}
	if nil == tree.root.child.child.child.child.handler {
		t.Errorf("")
	}
}

func TestInsertSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1"), nil)
	tree.Insert(generateStaticChunk("/path2"), nil)

	if "/path" != tree.root.prefix {
		t.Errorf("")
	}

	if "1" != tree.root.child.prefix {
		t.Errorf("")
	}

	if "2" != tree.root.child.sibling.prefix {
		t.Errorf("")
	}
}

func TestInsertSiblingNoCommon(t *testing.T) {
	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1"), nil)
	tree.Insert(generateStaticChunk("path2"), nil)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}

	if "path2" != tree.root.sibling.prefix {
		t.Errorf("")
	}

}

func TestInsertChildOnSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1"), nil)
	tree.Insert(generateStaticChunk("/path2"), nil)
	tree.Insert(generateStaticChunk("/path1/path3"), nil)

	if "/path" != tree.root.prefix {
		t.Errorf("")
	}

	if "1" != tree.root.child.prefix {
		t.Errorf("")
	}

	if "2" != tree.root.child.sibling.prefix {
		t.Errorf("")
	}

	if "/path3" != tree.root.child.child.prefix {
		t.Errorf("")
	}
}

func TestInsertSiblingOnSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert(generateStaticChunk("/path1"), nil)
	tree.Insert(generateStaticChunk("/path2"), nil)
	tree.Insert(generateStaticChunk("/path3"), nil)

	if "/path" != tree.root.prefix {
		t.Errorf("")
	}

	if "1" != tree.root.child.prefix {
		t.Errorf("")
	}

	if "2" != tree.root.child.sibling.prefix {
		t.Errorf("")
	}

	if "3" != tree.root.child.sibling.sibling.prefix {
		t.Errorf("")
	}
}

var testTree Tree
var flag *string

func TestInsertWithHandler(t *testing.T) {
	tree := Tree{}
	handler1, flag1 := generateHandler("/path1")
	handler2, _ := generateHandler("/path2")
	handler3, _ := generateHandler("/path3")
	handler4, flag4 := generateHandler("/path3/path4")
	handler5, _ := generateHandler("/path5/path4")
	tree.Insert(generateStaticChunk("/path1"), handler1)
	tree.Insert(generateStaticChunk("/path2"), handler2)
	tree.Insert(generateStaticChunk("/path3"), handler3)
	tree.Insert(generateStaticChunk("/path3/path4"), handler4)
	tree.Insert(generateStaticChunk("/path4/path5"), handler5)

	if nil != tree.root.handler {
		t.Errorf("")
	}

	handler := tree.root.child.handler
	if nil == handler {
		t.Errorf("")
	}
	handler(nil, nil)
	if "/path1" != *flag1 {
		t.Errorf("")
	}

	if nil == tree.root.child.sibling.handler {
		t.Errorf("")
	}

	if nil == tree.root.child.sibling.sibling.handler {
		t.Errorf("")
	}

	handler = tree.root.child.sibling.sibling.child.handler
	if nil == handler {
		t.Errorf("")
	}
	handler(nil, nil)
	if "/path3/path4" != *flag4 {
		t.Errorf("")
	}

	if nil == tree.root.child.sibling.sibling.sibling.handler {
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
	tree.Insert(generateStaticChunk("/path1"), handler1)
	tree.Insert(generateStaticChunk("/path2"), handler2)
	tree.Insert(generateStaticChunk("/path3"), handler3)
	tree.Insert(generateStaticChunk("/path3/path4"), handler4)
	tree.Insert(generateStaticChunk("/path4/path5"), handler5)

	handler, paramBag := tree.Find("/path3/path4")
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
	handler4, _ := generateHandler("/path1/{name}")
	handler5, flag5 := generateHandler("/{date}")
	handler6, flag6 := generateHandler("/path3/{slug}")
	parser := NewParser("/path1/{id}")
	parser.parse()
	tree.Insert(parser.chunks, handler1)
	parser = NewParser("/path1/{id}/path2")
	parser.parse()
	tree.Insert(parser.chunks, handler2)
	parser = NewParser("/path1/{id}-path2")
	parser.parse()
	tree.Insert(parser.chunks, handler3)
	parser = NewParser("/path1/{name}")
	parser.parse()
	tree.Insert(parser.chunks, handler4)
	parser = NewParser("/{date}")
	parser.parse()
	tree.Insert(parser.chunks, handler5)
	parser = NewParser("/path3/{slug}")
	parser.parse()
	tree.Insert(parser.chunks, handler6)

	data := []findResult{
		{path: "/path1/123", ok: true, f: flag1, schema: "/path1/{id}"},
		{path: "/path1/123/", ok: false, f: nil},
		{path: "/path1/123/path2", ok: true, f: flag2, schema: "/path1/{id}/path2"},
		{path: "/path1/123-path2", ok: true, f: flag3, schema: "/path1/{id}-path2"},
		{path: "/path1/pepe", ok: true, f: flag1, schema: "/path1/{id}"},
		{path: "/path1/pepe_path2", ok: true, f: flag1, schema: "/path1/{id}"}, //two siblings dynamic not allowed
		{path: "/2019-20-11", ok: true, f: flag5, schema: "/{date}"},
		{path: "/path3/123", ok: true, f: flag6, schema: "/path3/{slug}"},
		{path: "/path3/123/asdf", ok: true, f: flag6, schema: "/path3/{slug}"},
		{path: "/", ok: false, f: nil},
	}

	for _, item := range data {
		handler, _ := tree.Find(item.path)

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

	parser := NewParser("/path1/{id}")
	parser.parse()
	tree.Insert(parser.chunks, nil)
	parser = NewParser("/path1/{id}/path2/{slug}")
	parser.parse()
	tree.Insert(parser.chunks, nil)
	parser = NewParser("/path1")
	parser.parse()
	tree.Insert(parser.chunks, nil)
	parser = NewParser("/data1/{id}/data2/{id}")
	parser.parse()
	tree.Insert(parser.chunks, nil)

	data := []findResultWithParams{
		{path: "/path1/123", ok: true, keys: []string{ "id" }, values: []string{ "123" } },
		{path: "/path1/123/path2/this-is-a-slug", ok: true, keys: []string{ "id", "slug" }, values: []string{ "123", "this-is-a-slug" } },
		{path: "/path1", ok: true, keys: []string{}, values: []string{} },
		{path: "/data1/123/data2/456", keys: []string{ "id", "id" }, values: []string{ "123", "456" } },
		{path: "/", ok: false, keys: []string{}, values: []string{} },
	}

	for _, item := range data {
		_, paramBag := tree.Find(item.path)

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

func generateHandler(path string) (HandlerFunction, *string) {
	var flag string
	return func(response http.ResponseWriter, request *http.Request) {
		flag = path
	}, &flag

}
