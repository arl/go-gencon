package containers

const LockFreeQueue string = `

package {{.Package}}

import (
	"unsafe"
	"sync/atomic"
)

// private structure
type node struct {
	value *{{.Containee}}
	next *node
}

type {{.Container}} struct {
	dummy *node
	tail *node
}

{{if .Exported}}
// New{{.Container}} creates a new lock free queue of {{.Containee}}
func New{{.Container}}() *{{.Container}} {
{{else}}
// new{{.Container}} creates a new lock free queue of {{.Containee}}
func new{{.Container}}() *{{.Container}} {
{{end}}
	q := new({{.Container}})
	q.dummy = new(node)
	q.tail = q.dummy

	return q
}

// Enqueue places a new element at the back of the queue.
//
// This method is safe for concurrent use by multiple goroutines
func (q *{{.Container}}) Enqueue(v *{{.Containee}}) {
	var oldTail, oldTailNext *node

	newNode := new(node)
	newNode.value = v

	newNodeAdded := false

	for !newNodeAdded {
		oldTail = q.tail
		oldTailNext = oldTail.next

		if q.tail != oldTail {
			continue
		}

		if oldTailNext != nil {
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(oldTail), unsafe.Pointer(oldTailNext))
			continue
		}

		newNodeAdded = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&oldTail.next)), unsafe.Pointer(oldTailNext), unsafe.Pointer(newNode))
	}

	atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(oldTail), unsafe.Pointer(newNode))
}

// Dequeue dequeues the front element of the queue
//
// This method is safe for concurrent use by multiple goroutines
func (q *{{.Container}}) Dequeue() (*{{.Containee}}, bool) {
	var (
		temp              *{{.Containee}}
		oldDummy, oldHead *node
	)

	removed := false

	for !removed {
		oldDummy = q.dummy
		oldHead = oldDummy.next
		oldTail := q.tail

		if q.dummy != oldDummy {
			continue
		}

		if oldHead == nil {
			return nil, false
		}

		if oldTail == oldDummy {
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(oldTail), unsafe.Pointer(oldHead))
			continue
		}

		temp = oldHead.value
		removed = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.dummy)), unsafe.Pointer(oldDummy), unsafe.Pointer(oldHead))
	}

	return temp, true
}

func (q *{{.Container}}) Iterate() <-chan *{{.Containee}} {
	c := make(chan *{{.Containee}})
	go func() {
		for {
			item, ok := q.Dequeue()
			if !ok {
				break
			}

			c <- item
		}
		close(c)
	}()

	return c
}`
