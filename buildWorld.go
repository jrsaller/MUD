package main

import (
	"database/sql"
	"fmt"
	"log"
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

func ReadRoom(id int, rooms map[int]*Room) {
	fmt.Println(rooms[id].Name + "\n")
	fmt.Println(rooms[id].Description)
	fmt.Printf("Directions: [ ")
	for index, element := range rooms[id].Exits {
		if element.To != nil {
			fmt.Printf(numToDirection[index] + " ")
		}
	}
	fmt.Printf("]\n")
}
