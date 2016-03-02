/*
Package... errr... hum... guess what?!
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
*/
package main

import (
	"./config"
	"./config/parser"
	"./html"
	"./messageplexer"
	"./reqtraps"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// ProcessNewConnection handles the possible request with the specifc trap returned based on what is being requested.
func ProcessNewConnection(newConn net.Conn, roomName string, rooms *config.CherryRooms) {
	buf := make([]byte, 4096)
	bufLen, err := newConn.Read(buf)
	if err == nil {
		preprocessor := html.NewHTMLPreprocessor(rooms)
		httpPayload := string(buf[:bufLen])
		var trap reqtraps.RequestTrap
		trap = reqtraps.GetRequestTrap(httpPayload)
		trap().Handle(newConn, roomName, httpPayload, rooms, preprocessor)
	} else {
		newConn.Close()
	}
}

// Peer is the room listener.
func Peer(roomName string, c *config.CherryRooms) {
	port := c.GetListenPort(roomName)
	var portNum int64
	portNum, _ = strconv.ParseInt(port, 10, 16)
	var err error
	var room *config.RoomConfig
	room = c.GetRoomByPort(int16(portNum))
	room.MainPeer, err = net.Listen("tcp", c.GetServerName()+":"+port)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	defer room.MainPeer.Close()
	for {
		conn, err := room.MainPeer.Accept()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		go ProcessNewConnection(conn, roomName, c)
	}
}

// GetOption handles the command line options.
func GetOption(option, defaultValue string) string {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "--"+option+"=") {
			return arg[len(option)+3:]
		}
	}
	return defaultValue
}

func main() {
	var cherryRooms *config.CherryRooms
	var err *parser.CherryFileError
	var configPath string
	configPath = GetOption("config", "")
	if len(configPath) == 0 {
		fmt.Println("ERROR: --config option is missing.")
		os.Exit(1)
	}
	cherryRooms, err = parser.ParseCherryFile(configPath)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		rooms := cherryRooms.GetRooms()
		for ri, r := range rooms {
			go messageplexer.RoomMessagePlexer(r, cherryRooms)
			if ri < len(rooms)-1 {
				go Peer(r, cherryRooms)
			} else {
				Peer(r, cherryRooms)
			}
		}
	}
}
