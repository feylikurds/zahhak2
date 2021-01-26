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
	"errors"
	"reflect"
	"sync"
)

type EventManager struct {
	sync.RWMutex

	functions map[string]interface{}
}

func NewEventManager() *EventManager {
	return &EventManager{
		functions: map[string]interface{}{}}
}

func (em *EventManager) On(event string, function interface{}) {
	em.Lock()
	defer em.Unlock()

	em.functions[event] = function
}

func (em *EventManager) Fire(event string, params ...interface{}) ([]reflect.Value, error) {
	em.RLock()
	defer em.RUnlock()

	f := reflect.ValueOf(em.functions[event])

	if len(params) != f.Type().NumIn() {
		return nil, errors.New("Invalid number of parameters")
	}

	in := make([]reflect.Value, len(params))

	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result := f.Call(in)

	return result, nil
}

func (em *EventManager) FireBackground(event string, params ...interface{}) (chan []reflect.Value, error) {
	em.RLock()
	defer em.RUnlock()

	f := reflect.ValueOf(em.functions[event])

	if len(params) != f.Type().NumIn() {
		return nil, errors.New("Invalid number of parameters")
	}

	in := make([]reflect.Value, len(params))

	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	results := make(chan []reflect.Value)

	go func() {
		results <- f.Call(in)
	}()

	return results, nil
}

func (em *EventManager) Clear(event string) {
	em.Lock()
	defer em.Unlock()

	delete(em.functions, event)
}

func (em *EventManager) Empty() {
	em.Lock()
	defer em.Unlock()

	em.functions = map[string]interface{}{}
}

func (em *EventManager) HasEvent(event string) bool {
	em.RLock()
	defer em.RUnlock()

	_, ok := em.functions[event]

	return ok
}

func (em *EventManager) Events() []string {
	em.RLock()
	defer em.RUnlock()

	events := []string{}

	for event := range em.functions {
		events = append(events, event)
	}

	return events
}

func (em *EventManager) Count() int {
	em.RLock()
	defer em.RUnlock()

	return len(em.functions)
}
