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
    "strings"
)

func ProcessNewConnection(new_conn net.Conn, room *config.RoomConfig) {
    buf := make([]byte, 4096)
    buf_len, err := new_conn.Read(buf)
    var http_payload string
    if err == nil {
        http_payload = string(buf[:buf_len])
        if strings.HasPrefix(http_payload, "GET /join") {
            //  TODO(Santiago): Return the join form.
        } else if strings.HasPrefix(http_payload, "GET /brief") {
            //  TODO(Santiago): Return the brief for this room.
        } else if strings.HasPrefix(http_payload, "GET /top") {
            //  TODO(Santiago): Return the room's top frame.
        } else if strings.HasPrefix(http_payload, "GET /banner") {
            //  TODO(Santiago): Return the room's banner frame.
        } else if strings.HasPrefix(http_payload, "GET /body") {
            //  TODO(Santiago): Return the room's body frame and do not close this connection.
        } else if strings.HasPrefix(http_payload, "POST /exit") {
            //  TODO(Santiago): Clear user's information from room's context and return to him the exit document.
        } else if strings.HasPrefix(http_payload, "POST /join") {
            //  TODO(Santiago): Process the join request, adding (if allowed) the user inside room's context, and
            //                  return the room's basic struct [TOP/BODY/BANNER]
        } else if strings.HasPrefix(http_payload, "POST /banner") {
            //  TODO(Santiago): Enqueue user's message for future processing.
        } else if strings.HasPrefix(http_payload, "POST /query") {
            //  TODO(Santiago): Return the search results.
        }
    }
}

func MainPeer(room_name string, c *config.CherryRooms) {
    port := c.GetListenPort(room_name)
    var port_num int64
    port_num, _ = strconv.ParseInt(port, 10, 16)
    var err error
    var room *config.RoomConfig
    room = c.GetRoomByPort(int16(port_num))
    room.MainPeer, err = net.Listen("tcp", ":" + port)
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
        //  TODO(Santiago): Process the user's initial request.
        go ProcessNewConnection(conn, room)
    }
}

func main() {
    var cherry_rooms *config.CherryRooms
    var err *parser.CherryFileError
    cherry_rooms, err = parser.ParseCherryFile("config.cherry")
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("*** Configuration loaded!")
        MainPeer("foobar", cherry_rooms)
    }
}
