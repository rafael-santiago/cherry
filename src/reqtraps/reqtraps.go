package reqtraps

import (
    "net"
    "../config"
    "../html"
    "../rawhttp"
    "strings"
)

type RequestTrapInterface interface {
    Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor)
}

type RequestTrapHandleFunc func(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor)

func (h RequestTrapHandleFunc) Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    h(new_conn, room_name, http_payload, rooms, preprocessor)
}

type RequestTrap func() RequestTrapInterface

func BuildRequestTrap(handle RequestTrapHandleFunc) RequestTrap {
    return func () RequestTrapInterface {
        return RequestTrapHandleFunc(handle)
    }
}

func GetRequestTrap(http_payload string) RequestTrap {
    var http_method_part string
    var space_nr int = 0
    for _, h := range http_payload {
        if h == ' ' {
            space_nr++
        }
        if h == '\n' || h == '\r' || space_nr == 2 {
            break
        }
        http_method_part += string(h)
    }
    http_method_part += "$"
    if strings.HasPrefix(http_method_part, "GET /join$") {
        return BuildRequestTrap(GetJoin_Handle)
    }
    //if strings.HasPrefix(http_method_part, "GET /brief$") {
    //    return nil //  TODO(Santiago): Return the brief handle for this room.
    //}
    if strings.HasPrefix(http_method_part, "GET /top&") {
        return BuildRequestTrap(GetTop_Handle)
    }
    if strings.HasPrefix(http_method_part, "GET /banner&") {
        return BuildRequestTrap(GetBanner_Handle)
    }
    if strings.HasPrefix(http_method_part, "GET /body&") {
        return BuildRequestTrap(GetBody_Handle)
    }
    if strings.HasPrefix(http_method_part, "GET /exit&") {
        return BuildRequestTrap(GetExit_Handle)
    }
    if strings.HasPrefix(http_method_part, "POST /join$") {
        return BuildRequestTrap(PostJoin_Handle)
    }
    if strings.HasPrefix(http_method_part, "POST /banner&") {
        return BuildRequestTrap(PostBanner_Handle) //  TODO(Santiago): Enqueue user's message for further processing.
    }
    //if strins.HasPrefix(http_method_part, "POST /query&") {
    //    return BuildRequestTrap(PostQuery_Handle) //  TODO(Santiago): Return the search results.
    //}
    return BuildRequestTrap(BadAssError_Handle)
}

func GetJoin_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    //  INFO(Santiago): The form for room joining was requested, so we will flush it to client.
    var reply_buffer []byte
    reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetEntranceTemplate(room_name)), 200, true)
    new_conn.Write(reply_buffer)
    new_conn.Close()
}

func GetTop_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var user_data map[string]string
    user_data = rawhttp.GetFieldsFromGet(http_payload)
    var reply_buffer []byte
    if !rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"]) {
        reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetTopTemplate(room_name)), 200, true)
    }
    new_conn.Write(reply_buffer)
    new_conn.Close()
}

func GetBanner_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var user_data map[string]string
    var reply_buffer []byte
    user_data = rawhttp.GetFieldsFromGet(http_payload)
    if !rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"]) {
        reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetBannerTemplate(room_name)), 200, true)
    }
    new_conn.Write(reply_buffer)
    new_conn.Close()
}

func GetExit_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var user_data map[string]string
    var reply_buffer []byte
    user_data = rawhttp.GetFieldsFromGet(http_payload)
    if !rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"]) {
        reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        preprocessor.SetDataValue("{{.nickname}}", user_data["user"])
        preprocessor.SetDataValue("{{.session-id}}", user_data["id"])
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetExitTemplate(room_name)), 200, true)
    }
    new_conn.Write(reply_buffer)
    rooms.RemoveUser(room_name, user_data["user"])
    new_conn.Close()
}

func PostJoin_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    //  INFO(Santiago): Here, we need firstly parse the posted fields, check for "nickclash", if this is the case
    //                  flush the page informing it. Otherwise we add the user basic info and flush the room skeleton
    //                  [TOP/BODY/BANNER]. Then we finally close the connection.
    var user_data map[string]string
    var reply_buffer []byte
    user_data = rawhttp.GetFieldsFromPost(http_payload)
    if _, posted := user_data["user"]; !posted {
        new_conn.Close()
        return
    }
    if _, posted := user_data["color"]; !posted {
        new_conn.Close()
        return
    }
    if _, posted := user_data["says"]; !posted {
        new_conn.Close()
        return
    }
    preprocessor.SetDataValue("{{.nickname}}", user_data["user"])
    preprocessor.SetDataValue("{{.session-id}}", "0")
    if rooms.HasUser(room_name, user_data["user"]) {
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetNickclashTemplate(room_name)), 200, true)
    } else {
        rooms.AddUser(room_name, user_data["user"], user_data["color"], true)
        preprocessor.SetDataValue("{{.session-id}}", rooms.GetSessionId(user_data["user"], room_name))
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetSkeletonTemplate(room_name)), 200, true)
        //  INFO(Santiago): At this point the others and this user will get the join notification of him.
        //                  Yes, he/she could "hack" the join notification message just for fun :^)
        rooms.EnqueueMessage(room_name, user_data["user"], "", "", "", "", user_data["says"], "")
    }
    new_conn.Write(reply_buffer)
    new_conn.Close()
}

func GetBody_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var user_data map[string]string
    user_data = rawhttp.GetFieldsFromGet(http_payload)
    var valid_user bool
    valid_user = rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"])
    var reply_buffer []byte
    if !valid_user {
        reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetBodyTemplate(room_name)), 200, false)
    }
    new_conn.Write(reply_buffer)
    if valid_user {
        rooms.SetUserConnection(room_name, user_data["user"], new_conn)
    } else {
        new_conn.Close()
    }
}

func BadAssError_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    new_conn.Write(rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true))
    new_conn.Close()
}

func PostBanner_Handle(new_conn net.Conn, room_name, http_payload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var user_data map[string]string
    var reply_buffer []byte
    var invalid_request bool = false
    user_data = rawhttp.GetFieldsFromPost(http_payload)
    if _ , has := user_data["user"]; !has {
        invalid_request = true
    } else if _, has := user_data["id"]; !has {
        invalid_request = true
    } else if _, has := user_data["action"]; !has {
        invalid_request = true
    } else if _, has := user_data["whoto"]; !has {
        invalid_request = true
    } else if _, has := user_data["sound"]; !has {
        invalid_request = true
    } else if  _, has := user_data["image"]; !has {
        invalid_request = true
    } else if _, has := user_data["says"]; !has {
        invalid_request = true
    }
    if invalid_request || !rooms.IsValidUserRequest(room_name, user_data["user"], user_data["id"]) {
        reply_buffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        reply_buffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(room_name, rooms.GetBannerTemplate(room_name)), 200, true)
        var something_to_say bool =  (len(user_data["says"]) > 0 || len(user_data["image"]) > 0 || len(user_data["sound"]) > 0)
        if something_to_say {
            //  INFO(Santiago): Any further antiflood control would go from here.
            rooms.EnqueueMessage(room_name, user_data["user"], user_data["whoto"], user_data["action"], user_data["sound"], user_data["image"], user_data["says"], user_data["priv"])
        }
    }
    new_conn.Write(reply_buffer)
    new_conn.Close()
}
