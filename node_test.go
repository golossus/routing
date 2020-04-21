package routing

import (
	"regexp"
	"testing"
)

func TestRegexpToStringWorks(t *testing.T) {

	node1 := node{t: nodeTypeStatic}
	if node1.regexpToString() != "" {
		t.Errorf("Node Type static returns invalid string")
	}

	node2 := node{t: nodeTypeDynamic, regexp: nil}
	if node2.regexpToString() != "" {
		t.Errorf("Node without reg expression returns invalid string")
	}

	node3 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}
	if node3.regexpToString() != "[0-9]+" {
		t.Errorf("Node with regular expression returns invalid string")
	}

}

func TestIsCatchAllWorks(t *testing.T) {

	node1 := node{t: nodeTypeStatic}
	if node1.isCatchAll() {
		t.Errorf("Node Type static is catch all")
	}

	node2 := node{t: nodeTypeDynamic, regexp: nil}
	if node2.isCatchAll() {
		t.Errorf("Node without reg expression is catch all")
	}

	node3 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}
	if node3.isCatchAll() {
		t.Errorf("Node with no catch all regexp is catch all")
	}

	node4 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile(catchAllExpression)}
	if !node4.isCatchAll() {
		t.Errorf("Node with valid catch all expression is not catch all")
	}
}

func TestRegexpEqualsWorks(t *testing.T) {

	node1 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[a-z]+")}
	node2 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}
	node3 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}

	if node1.regexpEquals(&node2) {
		t.Errorf("Node1 is equal to node 2")
	}

	if !node2.regexpEquals(&node3) {
		t.Errorf("Node2 is not equal to node 3")
	}
}