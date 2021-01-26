/*
Zahhak2, a Golang multiplayer console game.
Copyright (C) 2016 Aryo Pehlewan aryopehlewan@hotmail.com
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package common

import (
	"sync"
)

type Ring struct {
	DefaultCapacity int
	head            int
	tail            int
	buff            []interface{}

	sync.RWMutex
}

func NewRing() *Ring {
	return &Ring{
		DefaultCapacity: 1024,
		head:            -1}
}

func (r *Ring) SetCapacity(size int) {
	r.Lock()
	defer r.Unlock()

	r.checkInit()
	r.extend(size)
}

func (r *Ring) Capacity() int {
	r.RLock()
	defer r.RUnlock()

	v := len(r.buff)

	return v
}

func (r *Ring) Enqueue(i interface{}) {
	r.Lock()
	defer r.Unlock()

	r.checkInit()
	r.set(r.head+1, i)
	old := r.head
	r.head = r.mod(r.head + 1)
	if old != -1 && r.head == r.tail {
		r.tail = r.mod(r.tail + 1)
	}
}

func (r *Ring) Dequeue() interface{} {
	r.Lock()
	defer r.Unlock()

	r.checkInit()
	if r.head == -1 {
		return nil
	}
	v := r.get(r.tail)
	if r.tail == r.head {
		r.head = -1
		r.tail = 0
	} else {
		r.tail = r.mod(r.tail + 1)
	}

	return v
}

func (r *Ring) Peek() interface{} {
	r.RLock()
	defer r.RUnlock()

	r.checkInit()

	if r.head == -1 {
		return nil
	}

	v := r.get(r.tail)

	return v
}

func (r *Ring) Values() []interface{} {
	r.RLock()
	defer r.RUnlock()

	if r.head == -1 {
		return []interface{}{}
	}
	c := len(r.buff)
	arr := make([]interface{}, 0, c)
	for i := 0; i < c; i++ {
		idx := r.mod(i + r.tail)
		arr = append(arr, r.get(idx))
		if idx == r.head {
			break
		}
	}

	return arr
}

func (r *Ring) Empty() []interface{} {
	r.Lock()
	defer r.Unlock()

	if r.head == -1 {
		return []interface{}{}
	}
	c := len(r.buff)
	arr := make([]interface{}, 0, c)
	for i := 0; i < c; i++ {
		idx := r.mod(i + r.tail)
		arr = append(arr, r.get(idx))
		if idx == r.head {
			break
		}
	}

	size := len(r.buff)
	r.buff = nil
	r.head, r.tail = -1, 0

	r.extend(size)

	return arr
}

func (r *Ring) Clear() {
	r.Lock()
	defer r.Unlock()

	size := len(r.buff)
	r.buff = nil
	r.head, r.tail = -1, 0

	r.extend(size)
}

func (r *Ring) set(p int, v interface{}) {
	r.buff[r.mod(p)] = v
}

func (r *Ring) get(p int) interface{} {
	v := r.buff[r.mod(p)]

	return v
}

func (r *Ring) mod(p int) int {
	v := p % len(r.buff)

	return v
}

func (r *Ring) checkInit() {
	if r.buff == nil {
		r.buff = make([]interface{}, r.DefaultCapacity)
		for i := range r.buff {
			r.buff[i] = nil
		}
		r.head, r.tail = -1, 0
	}
}

func (r *Ring) extend(size int) {
	if size == len(r.buff) {
		return
	} else if size < len(r.buff) {
		r.buff = r.buff[0:size]
	}
	newb := make([]interface{}, size-len(r.buff))
	for i := range newb {
		newb[i] = nil
	}
	r.buff = append(r.buff, newb...)
}
