// go gen-con - Go Generic Containers
// Copyright 2016 Aurélien Rainone. All rights reserved.
// Use of this source code is governed by a MIT license
// that can be found in the LICENSE file.

package containers

const Set string = `
package {{.Package}}

type {{.Container}} struct {
	set map[*{{.Containee}}]struct{}
}

{{if .Exported}}
// New{{.Container}} initializes and returns a new set of {{.Containee}}
func New{{.Container}}() *{{.Container}} { {{else}}
// new{{.Container}} initializes and returns a new set of {{.Containee}}
func new{{.Container}}() *{{.Container}} { {{end}}
	return &{{.Container}}{make(map[*{{.Containee}}]struct{})}
}

// Len returns the set length
func (s *{{.Container}}) Len() int {
	return len(s.set)
}

// Add adds a new element to the set
func (s *{{.Container}}) Add(i *{{.Containee}}) {
	s.set[i] = struct{}{}
}

// Contains returns true if element is contained in the set
func (s *{{.Container}}) Contains(i *{{.Containee}}) bool {
	_, found := s.set[i]
	return found
}

// Remove removes an element from the set
func (s *{{.Container}}) Remove(i *{{.Containee}}) {
	delete(s.set, i)
}

// Each runs a function for each element.
//
// If f() returns null, Each stops the iteration immeditely
func (s *{{.Container}}) Each(f func(*{{.Containee}}) bool) {
	for i := range s.set {
		if !f(i) {
			return
		}
	}
}

// Union adds all elements from another set
func (s *{{.Container}}) Union(other *{{.Container}}) {
	for i := range other.set {
		s.set[i] = struct{}{}
	}
}
`
