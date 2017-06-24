/*
Package messageplexer is the part responsible for messages delivering on Cherry.
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
*/
package messageplexer

import (
	"net"
	"github.com/rafael-santiago/cherry/src/pkg/config"
	"github.com/rafael-santiago/cherry/src/pkg/html"
	"runtime"
)

// RoomMessagePlexer performs all message delivering stuff.
func RoomMessagePlexer(roomName string, rooms *config.CherryRooms) {
	preprocessor := html.NewHTMLPreprocessor(rooms)
	var allUsers = rooms.GetAllUsersAlias(roomName)
	for {
		runtime.Gosched()
		currMessage := rooms.GetNextMessage(roomName)
		if len(currMessage.Say) == 0 && len(currMessage.Image) == 0 /*&& len(currMessage.Sound) == 0*/ {
			continue
		}
		var actionTemplate string
		if rooms.HasAction(roomName, currMessage.Action) {
			actionTemplate = rooms.GetRoomActionTemplate(roomName, currMessage.Action)
		}
		if len(actionTemplate) == 0 {
			actionTemplate = "<p>({{.hour}}:{{.minute}}:{{.second}}) <b>{{.message-colored-user}}</b>: {{.message-says}}" //  INFO(Santiago): A very basic action template.
		}
		message := preprocessor.ExpandData(roomName, actionTemplate)
		if currMessage.Priv != "1" {
			rooms.AddPublicMessage(roomName, message)
		}
		preprocessor.SetDataValue("{{.current-formatted-message}}", message)
		messageHighlighted := preprocessor.ExpandData(roomName, rooms.GetHighlightTemplate(roomName))
		preprocessor.UnsetDataValue("{{.current-formatted-message}}")
		users := rooms.GetRoomUsers(roomName)
		for _, user := range users {
			if currMessage.Priv == "1" &&
				user != currMessage.From &&
				user != currMessage.To &&
				currMessage.To != allUsers {
				continue
			}
			if rooms.IsIgnored(user, currMessage.From, roomName) {
				continue
			}
			var messageBuffer []byte
			if user == currMessage.From ||
				user == currMessage.To {
				messageBuffer = []byte(messageHighlighted)
			} else {
				messageBuffer = []byte(message)
			}
			var conn net.Conn
			conn = rooms.GetUserConnection(roomName, user)
			if conn == nil {
				continue
			}
			_, e := conn.Write(messageBuffer)
			if e != nil {
				rooms.EnqueueMessage(roomName, user, "", "", "", rooms.GetExitMessage(roomName), "")
				rooms.RemoveUser(roomName, user)
			}
		}
		rooms.DequeueMessage(roomName)
	}
}
