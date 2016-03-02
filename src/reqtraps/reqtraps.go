/*
Package reqtraps implements the traps that handle each relevant HTTP method.
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package reqtraps

import (
    "net"
    "../config"
    "../html"
    "../rawhttp"
    "strings"
)

// RequestTrapInterface is used for what it suggests [Hi lint! you are so stupid!].
type RequestTrapInterface interface {
    Handle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor)
}

// RequestTrapHandleFunc idem.
type RequestTrapHandleFunc func(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor)

// Handle idem.
func (h RequestTrapHandleFunc) Handle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    h(newConn, roomName, httpPayload, rooms, preprocessor)
}

// RequestTrap idem.
type RequestTrap func() RequestTrapInterface

// BuildRequestTrap creates a trap based on a user request.
func BuildRequestTrap(handle RequestTrapHandleFunc) RequestTrap {
    return func () RequestTrapInterface {
        return RequestTrapHandleFunc(handle)
    }
}

// GetRequestTrap returns the correct trap that should be used to handle the user request.
func GetRequestTrap(httpPayload string) RequestTrap {
    var httpMethodPart string
    var spaceNr int
    for _, h := range httpPayload {
        if h == ' ' {
            spaceNr++
        }
        if h == '\n' || h == '\r' || spaceNr == 2 {
            break
        }
        httpMethodPart += string(h)
    }
    httpMethodPart += "$"
    if strings.HasPrefix(httpMethodPart, "GET /join$") {
        return BuildRequestTrap(GetJoinHandle)
    }
    if strings.HasPrefix(httpMethodPart, "GET /brief$") {
        return BuildRequestTrap(GetBriefHandle)
    }
    if strings.HasPrefix(httpMethodPart, "GET /top&") {
        return BuildRequestTrap(GetTopHandle)
    }
    if strings.HasPrefix(httpMethodPart, "GET /banner&") {
        return BuildRequestTrap(GetBannerHandle)
    }
    if strings.HasPrefix(httpMethodPart, "GET /body&") {
        return BuildRequestTrap(GetBodyHandle)
    }
    if strings.HasPrefix(httpMethodPart, "GET /exit&") {
        return BuildRequestTrap(GetExitHandle)
    }
    if strings.HasPrefix(httpMethodPart, "POST /join$") {
        return BuildRequestTrap(PostJoinHandle)
    }
    if strings.HasPrefix(httpMethodPart, "POST /banner&") {
        return BuildRequestTrap(PostBannerHandle)
    }
    if strings.HasPrefix(httpMethodPart, "GET /find$") {
        return BuildRequestTrap(GetFindHandle)
    }
    if strings.HasPrefix(httpMethodPart, "POST /find$") {
        return BuildRequestTrap(PostFindHandle)
    }
    return BuildRequestTrap(BadAssErrorHandle)
}

// GetFindHandle implements the handle for the find document (GET).
func GetFindHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var replyBuffer []byte
    replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetFindBotTemplate(roomName)), 200, true)
    newConn.Write(replyBuffer)
    newConn.Close()
}

// PostFindHandle implements the handle for the find document (POST).
func PostFindHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var userData map[string]string
    userData = rawhttp.GetFieldsFromPost(httpPayload)
    var replyBuffer []byte
    if _, posted := userData["user"]; !posted {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        var result string
        result = preprocessor.ExpandData(roomName, rooms.GetFindResultsHeadTemplate(roomName))
        listing := rooms.GetFindResultsBodyTemplate(roomName)
        availRooms := rooms.GetRooms()
        user := strings.ToUpper(userData["user"])
        if len(user) > 0 {
            for _, r := range availRooms {
                users := rooms.GetRoomUsers(r)
                preprocessor.SetDataValue("{{.find-result-users-total}}", rooms.GetUsersTotal(r))
                preprocessor.SetDataValue("{{.find-result-room-name}}", r)
                for _, u := range users {
                    if strings.HasPrefix(strings.ToUpper(u), user) {
                        preprocessor.SetDataValue("{{.find-result-user}}", u)
                        result += preprocessor.ExpandData(roomName, listing)
                    }
                }
            }
        }
        result += preprocessor.ExpandData(roomName, rooms.GetFindResultsTailTemplate(roomName))
        replyBuffer = rawhttp.MakeReplyBuffer(result, 200, true)
    }
    newConn.Write(replyBuffer)
    newConn.Close()
}

// GetJoinHandle implements the handle for the join document (GET).
func GetJoinHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    //  INFO(Santiago): The form for room joining was requested, so we will flush it to client.
    var replyBuffer []byte
    replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetEntranceTemplate(roomName)), 200, true)
    newConn.Write(replyBuffer)
    newConn.Close()
}

// GetTopHandle implements the handle for the top document (GET).
func GetTopHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var userData map[string]string
    userData = rawhttp.GetFieldsFromGet(httpPayload)
    var replyBuffer []byte
    if !rooms.IsValidUserRequest(roomName, userData["user"], userData["id"]) {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetTopTemplate(roomName)), 200, true)
    }
    newConn.Write(replyBuffer)
    newConn.Close()
}

// GetBannerHandle implements the handle for the banner document (GET).
func GetBannerHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var userData map[string]string
    var replyBuffer []byte
    userData = rawhttp.GetFieldsFromGet(httpPayload)
    preprocessor.SetDataValue("{{.nickname}}", userData["user"])
    preprocessor.SetDataValue("{{.session-id}}", userData["id"])
    if !rooms.IsValidUserRequest(roomName, userData["user"], userData["id"]) {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetBannerTemplate(roomName)), 200, true)
    }
    newConn.Write(replyBuffer)
    newConn.Close()
}

// GetExitHandle implements the handle for the exit document (GET).
func GetExitHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var userData map[string]string
    var replyBuffer []byte
    userData = rawhttp.GetFieldsFromGet(httpPayload)
    if !rooms.IsValidUserRequest(roomName, userData["user"], userData["id"]) {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        preprocessor.SetDataValue("{{.nickname}}", userData["user"])
        preprocessor.SetDataValue("{{.session-id}}", userData["id"])
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetExitTemplate(roomName)), 200, true)
    }
    rooms.EnqueueMessage(roomName, userData["user"], "", "", "",  rooms.GetExitMessage(roomName), "")
    newConn.Write(replyBuffer)
    rooms.RemoveUser(roomName, userData["user"])
    newConn.Close()
}

// PostJoinHandle implements the handle for the join document (POST).
func PostJoinHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    //  INFO(Santiago): Here, we need firstly parse the posted fields, check for "nickclash", if this is the case
    //                  flush the page informing it. Otherwise we add the user basic info and flush the room skeleton
    //                  [TOP/BODY/BANNER]. Then we finally close the connection.
    var userData map[string]string
    var replyBuffer []byte
    userData = rawhttp.GetFieldsFromPost(httpPayload)
    if _, posted := userData["user"]; !posted {
        newConn.Close()
        return
    }
    if _, posted := userData["color"]; !posted {
        newConn.Close()
        return
    }
    preprocessor.SetDataValue("{{.nickname}}", userData["user"])
    preprocessor.SetDataValue("{{.session-id}}", "0")
    if rooms.HasUser(roomName, userData["user"]) || userData["user"] == rooms.GetAllUsersAlias(roomName) {
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetNickclashTemplate(roomName)), 200, true)
    } else {
        rooms.AddUser(roomName, userData["user"], userData["color"], true)
        preprocessor.SetDataValue("{{.session-id}}", rooms.GetSessionId(userData["user"], roomName))
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetSkeletonTemplate(roomName)), 200, true)
        rooms.EnqueueMessage(roomName, userData["user"], "", "", "", rooms.GetJoinMessage(roomName), "")
    }
    newConn.Write(replyBuffer)
    newConn.Close()
}

// GetBriefHandle implements the handle for the brief document (GET).
func GetBriefHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var replyBuffer []byte
    if rooms.IsAllowingBriefs(roomName) {
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetBriefTemplate(roomName)), 200, true)
    } else {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    }
    newConn.Write(replyBuffer)
    newConn.Close()
}

// GetBodyHandle implements the handle for the body document (GET).
func GetBodyHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var userData map[string]string
    userData = rawhttp.GetFieldsFromGet(httpPayload)
    var validUser bool
    validUser = rooms.IsValidUserRequest(roomName, userData["user"], userData["id"])
    var replyBuffer []byte
    if !validUser {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else {
        replyBuffer = rawhttp.MakeReplyBuffer(preprocessor.ExpandData(roomName, rooms.GetBodyTemplate(roomName)), 200, false)
    }
    newConn.Write(replyBuffer)
    if validUser {
        rooms.SetUserConnection(roomName, userData["user"], newConn)
    } else {
        newConn.Close()
    }
}

// BadAssErrorHandle implements the handle for the any unexpected HTTP request (GET/POST/Whatever).
func BadAssErrorHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    newConn.Write(rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true))
    newConn.Close()
}

// PostBannerHandle implements the handle for the banner document (POST).
func PostBannerHandle(newConn net.Conn, roomName, httpPayload string, rooms *config.CherryRooms, preprocessor *html.Preprocessor) {
    var userData map[string]string
    var replyBuffer []byte
    var invalidRequest = false
    userData = rawhttp.GetFieldsFromPost(httpPayload)
    if _ , has := userData["user"]; !has {
        invalidRequest = true
    } else if _, has := userData["id"]; !has {
        invalidRequest = true
    } else if _, has := userData["action"]; !has {
        invalidRequest = true
    } else if _, has := userData["whoto"]; !has {
        invalidRequest = true
    } else if  _, has := userData["image"]; !has {
        invalidRequest = true
    } else if _, has := userData["says"]; !has {
        invalidRequest = true
    }
    var restoreBanner = true
    if invalidRequest || !rooms.IsValidUserRequest(roomName, userData["user"], userData["id"]) {
        replyBuffer = rawhttp.MakeReplyBuffer(html.GetBadAssErrorData(), 404, true)
    } else if userData["action"] == rooms.GetIgnoreAction(roomName) {
        if userData["user"] != userData["whoto"] && ! rooms.IsIgnored(userData["user"], userData["whoto"], roomName) {
            rooms.AddToIgnoreList(userData["user"], userData["whoto"], roomName)
            rooms.EnqueueMessage(roomName, userData["user"], "", "", "", rooms.GetOnIgnoreMessage(roomName) + userData["whoto"], "1")
            restoreBanner = false
        }
    } else if userData["action"] == rooms.GetDeIgnoreAction(roomName) {
        if rooms.IsIgnored(userData["user"], userData["whoto"], roomName) {
            rooms.DelFromIgnoreList(userData["user"], userData["whoto"], roomName)
            rooms.EnqueueMessage(roomName, userData["user"], "", "", "", rooms.GetOnDeIgnoreMessage(roomName) + userData["whoto"], "1")
            restoreBanner = false
        }
    } else {
        var somethingToSay = (len(userData["says"]) > 0 || len(userData["image"]) > 0 || len(userData["sound"]) > 0)
        if somethingToSay {
            //  INFO(Santiago): Any further antiflood control would go from here.
            rooms.EnqueueMessage(roomName, userData["user"], userData["whoto"], userData["action"], userData["image"], userData["says"], userData["priv"])
        }
    }
    preprocessor.SetDataValue("{{.nickname}}", userData["user"])
    preprocessor.SetDataValue("{{.session-id}}", userData["id"])
    if userData["priv"] == "1" {
        preprocessor.SetDataValue("{{.priv}}", "checked")
    }
    tempBanner := preprocessor.ExpandData(roomName, rooms.GetBannerTemplate(roomName))
    if restoreBanner {
        tempBanner = strings.Replace(tempBanner,
                                      "<option value = \"" + userData["whoto"] + "\">",
                                      "<option value = \"" + userData["whoto"] + "\" selected>", -1)
        tempBanner = strings.Replace(tempBanner,
                                      "<option value = \"" + userData["action"] + "\">",
                                      "<option value = \"" + userData["action"] + "\" selected>", -1)
    }
    replyBuffer = rawhttp.MakeReplyBuffer(tempBanner, 200, true)
    newConn.Write(replyBuffer)
    newConn.Close()
}
