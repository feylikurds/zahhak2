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

package common

type ICreature interface {
	IGameObject

	Move(bool, int, int)
	Next(bool, int, int)
	GetNext() (int, int)
	Battle(func(ICreature))
	ChangeHealth(bool, int)
	ChangeStrength(bool, int)
	ChangeTreasure(bool, int)
	GetHealth() int
	GetStrength() int
	GetTreasure() int
	Dead() bool
	Stay(bool)
	Release(bool)
	Halted() bool
	LoadRooms(IRooms)
}
