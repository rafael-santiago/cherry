/*
Package main.
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
*/
package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"pkg/config"
	"pkg/config/parser"
	"pkg/html"
	"pkg/messageplexer"
	"pkg/reqtraps"
	"strconv"
	"strings"
	"syscall"
)

const cherryVersion = "1.2"

func processNewConnection(newConn net.Conn, roomName string, rooms *config.CherryRooms) {
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

func getListenPeer(c *config.CherryRooms, port string) (net.Listener, error) {
	var listenConn net.Listener
	var listenError error
	if c.GetCertificatePath() != "" && c.GetPrivateKeyPath() != "" {
		cert, err := tls.LoadX509KeyPair(c.GetCertificatePath(), c.GetPrivateKeyPath())
		if err != nil {
			return nil, err
		}
		secParams := &tls.Config{Certificates: []tls.Certificate{cert}}
		listenConn, listenError = tls.Listen("tcp", c.GetServerName()+":"+port, secParams)
	} else {
		listenConn, listenError = net.Listen("tcp", c.GetServerName()+":"+port)
	}
	return listenConn, listenError
}

func peer(roomName string, c *config.CherryRooms) {
	port := c.GetListenPort(roomName)
	var portNum int64
	portNum, _ = strconv.ParseInt(port, 10, 16)
	var err error
	var room *config.RoomConfig
	room = c.GetRoomByPort(int16(portNum))
	//room.MainPeer, err = net.Listen("tcp", c.GetServerName()+":"+port)
	room.MainPeer, err = getListenPeer(c, port)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(1)
	}
	defer room.MainPeer.Close()
	for {
		conn, err := room.MainPeer.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go processNewConnection(conn, roomName, c)
	}
}

func getOption(option, defaultValue string, flagOption ...bool) string {
	isFlagOption := false
	if len(flagOption) > 0 {
		isFlagOption = flagOption[0]
	}
	for _, arg := range os.Args {
		if !isFlagOption {
			if strings.HasPrefix(arg, "--"+option+"=") {
				return arg[len(option)+3:]
			}
		} else if strings.HasPrefix(arg, "--"+option) {
			return "1"
		}
	}
	return defaultValue
}

func cleanup() {
	fmt.Println("INFO: Aborting signal received. Now your Cherry tree is being uprooted...  ;) Goodbye!!")
}

func announceVersion() {
	fmt.Println("cherry-" + cherryVersion)
}

func offerHelp() {
	fmt.Println("usage: cherry [--config=<cherry config filepath> | --help | --version]")
}

func openRooms(configPath string) {
	var cherryRooms *config.CherryRooms
	var err *parser.CherryFileError
	cherryRooms, err = parser.ParseCherryFile(configPath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		rooms := cherryRooms.GetRooms()
		for _, r := range rooms {
			go messageplexer.RoomMessagePlexer(r, cherryRooms)
			go peer(r, cherryRooms)
		}
	}
	sigintWatchdog := make(chan os.Signal, 1)
	signal.Notify(sigintWatchdog, os.Interrupt)
	signal.Notify(sigintWatchdog, syscall.SIGINT|syscall.SIGTERM)
	<-sigintWatchdog
	cleanup()
}

func main() {
	versionInfo := getOption("version", "", true)
	if len(versionInfo) > 0 {
		announceVersion()
		os.Exit(0)
	}
	help := getOption("help", "", true)
	if len(help) > 0 {
		offerHelp()
		os.Exit(0)
	}
	configPath := getOption("config", "")
	if len(configPath) == 0 {
		offerHelp()
		os.Exit(1)
	}
	openRooms(configPath)
}
