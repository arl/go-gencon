// go gen-con - Go Generic Containers
// Copyright 2016 Aurélien Rainone. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package containers

const Stack string = `
package {{.Package}}

// {{.Container}} represents a stack of {{.Containee}}.
type {{.Container}} struct {
	top  *item
	size int
}

// internal item structure.
type item struct {
	value *{{.Containee}}
	next  *item
}

{{if .Exported}}
// New{{.Container}} initializes and returns a new stack of {{.Containee}}.
func New{{.Container}}() *{{.Container}} {
{{else}}
// new{{.Container}} initializes and returns a new stack of {{.Containee}}.
func new{{.Container}}() *{{.Container}} {
{{end}}	return &{{.Container}}{}
}

// Len returns the stack's length.
func (s *{{.Container}}) Len() int {
	return s.size
}

// Push pushes a new item on top of the stack.
func (s *{{.Container}}) Push(value *{{.Containee}}) {
	s.top = &item{value, s.top}
	s.size++
}

// Pop removes the topmost item from the stack and return its value.
//
// If the stack is empty, Pop returns nil.
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

// Peek returns the topmost item without removing it from the stack.
func (s *{{.Container}}) Peek() (value *{{.Containee}}, exists bool) {
	exists = false
	if s.size > 0 {
		value = s.top.value
		exists = true
	}
	return
}

// PeekN returns at max the N topmost item without removing them from the stack.
func (s *{{.Container}}) PeekN(n int) []*{{.Containee}} {
	var (
		N   []*{{.Containee}}
		cur *item
	)
	N = make([]*{{.Containee}}, 0, n)
	cur = s.top
	for len(N) < n {
		if cur == nil {
			break
		}
		N = append(N, cur.value)
		cur = cur.next
	}
	return N
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
