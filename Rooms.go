/*
Zahhak2, a Golang console game.
Copyright (C) 2021 Aryo Pehlewan aryopehlewan@hotmail.com
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
	z "./common"
	zgo "./gameobjects"
)

type Rooms struct {
	rooms [][]z.IRoom
}

func NewRooms(session string, broadcast chan *z.Message, capacity, width, height int) *Rooms {
	rooms := make([][]z.IRoom, width)

	for x := 0; x < width; x++ {
		rooms[x] = make([]z.IRoom, height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rooms[x][y] = zgo.NewRoom(broadcast, capacity, x, y)
		}
	}

	return &Rooms{
		rooms: rooms}
}

func (r *Rooms) HasRoom(x, y int) bool {
	return r.rooms[x][y].HasRoom()
}

func (r *Rooms) HasRoomForTwo(x, y int) bool {
	return r.rooms[x][y].HasRoomForTwo()
}

func (r *Rooms) Enter(broadcast bool, x, y int, igo z.IGameObject) {
	r.rooms[x][y].Enter(broadcast, igo)
}

func (r *Rooms) Leave(broadcast bool, x, y int, igo z.IGameObject) {
	r.rooms[x][y].Leave(broadcast, igo)
}

func (r *Rooms) GetCreatures(x, y int) []z.ICreature {
	return r.rooms[x][y].GetCreatures()
}

func (r *Rooms) GetHealths(x, y int) []z.IHealth {
	return r.rooms[x][y].GetHealths()
}

func (r *Rooms) GetStrengths(x, y int) []z.IStrength {
	return r.rooms[x][y].GetStrengths()
}

func (r *Rooms) GetTreasures(x, y int) []z.ITreasure {
	return r.rooms[x][y].GetTreasures()
}

func (r *Rooms) GetGameObjects(x, y int) []z.IGameObject {
	return r.rooms[x][y].GetGameObjects()
}
