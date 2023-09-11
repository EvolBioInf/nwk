package nwk

import (
	"os"
	"testing"
)

func TestNewick(t *testing.T) {
	in := "test1.nwk"
	f, err := os.Open(in)
	if err != nil {
		t.Errorf("couldn't open %q", in)
	}
	defer f.Close()
	sc := NewScanner(f)
	want := `(((A:0.2,B:0.3):0.3,(D:0.5,E:0.3):0.2):0.3,F:0.7);`
	var root *Node
	for sc.Scan() {
		root = sc.Tree()
		get := root.String()
		if get != want {
			t.Errorf("get:\n%s\nwant:\n%s\n",
				get, want)
		}
	}
	root2 := root.CopyClade()
	get := root2.String()
	if get != want {
		t.Errorf("get:\n%s\nwant\n%s\n", get, want)
	}
	root.UniformLabels("n")
	want = `(((n13:0.2,n14:0.3)n12:0.3,` +
		`(n16:0.5,n17:0.3)n15:0.2)n11` +
		`:0.3,n18:0.7)n10;`
	get = root.String()
	if get != want {
		t.Errorf("get:\n%s\nwant:\n%s",
			get, want)
	}
	n1 := root.Child.Child.Child
	n2 := root.Child.Child.Sib.Child
	l1 := n1.LCA(n2)
	l2 := n2.LCA(n1)
	if l1.Id != l2.Id || l1.Id != 11 {
		t.Errorf("get:\n%d\nwant:\n%d",
			l1.Id, 11)
	}
	ud := n1.UpDistance(root)
	if ud != 0.8 {
		t.Errorf("get:\n%g\nwant:\n%g", ud, 0.8)
	}
	ch := NewNode()
	ch.Label = "new"
	root.AddChild(ch)
	get = root.String()
	ot := want
	want = `(((n13:0.2,n14:0.3)n12:0.3,` +
		`(n16:0.5,n17:0.3)n15:0.2)n11` +
		`:0.3,n18:0.7,new)n10;`
	if get != want {
		t.Errorf("get:\n%s\nwant:\n%s", get, want)
	}
	root.RemoveChild(ch)
	get = root.String()
	want = ot
	if get != want {
		t.Errorf("get:\n%s\nwant:\n%s", get, want)
	}
	get = root.Print()
	want = "n10\n   n18\n   n11\n      n15\n         n17\n" +
		"         n16\n      n12\n         n14\n" +
		"         n13\n"
	if get != want {
		t.Errorf("get:\n%s\nwant:\n%s", get, want)
	}
	want = "n10$n11$n12$n13$n14$n15$n16$n17$n18"
	get = root.Key("$")
	if get != want {
		t.Errorf("get:\n%s\nwant:\n%s", get, want)
	}
	in = "test2.nwk"
	f, err = os.Open(in)
	if err != nil {
		t.Errorf("couldn't open %q", in)
	}
	defer f.Close()
	sc = NewScanner(f)
	sc.Scan()
	origRoot := sc.Tree()
	copyRoot := origRoot.CopyClade()
	want = `((T1:47)4:1,((T6:15,T7:11)5:20,` +
		`((T8:34,T9:37)6:3,T10:41)7:2)8:4)9;`
	v := copyRoot.Child.Child.Sib
	v.RemoveClade()
	get = copyRoot.String()
	if get != want {
		t.Errorf("get\n%s\nwant:\n%s\n", get, want)
	}
	want = `((T1:47,(T5:31)3:10)4:1,((T6:15,T7:11)5:20,` +
		`((T8:34,T9:37)6:3,T10:41)7:2)8:4)9;`
	copyRoot = origRoot.CopyClade()
	v = copyRoot.Child.Child.Sib.Child
	v.RemoveClade()
	get = copyRoot.String()
	if get != want {
		t.Errorf("get:\n%s\nwant:\n%s\n", get, want)
	}
}
