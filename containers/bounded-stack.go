package containers

const BoundedStack string = `
package {{.Package}}

// {{.Container}} represents a bounded stack of {{.Containee}}.
type {{.Container}} struct {
	top  *item
	size int
	max  int
}

// internal item structure
type item struct {
	value *{{.Containee}}
	next  *item
}

// New{{.Container}} initializes and returns a new bounded stack of {{.Containee}}
func New{{.Container}}(max int) *{{.Container}} {
	return &{{.Container}}{max: max}
}

// Len returns the stack's length
func (s *{{.Container}}) Len() int {
	return s.size
}

// Max returns the stack's maximum size
func (s *{{.Container}}) Max() int {
	return s.max
}

// Push pushes a new item on top of the stack.
//
// In case this operation would make the stack size greater than its maximum,
// the bottommost element is removed before pushing the new one.
func (s *{{.Container}}) Push(value *{{.Containee}}) {
	if s.size+1 > s.max {
		if last := s.PopLast(); last == nil {
			panic("Unexpected nil in stack")
		}
	}
	s.top = &item{value, s.top}
	s.size++
}

// Pop removes the topmost item from the stack and return its value
//
// If the stack is empty, Pop returns nil
func (s *{{.Container}}) Pop() (value *{{.Containee}}) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

// PopLast removes the bottommost item.
//
// PopLast does nothing if the stack does not contain at least 2 items.
func (s *{{.Container}}) PopLast() (value *{{.Containee}}) {
	if lastElem := s.popLast(s.top); s.size >= 2 && lastElem != nil {
		return lastElem.value
	}
	return nil
}

// Peek returns the topmost without removing it from the stack
func (s *{{.Container}}) Peek() (value *{{.Containee}}, exists bool) {
	exists = false
	if s.size > 0 {
		value = s.top.value
		exists = true
	}

	return
}

func (s *{{.Container}}) popLast(elem *item) *item {
	if elem == nil {
		return nil
	}
	// not last because it has next and a grandchild
	if elem.next != nil && elem.next.next != nil {
		return s.popLast(elem.next)
	}

	// current elem is second from bottom, as next elem has no child
	if elem.next != nil && elem.next.next == nil {
		last := elem.next
		// make current elem bottom of stack by removing its next item
		elem.next = nil
		s.size--
		return last
	}
	return nil
}
`
