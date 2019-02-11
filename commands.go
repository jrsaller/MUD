package main

import (
	"fmt"
	"strings"
)

//LOOKUP TABLES
var directionToNum = map[string]int{"n": 0, "e": 1, "w": 2, "s": 3, "u": 4, "d": 5}

var numToDirection = map[int]string{0: "n", 1: "e", 2: "w", 3: "s", 4: "u", 5: "d"}

var rooms map[int]*Room

var zones map[int]*Zone

//COMMANDS
var commands = make(map[string]func(string, *Player))

func addCommand(command string, action func(string, *Player)) {
	commands[command] = action
}

func look(line string, p *Player) {
	lookDirections := map[string]int{
		"n": 0, "no": 0, "nor": 0, "nort": 0, "north": 0,
		"e": 1, "ea": 1, "eas": 1, "east": 1,
		"w": 2, "we": 2, "wes": 2, "west": 2,
		"s": 3, "so": 3, "sou": 3, "sout": 3, "south": 3,
		"u": 4, "up": 4,
		"d": 5, "do": 5, "dow": 5, "down": 5}
	lineList := strings.Fields(line)
	if len(lineList) < 2 {
		ReadRoom(p.Location.ID, rooms)
		return
	}
	if direction, ok := lookDirections[string(lineList[1])]; ok {
		if p.Location.Exits[direction].Description != "" {
			fmt.Println(p.Location.Exits[direction].Description)
		} else {
			fmt.Println("You do not see anything interesting\n")
		}
	} else {
		fmt.Println("Look where?")
	}
}
func north(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[0].To != nil {
		player.Location = player.Location.Exits[0].To
		ReadRoom(player.Location.ID, rooms)
	} else {
		fmt.Println("You cant move that way")
	}
}

func east(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[1].To != nil {
		player.Location = player.Location.Exits[1].To
		ReadRoom(player.Location.ID, rooms)
	} else {
		fmt.Println("You cant move that way")
	}
}

func south(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[3].To != nil {
		player.Location = player.Location.Exits[3].To
		ReadRoom(player.Location.ID, rooms)
	} else {
		fmt.Println("You cant move that way")
	}
}

func west(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[2].To != nil {
		player.Location = player.Location.Exits[2].To
		ReadRoom(player.Location.ID, rooms)
	} else {
		fmt.Println("You cant move that way")
	}
}

func up(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[4].To != nil {
		player.Location = player.Location.Exits[4].To
		ReadRoom(player.Location.ID, rooms)
	} else {
		fmt.Println("You cant move that way")
	}
}

func down(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[5].To != nil {
		player.Location = player.Location.Exits[5].To
		ReadRoom(player.Location.ID, rooms)
	} else {
		fmt.Println("You cant move that way")
	}
}

func sigh(_ string, player *Player) {
	fmt.Println("You sigh")
}
func laugh(_ string, player *Player) {
	fmt.Println("You laugh")
}
func recall(_ string, player *Player) {
	fmt.Println("You kneel down to pray and your vision blurs...")
	player.Location = rooms[3001]
	ReadRoom(player.Location.ID, rooms)
}
