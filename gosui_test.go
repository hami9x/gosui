package gosui

import (
	// "container/list"
	"fmt"
	"image"
	chk "launchpad.net/gocheck"
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) { chk.TestingT(t) }

type MySuite struct{}

var _ = chk.Suite(&MySuite{})

type DummyBackend struct {
	c int
}

func (b *DummyBackend) DrawRect(rect image.Rectangle, radiis [4]image.Point, paint Paint) {
	// fmt.Printf("I'm drawing %v ^^\n", rect)
	b.c += 1
	// fmt.Print("*")
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
	r2 := NewRectElement(root, MakeRect(70, 70, 120, 120))
	r2.ZIndex = 1
	r3 := NewRectElement(root, MakeRect(0, 0, 300, 300))
	r3.ZIndex = -1
	NewRectElement(root, MakeRect(99, 99, 100, 100))
	backend := new(DummyBackend)
	r1.Redraw(backend, root)
	c.Check(backend.c, chk.Equals, 6) //r1, r2 and "r4" are redrawn
	//the 3rd aren't redrawn because it is behind r1
}

func BenchmarkDisplayEngine(t *testing.B) {
	rand.Seed(time.Now().Unix())
	root := NewRootElement()
	l := make([](*Element), 0)
	l = append(l, root)
	for i := 1; i < 1000; i += 1 {
		rect := NewRectElement(l[random(0, len(l)-1)],
			MakeRectWH(random(0, 500), random(0, 500), random(0, 500), random(0, 500)))
		l = append(l, rect)
	}

	backend := new(DummyBackend)
	root.Draw(backend)
	for i := 1; i < 100; i += 1 {
		l[random(0, len(l)-1)].Redraw(backend, root)
	}
	fmt.Printf("\n %v objects drawn\n", (backend.c-1)/2)
}
