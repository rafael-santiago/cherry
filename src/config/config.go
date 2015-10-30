/*
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package config

import (
    "sync"
//    "container/list"
    "net"
    "fmt"
)

type RoomMisc struct {
    listen_port int16
    join_message string
    exit_message string
    on_ignore_message string
    on_deignore_message string
    greeting_message string
    private_message_marker string
    max_users int
    allow_brief bool
    flooding_police bool
    max_flood_allowed_before_kick int
    all_users_alias string
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
    Sound string
    Image string
    Say string
    Priv string
}

type RoomUser struct {
    session_id string
    color string
    //ignorelist *list.List
    ignorelist []string
    kickout bool
    conn *net.Conn
}

type RoomConfig struct {
    mutex *sync.Mutex
    //message_queue *list.List
    message_queue []Message
    users map[string]*RoomUser
    templates map[string]string
    misc *RoomMisc
    actions map[string]*RoomAction
    images map[string]*RoomMediaResource
    sounds map[string]*RoomMediaResource
}

type CherryRooms struct {
    configs map[string]*RoomConfig
}

func NewCherryRooms() *CherryRooms {
    return &CherryRooms{make(map[string]*RoomConfig)}
}

func (c *CherryRooms) AddUser(room_name, nickname, id, color string, kickout bool, conn *net.Conn) {
    c.configs[room_name].mutex.Lock()
    c.configs[room_name].users[nickname] = &RoomUser{id, color, make([]string, 0)/*new(list.List)*/, kickout, conn}
    c.configs[room_name].mutex.Unlock()
}

func (c *CherryRooms) RemoveUser(room_name, nickname string) {
    c.configs[room_name].mutex.Lock()
    delete(c.configs[room_name].users, nickname)
    c.configs[room_name].mutex.Unlock()
}

func (c *CherryRooms) EnqueueMessage(room_name, from, to, action, sound, image, say, priv string) {
    c.configs[room_name].mutex.Lock()
    c.configs[room_name].message_queue = append(c.configs[room_name].message_queue, Message{from, to, action, sound, image, say, priv})
    //c.configs[room_name].message_queue.PushBack(Message{from, to, action, sound, image, say, priv})
    c.configs[room_name].mutex.Unlock()
}

func (c *CherryRooms) DequeueMessage(room_name string) {
    c.configs[room_name].mutex.Lock()
    //c.configs[room_name].message_queue.Remove(c.configs[room_name].message_queue.Front())
    c.configs[room_name].message_queue = c.configs[room_name].message_queue[1:]
    c.configs[room_name].mutex.Unlock()
}

func (c *CherryRooms) GetNextMessage(room_name string) Message {
    c.configs[room_name].mutex.Lock()
    var message Message
    message = c.configs[room_name].message_queue[0]
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetSessionId(from, room_name string) string {
    c.configs[room_name].mutex.Lock()
    var session_id string
    session_id = c.configs[room_name].users[from].session_id
    c.configs[room_name].mutex.Unlock()
    return session_id
}

func (c *CherryRooms) GetColor(from, room_name string) string {
    c.configs[room_name].mutex.Lock()
    var color string
    color = c.configs[room_name].users[from].color
    c.configs[room_name].mutex.Unlock()
    return color
}

func (c *CherryRooms) GetIgnoreList(from, room_name string) string {
    c.configs[room_name].mutex.Lock()
    var ignore_list string
    ignoring := c.configs[room_name].users[from].ignorelist
    last_index := len(ignoring) - 1
    for c, who := range ignoring {
        ignore_list += "\"" + who + "\""
        if c != last_index {
            who += ", "
        }
    }
    c.configs[room_name].mutex.Unlock()
    return ignore_list
}

func (c *CherryRooms) GetGreetingMessage(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var message string
    message = c.configs[room_name].misc.greeting_message
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetJoinMessage(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var message string
    message = c.configs[room_name].misc.join_message
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetExitMessage(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var message string
    message = c.configs[room_name].misc.exit_message
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetOnIgnoreMessage(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var message string
    message = c.configs[room_name].misc.on_ignore_message
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetOnDeIgnoreMessage(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var message string
    message = c.configs[room_name].misc.on_deignore_message
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetPrivateMessageMarker(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var message string
    message = c.configs[room_name].misc.private_message_marker
    c.configs[room_name].mutex.Unlock()
    return message
}

func (c *CherryRooms) GetMaxUsers(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var max string
    max = fmt.Sprintf("%d", c.configs[room_name].misc.max_users)
    c.configs[room_name].mutex.Unlock()
    return max
}

func (c *CherryRooms) GetAllUsersAlias(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var alias string
    alias = c.configs[room_name].misc.all_users_alias
    c.configs[room_name].mutex.Unlock()
    return alias
}

func (c *CherryRooms) GetActionList(room_name string) string {
    return "TODO(Santiago): what?"
}

func (c *CherryRooms) GetImageList(room_name string) string {
    return "TODO(Santiago): what?"
}

func (c *CherryRooms) GetSoundList(room_name string) string {
    return "TODO(Santiago): what?"
}

func (c *CherryRooms) GetUsersList(room_name string) string {
    return "TODO(Santiago): what?"
}

func (c *CherryRooms) GetTopTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["top"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetBodyTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["body"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetBannerTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["banner"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetHighlightTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["highlight"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetEntranceTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["entrance"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetExitTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["exit"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetNickclashTemplate(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var data string
    data = c.configs[room_name].templates["nickclash"]
    c.configs[room_name].mutex.Unlock()
    return data
}

func (c *CherryRooms) GetLastPublicMessages(room_name string) string {
    return "TODO(Santiago): what?"
}

func (c *CherryRooms) GetServername() string {
    return "TODO(Santiago): what?"
}

func (c *CherryRooms) GetListenPort(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var port string
    port = fmt.Sprintf("%d", c.configs[room_name].misc.listen_port)
    c.configs[room_name].mutex.Unlock()
    return port
}

func (c *CherryRooms) GetUsersTotal(room_name string) string {
    c.configs[room_name].mutex.Lock()
    var total string
    total = fmt.Sprintf("%d", len(c.configs[room_name].users))
    c.configs[room_name].mutex.Unlock()
    return total
}

func (c *CherryRooms) AddRoom(room_name string, listen_port int16) bool {
    if c.HasRoom(room_name) || c.PortBusyByAnotherRoom(listen_port) {
        return false
    }
    c.configs[room_name] = c.init_config()
    c.configs[room_name].misc.listen_port = listen_port
    return true;
}

func (c *CherryRooms) AddAction(room_name, id, label, template string) {
    c.configs[room_name].actions[id] = &RoomAction{label, template}
}

func (c *CherryRooms) AddImage(room_name, id, label, template, url string) {
    c.configs[room_name].images[id] = c.new_media_resource(label, template, url)
}

func (c *CherryRooms) AddSound(room_name, id, label, template, url string) {
    c.configs[room_name].sounds[id] = c.new_media_resource(label, template, url)
}

func (c *CherryRooms) new_media_resource(label, template, url string) *RoomMediaResource {
    return &RoomMediaResource{label, template, url}
}

func (c *CherryRooms) HasAction(room_name, id string) bool {
    _, ok := c.configs[room_name].actions[id]
    return ok
}

func (c *CherryRooms) HasImage(room_name, id string) bool {
    _, ok := c.configs[room_name].images[id]
    return ok
}

func (c *CherryRooms) HasSound(room_name, id string) bool {
    _, ok := c.configs[room_name].sounds[id]
    return ok
}

func (c *CherryRooms) HasRoom(room_name string) bool {
    _, ok := c.configs[room_name]
    return ok
}

func (c *CherryRooms) PortBusyByAnotherRoom(port int16) bool {
    for _, c := range c.configs {
        if c.misc.listen_port == port {
            return true
        }
    }
    return false
}

func (c *CherryRooms) GetRoomByPort(port int16) *RoomConfig {
    for _, r := range c.configs {
        if r.misc.listen_port == port {
            return r
        }
    }
    return nil
}

func (c *CherryRooms) init_config() *RoomConfig {
    var room_config *RoomConfig
    room_config = new(RoomConfig)
    room_config.misc = &RoomMisc{}
    room_config.message_queue = make([]Message, 0)//list.New()
    room_config.users = make(map[string]*RoomUser)
    room_config.templates = make(map[string]string)
    room_config.actions = make(map[string]*RoomAction)
    room_config.images = make(map[string]*RoomMediaResource)
    room_config.sounds = make(map[string]*RoomMediaResource)
    room_config.mutex = new(sync.Mutex)
    return room_config
}

func (c *CherryRooms) AddTemplate(room_name, id, template string) {
    c.configs[room_name].templates[id] = template
}

func (c *CherryRooms) HasTemplate(room_name, id string) bool {
    _, ok := c.configs[room_name].templates[id]
    return ok
}

func (c *CherryRooms) SetJoinMessage(room_name, message string) {
    c.configs[room_name].misc.join_message = message
}

func (c *CherryRooms) SetExitMessage(room_name, message string) {
    c.configs[room_name].misc.exit_message = message
}

func (c *CherryRooms) SetOnIgnoreMessage(room_name, message string) {
    c.configs[room_name].misc.on_ignore_message = message
}

func (c *CherryRooms) SetOnDeIgnoreMessage(room_name, message string) {
    c.configs[room_name].misc.on_deignore_message = message
}

func (c *CherryRooms) SetGreetingMessage(room_name, message string) {
    c.configs[room_name].misc.greeting_message = message
}

func (c *CherryRooms) SetPrivateMessageMarker(room_name, marker string) {
    c.configs[room_name].misc.private_message_marker = marker
}

func (c *CherryRooms) SetMaxUsers(room_name string, value int) {
    c.configs[room_name].misc.max_users = value
}

func (c *CherryRooms) SetAllowBrief(room_name string, value bool) {
    c.configs[room_name].misc.allow_brief = value
}

func (c *CherryRooms) SetFloodingPolice(room_name string, value bool) {
    c.configs[room_name].misc.flooding_police = value
}

func (c *CherryRooms) SetMaxFloodAllowedBeforeKick(room_name string, value int) {
    c.configs[room_name].misc.max_flood_allowed_before_kick = value
}

func (c *CherryRooms) SetAllUsersAlias(room_name, alias string) {
    c.configs[room_name].misc.all_users_alias = alias
}

func (c *CherryRooms) Lock(room_name string) {
    c.configs[room_name].mutex.Lock()
}

func (c *CherryRooms) Unlock(room_name string) {
    c.configs[room_name].mutex.Unlock()
}
