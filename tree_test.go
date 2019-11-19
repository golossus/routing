package hw14_go

import (
	"net/http"
	"testing"
)

func TestInsertOnEmptyTree(t *testing.T) {
	tree := Tree{}
	tree.Insert("/path1", nil)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}
}

func TestInsertChild(t *testing.T) {
	tree := Tree{}
	tree.Insert("/path1", nil)
	tree.Insert("/path1/path2", nil)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}

	if "/path2" != tree.root.child.prefix {
		t.Errorf("")
	}
}
func TestInsertSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert("/path1", nil)
	tree.Insert("/path2", nil)

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
	tree.Insert("/path1", nil)
	tree.Insert("path2", nil)

	if "/path1" != tree.root.prefix {
		t.Errorf("")
	}

	if "path2" != tree.root.sibling.prefix {
		t.Errorf("")
	}

}

func TestInsertChildOnSibling(t *testing.T) {
	tree := Tree{}
	tree.Insert("/path1", nil)
	tree.Insert("/path2", nil)
	tree.Insert("/path1/path3", nil)

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
	tree.Insert("/path1", nil)
	tree.Insert("/path2", nil)
	tree.Insert("/path3", nil)

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
	tree.Insert("/path1", handler1)
	tree.Insert("/path2", handler2)
	tree.Insert("/path3", handler3)
	tree.Insert("/path3/path4", handler4)
	tree.Insert("/path5/path4", handler5)

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
	tree.Insert("/path1", handler1)
	tree.Insert("/path2", handler2)
	tree.Insert("/path3", handler3)
	tree.Insert("/path3/path4", handler4)
	tree.Insert("/path5/path4", handler5)

	handler := tree.Find("/path3/path4")
	handler(nil,nil)

	if *flag4 != "/path3/path4" {
		t.Errorf("")
	}
}

func generateHandler(path string) (HandlerFunction, *string) {
	var flag string
	return func(response http.ResponseWriter, request *http.Request) {
		flag = path
	}, &flag

}
