package messageplexer

import (
    "../config"
    "../html"
    "net"
)

func RoomMessagePlexer(room_name string, rooms *config.CherryRooms) {
    preprocessor := html.NewHtmlPreprocessor(rooms)
    for {
        curr_message := rooms.GetNextMessage(room_name)
        if len(curr_message.Say) == 0 {
            continue
        }
        var action_template string
        if rooms.HasAction(room_name, curr_message.Action) {
            action_template = rooms.GetRoomActionTemplate(room_name, curr_message.Action)
        }
        message := preprocessor.ExpandData(room_name, action_template)
        users := rooms.GetRoomUsers(room_name)
        for _, user := range users {
            var conn net.Conn
            conn = rooms.GetUserConnection(room_name, user)
            if conn == nil {
                continue
            }
            conn.Write([]byte(message))
        }
        rooms.DequeueMessage(room_name)
    }
}
