package chat

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

const separator string = ", "

// >> string: command name, int: arg count
var commands = map[string]int{
	"JOIN":   2,
	"SWITCH": 1,
	"NAME":   1,
}

type user struct {
	UserIP      string
	Name        string
	CurrentRoom int
	Connection  net.Conn
}

// >> sets the username
func (u *user) setName(name string) {
	u.Name = name
}

// >> Lets user join room
// note that the address is for the server
// and should be gotten from the CLI
func (u *user) join(address string, roomID int) error {
	var err error
	if roomID < 0 || roomID > maxRooms {
		err = fmt.Errorf("error: user attempted to join room id: %d\nmax: %d", roomID, maxRooms)
		log.Println(err)
		return err
	}

	u.Connection, err = net.Dial("tcp4", address)
	if err != nil {
		log.Printf("error: failed to join %s\n", address)
		log.Println(err)
		return err
	}

	return nil
}

// >> Gets the room ID
// if raw room ID is given, then simply returns an int of that
// otherwise if not numeric (room name is provided by user):
// 		- client asks server if provided name is available
//		- If yes: server replyes with room ID
//		- If no:  server replyes with error code 401
func (u *user) getRoomID(argString string) (int, error) {
	room, err := strconv.Atoi(argString)
	if err != nil {
		var rErr error
		room, rErr = u.ResolveRoomName(argString)
		if rErr != nil {
			return -1, rErr
		}
	}
	return room, nil
}

// >> Checks if the required argument count is correct
func checkArgs(argCount int, actual int) error {
	if argCount != actual {
		return fmt.Errorf("error: command arg count incorrect\nExpected: %d, Actual: %d", len(args), actual)
	}
	return nil
}

// >> Handles user input
func (u *user) handleUserInput(command string) error {
	// >> Splitting the command into it's components
	args := strings.Split(command, separator)
	argCount := len(args)

	// >> If command is actually a message
	if args[0][0] != ':' {
		u.sendMessage(command)
	}

	// >> checking that arguments provided are the same as expected for the command
	commandType := args[0][:1]
	err := checkArgs(argCount, commands[commandType])
	if err != nil {
		return err
	}

	// >> Executing command depending on types
	switch commandType {
	case "JOIN":
		room, err := u.getRoomID(args[2])
		if err != nil {
			log.Println("error: can't find room %s on server", args[2])
			return err
		}
		u.join(args[1], room)
		break

	case "SWITCH":
		room, err := u.getRoomID(args[2])
		if err != nil {
			log.Println("error: can't find room %s on server", args[1])
			return err
		}
		u.join(u.Connection.LocalAddr().String(), room)
		break

	case "NAME":
		u.setName(args[1])
		break

	}

	return nil
}
