package routing

import (
	"regexp"
	"testing"
)

func TestNode_RegexpToString_Works(t *testing.T) {

	node1 := node{t: nodeTypeStatic}
	if node1.regexpToString() != "" {
		t.Errorf("node Type static returns invalid string")
	}

	node2 := node{t: nodeTypeDynamic, regexp: nil}
	if node2.regexpToString() != "" {
		t.Errorf("node without reg expression returns invalid string")
	}

	node3 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}
	if node3.regexpToString() != "[0-9]+" {
		t.Errorf("node with regular expression returns invalid string")
	}

}

func TestNode_IsCatchAll_Works(t *testing.T) {

	node1 := node{t: nodeTypeStatic}
	if node1.isCatchAll() {
		t.Errorf("node Type static is catch all")
	}

	node2 := node{t: nodeTypeDynamic, regexp: nil}
	if node2.isCatchAll() {
		t.Errorf("node without reg expression is catch all")
	}

	node3 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}
	if node3.isCatchAll() {
		t.Errorf("node with no catch all regexp is catch all")
	}

	node4 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile(catchAllExpression)}
	if !node4.isCatchAll() {
		t.Errorf("node with valid catch all expression is not catch all")
	}
}

func TestNode_RegexpEquals_Works(t *testing.T) {

	node1 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[a-z]+")}
	node2 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}
	node3 := node{t: nodeTypeDynamic, regexp: regexp.MustCompile("[0-9]+")}

	if node1.regexpEquals(&node2) {
		t.Errorf("node1 is equal to node 2")
	}

	if !node2.regexpEquals(&node3) {
		t.Errorf("node2 is not equal to node 3")
	}
}

func TestNode_HasParameters_Works(t *testing.T) {

	node1 := node{t: nodeTypeStatic}
	node2 := node{t: nodeTypeDynamic, parent: &node1, regexp: regexp.MustCompile("[0-9]+")}
	node3 := node{t: nodeTypeStatic, parent: &node2}

	if node1.hasParameters() != false {
		t.Errorf("node1 has parameters")
	}
	if node2.hasParameters() != true {
		t.Errorf("node2 has no parameters")
	}
	if node3.hasParameters() != true {
		t.Errorf("node3 has no parameters")
	}
}
