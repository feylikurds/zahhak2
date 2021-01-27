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
	"strconv"

	z "./common"
)

func (g *Game) igo(broadcast bool, event, class, id string, args []string) {
	defer g.recover()

	function := g.action(event)
	gom := g.classGOM(class)

	g.findDoGameObject(broadcast, gom, id, function, args)
}

func (g *Game) action(event string) func(bool, *GameObjectMap, z.IGameObject, []string) {
	var function func(bool, *GameObjectMap, z.IGameObject, []string)

	switch event {
	case "Start":
		function = g.startGameObject

	case "Stop":
		function = g.stopGameObject

	case "Run":
		function = g.runGameObject

	case "Delete":
		function = g.deleteGameObject

	case "SetName":
		function = g.nameGameObject

	case "SetID":
		function = g.idGameObject

	case "SetPosition":
		function = g.positionGameObject

	case "Enter":
		function = g.enterGameObject

	case "Leave":
		function = g.leaveGameObject

	case "Stay":
		function = g.stayCreature

	case "Release":
		function = g.releaseCreature

	case "Next":
		function = g.nextCreature

	case "ChangeHealth":
		function = g.healthCreature

	case "ChangeStrength":
		function = g.strengthCreature

	case "ChangeTreasure":
		function = g.treasureCreature

	default:
		function = nil
	}

	return function
}

func (g *Game) classGOM(class string) *GameObjectMap {
	var gom *GameObjectMap

	switch class {
	case "Player":
		gom = g.players

	case "Monster":
		gom = g.monsters

	case "Bomb":
		gom = g.bombs

	case "Portal":
		gom = g.portals

	case "Missle":
		gom = g.missles

	default:
		gom = nil
	}

	return gom
}

func (g *Game) loadRooms(gom *GameObjectMap) {
	if gom == nil {
		return
	}

	gos := gom.GetValues()

	for _, igo := range gos {
		c := igo.(z.ICreature)
		c.LoadRooms(g.rooms)
	}
}

func (g *Game) prepareMonsters(gom *GameObjectMap) {
	if gom == nil {
		return
	}

	gos := gom.GetValues()

	for _, igo := range gos {
		c := igo.(z.IMonster)
		c.LoadPlayers(g.players)
		c.LoadPlayer()

		x, y := c.GetPosition()
		g.rooms.Enter(false, x, y, igo)

		go c.Blink()
	}
}

func (g *Game) prepareMissles(gom *GameObjectMap) {
	if gom == nil {
		return
	}

	gos := gom.GetValues()

	for _, igo := range gos {
		c := igo.(z.IMissle)

		x, y := c.GetPosition()
		g.rooms.Enter(false, x, y, igo)

		go c.Flash()
	}
}

func (g *Game) clearWorld(broadcast bool) {
	for _, o := range g.players.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.players, o, nil)
		}
	}

	for _, o := range g.monsters.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.monsters, o, nil)
		}
	}

	for _, o := range g.bombs.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.bombs, o, nil)
		}
	}

	for _, o := range g.portals.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.portals, o, nil)
		}
	}

	for _, o := range g.missles.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.missles, o, nil)
		}
	}

	for _, o := range g.healths.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.healths, o, nil)
		}
	}

	for _, o := range g.strengths.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.strengths, o, nil)
		}
	}

	for _, o := range g.treasures.GetValues() {
		if o.Deleted() {
			go g.clearGameObject(broadcast, g.treasures, o, nil)
		}
	}
}

func (g *Game) startCreatures(broadcast bool) {
	for _, c := range g.players.GetValues() {
		c.Start(broadcast)
	}

	for _, c := range g.monsters.GetValues() {
		c.Start(broadcast)
	}

	for _, c := range g.bombs.GetValues() {
		c.Start(broadcast)
	}

	for _, c := range g.portals.GetValues() {
		c.Start(broadcast)
	}

	for _, c := range g.missles.GetValues() {
		c.Start(broadcast)
	}
}

func (g *Game) stopCreatures(broadcast bool) {
	for _, c := range g.players.GetValues() {
		c.Stop(broadcast)
	}

	for _, c := range g.monsters.GetValues() {
		c.Stop(broadcast)
	}

	for _, c := range g.bombs.GetValues() {
		c.Stop(broadcast)
	}

	for _, c := range g.portals.GetValues() {
		c.Stop(broadcast)
	}

	for _, c := range g.missles.GetValues() {
		c.Stop(broadcast)
	}
}

func (g *Game) runCreatures(broadcast bool, exclude z.IGameObject) {
	for _, c := range g.players.GetValues() {
		if c == exclude {
			return
		}

		c.Run(broadcast)
	}

	for _, c := range g.monsters.GetValues() {
		if c == exclude {
			return
		}

		c.Run(broadcast)
	}

	for _, c := range g.bombs.GetValues() {
		if c == exclude {
			return
		}

		c.Run(broadcast)
	}

	for _, c := range g.portals.GetValues() {
		if c == exclude {
			return
		}

		c.Run(broadcast)
	}

	for _, c := range g.missles.GetValues() {
		if c == exclude {
			return
		}

		c.Run(broadcast)
	}
}

func (g *Game) findDoGameObject(broadcast bool, gom *GameObjectMap, key string, function func(bool, *GameObjectMap, z.IGameObject, []string), args []string) {
	if gom == nil {
		return
	}

	igo, e := gom.Get(key)

	if e == nil {
		function(broadcast, gom, igo, args)
	}
}

func (g *Game) startGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	igo.Start(broadcast)
}

func (g *Game) stopGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	igo.Stop(broadcast)
}

func (g *Game) runGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	igo.Run(broadcast)
}

func (g *Game) deleteGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	igo.Delete(broadcast)
}

func (g *Game) nameGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	name := args[0]

	igo.SetName(broadcast, name)
}

func (g *Game) idGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	id := args[0]

	igo.SetID(broadcast, id)
}

func (g *Game) positionGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	x, _ := strconv.Atoi(args[0])
	y, _ := strconv.Atoi(args[1])

	igo.SetPosition(broadcast, x, y)
}

func (g *Game) enterGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	x, _ := strconv.Atoi(args[0])
	y, _ := strconv.Atoi(args[1])

	g.rooms.Enter(broadcast, x, y, igo)
}

func (g *Game) leaveGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	x, _ := strconv.Atoi(args[0])
	y, _ := strconv.Atoi(args[1])

	g.rooms.Leave(broadcast, x, y, igo)
}

func (g *Game) clearGameObject(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	x, y := igo.GetPosition()
	id := igo.GetID()

	igo.Stop(broadcast)
	g.rooms.Leave(broadcast, x, y, igo)

	gom.Delete(id)
	igo = nil
}

func (g *Game) stayCreature(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	if c, ok := igo.(z.ICreature); ok {
		c.Stay(broadcast)
	}
}

func (g *Game) releaseCreature(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	if c, ok := igo.(z.ICreature); ok {
		c.Release(broadcast)
	}
}

func (g *Game) nextCreature(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	if c, ok := igo.(z.ICreature); ok {
		x, _ := strconv.Atoi(args[0])
		y, _ := strconv.Atoi(args[1])

		c.Next(broadcast, x, y)
	}
}

func (g *Game) healthCreature(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	if c, ok := igo.(z.ICreature); ok {
		points, _ := strconv.Atoi(args[0])

		c.ChangeHealth(broadcast, points)
	}
}

func (g *Game) strengthCreature(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	if c, ok := igo.(z.ICreature); ok {
		points, _ := strconv.Atoi(args[0])

		c.ChangeStrength(broadcast, points)
	}
}

func (g *Game) treasureCreature(broadcast bool, gom *GameObjectMap, igo z.IGameObject, args []string) {
	if igo == nil {
		return
	}

	if c, ok := igo.(z.ICreature); ok {
		points, _ := strconv.Atoi(args[0])

		c.ChangeTreasure(broadcast, points)
	}
}
