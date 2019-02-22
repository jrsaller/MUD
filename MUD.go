package main

import (
	"bufio"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"

	_ "github.com/mattn/go-sqlite3"
)

type InputEvent struct {
	player  *Player
	command []string
	Login   bool
}

type OutputEvent struct {
	player *Player
	Text   string
}

var ActivePlayers map[string]*Player

func acceptConnections(inevent chan InputEvent) {
	//MAKE THE ACTIVE PLAYERS MAP DATA STRUCTURE
	ActivePlayers = make(map[string]*Player)
	//START LISTENING ON :9001 FOR CONNECTIONS
	ln, err := net.Listen("tcp", ":9001")
	//CHECK LISTEN FOR ERRORS
	if err != nil {
		log.Fatalf("Unable to listen: %v", err)
	}
	//WRITE TO CONSOLE WHERE TO ACCEPT CONNECTIONS
	fmt.Println("accepting connections on " + ln.Addr().String())
	for {
		//ATTEMPT TO ACCEPT A CONNECTION
		conn, err := ln.Accept()
		go func() {
			if err != nil {
				log.Fatalf("Unable to accept connection: %v", err)
			}
			//WRITE TO CONSOLE ABOUT NEW CONNECTION
			fmt.Println("NEW CONNECTION DETECTED AT: ", conn.RemoteAddr())

			//CREATE SCANNER FOR NEW CONNECTION
			scan := bufio.NewScanner(conn)
			//CREATE EMPTY NAME AND PASSWORD STRING
			var name, pass string

			//GET NAME FROM CONNECTION
			fmt.Fprintf(conn, "What is your name?\n>")
			for scan.Scan() {
				name = strings.Trim(scan.Text(), " ")
				if name == "" {
					fmt.Fprintf(conn, "Name cannot be empty")
					conn.Close()
					fmt.Println("Connection closed for", conn.RemoteAddr())
					return
				}
				break
			}

			//GET PASSWORD FROM CONNECTION
			fmt.Fprintf(conn, "What is your pasword?\n>")
			for scan.Scan() {
				pass = scan.Text()
				if len(pass) < 8 {
					fmt.Fprintln(conn, "Password must have more than 8 characters")
					conn.Close()
					fmt.Println("Connection closed for", conn.RemoteAddr())
					return
				}
				break
			}
			//MAKE BOOLS FOR IF SALT NEEDS TO BE CREATED AND IF THEY ARE JOINING THE GAME
			var makeSalt bool
			//CHECK IF PLAYER EXISTS IN EXISTING PLAYER LIST
			if item, ok := players[name]; ok {
				//GET THE BINARIES FROM SALT
				s, err := base64.StdEncoding.DecodeString(item.Salt)
				//GET BINARIES OF HASH
				hash2, err2 := base64.StdEncoding.DecodeString(item.Hash)
				if err2 != nil {
					log.Fatal(err)
				}
				//GET HASH OF TYPED PASSWORD USING SALT FROM PREVIOUS LOGINS
				enteredPass := pbkdf2.Key([]byte(pass), s, 64*1024, 32, sha256.New)

				//COMPARE ENTERED HASH AND STORED HASH
				if subtle.ConstantTimeCompare(hash2, enteredPass) != 1 {
					fmt.Fprintln(conn, "Incorrect password.")
					conn.Close()
					fmt.Println("Incorrect login for", name)
					fmt.Println("Connection closed for: ", conn.RemoteAddr())
					//DO NOT MAKE SALT AND DO NOT JOIN MUD
					return
				} else {
					//ENTERED HASH WAS CORRECT, WELCOME BACK uSER
					fmt.Fprintln(conn, "Welcome back "+name)
					//DO NOT MAKE A NEW SALT, BUT DO JOIN THE MUD
					makeSalt = false
				}
			} else {
				//PLAYER DOES NOT EXIST YET
				fmt.Fprintf(conn, "Welcome New Player!\n")
				//Do MAKE SALT AND JOIN THE MUD
				makeSalt = true
			}
			//MAKE OUTPUT CHANNEL FOR PLAYER, AND GIVE TO PLAYER
			outputEvents := make(chan OutputEvent, 100)
			p := Player{rooms[3001], outputEvents, name}
			//IF WE NEED TO MAKE THE SALT, DONE HERE
			if makeSalt {
				//CREATED SALT
				salt := make([]byte, 32)
				_, salterr := crand.Read(salt)
				if salterr != nil {
					log.Fatal(salterr)
				}
				//GENERATE READABLE VERSION OF SALT
				salt64 := base64.StdEncoding.EncodeToString(salt)
				//GENERATE HASH WITH SALT
				hash1 := pbkdf2.Key([]byte(pass), salt, 64*1024, 32, sha256.New)
				//GENERATE READABLE HASH
				hash164 := base64.StdEncoding.EncodeToString(hash1)
				//CREATE NEW LOGIN EVENT,PASSED TO INPUT EVENT CHANNEL
				inevent <- InputEvent{&p, []string{name, salt64, hash164, strconv.Itoa(len(players))}, true}
			}
			//JOINING THE MUD
			//CHECK IF THE PLAYER WAS ACTIVE BEFORE
			if item, ok := ActivePlayers[name]; ok {
				//PLAYER WAS ACTIVE
				//SET NEW PLAYER LOCATION TO WHERE THE PLAYER WAS BEFORE
				p.Location = item.Location
				//SAY WHERE THE NEW PLAYER CONNECTED FROM
				fmt.Println(p.Name, "has joined from a new location,", conn.RemoteAddr())
				//CLOSE THE PREVIOUS PLAYERS OUTPUT CHANNEL
				close(item.Outputs)
				//SET THE OLD PLAYERS CHANNEL TO NIL AFTER ITS CLOSED
				item.Outputs = nil
			}
			//SET THE PLAYERS AS ACTIVE, WITH NAME AS KEY AND VALUE AS THE PLAYER OBJET"S REFERENCE
			ActivePlayers[name] = &p
			//SPIN UP THE CONNECTION HANDLER, PASS THE INPUT EVENT, CONNECTION AND PLAYER OBJECT
			go handleConnection(inevent, conn, &p)

		}()
	}
}

func handleConnection(inevents chan InputEvent, conn net.Conn, p *Player) {
	//MAKE SCANNER FOR THE CONNECTION
	scanner := bufio.NewScanner(conn)
	//SAY WHERE THE PLAYER HAS JOINED FROM
	fmt.Println(p.Name, "has joined the game from ", conn.RemoteAddr())
	//FORCE THE PLAYER TO LOOK IN THE FIRST ROOM
	look("look", p)

	//goroutine that accepts events from the main goroutine

	go func() {
		//WAIT FOR TEXT TO BE PASSED TO THE PLAYER'S OUTPUT CHANNEL
		for event := range p.Outputs {
			//PRINT THE TEXT OF THE EVENT TO THE CONNECTION
			fmt.Fprintf(conn, event.Text)
		}
		//ARRIVE HERE WHEN THE OUTPUTS CHANNEL HAS BEEN CLOSED
		//CLOSE THE CONNECTION
		conn.Close()
		//ANNOUNCE THAT THE PLAYER HAS DISCONNECTED
		fmt.Println("Connection closed at", conn.RemoteAddr())
		fmt.Println("Output goroutine closing for", conn.RemoteAddr())
		//EXIT FROM THIS GOROUTINE
		return
	}()

	//GET INPUTS FROM THE CONNETION USING THE SCANNER
	for scanner.Scan() {
		//READ A LINE FROM THE USER
		line := scanner.Text()
		//CONVERT THE LINE TO A SLICE OF STRINGS
		lineList := strings.Fields(line)
		//CHECK IF THE YHIT ENTER WITH NO COMMANDS
		if len(lineList) == 0 {
			fmt.Fprintln(conn, "You must enter a command")
			//CHECK HERE TO ENSURE THE COMMAND THEY TYPED IS AN AVAILABLE COMMAND
		} else if _, ok := commands[lineList[0]]; ok == true {
			//CREATE THE INPUT EVENT THE USER HAS TYPED, SEND IT TO THE MAIN LOOP THROUGH THE INPUT EVENT CHANNEL
			inputEvent := InputEvent{p, lineList, false}
			inevents <- inputEvent
		} else {
			p.Printf("Huh?\n")
		}
	}
	//SCANNER HAS BEEN CLOSED FOR THE CONNECTION, ANNOUNCE SUCH
	fmt.Println("(Goroutine closing) End of scanner for", conn.RemoteAddr())
	//CREATE AN EVENT THAT WILL CAUSE THE CONNECTION TO QUIT,SEND IT THROUGH THE INPUT EVENTS CHANNEL
	exitgame := InputEvent{p, strings.Fields("quit"), false}
	inevents <- exitgame
	return

}
func main() {
	//ADD COMMANDS
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

	addCommand("gossip", gossip)
	addCommand("say", say)
	addCommand("tell", tell)
	addCommand("shout", shout)
	addCommand("where", where)

	addCommand("quit", quit)
	fmt.Println("COMMANDS INSTALLED")

	//LOADING THE WORLD FROM THE DATABASE

	fmt.Println("READING WORLD FILE")
	db, err := sql.Open("sqlite3", "./world.db")

	//ZONES TRANSACTION

	zonetx, err := db.Begin()
	zones := ReadZones(zonetx)
	if err != nil {
		_ = zonetx.Rollback()
		log.Fatalf("Failure at zone reading %v", err)
	}
	zonetx.Commit()

	//ROOMS TRANSACTION

	roomstx, err := db.Begin()
	rooms = ReadRooms(roomstx, zones)
	if err != nil {
		_ = roomstx.Rollback()
		log.Fatalf("failed reading a room from the database: %v", err)
	}
	roomstx.Commit()
	fmt.Println("Read", len(zones), "zones &", len(rooms), "rooms\n")

	//EXITS TRANSACTION

	exitstx, err := db.Begin()
	ReadExits(exitstx, rooms)
	if err != nil {
		_ = exitstx.Rollback()
		log.Fatalf("Error reading an exit from the database: %v", err)
	}
	exitstx.Commit()

	playertx, err := db.Begin()
	players = ReadPlayers(playertx)
	if err != nil {
		_ = playertx.Rollback()
		log.Fatalf("Failed to read the players from the database: %v", err)
	}
	playertx.Commit()

	//CREATE THE MAIN EVENT CHANNEL

	inputEvents := make(chan InputEvent)

	//INITIATE CONNECTION LISTENER
	go acceptConnections(inputEvents)

	//MAIN EVENT LOOP

	for event := range inputEvents {
		//MAKE SURE THE PLAYERS OUTPUT CHANNEL HASNT BEEN CLOSED(WE SET THE CHANNEL TO NIL DIRECTLY AFTER CLOSING IT)
		if event.player.Outputs != nil {
			//CHECK IF THIS IS A LOGIN EVENT
			if event.Login {
				//MAKE THE PLAYER ENTRY TRANSACTION
				createPlayerTX, err := db.Begin()
				//CALL THE MAKE PLAYER ENTRY FUNCTION
				makePlayerEntry(createPlayerTX, event.command[0], event.command[1], event.command[2], event.command[3])
				//CHECK THAT ENTER PlAYER WAS SUCCESSFUL
				if err != nil {
					_ = createPlayerTX.Rollback()
					log.Fatalf("Failed to enter a player into the database: %v", err)
				}
				//COMMIT THE TRANSACTION
				createPlayerTX.Commit()
				fmt.Println("Inserted new player", event.player.Name)

				//NOT A LOGIN EVENT
			} else {
				//RUN THE COMMAND THAT WE PULL FROM THE COMMANDS MAP,USING THE COMMAND FROM THE INPUT EVENT
				//AND THE PLAYER FROM THE INPUT EVENT AS PARAMETERS
				commands[event.command[0]](strings.Join(event.command, " "), event.player)

			}
		}
	}

}
