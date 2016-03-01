/*
Package config "classify" (:P) a plain cherry file.
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package config

import (
    "sync"
    "net"
    "fmt"
    "crypto/md5"
    "io"
    "sort"
)

type RoomMisc struct {
    listenPort int16
    joinMessage string
    exitMessage string
    onIgnoreMessage string
    onDeIgnoreMessage string
    greetingMessage string
    privateMessageMarker string
    maxUsers int
    allowBrief bool
    floodingPolice bool
    maxFloodAllowedBeforeKick int
    allUsersAlias string
}

type RoomAction struct {
    label string
    template string
}

type RoomMediaResource struct {
    label string
    template string
    url string
}

type Message struct {
    From string
    To string
    Action string
//    Sound string
    Image string
    Say string
    Priv string
}

type RoomUser struct {
    sessionId string
    color string
    ignoreList []string
    kickout bool
    conn net.Conn
}

type RoomConfig struct {
    mutex *sync.Mutex
    MainPeer net.Listener
    messageQueue []Message
    publicMessages []string
    users map[string]*RoomUser
    templates map[string]string
    misc *RoomMisc
    actions map[string]*RoomAction
    images map[string]*RoomMediaResource
    //sounds map[string]*RoomMediaResource
    ignoreAction string
    deignoreAction string
}

type CherryRooms struct {
    configs map[string]*RoomConfig
    servername string
}

func NewCherryRooms() *CherryRooms {
    return &CherryRooms{make(map[string]*RoomConfig), "localhost"}
}

func (c *CherryRooms) GetRoomActionLabel(roomName, action string) string {
    c.Lock(roomName)
    var label string
    label = c.configs[roomName].actions[action].label
    c.Unlock(roomName)
    return label
}

func (c *CherryRooms) GetRoomUsers(roomName string) []string {
    var users []string
    users = make([]string, 0)
    c.Lock(roomName)
    for user, _ := range c.configs[roomName].users {
        users = append(users, user)
    }
    c.Unlock(roomName)
    return users
}

func (c *CherryRooms) GetRooms() []string {
    var rooms []string
    rooms = make([]string, 0)
    for room, _ := range c.configs {
        rooms = append(rooms, room)
    }
    return rooms
}

func (c *CherryRooms) GetUserConnection(roomName, user string) net.Conn {
    var conn net.Conn
    c.Lock(roomName)
    conn = c.configs[roomName].users[user].conn
    c.Unlock(roomName)
    return conn
}

func (c *CherryRooms) GetRoomActionTemplate(roomName, action string) string {
    c.Lock(roomName)
    var template string
    template = c.configs[roomName].actions[action].template
    c.Unlock(roomName)
    return template
}

func (c *CherryRooms) AddUser(roomName, nickname, color string, kickout bool) {
    c.configs[roomName].mutex.Lock()
    md := md5.New()
    io.WriteString(md, roomName + nickname + color)
    id := fmt.Sprintf("%x", md.Sum(nil))
    c.configs[roomName].users[nickname] = &RoomUser{id, color, make([]string, 0), kickout, nil}
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) RemoveUser(roomName, nickname string) {
    c.configs[roomName].mutex.Lock()
    delete(c.configs[roomName].users, nickname)
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) EnqueueMessage(roomName, from, to, action, image, say, priv string) {
    c.configs[roomName].mutex.Lock()
    c.configs[roomName].messageQueue = append(c.configs[roomName].messageQueue, Message{from, to, action, image, say, priv})
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) DequeueMessage(roomName string) {
    c.configs[roomName].mutex.Lock()
    if len(c.configs[roomName].messageQueue) >= 1 {
        c.configs[roomName].messageQueue = c.configs[roomName].messageQueue[1:]
    }
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) GetNextMessage(roomName string) Message {
    c.configs[roomName].mutex.Lock()
    var message Message
    if len(c.configs[roomName].messageQueue) > 0 {
        message = c.configs[roomName].messageQueue[0]
    } else {
        message = Message{}
    }
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetSessionId(from, roomName string) string {
    if len(from) == 0 || !c.HasUser(roomName, from) {
        return ""
    }
    c.configs[roomName].mutex.Lock()
    var sessionId string
    sessionId = c.configs[roomName].users[from].sessionId
    c.configs[roomName].mutex.Unlock()
    return sessionId
}

func (c *CherryRooms) GetColor(from, roomName string) string {
    if len(from) == 0 || !c.HasUser(roomName, from) {
        return ""
    }
    c.configs[roomName].mutex.Lock()
    var color string
    color = c.configs[roomName].users[from].color
    c.configs[roomName].mutex.Unlock()
    return color
}

func (c *CherryRooms) GetIgnoreList(from, roomName string) string {
    if len(from) == 0 || !c.HasUser(roomName, from) {
        return ""
    }
    c.configs[roomName].mutex.Lock()
    var ignoreList string
    ignoring := c.configs[roomName].users[from].ignoreList
    lastIndex := len(ignoring) - 1
    for c, who := range ignoring {
        ignoreList += "\"" + who + "\""
        if c != lastIndex {
            who += ", "
        }
    }
    c.configs[roomName].mutex.Unlock()
    return ignoreList
}

func (c *CherryRooms) AddToIgnoreList(from, to, roomName string) {
    if len(from) == 0 || len(to) == 0 || !c.HasUser(roomName, from) || !c.HasUser(roomName, to) {
        return
    }
    c.configs[roomName].mutex.Lock()
    for _, t := range c.configs[roomName].users[from].ignoreList {
        if t == to {
            c.configs[roomName].mutex.Unlock()
            return
        }
    }
    c.configs[roomName].users[from].ignoreList = append(c.configs[roomName].users[from].ignoreList, to)
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) DelFromIgnoreList(from, to, roomName string) {
    if len(from) == 0 || len(to) == 0 || !c.HasUser(roomName, from) || !c.HasUser(roomName, to) {
        return
    }
    var index int = -1
    c.configs[roomName].mutex.Lock()
    for it, t := range c.configs[roomName].users[from].ignoreList {
        if t == to {
            index = it
            break
        }
    }
    if index != -1 {
        c.configs[roomName].users[from].ignoreList = append(c.configs[roomName].users[from].ignoreList[:index], c.configs[roomName].users[from].ignoreList[index+1:]...)
    }
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) IsIgnored(from, to, roomName string) bool {
    if len(from) == 0 || len(to) == 0 || !c.HasUser(roomName, from) || !c.HasUser(roomName, to) {
        return false
    }
    var retval bool = false
    c.configs[roomName].mutex.Lock()
    for _, t := range c.configs[roomName].users[from].ignoreList {
        if t == to {
            retval = true
            break
        }
    }
    c.configs[roomName].mutex.Unlock()
    return retval
}

func (c *CherryRooms) GetGreetingMessage(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var message string
    message = c.configs[roomName].misc.greetingMessage
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetJoinMessage(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var message string
    message = c.configs[roomName].misc.joinMessage
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetExitMessage(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var message string
    message = c.configs[roomName].misc.exitMessage
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetOnIgnoreMessage(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var message string
    message = c.configs[roomName].misc.onIgnoreMessage
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetOnDeIgnoreMessage(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var message string
    message = c.configs[roomName].misc.onDeIgnoreMessage
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetPrivateMessageMarker(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var message string
    message = c.configs[roomName].misc.privateMessageMarker
    c.configs[roomName].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetMaxUsers(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var max string
    max = fmt.Sprintf("%d", c.configs[roomName].misc.maxUsers)
    c.configs[roomName].mutex.Unlock()
    return max
}

func (c *CherryRooms) GetAllUsersAlias(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var alias string
    alias = c.configs[roomName].misc.allUsersAlias
    c.configs[roomName].mutex.Unlock()
    return alias
}

func (c *CherryRooms) GetActionList(roomName string) string {
    c.Lock(roomName)
    var actionList string = ""
    var actions []string
    actions = make([]string, 0)
    for action, _ := range c.configs[roomName].actions {
        actions = append(actions, action)
    }
    sort.Strings(actions)
    for _, action := range actions {
        actionList += "<option value = \"" + action + "\">" + c.configs[roomName].actions[action].label + "\n"
    }
    c.Unlock(roomName)
    return actionList
}

func(c *CherryRooms) getMediaResourceList(roomName string, mediaResource map[string]*RoomMediaResource) string {
    c.Lock(roomName)
    var mediaRsrcList string = ""
    var resources []string
    resources = make([]string, 0)
    for resource, _ := range mediaResource {
        resources = append(resources, resource)
    }
    sort.Strings(resources)
    for _, resource := range resources {
        mediaRsrcList += "<option value = \"" + resource + "\">" + mediaResource[resource].label + "\n"
    }
    c.Unlock(roomName)
    return mediaRsrcList
}

func (c *CherryRooms) GetImageList(roomName string) string {
    return c.getMediaResourceList(roomName, c.configs[roomName].images)
}

//func (c *CherryRooms) GetSoundList(room_name string) string {
//    return c.getMediaResourceList(room_name, c.configs[room_name].sounds)
//}

func (c *CherryRooms) GetUsersList(roomName string) string {
    c.Lock(roomName)
    var users []string
    users = make([]string, 0)
    for user, _ := range c.configs[roomName].users {
        users = append(users, user)
    }
    //  WARN(Santiago): Already locked, we can acquire this piece of information directly... otherwise we got a deadlock.
    allUsersAlias := c.configs[roomName].misc.allUsersAlias
    var usersList string = "<option value = \"" + allUsersAlias + "\">" + allUsersAlias + "\n"
    sort.Strings(users)
    for _, user := range users {
        usersList += "<option value = \"" + user + "\">" + user + "\n"
    }
    c.Unlock(roomName)
    return usersList
}

func (c *CherryRooms) getRoomTemplate(roomName, template string) string {
    c.configs[roomName].mutex.Lock()
    var data string
    data = c.configs[roomName].templates[template]
    c.configs[roomName].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetTopTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "top")
}

func (c *CherryRooms) GetBodyTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "body")
}

func (c *CherryRooms) GetBannerTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "banner")
}

func (c *CherryRooms) GetHighlightTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "highlight")
}

func (c *CherryRooms) GetEntranceTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "entrance")
}

func (c *CherryRooms) GetExitTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "exit")
}

func (c *CherryRooms) GetNickclashTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "nickclash")
}

func (c *CherryRooms) GetSkeletonTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "skeleton")
}

func (c *CherryRooms) GetBriefTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "brief")
}

func (c *CherryRooms) GetFindResultsHeadTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "find-results-head")
}

func (c *CherryRooms) GetFindResultsBodyTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "find-results-body")
}

func (c *CherryRooms) GetFindResultsTailTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "find-results-tail")
}

func (c *CherryRooms) GetFindBotTemplate(roomName string) string {
    return c.getRoomTemplate(roomName, "find-bot")
}

func (c *CherryRooms) GetLastPublicMessages(roomName string) string {
    if !c.HasRoom(roomName) {
        return ""
    }
    var retval string
    c.Lock(roomName)
    msgs := c.configs[roomName].publicMessages
    c.Unlock(roomName)
    for _, m := range msgs {
        retval += m
    }
    return retval
}

func (c *CherryRooms) AddPublicMessage(roomName, message string) {
    if !c.HasRoom(roomName) {
        return
    }
    c.Lock(roomName)
    if (len(c.configs[roomName].publicMessages) == 10) {
        c.configs[roomName].publicMessages = c.configs[roomName].publicMessages[1:len(c.configs[roomName].publicMessages)-1]
    }
    c.configs[roomName].publicMessages = append(c.configs[roomName].publicMessages, message)
    c.Unlock(roomName)
}

func (c *CherryRooms) GetListenPort(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var port string
    port = fmt.Sprintf("%d", c.configs[roomName].misc.listenPort)
    c.configs[roomName].mutex.Unlock()
    return port
}

func (c *CherryRooms) GetUsersTotal(roomName string) string {
    c.configs[roomName].mutex.Lock()
    var total string
    total = fmt.Sprintf("%d", len(c.configs[roomName].users))
    c.configs[roomName].mutex.Unlock()
    return total
}

func (c *CherryRooms) AddRoom(roomName string, listenPort int16) bool {
    if c.HasRoom(roomName) || c.PortBusyByAnotherRoom(listenPort) {
        return false
    }
    c.configs[roomName] = c.initConfig()
    c.configs[roomName].misc.listenPort = listenPort
    return true
}

func (c *CherryRooms) AddAction(roomName, id, label, template string) {
    c.configs[roomName].actions[id] = &RoomAction{label, template}
}

func (c *CherryRooms) AddImage(roomName, id, label, template, url string) {
    c.configs[roomName].images[id] = c.newMediaResource(label, template, url)
}

//func (c *CherryRooms) AddSound(room_name, id, label, template, url string) {
//    c.configs[room_name].sounds[id] = c.newMediaResource(label, template, url)
//}

func (c *CherryRooms) newMediaResource(label, template, url string) *RoomMediaResource {
    return &RoomMediaResource{label, template, url}
}

func (c *CherryRooms) HasAction(roomName, id string) bool {
    _, ok := c.configs[roomName].actions[id]
    return ok
}

func (c *CherryRooms) HasImage(roomName, id string) bool {
    _, ok := c.configs[roomName].images[id]
    return ok
}

//func (c *CherryRooms) HasSound(room_name, id string) bool {
//    _, ok := c.configs[room_name].sounds[id]
//    return ok
//}

func (c *CherryRooms) HasRoom(roomName string) bool {
    _, ok := c.configs[roomName]
    return ok
}

func (c *CherryRooms) PortBusyByAnotherRoom(port int16) bool {
    for _, c := range c.configs {
        if c.misc.listenPort == port {
            return true
        }
    }
    return false
}

func (c *CherryRooms) GetRoomByPort(port int16) *RoomConfig {
    for _, r := range c.configs {
        if r.misc.listenPort == port {
            return r
        }
    }
    return nil
}

func (c *CherryRooms) initConfig() *RoomConfig {
    var room_config *RoomConfig
    room_config = new(RoomConfig)
    room_config.misc = &RoomMisc{}
    room_config.messageQueue = make([]Message, 0)
    room_config.publicMessages = make([]string, 0)
    room_config.users = make(map[string]*RoomUser)
    room_config.templates = make(map[string]string)
    room_config.actions = make(map[string]*RoomAction)
    room_config.images = make(map[string]*RoomMediaResource)
    //room_config.sounds = make(map[string]*RoomMediaResource)
    room_config.mutex = new(sync.Mutex)
    return room_config
}

func (c *CherryRooms) AddTemplate(roomName, id, template string) {
    c.configs[roomName].templates[id] = template
}

func (c *CherryRooms) HasTemplate(roomName, id string) bool {
    _, ok := c.configs[roomName].templates[id]
    return ok
}

func (c *CherryRooms) SetJoinMessage(roomName, message string) {
    c.configs[roomName].misc.joinMessage = message
}

func (c *CherryRooms) SetExitMessage(roomName, message string) {
    c.configs[roomName].misc.exitMessage = message
}

func (c *CherryRooms) SetOnIgnoreMessage(roomName, message string) {
    c.configs[roomName].misc.onIgnoreMessage = message
}

func (c *CherryRooms) SetOnDeIgnoreMessage(roomName, message string) {
    c.configs[roomName].misc.onDeIgnoreMessage = message
}

func (c *CherryRooms) SetGreetingMessage(roomName, message string) {
    c.configs[roomName].misc.greetingMessage = message
}

func (c *CherryRooms) SetPrivateMessageMarker(roomName, marker string) {
    c.configs[roomName].misc.privateMessageMarker = marker
}

func (c *CherryRooms) SetMaxUsers(roomName string, value int) {
    c.configs[roomName].misc.maxUsers = value
}

func (c *CherryRooms) SetAllowBrief(roomName string, value bool) {
    c.configs[roomName].misc.allowBrief = value
}

func (c *CherryRooms) IsAllowingBriefs(roomName string) bool {
    return c.configs[roomName].misc.allowBrief
}

func (c *CherryRooms) SetFloodingPolice(roomName string, value bool) {
    c.configs[roomName].misc.floodingPolice = value
}

func (c *CherryRooms) SetMaxFloodAllowedBeforeKick(roomName string, value int) {
    c.configs[roomName].misc.maxFloodAllowedBeforeKick = value
}

func (c *CherryRooms) SetAllUsersAlias(roomName, alias string) {
    c.configs[roomName].misc.allUsersAlias = alias
}

func (c *CherryRooms) Lock(roomName string) {
    c.configs[roomName].mutex.Lock()
}

func (c *CherryRooms) Unlock(roomName string) {
    c.configs[roomName].mutex.Unlock()
}

func (c *CherryRooms) GetServername() string {
    return c.servername
}

func (c *CherryRooms) SetServername(servername string) {
    c.servername = servername
}

func (c *CherryRooms) HasUser(roomName, user string) bool {
    _, ok := c.configs[roomName]
    if !ok {
        return false
    }
    _, ok = c.configs[roomName].users[user]
    return ok
}

func (c *CherryRooms) IsValidUserRequest(roomName, user, id string) bool {
    var valid bool = false
    if (c.HasUser(roomName, user)) {
        valid = (id == c.GetSessionId(user, roomName))
    }
    return valid
}

func (c *CherryRooms) SetIgnoreAction(roomName, action string) {
    c.Lock(roomName)
    c.configs[roomName].ignoreAction = action
    c.Unlock(roomName)
}

func (c *CherryRooms) SetDeIgnoreAction(roomName, action string) {
    c.Lock(roomName)
    c.configs[roomName].deignoreAction = action
    c.Unlock(roomName)
}

func (c *CherryRooms) GetIgnoreAction(roomName string) string {
    c.Lock(roomName)
    var retval string
    retval = c.configs[roomName].ignoreAction

    c.Unlock(roomName)
    return retval
}

func (c *CherryRooms) GetDeIgnoreAction(roomName string) string {
    c.Lock(roomName)
    var retval string
    retval = c.configs[roomName].deignoreAction
    c.Unlock(roomName)
    return retval
}

func (c *CherryRooms) SetUserConnection(roomName, user string, conn net.Conn) {
    c.Lock(roomName)
    c.configs[roomName].users[user].conn = conn
    c.Unlock(roomName)
}

func (c *CherryRooms) GetServerName() string {
    return c.servername
}
