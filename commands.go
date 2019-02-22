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

var players map[string]playerInfo

//COMMANDS
var commands = make(map[string]func(string, *Player))

func addCommand(command string, action func(string, *Player)) {
	commands[command] = action
}

func look(line string, player *Player) {
	lookDirections := map[string]int{
		"n": 0, "no": 0, "nor": 0, "nort": 0, "north": 0,
		"e": 1, "ea": 1, "eas": 1, "east": 1,
		"w": 2, "we": 2, "wes": 2, "west": 2,
		"s": 3, "so": 3, "sou": 3, "sout": 3, "south": 3,
		"u": 4, "up": 4,
		"d": 5, "do": 5, "dow": 5, "down": 5}
	lineList := strings.Fields(line)
	if len(lineList) < 2 {
		//fmt.Println("LOOK")
		ReadRoom(player.Location.ID, rooms, player)
		//fmt.Println("HAVE READ THE ROOM")
		return
	}
	if direction, ok := lookDirections[string(lineList[1])]; ok {
		if player.Location.Exits[direction].Description != "" {
			player.Printf(player.Location.Exits[direction].Description)
		} else {
			player.Printf("You do not see anything interesting\n")
		}
	} else {
		player.Printf("Look where?")
	}
}
func north(line string, player *Player) {
	//fmt.Println(player.Location.Exits)

	if player.Location.Exits[0].To != nil {
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " leaves to the north\n")
			}
		}
		player.Location = player.Location.Exits[0].To
		ReadRoom(player.Location.ID, rooms, player)
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " enters from the south\n")
			}
		}
	} else {
		player.Printf("You cant move that way\n")
	}
}

func east(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[1].To != nil {
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " leaves to the east\n")
			}
		}
		player.Location = player.Location.Exits[1].To
		ReadRoom(player.Location.ID, rooms, player)
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " enters from the west\n")
			}
		}
	} else {
		player.Printf("You cant move that way\n")
	}
}

func south(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[3].To != nil {
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " leaves to the south\n")
			}
		}
		player.Location = player.Location.Exits[3].To
		ReadRoom(player.Location.ID, rooms, player)
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " enters from the north\n")
			}
		}
	} else {
		player.Printf("You cant move that way\n")
	}
}

func west(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[2].To != nil {
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " leaves to the west\n")
			}
		}
		player.Location = player.Location.Exits[2].To
		ReadRoom(player.Location.ID, rooms, player)
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " enters from the east\n")
			}
		}
	} else {
		player.Printf("You cant move that way\n")
	}
}

func up(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[4].To != nil {
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " leaves upward\n")
			}
		}
		player.Location = player.Location.Exits[4].To
		ReadRoom(player.Location.ID, rooms, player)
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " enters from below\n")
			}
		}
	} else {
		player.Printf("You cant move that way\n")
	}
}

func down(line string, player *Player) {
	//fmt.Println(player.Location.Exits)
	if player.Location.Exits[5].To != nil {
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " leaves downward\n")
			}
		}
		player.Location = player.Location.Exits[5].To
		ReadRoom(player.Location.ID, rooms, player)
		for _, item := range ActivePlayers {
			if item != player && item.Location == player.Location {
				item.Printf(player.Name + " enters from above\n")
			}
		}
	} else {
		player.Printf("You cant move that way\n")
	}
}

func sigh(_ string, player *Player) {
	player.Printf("You sigh\n")
	for _, item := range ActivePlayers {
		if item.Location == player.Location && item != player {
			item.Printf(player.Name + " sighs\n")
		}
	}
}
func laugh(_ string, player *Player) {
	player.Printf("You laugh\n")
	for _, item := range ActivePlayers {
		if item.Location == player.Location && item != player {
			item.Printf(player.Name + " laughs\n")
		}
	}
}
func recall(_ string, player *Player) {
	player.Printf("You kneel down to pray and your vision blurs...\n")
	player.Location = rooms[3001]
	ReadRoom(player.Location.ID, rooms, player)
	for _, item := range ActivePlayers {
		if item != player && item.Location == player.Location {
			item.Printf(player.Name + " appears in the center of the room\n")
		}
	}
}

func quit(_ string, player *Player) {
	player.Printf("Goodbye!\n")
	//CLOSE THE OUTPUTS CHANNEL FOR THE PLAYER
	close(player.Outputs)
	fmt.Println("Closed output channel for", player.Name)
	//SET PLAYERS OUTPUT CHANNEL TO NIL
	player.Outputs = nil
	//REMOVE THE PLAYER FROM THE ACTIVE PLAYERS MAP
	delete(ActivePlayers, player.Name)
	//ANNOUNCE THE OUTPUT CHANNEL HAS BEEN CLOSED

	for _, item := range ActivePlayers {
		if item.Location == player.Location {
			item.Printf(player.Name + " has left the game\n")
		}
	}
}

func gossip(message string, player *Player) {
	for _, item := range ActivePlayers {

		item.Printf(player.Name + " gossips: " + strings.TrimPrefix(message, "gossip") + "\n")

	}
}

func say(message string, player *Player) {
	for _, item := range ActivePlayers {
		if item.Location == player.Location {
			item.Printf(player.Name + " says:" + strings.TrimPrefix(message, "say") + "\n")
		}
	}
}

func tell(message string, player *Player) {
	field := strings.Fields(message)
	for _, item := range ActivePlayers {
		if item.Name == field[1] {
			item.Printf(player.Name + " tells you:" + strings.TrimPrefix(message, "tell "+item.Name) + "\n")
			return
		}
	}
}

func shout(message string, player *Player) {
	for _, item := range ActivePlayers {
		if item.Location.Zone == player.Location.Zone {
			item.Printf(player.Name + " shouts:" + strings.TrimPrefix(message, "shout") + "\n")
		}
	}
}

func where(message string, player *Player) {
	for _, item := range ActivePlayers {
		if item != player && item.Location.Zone == player.Location.Zone {
			player.Printf(item.Name + ":" + item.Location.Name + "\n")
		}
	}
}

//func
