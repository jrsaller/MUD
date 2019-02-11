package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//MAIN OPERATION

func main() {
	//var zones map[int]*Zone
	//addCommand("look", look)
	//fmt.Println("Adding commands...")
	addCommand("n", north)
	addCommand("no", north)
	addCommand("nor", north)
	addCommand("nort", north)
	addCommand("north", north)

	addCommand("s", south)
	addCommand("so", south)
	addCommand("sou", south)
	addCommand("sout", south)
	addCommand("south", south)

	addCommand("e", east)
	addCommand("ea", east)
	addCommand("eas", east)
	addCommand("east", east)

	addCommand("w", west)
	addCommand("we", west)
	addCommand("wes", west)
	addCommand("west", west)

	addCommand("u", up)
	addCommand("up", up)

	addCommand("d", down)
	addCommand("do", down)
	addCommand("dow", down)
	addCommand("down", down)

	addCommand("l", look)
	addCommand("lo", look)
	addCommand("loo", look)
	addCommand("look", look)
	addCommand("sigh", sigh)
	addCommand("laugh", laugh)
	addCommand("recall", recall)
	fmt.Println("COMMANDS INSTALLED")

	db, err := sql.Open("sqlite3", "./world.db")
	//ZONES TRANSACTION
	zonetx, err := db.Begin()
	zones = ReadZones(zonetx)
	if err != nil {
		_ = zonetx.Rollback()
		log.Fatal(err)
	}
	zonetx.Commit()

	//ROOMS TRANSACTION
	roomstx, err := db.Begin()
	rooms = ReadRooms(roomstx, zones)
	if err != nil {
		_ = roomstx.Rollback()
		log.Fatal(err)
	}
	roomstx.Commit()
	fmt.Println("Read", len(zones), "Zones & ", len(rooms), "Rooms\n")
	//EXITS TRANSACTION
	exitstx, err := db.Begin()
	ReadExits(exitstx, rooms)
	if err != nil {
		_ = exitstx.Rollback()
		log.Fatal(err)
	}
	exitstx.Commit()

	player := Player{rooms[3001]}
	scanner := bufio.NewScanner(os.Stdin)
	ReadRoom(player.Location.ID, rooms)
	fmt.Print("What would you like to do?\n>")

	for scanner.Scan() {
		line := scanner.Text()
		lineList := strings.Fields(line)
		if len(lineList) == 0 {
			fmt.Println("You must enter a command")
		} else if _, ok := commands[lineList[0]]; ok == true {
			//param := strings.TrimLeft(line+" ", lineList[0])
			commands[lineList[0]](line, &player)
		} else {
			fmt.Println("Huh?")
		}

		fmt.Print("What would you like to do?\n>")

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner : %v", err)
	}
}
