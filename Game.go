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
	"errors"
	"fmt"
	"sync"
	"time"

	tb "github.com/nsf/termbox-go"

	zc "./canvas"
	z "./common"
	zgo "./gameobjects"
	zm "./music"
)

type Game struct {
	sync.RWMutex

	config      *z.Config
	id          string
	session     string
	broadcast   chan *z.Message
	gameManager *GameManager

	rooms  *Rooms
	player z.IPlayer

	players   *GameObjectMap
	monsters  *GameObjectMap
	bombs     *GameObjectMap
	portals   *GameObjectMap
	missles   *GameObjectMap
	healths   *GameObjectMap
	strengths *GameObjectMap
	treasures *GameObjectMap

	statuses z.IRing
	canvas   *zc.Canvas
	music    *zm.Music

	display bool
	paused  bool
}

func NewGame(config *z.Config) *Game {
	return &Game{
		paused:  true,
		config:  config,
		id:      z.UUID(),
		session: z.UUID(),
	}
}

func (g *Game) Event(action string) *z.Message {
	m := z.NewMessage("Game", g.id, action)

	return m
}

func (g *Game) Start() {
	defer g.recover()

	g.config.Init()

	g.init()

	g.initMultiplayer()
}

func (g *Game) init() {
	g.broadcast = make(chan *z.Message, 1024)

	g.gameManager = NewGameManager(g)
	g.gameManager.Run()
}

func (g *Game) initMultiplayer() {
	if g.config.Multiplayer {
		if g.config.Server {
			g.display = true

			g.initGame()

			g.runCreatures(false, nil)

			g.announce(false, "Multiplayer mode", z.BoldColorWhite)
			g.announce(false, "Hosting game", z.BoldColorWhite)
			g.announce(false, "Port: "+g.config.Port, z.BoldColorWhite)
			g.announce(false, "IP: "+g.config.Host, z.BoldColorWhite)
			g.announce(false, "Waiting for others", z.BoldColorWhite)
			g.announce(false, "Press enter to begin", z.BoldColorWhite)
		} else {
			g.client()
		}
	} else {
		g.display = true

		g.initGame()

		g.announce(false, "Single player mode", z.BoldColorWhite)

		g.runCreatures(false, nil)

		g.pause(false, false)
	}
}

func (g *Game) announce(broadcast bool, text string, color tb.Attribute) {
	if broadcast {
		msg := g.Event("Announce")
		msg.Params["Status"] = text
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		g.broadcast <- msg
	}

	status := &z.Status{Text: text, Color: color}

	go g.statuses.Enqueue(status)
}

func (g *Game) initGame() {
	g.createWorld()

	g.initPlayer()

	g.initHealths()

	g.initStrengths()

	g.initTreasures()

	g.initBombs()

	g.initPortals()

	g.initMonsters()
}

func (g *Game) createWorld() {
	g.statuses = z.NewRing()
	g.statuses.SetCapacity(g.config.NumMsgsDisplay)

	for i := 0; i < g.config.NumMsgsDisplay; i++ {
		g.announce(false, "", z.BoldColorWhite)
	}

	g.music = zm.NewMusic()
	g.music.Run()
	g.music.Background()

	g.announce(false, "Creating world", z.BoldColorWhite)

	g.rooms = NewRooms(g.session, g.broadcast, g.config.Capacity, g.config.WorldWidth, g.config.WorldHeight)

	g.canvas = zc.NewCanvas(g.config, g.rooms, g.statuses)
	g.players = NewGameObjectMap()
	g.healths = NewGameObjectMap()
	g.strengths = NewGameObjectMap()
	g.treasures = NewGameObjectMap()
	g.bombs = NewGameObjectMap()
	g.portals = NewGameObjectMap()
	g.monsters = NewGameObjectMap()
	g.missles = NewGameObjectMap()
}

func (g *Game) initPlayer() {
	g.announce(false, "Initializing player", z.BoldColorWhite)

	id := z.UUID()
	g.newPlayer(false, g.config.Name, id, false)
	igo, _ := g.players.Get(id)
	g.player = igo.(*zgo.Player)

	x, y := g.config.WorldWidth/2, g.config.WorldHeight/2
	g.player.Move(false, x, y)
	g.player.Next(false, 0, 0)
}

func (g *Game) initHealths() {
	s := fmt.Sprintf("Initializing %d healths", g.config.NumHealths)
	g.announce(false, s, z.BoldColorWhite)

	for i := 0; i < g.config.NumHealths; i++ {
		health := zgo.NewHealth(g.broadcast)
		g.healths.Set(health.GetID(), health)
		g.placeRandomly(false, health)
	}
}

func (g *Game) initStrengths() {
	s := fmt.Sprintf("Initializing %d strengths", g.config.NumStrengths)
	g.announce(false, s, z.BoldColorWhite)

	for i := 0; i < g.config.NumStrengths; i++ {
		strength := zgo.NewStrength(g.broadcast)
		g.strengths.Set(strength.GetID(), strength)
		g.placeRandomly(false, strength)
	}
}

func (g *Game) initTreasures() {
	s := fmt.Sprintf("Initializing %d treasures", g.config.NumTreasures)
	g.announce(false, s, z.BoldColorWhite)

	for i := 0; i < g.config.NumTreasures; i++ {
		treasure := zgo.NewTreasure(g.broadcast)
		g.treasures.Set(treasure.GetID(), treasure)
		g.placeRandomly(false, treasure)
	}
}

func (g *Game) initBombs() {
	s := fmt.Sprintf("Initializing %d bombs", g.config.NumBombs)
	g.announce(false, s, z.BoldColorWhite)

	for i := 0; i < g.config.NumBombs; i++ {
		bomb := zgo.NewBomb(g.broadcast, g.rooms)
		g.bombs.Set(bomb.GetID(), bomb)
		g.placeRandomly(false, bomb)
	}
}

func (g *Game) initPortals() {
	s := fmt.Sprintf("Initializing %d portals", g.config.NumPortals)
	g.announce(false, s, z.BoldColorWhite)

	for i := 0; i < g.config.NumPortals; i++ {
		portal := zgo.NewPortal(g.broadcast, g.config.WorldWidth, g.config.WorldHeight, g.rooms)
		g.portals.Set(portal.GetID(), portal)
		g.placeRandomly(false, portal)
	}
}

func (g *Game) initMonsters() {
	s := fmt.Sprintf("Initializing %d monsters", g.config.NumMonsters)
	g.announce(false, s, z.BoldColorWhite)

	for i := 0; i < g.config.NumMonsters; i++ {
		monster := zgo.NewMonster(g.broadcast, g.config.WorldWidth, g.config.WorldHeight, g.rooms, g.players, g.config.Difficulty)
		g.monsters.Set(monster.GetID(), monster)
		g.placeRandomly(false, monster)
	}
}

func (g *Game) placeRandomly(broadcast bool, igo z.IGameObject) {
	class := igo.GetClass()
	x, y := g.randomFreePlace()

	if class == "Player" || class == "Monster" {
		c, _ := igo.(z.ICreature)
		c.Move(broadcast, x, y)
		c.Next(broadcast, 0, 0)
	} else {
		g.rooms.Enter(broadcast, x, y, igo)
		igo.SetPosition(broadcast, x, y)
	}
}

func (g *Game) randomFreePlace() (int, int) {
	x, y := 0, 0

	for {
		x = z.RandomNumber(0, g.config.WorldWidth-1)
		y = z.RandomNumber(0, g.config.WorldHeight-1)

		if !g.rooms.HasRoomForTwo(x, y) {
			continue
		}

		break
	}

	return x, y
}

func (g *Game) Play() {
	go func() {
		defer g.recover()

		for {
			if g.config.Multiplayer {
				g.clearWorld(g.config.Server)
			} else {
				g.clearWorld(false)
			}

			time.Sleep(1000 * time.Millisecond)
		}
	}()
}

func (g *Game) Display() {
	defer g.recover()

	if g.display {
		n := g.config.NumTreasures
		t := g.treasures.Len()
		h := 0
		s := 0

		if g.player != nil {
			h = g.player.GetHealth()
			s = g.player.GetStrength()
		}

		go g.canvas.Draw(h, s, n-t, n)
	}
}

func (g *Game) MoveKey(x int, y int) {
	defer g.recover()

	if g.paused {
		return
	}

	go g.player.Next(g.config.Multiplayer, x, y)
}

func (g *Game) MoveMouse(x int, y int) {
	defer g.recover()

	if g.paused {
		return
	}

	nextX, nextY := z.MouseToRelative(g.player, x, y)

	go g.player.Next(g.config.Multiplayer, nextX, nextY)
}

func (g *Game) Fire() {
	defer g.recover()

	if g.paused || g.player.GetStrength() <= 0 {
		return
	}

	g.player.ChangeStrength(g.config.Multiplayer, z.STRENGTH_LOST)

	id := z.UUID()
	x, y := g.player.GetPosition()
	nextX, nextY := g.player.GetNext()

	g.newMissle(g.config.Multiplayer, g.config.Name, id, x, y, nextX, nextY)
	g.igo(g.config.Multiplayer, "Start", "Missle", id, []string{})

	if g.config.Multiplayer && !g.config.Server {
		msg := g.Event("Run")
		msg.Class = "Missle"
		msg.ID = id
		g.broadcast <- msg
	} else {
		g.igo(g.config.Multiplayer, "Run", "Missle", id, []string{})
	}
}

func (g *Game) Pause() {
	paused := !g.paused

	if paused {
		msg := g.Event("Announce")
		msg.Params["Status"] = "Pausing game"
		msg.Params["Color"] = fmt.Sprintf("%d", z.BoldColorWhite)
		g.broadcast <- msg
	} else {
		msg := g.Event("Announce")
		msg.Params["Status"] = "Resuming game"
		msg.Params["Color"] = fmt.Sprintf("%d", z.BoldColorWhite)
		g.broadcast <- msg
	}

	if g.config.Multiplayer {
		if g.config.Server {
			g.pause(false, paused)
		}
	} else {
		g.pause(false, paused)
	}
}

func (g *Game) pause(broadcast bool, state bool) {
	defer g.recover()

	if state {
		g.paused = true
		g.stopCreatures(broadcast)
	} else {
		g.paused = false
		g.startCreatures(broadcast)
	}
}

func (g *Game) sfx(broadcast bool, effect string) {
	defer g.recover()

	if broadcast {
		msg := g.Event("Sfx")
		msg.Params["Effect"] = effect
		g.broadcast <- msg
	}

	go g.music.Play(effect)
}

func (g *Game) recover() {
	mode := "Game: "

	if g.config.Multiplayer {
		if g.config.Server {
			mode += "Server: "
		} else {
			mode += "Client: "
		}
	}

	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New(mode + "Recover(): " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New(mode + "Recover(): " + e.Error())
			z.LogPanic(e)
		}
	}
}
