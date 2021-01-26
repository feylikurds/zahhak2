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
	"fmt"
	"strconv"
	"time"

	tb "github.com/nsf/termbox-go"

	z "./common"
)

type GameManager struct {
	mode         string
	eventManager *z.EventManager

	multiplayer    bool
	server         bool
	gameID         string
	session        string
	broadcast      chan *z.Message
	host           string
	port           string
	networkManager *NetworkManager
}

func NewGameManager(g *Game) *GameManager {
	multiplayer := g.config.Multiplayer
	server := g.config.Server
	mode := "GameManager: "

	if multiplayer {
		if server {
			mode += "Server: "
		} else {
			mode += "Client: "
		}
	}

	eventManager := z.NewEventManager()

	gameID := g.id
	session := g.session
	broadcast := g.broadcast
	host := g.config.Host
	port := g.config.Port

	eventManager.On("Announce", g.announce)
	eventManager.On("Sfx", g.sfx)
	eventManager.On("Pause", g.pause)
	eventManager.On("NewClient", g.newClient)
	eventManager.On("NewPlayer", g.newPlayer)
	eventManager.On("IGO", g.igo)

	return &GameManager{
		mode:         mode,
		eventManager: eventManager,

		multiplayer: multiplayer,
		server:      server,
		gameID:      gameID,
		session:     session,
		broadcast:   broadcast,
		host:        host,
		port:        port,
	}
}

func (gm *GameManager) Run() {
	if gm.multiplayer {
		gm.networkManager = NewNetworkManager(gm.gameID, gm.session, gm.eventManager, gm.host, gm.port, gm.server)
		gm.networkManager.Run()
	}

	go func() {
		defer gm.recover()

		var e error
		m := &z.Message{}
		ok := true

		for {
			m, ok = <-gm.broadcast

			if !ok {
				e = errors.New(gm.mode + "Outgoing: channel error")
				gm.error(e)

				time.Sleep(1 * time.Second)

				continue
			}

			if gm.multiplayer {
				m.Params["Session"] = gm.session
				m.Params["GameID"] = gm.gameID

				gm.send(m)
			} else {
				action := m.Action

				switch action {
				case "Error":
					status := m.Params["Exception"]
					color := z.BoldColorRed

					gm.eventManager.Fire("Announce", false, status, color)

				case "Announce":
					status := m.Params["Status"]
					s := m.Params["Color"]
					i, _ := strconv.Atoi(s)
					color := tb.Attribute(i)

					gm.eventManager.Fire("Announce", false, status, color)

				case "Sfx":
					effect := m.Params["Effect"]

					gm.eventManager.Fire("Sfx", false, effect)

				default:
				}
			}
		}
	}()

	go func() {
		defer gm.recover()

		if !gm.multiplayer {
			return
		}

		for {
			m, ok := <-gm.networkManager.Messages

			if !ok {
				e := errors.New(gm.mode + "Incoming: channel error")
				gm.error(e)

				time.Sleep(1 * time.Second)

				continue
			}

			session := m.Params["Session"]
			gameID := m.Params["GameID"]
			action := m.Action

			if session != gm.session {
				continue
			}

			if gm.server {
				if action != "Error" {
					gm.send(m)
				}
			}

			switch action {
			case "Error":
				status := m.Params["Exception"]
				color := z.BoldColorRed

				gm.eventManager.Fire("Announce", false, status, color)

			case "NewPlayer":
				name := m.Params["Name"]
				id := m.Params["ID"]

				gm.eventManager.Fire("NewPlayer", false, name, id, true)

			case "Start", "Stop", "Delete", "Stay", "Release":
				class := m.Class
				id := m.ID

				gm.eventManager.Fire("IGO", false, action, class, id, []string{})

			case "SetName", "SetID":
				class := m.Class
				id := m.ID
				prop := action[3:]
				val := m.Params[prop]

				gm.eventManager.Fire("IGO", false, action, class, id, []string{val})

			case "ChangeHealth", "ChangeStrength", "ChangeTreasure":
				class := m.Class
				id := m.ID
				points := m.Params["Points"]

				gm.eventManager.Fire("IGO", false, action, class, id, []string{points})

			case "Next", "SetPosition":
				class := m.Class
				id := m.ID
				x := m.Params["X"]
				y := m.Params["Y"]

				gm.eventManager.Fire("IGO", false, action, class, id, []string{x, y})

			case "Enter", "Leave":
				class := m.Params["Class"]
				id := m.Params["ID"]
				x := m.Params["X"]
				y := m.Params["Y"]

				gm.eventManager.Fire("IGO", false, action, class, id, []string{x, y})

			default:
			}

			if gm.server {
				switch action {
				case "Run":
					class := m.Class
					id := m.ID

					gm.eventManager.Fire("IGO", false, action, class, id, []string{})

				default:
				}
			}

			if gm.server && gm.gameID == gameID {
				continue
			} else {
				switch action {
				case "Announce":
					status := m.Params["Status"]
					s := m.Params["Color"]
					i, _ := strconv.Atoi(s)
					color := tb.Attribute(i)

					gm.eventManager.Fire("Announce", false, status, color)

				case "Sfx":
					effect := m.Params["Effect"]

					gm.eventManager.Fire("Sfx", false, effect)

				case "Pause":
					s := m.Params["State"]
					state := true

					if s == "False" {
						state = false
					}
					gm.eventManager.Fire("Pause", false, state)

				default:
				}
			}
		}
	}()
}

func (gm *GameManager) send(message *z.Message) {
	e := gm.networkManager.Send(message)

	if e != nil {
		if gm.server && e.Error() == "No connection" {
			return
		}

		z.LogError(errors.New(gm.mode + "Send(): " + e.Error()))
		gm.error(e)
	}
}

func (gm *GameManager) SetSession(session string) {
	gm.session = session
	gm.networkManager.session = session
}

func (gm *GameManager) error(e error) {
	gm.eventManager.Fire("Announce", false, e.Error(), z.BoldColorRed)
}

func (gm *GameManager) recover() {
	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New(gm.mode + "Recover(): " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New(gm.mode + "Recover(): " + e.Error())
			z.LogPanic(e)
		}
	}
}
