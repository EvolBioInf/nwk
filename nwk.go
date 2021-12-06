// Package nwk implements data structures, methods, and functions to manipulate phylogenies in Newick format.
package nwk

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"text/scanner"
)

// A Node is the basic uni of a Newick tree.
type Node struct {
	Id                 int
	Child, Sib, Parent *Node
	Label              string
	Length             float64
	HasLength          bool
	marked             bool
}

// Scanner scans an input file one tree at a time.
type Scanner struct {
	r    *bufio.Reader
	text string
}

var nodeId = 1

// Method AddChild adds a child node to a Node.  Inside
func (n *Node) AddChild(v *Node) {
	v.Parent = n
	if n.Child == nil {
		n.Child = v
	} else {
		w := n.Child
		for w.Sib != nil {
			w = w.Sib
		}
		w.Sib = v
	}
}

//  The method RemoveChild removes a direct child node, if  present. If not, it returns an error.
func (v *Node) RemoveChild(c *Node) error {
	if v.Child == nil {
		return errors.New("no children")
	}
	if v.Child.Id == c.Id {
		v.Child = v.Child.Sib
		return nil
	}
	w := v.Child
	for w.Sib != nil && w.Sib.Id != c.Id {
		w = w.Sib
	}
	if w.Sib == nil {
		return errors.New("child not found")
	} else {
		w.Sib = w.Sib.Sib
	}
	return nil
}

// The method LCA returns the lowest common ancestor of two nodes or nil, if none could be found.
func (v *Node) LCA(w *Node) *Node {
	clearPath(v)
	clearPath(w)
	markPath(v)
	for w != nil && !w.marked {
		w = w.Parent
	}
	return w
}

//  The method UpDistance returns the distance between the node  and one of its ancestors.
func (v *Node) UpDistance(w *Node) float64 {
	s := 0.0
	x := v
	for x != nil && x.Id != w.Id {
		s += x.Length
		x = x.Parent
	}
	if x == nil {
		log.Fatal("can't find ancestor")
	}
	return s
}

//  The method UniformLabels labels all nodes in the subtree with  a prefix followed by the node ID.
func (v *Node) UniformLabels(pre string) {
	label(v, pre)
}

// String turns a tree into its Newick string.
func (n *Node) String() string {
	w := new(bytes.Buffer)
	writeTree(n, w)
	return w.String()
}

// The method Scan advances the scanner by one tree.
func (s *Scanner) Scan() bool {
	var err error
	s.text, err = s.r.ReadString(';')
	if err == nil {
		return true
	}
	return false
}

// The method Tree returns the most recent tree scanned.
func (s *Scanner) Tree() *Node {
	var root *Node
	var tokens []string
	tree := s.Text()
	tree = strings.ReplaceAll(tree, "[", "/*")
	tree = strings.ReplaceAll(tree, "]", "*/")
	tree = strings.ReplaceAll(tree, "'", "\"")
	tree = strings.ReplaceAll(tree, "\"\"", "'")
	c1 := []rune(tree)
	var c2 []rune
	isNum := false
	for _, r := range c1 {
		if r == ':' {
			isNum = true
			c2 = append(c2, '"')
		}
		if isNum && (r == ',' || r == ';' || r == ' ' || r == ')') {
			isNum = false
			c2 = append(c2, '"')
		}
		c2 = append(c2, r)
	}
	tree = string(c2)
	var tsc scanner.Scanner
	tsc.Init(strings.NewReader(tree))
	for t := tsc.Scan(); t != scanner.EOF; t = tsc.Scan() {
		text := tsc.TokenText()
		if text[0] == '"' {
			var err error
			text, err = strconv.Unquote(text)
			if err != nil {
				log.Fatalf("couldn't unquote %q\n", text)
			}
		} else {
			text = strings.ReplaceAll(text, "_", " ")
		}
		tokens = append(tokens, text)
	}
	i := 0
	v := root
	for i < len(tokens) {
		t := tokens[i]
		if t == "(" {
			if v == nil {
				v = NewNode()
			}
			v.AddChild(NewNode())
			v = v.Child
		}
		if t == ")" {
			v = v.Parent
		}
		if t == "," {
			s := NewNode()
			s.Parent = v.Parent
			v.Sib = s
			v = v.Sib
		}
		if t[0] == ':' {
			l, err := strconv.ParseFloat(t[1:], 64)
			if err != nil {
				log.Fatalf("didn't understand %q\n", t[1:])
			}
			v.Length = l
			v.HasLength = true
		}
		if t == ";" {
			break
		}
		if strings.IndexAny(t[:1], ")(,:;") == -1 {
			v.Label = t
		}
		i++
	}
	root = v
	return root
}

// The method Text returns the text scanned most recently.
func (s *Scanner) Text() string {
	return s.text
}

// NewNode returns a new node with a unique Id.
func NewNode() *Node {
	n := new(Node)
	n.Id = nodeId
	nodeId++
	return n
}
func clearPath(v *Node) {
	for v != nil {
		v.marked = false
		v = v.Parent
	}
}
func markPath(v *Node) {
	for v != nil {
		v.marked = true
		v = v.Parent
	}
}
func label(v *Node, pre string) {
	if v == nil {
		return
	}
	label(v.Child, pre)
	label(v.Sib, pre)
	v.Label = pre + strconv.Itoa(v.Id)
}
func writeTree(v *Node, w *bytes.Buffer) {
	if v == nil {
		return
	}
	if v.Parent != nil && v.Parent.Child.Id != v.Id {
		fmt.Fprint(w, ",")
	}
	if v.Child != nil {
		fmt.Fprint(w, "(")
	}
	writeTree(v.Child, w)
	printLabel(w, v)
	writeTree(v.Sib, w)
	if v.Parent != nil && v.Sib == nil {
		fmt.Fprint(w, ")")
	}
	if v.Parent == nil {
		fmt.Fprint(w, ";")
	}
}
func printLabel(w *bytes.Buffer, v *Node) {
	label := v.Label
	if strings.IndexAny(label, "(),") != -1 {
		label = strings.ReplaceAll(label, "'", "''")
		label = fmt.Sprintf("'%s'", label)
	} else {
		label = strings.ReplaceAll(label, " ", "_")
	}
	fmt.Fprintf(w, "%s", label)
	if v.HasLength && v.Parent != nil {
		fmt.Fprintf(w, ":%.3g", v.Length)
	}
}

//  NewScanner returns a scanner for scanning Newick-formatted  phylogenies.
func NewScanner(r io.Reader) *Scanner {
	sc := new(Scanner)
	sc.r = bufio.NewReader(r)
	return sc
}
