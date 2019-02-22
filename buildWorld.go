package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

//STRUCTS WORLD BUILDING
type Zone struct {
	ID    int
	Name  string
	Rooms []*Room
}

type Room struct {
	ID          int
	Zone        *Zone
	Name        string
	Description string
	Exits       [6]Exit
}

type Exit struct {
	To          *Room
	Description string
}

type Player struct {
	Location *Room
	Outputs  chan OutputEvent
	Name     string
}

//READ WORLD COMMANDS

func ReadZones(tx *sql.Tx) map[int]*Zone {

	zones = make(map[int]*Zone)
	rows, err := tx.Query("SELECT * FROM zones")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var zoneid int
		var zonename string
		var z Zone
		rows.Scan(&zoneid, &zonename)
		z.ID = zoneid
		z.Name = zonename
		zones[zoneid] = &z
		//fmt.Println(zones[zoneid])
	}
	return zones

}

func ReadRooms(tx *sql.Tx, zones map[int]*Zone) map[int]*Room {

	rooms = make(map[int]*Room)
	rows, err := tx.Query("SELECT * FROM rooms")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var roomid, zoneid int
		var roomname, description string
		var r Room
		rows.Scan(&roomid, &zoneid, &roomname, &description)
		r.ID = roomid
		r.Zone = zones[zoneid]
		r.Name = roomname
		r.Description = description
		rooms[roomid] = &r

	}
	return rooms
}

func ReadExits(tx *sql.Tx, rooms map[int]*Room) {
	rows, err := tx.Query("SELECT * FROM exits")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var from, to int
		var direction, description string
		var e Exit
		rows.Scan(&from, &to, &direction, &description)
		e.To = rooms[to]
		e.Description = description
		rooms[from].Exits[directionToNum[direction]] = Exit{rooms[to], description}

	}
}

func ReadRoom(id int, rooms map[int]*Room, p *Player) {
	description := rooms[id].Name + "\n"
	description += rooms[id].Description + "\n"
	description += "Directions: [ "
	for index, element := range rooms[id].Exits {
		if element.To != nil {
			description += numToDirection[index] + " "
		}
	}
	//fmt.Println("PRINTING TO PLAYER OUTPUTS")
	p.Printf(description + "]\n")
	//fmt.Println("END OF READ ROOM")
	for _, item := range ActivePlayers {
		if item.Location == rooms[id] && item != p {
			p.Printf("You see " + item.Name + " standing here\n")
		}
	}
}

type playerInfo struct {
	Name string
	Salt string
	Hash string
}

func ReadPlayers(tx *sql.Tx) map[string]playerInfo {
	players = make(map[string]playerInfo)
	pQuery, err := tx.Query("SELECT * FROM players")
	if err != nil {
		log.Fatal(err)
	}
	//var names []string
	for pQuery.Next() {
		var id int
		var name string
		var salt string
		var hash string
		pQuery.Scan(&id, &name, &salt, &hash)
		players[name] = playerInfo{name, salt, hash}
	}
	return players
}

func (p *Player) Printf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	p.Outputs <- OutputEvent{p, msg}
}

func makePlayerEntry(tx *sql.Tx, name string, salt string, hash string, numPlayers string) {
	id, err := strconv.Atoi(numPlayers)
	if err != nil {
		log.Fatal(err)
	}
	id += 1
	newid := strconv.Itoa(id)
	queryString := "INSERT INTO players VALUES(" + newid + ",\"" + name + "\",\"" + salt + "\",\"" + hash + "\")"
	tx.Exec(queryString)
	players[name] = playerInfo{name, salt, hash}
}
