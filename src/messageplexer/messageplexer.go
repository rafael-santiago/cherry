package messageplexer

import (
    "../config"
    "../html"
    "net"
)

func RoomMessagePlexer(room_name string, rooms *config.CherryRooms) {
    preprocessor := html.NewHtmlPreprocessor(rooms)
    var all_users string = rooms.GetAllUsersAlias(room_name)
    for {
        curr_message := rooms.GetNextMessage(room_name)
        if len(curr_message.Say) == 0 && len(curr_message.Image) == 0 /*&& len(curr_message.Sound) == 0*/ {
            continue
        }
        var action_template string
        if rooms.HasAction(room_name, curr_message.Action) {
            action_template = rooms.GetRoomActionTemplate(room_name, curr_message.Action)
        }
        if len(action_template) == 0 {
            action_template = "<p>({{.hour}}:{{.minute}}:{{.second}}) <b>{{.message-colored-user}}</b>: {{.message-says}}" //  INFO(Santiago): A very basic action template.
        }
        message := preprocessor.ExpandData(room_name, action_template)
        if curr_message.Priv != "1" {
            rooms.AddPublicMessage(room_name, message)
        }
        preprocessor.SetDataValue("{{.current-formatted-message}}", message)
        message_highlighted := preprocessor.ExpandData(room_name, rooms.GetHighlightTemplate(room_name))
        preprocessor.UnsetDataValue("{{.current-formatted-message}}")
        users := rooms.GetRoomUsers(room_name)
        for _, user := range users {
            if curr_message.Priv == "1" &&
               user != curr_message.From &&
               user != curr_message.To &&
               curr_message.To != all_users {
                continue
            }
            if rooms.IsIgnored(user, curr_message.From, room_name) {
                continue
            }
            var message_buffer []byte
            if user == curr_message.From ||
               user == curr_message.To {
                message_buffer = []byte(message_highlighted)
            } else {
                message_buffer = []byte(message)
            }
            var conn net.Conn
            conn = rooms.GetUserConnection(room_name, user)
            if conn == nil {
                continue
            }
            _, e := conn.Write(message_buffer)
            if e != nil {
                rooms.EnqueueMessage(room_name, user, "", "", "", rooms.GetExitMessage(room_name), "")
                rooms.RemoveUser(room_name, user)
            }
        }
        rooms.DequeueMessage(room_name)
    }
}
