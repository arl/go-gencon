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
`
