package reqtraps

import (
    "net"
    "../config"
    "../html"
    "../rawhttp"
)

type RequestTrapInterface interface {
    Build(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms)
    Handle()
}

type RequestTrapBase struct {
    new_conn net.Conn
    room_name string
    http_payload string
    rooms *config.CherryRooms
    preprocessor *html.Preprocessor
    reply_buffer []byte
}

func (t *RequestTrapBase) Build(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms) {
    t.new_conn = new_conn
    t.room_name = room_name
    t.http_payload = http_payload
    t.rooms = rooms
    t.preprocessor = html.NewHtmlPreprocessor(t.rooms)
}

func (t *RequestTrapBase) Handle() {
    t.reply_buffer = rawhttp.MakeReplyBuffer("<html></html>", 200, true)
    t.write()
}

func (t *RequestTrapBase) write() {
    t.new_conn.Write(t.reply_buffer)
    t.new_conn.Close()
}

type GetJoinRequestTrap struct {
    RequestTrapBase
}

func (t *GetJoinRequestTrap) Handle() {
    //  INFO(Santiago): The form for room joining was requested, so we will flush it to client.
    t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetEntranceTemplate(t.room_name)), 200, true)
    t.write()
}

type GetTopRequestTrap struct {
    RequestTrapBase
}

func (t *GetTopRequestTrap) Handle() {
    user_data := rawhttp.GetFieldsFromGet(t.http_payload)
    if !t.rooms.IsValidUserRequest(t.room_name, user_data["user"], user_data["id"]) {
        t.reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetTopTemplate(t.room_name)), 200, true)
    }
    t.write()
}

type GetBannerRequestTrap struct {
    RequestTrapBase
}

func (t *GetBannerRequestTrap) Handle() {
    user_data := rawhttp.GetFieldsFromGet(t.http_payload)
    if !t.rooms.IsValidUserRequest(t.room_name, user_data["user"], user_data["id"]) {
        t.reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetBannerTemplate(t.room_name)), 200, true)
    }
    t.write()
}

type GetExitRequestTrap struct {
    RequestTrapBase
    user_data map[string]string
}

func (t *GetExitRequestTrap) Handle() {
    t.user_data = rawhttp.GetFieldsFromGet(t.http_payload)
    if !t.rooms.IsValidUserRequest(t.room_name, t.user_data["user"], t.user_data["id"]) {
        t.reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        t.preprocessor.SetDataValue("{{.nickname}}", t.user_data["user"])
        t.preprocessor.SetDataValue("{{.session-id}}", t.user_data["id"])
        t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetExitTemplate(t.room_name)), 200, true)
    }
    t.write()
}

func (t *GetExitRequestTrap) write() {
    t.new_conn.Write(t.reply_buffer)
    t.rooms.RemoveUser(t.room_name, t.user_data["user"])
    t.new_conn.Close()
}

type PostJoinRequestTrap struct {
    RequestTrapBase
}

func (t *PostJoinRequestTrap) Handle() {
    //  INFO(Santiago): Here, we need firstly parse the posted fields, check for "nickclash", if this is the case
    //                  flush the page informing it. Otherwise we add the user basic info and flush the room skeleton
    //                  [TOP/BODY/BANNER]. Then we finally close the connection.
    user_data := rawhttp.GetFieldsFromPost(t.http_payload)
    if _, posted := user_data["user"]; !posted {
        t.new_conn.Close()
    }
    if _, posted := user_data["color"]; !posted {
        t.new_conn.Close()
    }
    if _, posted := user_data["says"]; !posted {
        t.new_conn.Close()
    }
    t.preprocessor.SetDataValue("{{.nickname}}", user_data["user"])
    t.preprocessor.SetDataValue("{{.session-id}}", "0")
    if t.rooms.HasUser(t.room_name, user_data["user"]) {
        t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetNickclashTemplate(t.room_name)), 200, true)
    } else {
        t.rooms.AddUser(t.room_name, user_data["user"], user_data["color"], true)
        t.preprocessor.SetDataValue("{{.session-id}}", t.rooms.GetSessionId(user_data["user"], t.room_name))
        t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetSkeletonTemplate(t.room_name)), 200, true)
        //  INFO(Santiago): At this point the others and this user will get the join notification of him.
        //                  Yes, he/she could "hack" the join notification message for fun :^)
        t.rooms.EnqueueMessage(t.room_name, user_data["user"], "", "", "", "", user_data["says"], "")
    }
    t.write()
}

type GetBodyRequestTrap struct {
    RequestTrapBase
    valid_user bool
    user_data map[string]string
}

func (t *GetBodyRequestTrap) Handle() {
    t.user_data = rawhttp.GetFieldsFromGet(t.http_payload)
    t.valid_user = t.rooms.IsValidUserRequest(t.room_name, t.user_data["user"], t.user_data["id"])
    if !t.valid_user {
        t.reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        t.reply_buffer = rawhttp.MakeReplyBuffer(t.preprocessor.ExpandData(t.room_name, t.rooms.GetBodyTemplate(t.room_name)), 200, false)
    }
    t.write()
}

func (t *GetBodyRequestTrap) write() {
    t.new_conn.Write(t.reply_buffer)
    if t.valid_user {
        t.rooms.SetUserConnection(t.room_name, t.user_data["user"], t.new_conn) //  INFO(Santiago): This connection will closed only on exit request.
    } else {
        t.new_conn.Close()
    }
}
