/*
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package config

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

type RoomConfig struct {
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
    room_config.actions = make(map[string]*RoomAction)
    room_config.images = make(map[string]*RoomMediaResource)
    room_config.sounds = make(map[string]*RoomMediaResource)
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
