/*
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
    "fmt"
    "os"
    "net"
    "strconv"
    "./html"
    "./reqtraps"
    "./messageplexer"
    "strings"
)

func ProcessNewConnection(new_conn net.Conn, room_name string, rooms *config.CherryRooms) {
    buf := make([]byte, 4096)
    buf_len, err := new_conn.Read(buf)
    if err == nil {
        preprocessor := html.NewHtmlPreprocessor(rooms)
        http_payload := string(buf[:buf_len])
        var trap reqtraps.RequestTrap
        trap = reqtraps.GetRequestTrap(http_payload)
        trap().Handle(new_conn, room_name, http_payload, rooms, preprocessor)
    } else {
        new_conn.Close()
    }
}

func MainPeer(room_name string, c *config.CherryRooms) {
    port := c.GetListenPort(room_name)
    var port_num int64
    port_num, _ = strconv.ParseInt(port, 10, 16)
    var err error
    var room *config.RoomConfig
    room = c.GetRoomByPort(int16(port_num))
    room.MainPeer, err = net.Listen("tcp", c.GetServerName() + ":" + port)
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
        go ProcessNewConnection(conn, room_name, c)
    }
}

func GetOption(option, default_value string) string {
    for _, arg := range os.Args {
        if strings.HasPrefix(arg, "--" + option + "=") {
            return arg[len(option) + 3:]
        }
    }
    return default_value
}

func main() {
    var cherry_rooms *config.CherryRooms
    var err *parser.CherryFileError
    var config_path string
    config_path = GetOption("config", "")
    if len(config_path) == 0 {
        fmt.Println("ERROR: --config option is missing.")
        os.Exit(1)
    }
    cherry_rooms, err = parser.ParseCherryFile(config_path)
    if err != nil {
        fmt.Println(err.Error())
    } else {
        rooms := cherry_rooms.GetRooms()
        for _, r := range rooms {
            go messageplexer.RoomMessagePlexer(r, cherry_rooms)
        }
        MainPeer("foobar", cherry_rooms)
    }
}
