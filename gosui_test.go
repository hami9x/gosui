package gosui

import (
	// "container/list"
	"fmt"
	"image"
	"math/rand"
	"testing"
	"time"

	chk "launchpad.net/gocheck"
)

func Test(t *testing.T) { chk.TestingT(t) }

type MySuite struct{}

var _ = chk.Suite(&MySuite{})

type DummyBackend struct {
	c int
}

func (b *DummyBackend) DrawRect(rect image.Rectangle, radiis [4]int, paint Paint) {
	// fmt.Printf("I'm drawing %v ^^\n", rect)
	b.c += 1
	// fmt.Print("*")
}

func (b *DummyBackend) Init(w, h int) {}

func (b *DummyBackend) DrawElementsInArea(l DrawPriorityList, area image.Rectangle) {
	for _, o := range l {
		o.Draw(b)
	}
}

func random(min, max int) int {
	if max == min {
		return min
	}
	return rand.Intn(max-min) + min
}

func (s *MySuite) TestDraw(c *chk.C) {
	root := NewRootElement()
	NewRectElement(root, MakeRect(0, 0, 100, 100))
	NewRectElement(root, MakeRect(70, 70, 120, 120))
	NewRectElement(root, MakeRect(200, 200, 300, 300))
	backend := new(DummyBackend)

	root.Draw(backend)
	c.Check(backend.c, chk.Equals, 3)
}

func (s *MySuite) TestRedraw(c *chk.C) {
	root := NewRootElement()
	r1 := NewRectElement(root, MakeRect(0, 0, 100, 100))
	NewRectElement(root, MakeRect(70, 70, 120, 120))
	NewRectElement(root, MakeRect(0, 0, 300, 300))
	NewRectElement(root, MakeRect(101, 101, 102, 102))
	backend := new(DummyBackend)
	Redraw(r1, backend, root)
	c.Check(backend.c, chk.Equals, 3) //r1, r2, r3 are redrawn
}

func BenchmarkDisplayEngine(t *testing.B) {
	rand.Seed(time.Now().Unix())
	root := NewRootElement()
	l := make([](*AbstractElement), 0)
	l = append(l, root)
	for i := 1; i < 1000; i += 1 {
		p := NewAbstractElement(l[random(0, len(l)-1)])
		NewRectElement(p,
			MakeRectWH(random(0, 500), random(0, 500), random(0, 500), random(0, 500)))
		l = append(l, p)
	}

	backend := new(DummyBackend)
	root.Draw(backend)
	for i := 1; i < 1000; i += 1 {
		Redraw(l[random(0, len(l)-1)], backend, root)
	}
	fmt.Printf("\n %v objects drawn\n", (backend.c-1)/2)
}
