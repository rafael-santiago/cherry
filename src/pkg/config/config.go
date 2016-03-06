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
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
	"sync"
)

// RoomMisc gathers the misc options for a room.
type RoomMisc struct {
	listenPort                int16
	joinMessage               string
	exitMessage               string
	onIgnoreMessage           string
	onDeIgnoreMessage         string
	greetingMessage           string
	privateMessageMarker      string
	maxUsers                  int
	allowBrief                bool
	floodingPolice            bool
	maxFloodAllowedBeforeKick int
	allUsersAlias             string
}

// RoomAction gathers the label and the template (data) from an action.
type RoomAction struct {
	label    string
	template string
}

// RoomMediaResource gathers the label, template (data) and uri from a generical resource (until now, images only).
type RoomMediaResource struct {
	label    string
	template string
	url      string
}

// Message gathers a message that must be formatted and delivered.
type Message struct {
	From   string
	To     string
	Action string
	//    Sound string
	Image string
	Say   string
	Priv  string
}

// RoomUser is the user context.
type RoomUser struct {
	sessionID  string
	color      string
	ignoreList []string
	kickout    bool
	conn       net.Conn
	addr       string
}

// RoomConfig represents in memory a defined room loaded from a cherry file.
type RoomConfig struct {
	mutex          *sync.Mutex
	MainPeer       net.Listener
	messageQueue   []Message
	publicMessages []string
	users          map[string]*RoomUser
	templates      map[string]string
	misc           *RoomMisc
	actions        map[string]*RoomAction
	images         map[string]*RoomMediaResource
	//sounds map[string]*RoomMediaResource
	ignoreAction   string
	deignoreAction string
}

// CherryRooms represents your cherry tree... I mean your cherry server.
type CherryRooms struct {
	configs    map[string]*RoomConfig
	servername string
}

// NewCherryRooms creates a new server container.
func NewCherryRooms() *CherryRooms {
	return &CherryRooms{make(map[string]*RoomConfig), "localhost"}
}

// GetRoomActionLabel spits a room action label.
func (c *CherryRooms) GetRoomActionLabel(roomName, action string) string {
	c.Lock(roomName)
	var label string
	label = c.configs[roomName].actions[action].label
	c.Unlock(roomName)
	return label
}

// GetRoomUsers spits all users connected in the specified room.
func (c *CherryRooms) GetRoomUsers(roomName string) []string {
	var users []string
	users = make([]string, 0)
	c.Lock(roomName)
	for user := range c.configs[roomName].users {
		users = append(users, user)
	}
	c.Unlock(roomName)
	return users
}

// GetRooms spits all opened rooms.
func (c *CherryRooms) GetRooms() []string {
	var rooms []string
	rooms = make([]string, 0)
	for room := range c.configs {
		rooms = append(rooms, room)
	}
	return rooms
}

// GetUserConnection returns the net.Conn that represents an "user peer".
func (c *CherryRooms) GetUserConnection(roomName, user string) net.Conn {
	var conn net.Conn
	c.Lock(roomName)
	conn = c.configs[roomName].users[user].conn
	c.Unlock(roomName)
	return conn
}

// GetRoomActionTemplate returns the template data from an action.
func (c *CherryRooms) GetRoomActionTemplate(roomName, action string) string {
	c.Lock(roomName)
	var template string
	template = c.configs[roomName].actions[action].template
	c.Unlock(roomName)
	return template
}

// AddUser does what it is saying. BELIEVE or NOT!!!
func (c *CherryRooms) AddUser(roomName, nickname, color string, kickout bool) {
	c.configs[roomName].mutex.Lock()
	md := md5.New()
	io.WriteString(md, roomName+nickname+color)
	id := fmt.Sprintf("%x", md.Sum(nil))
	c.configs[roomName].users[nickname] = &RoomUser{id, color, make([]string, 0), kickout, nil, ""}
	c.configs[roomName].mutex.Unlock()
}

// RemoveUser removes a user...
func (c *CherryRooms) RemoveUser(roomName, nickname string) {
	c.configs[roomName].mutex.Lock()
	delete(c.configs[roomName].users, nickname)
	c.configs[roomName].mutex.Unlock()
}

// EnqueueMessage adds to the queue an user message.
func (c *CherryRooms) EnqueueMessage(roomName, from, to, action, image, say, priv string) {
	c.configs[roomName].mutex.Lock()
	c.configs[roomName].messageQueue = append(c.configs[roomName].messageQueue, Message{from, to, action, image, say, priv})
	c.configs[roomName].mutex.Unlock()
}

// DequeueMessage removes from the queue the oldest user message.
func (c *CherryRooms) DequeueMessage(roomName string) {
	c.configs[roomName].mutex.Lock()
	if len(c.configs[roomName].messageQueue) >= 1 {
		c.configs[roomName].messageQueue = c.configs[roomName].messageQueue[1:]
	}
	c.configs[roomName].mutex.Unlock()
}

// GetNextMessage returns the next message that should be processed.
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

// GetSessionID returns the user's session ID.
func (c *CherryRooms) GetSessionID(from, roomName string) string {
	if len(from) == 0 || !c.HasUser(roomName, from) {
		return ""
	}
	c.configs[roomName].mutex.Lock()
	var sessionID string
	sessionID = c.configs[roomName].users[from].sessionID
	c.configs[roomName].mutex.Unlock()
	return sessionID
}

// GetColor returns the user's nickname color.
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

// GetIgnoreList returns all users ignored by an user.
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

// AddToIgnoreList add to the user context some user to be ignored.
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

// DelFromIgnoreList removes from the user context a previous ignored user.
func (c *CherryRooms) DelFromIgnoreList(from, to, roomName string) {
	if len(from) == 0 || len(to) == 0 || !c.HasUser(roomName, from) || !c.HasUser(roomName, to) {
		return
	}
	var index = -1
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

// IsIgnored returns "true" if the user U is ignoring the asshole A, otherwise guess what.
func (c *CherryRooms) IsIgnored(from, to, roomName string) bool {
	if len(from) == 0 || len(to) == 0 || !c.HasUser(roomName, from) || !c.HasUser(roomName, to) {
		return false
	}
	var retval = false
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

// GetGreetingMessage returns the pre-configurated greeting message.
func (c *CherryRooms) GetGreetingMessage(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var message string
	message = c.configs[roomName].misc.greetingMessage
	c.configs[roomName].mutex.Unlock()
	return message
}

// GetJoinMessage returns the pre-configurated join message.
func (c *CherryRooms) GetJoinMessage(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var message string
	message = c.configs[roomName].misc.joinMessage
	c.configs[roomName].mutex.Unlock()
	return message
}

// GetExitMessage returns the pre-configurated exit message.
func (c *CherryRooms) GetExitMessage(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var message string
	message = c.configs[roomName].misc.exitMessage
	c.configs[roomName].mutex.Unlock()
	return message
}

// GetOnIgnoreMessage returns the pre-configurated "on ignore" message.
func (c *CherryRooms) GetOnIgnoreMessage(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var message string
	message = c.configs[roomName].misc.onIgnoreMessage
	c.configs[roomName].mutex.Unlock()
	return message
}

// GetOnDeIgnoreMessage returns the pre-configurated "on deignore" message.
func (c *CherryRooms) GetOnDeIgnoreMessage(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var message string
	message = c.configs[roomName].misc.onDeIgnoreMessage
	c.configs[roomName].mutex.Unlock()
	return message
}

// GetPrivateMessageMarker returns the private message marker.
func (c *CherryRooms) GetPrivateMessageMarker(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var message string
	message = c.configs[roomName].misc.privateMessageMarker
	c.configs[roomName].mutex.Unlock()
	return message
}

// GetMaxUsers returns the total of users allowed in a room.
func (c *CherryRooms) GetMaxUsers(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var max string
	max = fmt.Sprintf("%d", c.configs[roomName].misc.maxUsers)
	c.configs[roomName].mutex.Unlock()
	return max
}

// GetAllUsersAlias returns the "all users" alias.
func (c *CherryRooms) GetAllUsersAlias(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var alias string
	alias = c.configs[roomName].misc.allUsersAlias
	c.configs[roomName].mutex.Unlock()
	return alias
}

// GetActionList returns a well-formatted "HTML combo" containing all actions.
func (c *CherryRooms) GetActionList(roomName string) string {
	c.Lock(roomName)
	var actionList = ""
	var actions []string
	actions = make([]string, 0)
	for action := range c.configs[roomName].actions {
		actions = append(actions, action)
	}
	sort.Strings(actions)
	for _, action := range actions {
		actionList += "<option value = \"" + action + "\">" + c.configs[roomName].actions[action].label + "\n"
	}
	c.Unlock(roomName)
	return actionList
}

func (c *CherryRooms) getMediaResourceList(roomName string, mediaResource map[string]*RoomMediaResource) string {
	c.Lock(roomName)
	var mediaRsrcList = ""
	var resources []string
	resources = make([]string, 0)
	for resource := range mediaResource {
		resources = append(resources, resource)
	}
	sort.Strings(resources)
	for _, resource := range resources {
		mediaRsrcList += "<option value = \"" + resource + "\">" + mediaResource[resource].label + "\n"
	}
	c.Unlock(roomName)
	return mediaRsrcList
}

// GetImageList returns a well-formatted "HTML combo" containing all images.
func (c *CherryRooms) GetImageList(roomName string) string {
	return c.getMediaResourceList(roomName, c.configs[roomName].images)
}

//func (c *CherryRooms) GetSoundList(room_name string) string {
//    return c.getMediaResourceList(room_name, c.configs[room_name].sounds)
//}

// GetUsersList returns a well-formatted "HTML combo" containing all users connected on a room.
func (c *CherryRooms) GetUsersList(roomName string) string {
	c.Lock(roomName)
	var users []string
	users = make([]string, 0)
	for user := range c.configs[roomName].users {
		users = append(users, user)
	}
	//  WARN(Santiago): Already locked, we can acquire this piece of information directly... otherwise we got a deadlock.
	allUsersAlias := c.configs[roomName].misc.allUsersAlias
	var usersList = "<option value = \"" + allUsersAlias + "\">" + allUsersAlias + "\n"
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

// GetTopTemplate spits the top template data.
func (c *CherryRooms) GetTopTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "top")
}

// GetBodyTemplate spits the body template data.
func (c *CherryRooms) GetBodyTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "body")
}

// GetBannerTemplate spits the banner template data.
func (c *CherryRooms) GetBannerTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "banner")
}

// GetHighlightTemplate spits the highlight template data.
func (c *CherryRooms) GetHighlightTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "highlight")
}

// GetEntranceTemplate spits the entrance template data.
func (c *CherryRooms) GetEntranceTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "entrance")
}

// GetExitTemplate spits the exit template data.
func (c *CherryRooms) GetExitTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "exit")
}

// GetNickclashTemplate spits the nickclash template data.
func (c *CherryRooms) GetNickclashTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "nickclash")
}

// GetSkeletonTemplate spits the skeleton template data.
func (c *CherryRooms) GetSkeletonTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "skeleton")
}

// GetBriefTemplate spits the brief template data.
func (c *CherryRooms) GetBriefTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "brief")
}

// GetFindResultsHeadTemplate spits the find template data (HEAD).
func (c *CherryRooms) GetFindResultsHeadTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "find-results-head")
}

// GetFindResultsBodyTemplate spits the find template data (BODY).
func (c *CherryRooms) GetFindResultsBodyTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "find-results-body")
}

// GetFindResultsTailTemplate spits the find template data (TAIL).
func (c *CherryRooms) GetFindResultsTailTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "find-results-tail")
}

// GetFindBotTemplate spits the find bot template data.
func (c *CherryRooms) GetFindBotTemplate(roomName string) string {
	return c.getRoomTemplate(roomName, "find-bot")
}

// GetLastPublicMessages spits the last public messages (well-formatted in HTML).
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

// AddPublicMessage adds a public message to the pool of public messages that will be used for the room brief's composing.
func (c *CherryRooms) AddPublicMessage(roomName, message string) {
	if !c.HasRoom(roomName) {
		return
	}
	c.Lock(roomName)
	if len(c.configs[roomName].publicMessages) == 10 {
		c.configs[roomName].publicMessages = c.configs[roomName].publicMessages[1 : len(c.configs[roomName].publicMessages)-1]
	}
	c.configs[roomName].publicMessages = append(c.configs[roomName].publicMessages, message)
	c.Unlock(roomName)
}

// GetListenPort returns the port that is being used for the room serving.
func (c *CherryRooms) GetListenPort(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var port string
	port = fmt.Sprintf("%d", c.configs[roomName].misc.listenPort)
	c.configs[roomName].mutex.Unlock()
	return port
}

// GetUsersTotal returns the total of users currently talking in a room.
func (c *CherryRooms) GetUsersTotal(roomName string) string {
	c.configs[roomName].mutex.Lock()
	var total string
	total = fmt.Sprintf("%d", len(c.configs[roomName].users))
	c.configs[roomName].mutex.Unlock()
	return total
}

// AddRoom adds a room to the "memory".
func (c *CherryRooms) AddRoom(roomName string, listenPort int16) bool {
	if c.HasRoom(roomName) || c.PortBusyByAnotherRoom(listenPort) {
		return false
	}
	c.configs[roomName] = c.initConfig()
	c.configs[roomName].misc.listenPort = listenPort
	return true
}

// AddAction adds an action to the "memory".
func (c *CherryRooms) AddAction(roomName, id, label, template string) {
	c.configs[roomName].actions[id] = &RoomAction{label, template}
}

// AddImage adds an image (data that represents an image) to the "memory".
func (c *CherryRooms) AddImage(roomName, id, label, template, url string) {
	c.configs[roomName].images[id] = c.newMediaResource(label, template, url)
}

//func (c *CherryRooms) AddSound(room_name, id, label, template, url string) {
//    c.configs[room_name].sounds[id] = c.newMediaResource(label, template, url)
//}

func (c *CherryRooms) newMediaResource(label, template, url string) *RoomMediaResource {
	return &RoomMediaResource{label, template, url}
}

// HasAction verifies if an action really exists for the indicated room.
func (c *CherryRooms) HasAction(roomName, id string) bool {
	_, ok := c.configs[roomName].actions[id]
	return ok
}

// HasImage verifies if an image really exists for the indicated room.
func (c *CherryRooms) HasImage(roomName, id string) bool {
	_, ok := c.configs[roomName].images[id]
	return ok
}

//func (c *CherryRooms) HasSound(room_name, id string) bool {
//    _, ok := c.configs[room_name].sounds[id]
//    return ok
//}

// HasRoom verifies if a room really exists in this server.
func (c *CherryRooms) HasRoom(roomName string) bool {
	_, ok := c.configs[roomName]
	return ok
}

// PortBusyByAnotherRoom verifies if there is some port clash between rooms.
func (c *CherryRooms) PortBusyByAnotherRoom(port int16) bool {
	for _, c := range c.configs {
		if c.misc.listenPort == port {
			return true
		}
	}
	return false
}

// GetRoomByPort returns a room (all configuration from it) given a port.
func (c *CherryRooms) GetRoomByPort(port int16) *RoomConfig {
	for _, r := range c.configs {
		if r.misc.listenPort == port {
			return r
		}
	}
	return nil
}

func (c *CherryRooms) initConfig() *RoomConfig {
	var roomConfig *RoomConfig
	roomConfig = new(RoomConfig)
	roomConfig.misc = &RoomMisc{}
	roomConfig.messageQueue = make([]Message, 0)
	roomConfig.publicMessages = make([]string, 0)
	roomConfig.users = make(map[string]*RoomUser)
	roomConfig.templates = make(map[string]string)
	roomConfig.actions = make(map[string]*RoomAction)
	roomConfig.images = make(map[string]*RoomMediaResource)
	//room_config.sounds = make(map[string]*RoomMediaResource)
	roomConfig.mutex = new(sync.Mutex)
	return roomConfig
}

// AddTemplate adds a template based on room name, ID.
func (c *CherryRooms) AddTemplate(roomName, id, template string) {
	c.configs[roomName].templates[id] = template
}

// HasTemplate verifies if a template really exists for a room.
func (c *CherryRooms) HasTemplate(roomName, id string) bool {
	_, ok := c.configs[roomName].templates[id]
	return ok
}

// SetJoinMessage sets the join message.
func (c *CherryRooms) SetJoinMessage(roomName, message string) {
	c.configs[roomName].misc.joinMessage = message
}

// SetExitMessage sets the exit message.
func (c *CherryRooms) SetExitMessage(roomName, message string) {
	c.configs[roomName].misc.exitMessage = message
}

// SetOnIgnoreMessage sets the "on ignore" message.
func (c *CherryRooms) SetOnIgnoreMessage(roomName, message string) {
	c.configs[roomName].misc.onIgnoreMessage = message
}

// SetOnDeIgnoreMessage sets the "on deignore" message.
func (c *CherryRooms) SetOnDeIgnoreMessage(roomName, message string) {
	c.configs[roomName].misc.onDeIgnoreMessage = message
}

// SetGreetingMessage sets the greeting message.
func (c *CherryRooms) SetGreetingMessage(roomName, message string) {
	c.configs[roomName].misc.greetingMessage = message
}

// SetPrivateMessageMarker sets the private message marker.
func (c *CherryRooms) SetPrivateMessageMarker(roomName, marker string) {
	c.configs[roomName].misc.privateMessageMarker = marker
}

// SetMaxUsers sets the maximum of users allowed in a room.
func (c *CherryRooms) SetMaxUsers(roomName string, value int) {
	c.configs[roomName].misc.maxUsers = value
}

// SetAllowBrief sets the allow brief option.
func (c *CherryRooms) SetAllowBrief(roomName string, value bool) {
	c.configs[roomName].misc.allowBrief = value
}

// IsAllowingBriefs verifies if briefs are allowed for a room.
func (c *CherryRooms) IsAllowingBriefs(roomName string) bool {
	return c.configs[roomName].misc.allowBrief
}

//func (c *CherryRooms) SetFloodingPolice(roomName string, value bool) {
//    c.configs[roomName].misc.floodingPolice = value
//}

//func (c *CherryRooms) SetMaxFloodAllowedBeforeKick(roomName string, value int) {
//    c.configs[roomName].misc.maxFloodAllowedBeforeKick = value
//}

// SetAllUsersAlias sets all users alias.
func (c *CherryRooms) SetAllUsersAlias(roomName, alias string) {
	c.configs[roomName].misc.allUsersAlias = alias
}

// Lock acquire the room mutex.
func (c *CherryRooms) Lock(roomName string) {
	c.configs[roomName].mutex.Lock()
}

// Unlock dispose the room mutex.
func (c *CherryRooms) Unlock(roomName string) {
	c.configs[roomName].mutex.Unlock()
}

// GetServername spits the server name.
func (c *CherryRooms) GetServername() string {
	return c.servername
}

// SetServername sets the server name.
func (c *CherryRooms) SetServername(servername string) {
	c.servername = servername
}

// HasUser verifies if the user is connected in the room.
func (c *CherryRooms) HasUser(roomName, user string) bool {
	_, ok := c.configs[roomName]
	if !ok {
		return false
	}
	_, ok = c.configs[roomName].users[user]
	return ok
}

// IsValidUserRequest verifies if the session ID really matches with the previously defined.
func (c *CherryRooms) IsValidUserRequest(roomName, user, id string, userConn net.Conn) bool {
	var valid = false
	if c.HasUser(roomName, user) {
		valid = (id == c.GetSessionID(user, roomName))
		if valid {
			c.Lock(roomName)
			userAddr := strings.Split(userConn.RemoteAddr().String(), ":")
			realAddr := c.configs[roomName].users[user].addr
			c.Unlock(roomName)
			if len(realAddr) > 0 && len(userAddr) > 0 {
				valid = (realAddr == userAddr[0])
			}
		}
	}
	return valid
}

// SetIgnoreAction sets the action that will be used for ignoring.
func (c *CherryRooms) SetIgnoreAction(roomName, action string) {
	c.Lock(roomName)
	c.configs[roomName].ignoreAction = action
	c.Unlock(roomName)
}

// SetDeIgnoreAction sets the action that will be used for "deignoring".
func (c *CherryRooms) SetDeIgnoreAction(roomName, action string) {
	c.Lock(roomName)
	c.configs[roomName].deignoreAction = action
	c.Unlock(roomName)
}

// GetIgnoreAction returns the action that represents the ignoring.
func (c *CherryRooms) GetIgnoreAction(roomName string) string {
	c.Lock(roomName)
	var retval string
	retval = c.configs[roomName].ignoreAction
	c.Unlock(roomName)
	return retval
}

// GetDeIgnoreAction returns the action that represents the "deignoring".
func (c *CherryRooms) GetDeIgnoreAction(roomName string) string {
	c.Lock(roomName)
	var retval string
	retval = c.configs[roomName].deignoreAction
	c.Unlock(roomName)
	return retval
}

// SetUserConnection registers a connection for a user recently enrolled in a room.
func (c *CherryRooms) SetUserConnection(roomName, user string, conn net.Conn) {
	c.Lock(roomName)
	c.configs[roomName].users[user].conn = conn
	remoteAddr := strings.Split(conn.RemoteAddr().String(), ":")
	if len(remoteAddr) > 0 {
		c.configs[roomName].users[user].addr = remoteAddr[0]
	}
	c.Unlock(roomName)
}

// GetServerName spits the server name.
func (c *CherryRooms) GetServerName() string {
	return c.servername
}
