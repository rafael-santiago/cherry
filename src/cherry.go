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
    "./rawhttp"
    "./html"
)

func ProcessNewConnection(new_conn net.Conn, room_name string, rooms *config.CherryRooms) {
    buf := make([]byte, 4096)
    buf_len, err := new_conn.Read(buf)
    var http_payload string
    var reply_buffer []byte
    preprocessor := html.NewHtmlPreprocessor(rooms)
    if err == nil {
        http_payload = string(buf[:buf_len])
        if strings.HasPrefix(http_payload, "GET /join") {
            //  INFO(Santiago): The form for room joining was requested, so we will flush it to the client.
            reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetEntranceTemplate(room_name)),
                                                   200,
                                                   true)
            new_conn.Write(reply_buffer)
            new_conn.Close()
        } else if strings.HasPrefix(http_payload, "GET /brief") {
            //  TODO(Santiago): Return the brief for this room.
        } else if strings.HasPrefix(http_payload, "GET /top") {
            //  TODO(Santiago): Return the room's top frame.
            user_data := rawhttp.GetFieldsFromGet(http_payload)
            if !rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"]) {
                reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
            } else {
                reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetTopTemplate(room_name)),
                                                       200,
                                                       true)
            }
            new_conn.Write(reply_buffer)
            new_conn.Close()
        } else if strings.HasPrefix(http_payload, "GET /banner") {
            user_data := rawhttp.GetFieldsFromGet(http_payload)
            if !rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"]) {
                reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
            } else {
                reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetBannerTemplate(room_name)),
                                                       200,
                                                       true)
            }
            new_conn.Write(reply_buffer)
            new_conn.Close()
        } else if strings.HasPrefix(http_payload, "GET /body") {
            //  TODO(Santiago): Return the room's body frame and do not close this connection.
        } else if strings.HasPrefix(http_payload, "POST /exit") {
            //  TODO(Santiago): Clear user's information from room's context and return to him the exit document.
        } else if strings.HasPrefix(http_payload, "POST /join") {
            //  INFO(Santiago): Here, we need firstly parse the posted fields, check for "nickclash", if this is the case
            //                  flush the page informing it. Otherwise we add the user basic info and flush the room skeleton
            //                  [TOP/BODY/BANNER]. Then we finally close the connection.
            user_data := rawhttp.GetFieldsFromPost(http_payload)
            if _, posted := user_data["user"]; !posted {
                new_conn.Close()
            }
            if _, posted := user_data["color"]; !posted {
                new_conn.Close()
            }
            if _, posted := user_data["says"]; !posted {
                new_conn.Close()
            }
            preprocessor.SetDataValue("{{.nickname}}", user_data["user"])
            preprocessor.SetDataValue("{{.session-id}}", "0")
            if rooms.HasUser(room_name, user_data["user"]) {
                reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetNickclashTemplate(room_name)), 200, true)
            } else {
                rooms.AddUser(room_name, user_data["user"], user_data["color"], false)
                preprocessor.SetDataValue("{{.session-id}}", rooms.GetSessionId(user_data["user"], room_name))
                reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetSkeletonTemplate(room_name)), 200, true)
                //  INFO(Santiago): At this point the others and this user will get the join notification of him.
                //                  Yes, he/she could "hack" the join notification message for fun :^)
                rooms.EnqueueMessage(room_name, user_data["user"], "", "", "", "", user_data["says"], "")
            }
            new_conn.Write(reply_buffer)
            new_conn.Close()
        } else if strings.HasPrefix(http_payload, "POST /banner") {
            //  TODO(Santiago): Enqueue user's message for future processing.
        } else if strings.HasPrefix(http_payload, "POST /query") {
            //  TODO(Santiago): Return the search results.
        } else {
            new_conn.Write(rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true))
            new_conn.Close()
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
    room.MainPeer, err = net.Listen("tcp", "192.30.70.3:" + port)
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
        go ProcessNewConnection(conn, room_name, c)
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
