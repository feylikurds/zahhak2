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

package menu

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gv "github.com/asaskevich/govalidator"
	tb "github.com/nsf/termbox-go"

	z "../common"
	zn "../networking"
)

type Menu struct {
	g *Graphics

	Quit        bool
	Multiplayer bool
	Server      bool
	Host        string
	Port        string
	Name        string
}

func NewMenu() *Menu {
	return &Menu{}
}

func (m *Menu) Options() {
	m.Host = z.HOST_IP
	m.Port = z.PORT_NUM
	m.Name = z.NAME

	m.main()

	if m.Quit || !m.Multiplayer {
		return
	}

	m.multiplayer()

	if m.Quit {
		return
	}

	if m.Server {
		m.host()

		return
	}

	m.client()
}

func (m *Menu) background() {
	m.g = NewGraphics(3.0)

	m.g.BlankScreen()

	m.g.DrawBackground("*", z.BoldColorWhite, 2.5)

	m.g.DrawBackground("*", z.ColorWhite, 3)

	m.g.DrawBackground("*", z.ColorBlack, 4)
}

func (m *Menu) main() {
	m.g.BlankScreen()
	m.background()
	m.g.Resize(10)

	cX, cY := m.g.GetCenter()
	uX, _ := m.g.GetUnits()

	m.g.Print(cX-2*uX, cY-6, "1. Single player", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY-4, "2. Multi-player", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY-2, "3. Quit", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY+3, "Enter 1-3 (default 1):", z.BoldColorCyan)

	m.g.SetCursor(cX-2*uX+len("Enter 1-3 (default 1):"), cY+3)

	m.g.Flush()

	c, _ := m.g.ReadChar()
	s := fmt.Sprintf("%c", c)
	n, _ := strconv.Atoi(s)

	switch n {
	case 2:
		m.Multiplayer = true
	case 3:
		m.Quit = true
	}

	tb.HideCursor()

}

func (m *Menu) multiplayer() {
	m.g.BlankScreen()
	m.background()
	m.g.Resize(10)

	cX, cY := m.g.GetCenter()
	uX, _ := m.g.GetUnits()

	m.g.Print(cX-2*uX, cY-6, "1. Host game (server)", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY-4, "2. Connect to a game (client)", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY-2, "3. Quit", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY+3, "Enter 1-3 (default 1):", z.BoldColorCyan)

	m.g.SetCursor(cX-2*uX+len("Enter 1-3 (default 1):"), cY+3)

	m.g.Flush()

	c, _ := m.g.ReadChar()
	n, _ := strconv.Atoi(string(c))

	switch n {
	case 2:
		m.Server = false
	case 3:
		m.Quit = true
	default:
		m.Server = true
	}

	tb.HideCursor()
}

func (m *Menu) host() {
	m.g.BlankScreen()
	m.background()
	m.g.Resize(10)

	n, e := strconv.Atoi(m.Port)

	if e != nil || !zn.IsTCPPortAvailable(n) {
		tb.Close()

		println("Error: server's port " + m.Port + " can not be opened.")
		os.Exit(1)
	}

	cX, cY := m.g.GetCenter()
	uX, _ := m.g.GetUnits()

	remoteHostnames, e := zn.RemoteHostnames()

	if e == nil {
		m.g.Print(cX-2*uX, cY-6, "Internet hostname", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY-5, remoteHostnames[0], z.BoldColorYellow)
	} else {
		m.g.Print(cX-2*uX, cY-6, "Internet hostname", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY-5, "DNS unable to resolve", z.BoldColorYellow)
	}

	remoteIP, e := zn.RemoteIP()

	if e == nil {
		m.g.Print(cX-2*uX, cY-4, "Internet IP", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY-3, remoteIP, z.BoldColorYellow)
	} else {
		m.g.Print(cX-2*uX, cY-4, "Internet IP", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY-3, e.Error(), z.BoldColorYellow)
	}

	localHostnames, e := zn.LocalHostnames()

	if e == nil {
		m.g.Print(cX-2*uX, cY-2, "Local hostname", z.BoldColorYellow|z.AttrUnderline)

		if len(localHostnames) > 1 {
			m.g.Print(cX-2*uX, cY-1, localHostnames[0], z.BoldColorYellow)
		} else {
			m.g.Print(cX-2*uX, cY-1, "DNS unable to resolve", z.BoldColorYellow)
		}
	} else {
		m.g.Print(cX-2*uX, cY-2, "Local hostname", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY-1, e.Error(), z.BoldColorYellow)
	}

	localIP, e := zn.LocalIP()

	if e == nil {
		m.g.Print(cX-2*uX, cY, "Local IP", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY+1, localIP, z.BoldColorYellow)

		m.Host = localIP
	} else {
		m.g.Print(cX-2*uX, cY, "Local IP", z.BoldColorYellow|z.AttrUnderline)
		m.g.Print(cX-2*uX, cY+1, e.Error(), z.BoldColorYellow)
	}

	m.g.Print(cX-2*uX, cY+3, "Press any key to start:", z.BoldColorCyan)

	m.g.SetCursor(cX-2*uX+len("Press any key to start:"), cY+3)

	m.g.Flush()

	m.g.ReadChar()

	m.Name = "Host"

	tb.HideCursor()
}

func (m *Menu) client() {
	m.g.BlankScreen()
	m.background()
	m.g.Resize(10)

	cX, cY := m.g.GetCenter()
	uX, _ := m.g.GetUnits()

	m.g.Print(cX-2*uX, cY-6, "Enter server's address", z.BoldColorYellow|z.AttrUnderline)
	m.g.Print(cX-2*uX, cY-4, "Format: address [optional port]", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY-3, "e.x. 192.168.0.100 1947", z.BoldColorYellow)
	m.g.Print(cX-2*uX, cY-1, "Press enter when done:", z.BoldColorCyan)

	x := cX - 2*uX

	s := m.g.Readline(x, cY)
	i := strings.Index(s, " ")

	if i > 0 {
		m.Host = s[:i]
		s := s[i+1:]
		n, e := strconv.Atoi(s)

		if e == nil && n > 0 && n < 65535 {
			m.Port = s
		}
	} else {
		m.Host = s
	}

	m.g.Print(cX-2*uX, cY+2, "Type your name:", z.BoldColorCyan)

	m.Name = strings.TrimSpace(m.g.Readline(x, cY+3))

	m.test()

	tb.HideCursor()
}

func (m *Menu) test() {
	if m.Name == "" {
		m.Name = "Opponent"
	}

	if m.Host == "" {
		tb.Close()

		println("Error: server's address is blank")
		os.Exit(1)
	}

	if !gv.IsHost(m.Host) {
		tb.Close()

		println("Error: server's address " + m.Host + " is not valid")
		os.Exit(1)
	}

	if !zn.IsRemotePortOpen(m.Host, m.Port) {
		tb.Close()

		println("Error: server's address " + m.Host + " and port " + m.Port + " is not accessible")
		os.Exit(1)
	}
}
