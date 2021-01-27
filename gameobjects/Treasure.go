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

package gameobjects

import (
	z "../common"
)

type Treasure struct {
	*Item
}

func NewTreasure(broadcast chan *z.Message) *Treasure {
	return &Treasure{
		Item: &Item{
			GameObject: &GameObject{
				Class:     "Treasure",
				Name:      "Treasure",
				Symbol:    'T',
				Color:     z.BoldColorBlue,
				ID:        z.UUID(),
				broadcast: broadcast,
				Paused:    true,
			},
			Points: 1,
		},
	}
}
