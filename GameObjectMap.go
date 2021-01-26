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

package main

import (
	"errors"
	"sync"

	z "./common"
)

type GameObjectMap struct {
	sync.RWMutex

	kvs map[string]z.IGameObject
}

func NewGameObjectMap() *GameObjectMap {
	return &GameObjectMap{
		kvs: map[string]z.IGameObject{}}
}

func (g *GameObjectMap) Set(key string, val z.IGameObject) {
	g.Lock()
	defer g.Unlock()

	g.kvs[key] = val
}

func (g *GameObjectMap) Get(key string) (z.IGameObject, error) {
	g.RLock()
	defer g.RUnlock()

	val, ok := g.kvs[key]

	if ok {
		return val, nil
	}

	return nil, errors.New("Key does not exist")
}

func (g *GameObjectMap) Delete(key string) {
	g.Lock()
	defer g.Unlock()

	delete(g.kvs, key)
}

func (g *GameObjectMap) Len() int {
	g.RLock()
	defer g.RUnlock()

	return len(g.kvs)
}

func (g *GameObjectMap) GetMap() map[string]z.IGameObject {
	g.RLock()
	defer g.RUnlock()

	newKvs := map[string]z.IGameObject{}

	for k, v := range g.kvs {
		newKvs[k] = v
	}

	return newKvs
}

func (g *GameObjectMap) GetKeys() []string {
	g.RLock()
	defer g.RUnlock()

	keys := []string{}

	for key := range g.kvs {
		keys = append(keys, key)
	}

	return keys
}

func (g *GameObjectMap) GetValues() []z.IGameObject {
	g.RLock()
	defer g.RUnlock()

	vals := []z.IGameObject{}

	for _, val := range g.kvs {
		vals = append(vals, val)
	}

	return vals
}

func (g *GameObjectMap) GetRandomKey() (string, error) {
	g.RLock()
	defer g.RUnlock()

	keys := []string{}
	var key string

	for k := range g.kvs {
		key = k
		keys = append(keys, key)
	}

	l := len(g.kvs) - 1
	e := errors.New("Map is empty")

	if l == 0 {
		e = nil
	} else if l > 0 {
		i := z.RandomNumber(0, l)
		key = keys[i]
		e = nil
	}

	return key, e
}

func (g *GameObjectMap) GetRandomValue() (z.IGameObject, error) {
	g.RLock()
	defer g.RUnlock()

	vals := []z.IGameObject{}
	var val z.IGameObject

	for _, v := range g.kvs {
		val = v
		vals = append(vals, val)
	}

	l := len(g.kvs) - 1
	e := errors.New("Map is empty")

	if l == 0 {
		e = nil
	} else if l > 0 {
		i := z.RandomNumber(0, l)
		val = vals[i]
		e = nil
	}

	return val, e
}
