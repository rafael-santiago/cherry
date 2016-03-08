/*
Package html implements expanders based on the special makers recognized by Cherry.
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
*/
package html

import (
	"fmt"
	"pkg/config"
	"strings"
	"time"
)

// Preprocessor is returned by NewHTMLPreprocessor and it actually performs
// the content expansion.
type Preprocessor struct {
	rooms        *config.CherryRooms
	dataExpander map[string]func(*Preprocessor, string, string, string) string
	dataValue    map[string]string
}

// NewHTMLPreprocessor creates a new HTML preprocessor.
func NewHTMLPreprocessor(rooms *config.CherryRooms) *Preprocessor {
	var preprocessor *Preprocessor
	preprocessor = new(Preprocessor)
	preprocessor.Init(rooms)
	return preprocessor
}

// SetDataValue sets a statical content for a special marker.
func (p *Preprocessor) SetDataValue(field, data string) {
	p.dataValue[field] = data
}

// UnsetDataValue removes a previous data set by SetDataValue.
func (p *Preprocessor) UnsetDataValue(field string) {
	p.dataValue[field] = ""
}

// Init sets all default expanders.
func (p *Preprocessor) Init(rooms *config.CherryRooms) {
	p.rooms = rooms
	p.dataValue = make(map[string]string)
	p.dataExpander = make(map[string]func(*Preprocessor, string, string, string) string)
	p.dataExpander["{{.nickname}}"] = nicknameExpander
	p.dataExpander["{{.session-id}}"] = sessionIDExpander
	p.dataExpander["{{.color}}"] = colorExpander
	p.dataExpander["{{.ignore-list}}"] = ignoreListExpander
	p.dataExpander["{{.hour}}"] = hourExpander
	p.dataExpander["{{.minute}}"] = minuteExpander
	p.dataExpander["{{.second}}"] = secondExpander
	//    p.dataExpander["{{.month}}"] = month_expander
	//    p.dataExpander["{{.day}}"] = day_expander
	//    p.dataExpander["{{.year}}"] = year_expander
	p.dataExpander["{{.greeting-message}}"] = greetingMessageExpander
	p.dataExpander["{{.join-message}}"] = joinMessageExpander
	p.dataExpander["{{.exit-message}}"] = exitMessageExpander
	p.dataExpander["{{.on-ignore-message}}"] = onIgnoreMessageExpander
	p.dataExpander["{{.on-deignore-message}}"] = onDeIgnoreMessageExpander
	p.dataExpander["{{.max-users}}"] = maxUsersExpander
	p.dataExpander["{{.all-users-alias}}"] = allUsersAliasExpander
	p.dataExpander["{{.action-list}}"] = actionListExpander
	p.dataExpander["{{.image-list}}"] = imageListExpander
	//    p.dataExpander["{{.sound-list}}"] = sound_list_expander
	p.dataExpander["{{.users-list}}"] = usersListExpander
	p.dataExpander["{{.top-template}}"] = topTemplateExpander
	p.dataExpander["{{.body-template}}"] = bodyTemplateExpander
	p.dataExpander["{{.banner-template}}"] = bannerTemplateExpander
	p.dataExpander["{{.highlight-template}}"] = highlightTemplateExpander
	p.dataExpander["{{.entrance-template}}"] = entranceTemplateExpander
	p.dataExpander["{{.exit-template}}"] = exitTemplateExpander
	p.dataExpander["{{.nickclash-template}}"] = nickclashTemplateExpander
	p.dataExpander["{{.last-public-messages}}"] = lastPublicMessagesExpander
	p.dataExpander["{{.servername}}"] = servernameExpander
	p.dataExpander["{{.listen-port}}"] = listenPortExpander
	p.dataExpander["{{.room-name}}"] = roomNameExpander
	p.dataExpander["{{.users-total}}"] = usersTotalExpander
	p.dataExpander["{{.message-action-label}}"] = messageActionLabelExpander
	p.dataExpander["{{.message-whoto}}"] = messageWhotoExpander
	p.dataExpander["{{.message-user}}"] = nicknameExpander
	p.dataExpander["{{.message-colored-user}}"] = coloredNicknameExpander
	p.dataExpander["{{.message-says}}"] = messageSaysExpander
	//    p.dataExpander["{{.message-sound}}"] = message_sound_expander
	p.dataExpander["{{.message-image}}"] = messageImageExpander
	p.dataExpander["{{.message-private-marker}}"] = messagePrivateMarkerExpander
	p.dataExpander["{{.current-formatted-message}}"] = nil
	p.dataExpander["{{.priv}}"] = nil
	p.dataExpander["{{.brief-last-public-messages}}"] = briefLastPublicMessagesExpander
	p.dataExpander["{{.brief-who-are-talking}}"] = briefWhoAreTalkingExpander
	p.dataExpander["{{.brief-users-total}}"] = briefUsersTotalExpander
	p.dataExpander["{{.find-result-user}}"] = nil
	p.dataExpander["{{.find-result-room-name}}"] = nil
	p.dataExpander["{{.find-result-users-total}}"] = nil
}

// ExpandData gives preference for statical data if it does not exist the data is processed by expanders.
func (p *Preprocessor) ExpandData(roomName, data string) string {
	if p.rooms.HasRoom(roomName) {
		for varName, expander := range p.dataExpander {
			localValue, exists := p.dataValue[varName]
			if exists && len(localValue) > 0 {
				data = strings.Replace(data, varName, localValue, -1)
			} else {
				if expander == nil {
					continue
				}
				data = expander(p, roomName, varName, data)
			}
		}
	}
	return data
}

func expandImageRefs(data string) string {
	var retData string
	dataLen := len(data)
	for d := 0; d < dataLen; {
		if data[d] == '[' {
			var uri string
			d++
			for d < dataLen && data[d] != ']' {
				uri += string(data[d])
				d++
			}
			if strings.HasSuffix(uri, ".gif") ||
				strings.HasSuffix(uri, ".jpg") ||
				strings.HasSuffix(uri, ".jpeg") ||
				strings.HasSuffix(uri, ".png") ||
				strings.HasSuffix(uri, ".bmp") {
				retData += "<img src = \"" + uri + "\">"
			}
		} else {
			retData += string(data[d])
		}
		d++
	}
	return retData
}

// GetBadAssErrorData spits the default 404 Cherry's document.
func GetBadAssErrorData() string {
	return "<html><h1>404 Bad ass error</h1><h3>No cherry for you!</h3></html>"
}

func briefUsersTotalExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetUsersTotal(roomName), -1)
}

func briefWhoAreTalkingExpander(p *Preprocessor, roomName, varName, data string) string {
	var users = p.rooms.GetRoomUsers(roomName)
	var tableData string
	tableData = "<table border = 0>"
	for _, u := range users {
		tableData += "\n\t<tr><td>" + u + "</td></tr>"
	}
	tableData += "\n</table>"
	return strings.Replace(data, varName, tableData, -1)
}

func briefLastPublicMessagesExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetLastPublicMessages(roomName), -1)
}

func messageActionLabelExpander(p *Preprocessor, roomName, varName, data string) string {
	action := p.rooms.GetNextMessage(roomName).Action
	if !p.rooms.HasAction(roomName, action) {
		return data
	}
	return strings.Replace(data, varName, p.rooms.GetRoomActionLabel(roomName, action), -1)
}

func messageWhotoExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetNextMessage(roomName).To, -1)
}

func messageSaysExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, expandImageRefs(p.rooms.GetNextMessage(roomName).Say), -1)
}

//func message_sound_expander(p *Preprocessor, roomName, varName, data string) string {
//    sound := p.rooms.GetNextMessage(roomName).Sound
//    if len(sound) > 0 {
//    }
//    return strings.Replace(data, varName, sound, -1)
//}

func messageImageExpander(p *Preprocessor, roomName, varName, data string) string {
	image := p.rooms.GetNextMessage(roomName).Image
	if len(image) > 0 {
		image = "<br><img src = \"" + image + "\">"
	}
	return strings.Replace(data, varName, image, -1)
}

func nicknameExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetNextMessage(roomName).From, -1)
}

func getHexColor(clKey string) string {
	var hexColors = make(map[string]string)
	hexColors["0"] = "#000000"
	hexColors["1"] = "#d10019"
	hexColors["2"] = "#0d7000"
	hexColors["3"] = "#c0c0c0"
	hexColors["4"] = "#b533ff"
	hexColors["5"] = "#ff3db5"
	hexColors["6"] = "#0019d1"
	hexColors["7"] = "#3de5ff"
	return hexColors[clKey]
}

func coloredNicknameExpander(p *Preprocessor, roomName, varName, data string) string {
	color := p.rooms.GetColor(p.rooms.GetNextMessage(roomName).From, roomName)
	coloredNickname := "<font color = \"" + getHexColor(color) + "\">" + p.rooms.GetNextMessage(roomName).From + "</font>"
	return strings.Replace(data, varName, coloredNickname, -1)
}

func sessionIDExpander(p *Preprocessor, roomName, varName, data string) string {
	from := p.rooms.GetNextMessage(roomName).From
	return strings.Replace(data, varName, p.rooms.GetSessionID(from, roomName), -1)
}

func colorExpander(p *Preprocessor, roomName, varName, data string) string {
	from := p.rooms.GetNextMessage(roomName).From
	return strings.Replace(data, varName, p.rooms.GetColor(from, roomName), -1)
}

func ignoreListExpander(p *Preprocessor, roomName, varName, data string) string {
	from := p.rooms.GetNextMessage(roomName).From
	return strings.Replace(data, varName, p.rooms.GetIgnoreList(from, roomName), -1)
}

func hourExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, fmt.Sprintf("%.2d", time.Now().Hour()), -1)
}

func minuteExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, fmt.Sprintf("%.2d", time.Now().Minute()), -1)
}

func secondExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, fmt.Sprintf("%.2d", time.Now().Second()), -1)
}

func greetingMessageExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetGreetingMessage(roomName), -1)
}

func joinMessageExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetJoinMessage(roomName), -1)
}

func exitMessageExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetExitMessage(roomName), -1)
}

func onIgnoreMessageExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetOnIgnoreMessage(roomName), -1)
}

func onDeIgnoreMessageExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetOnDeIgnoreMessage(roomName), -1)
}

func messagePrivateMarkerExpander(p *Preprocessor, roomName, varName, data string) string {
	var privateMarker string
	if p.rooms.GetNextMessage(roomName).Priv == "1" {
		privateMarker = p.rooms.GetPrivateMessageMarker(roomName)
	}
	return strings.Replace(data, varName, privateMarker, -1)
}

func maxUsersExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetMaxUsers(roomName), -1)
}

func allUsersAliasExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetAllUsersAlias(roomName), -1)
}

func actionListExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetActionList(roomName), -1)
}

func imageListExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetImageList(roomName), -1)
}

//func sound_list_expander(p *Preprocessor, roomName, varName, data string) string {
//    return strings.Replace(data, varName, p.rooms.GetSoundList(roomName), -1)
//}

func usersListExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetUsersList(roomName), -1)
}

func topTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetTopTemplate(roomName), -1)
}

func bodyTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetBodyTemplate(roomName), -1)
}

func bannerTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetBannerTemplate(roomName), -1)
}

func highlightTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetHighlightTemplate(roomName), -1)
}

func entranceTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetEntranceTemplate(roomName), -1)
}

func exitTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetExitTemplate(roomName), -1)
}

func nickclashTemplateExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetNickclashTemplate(roomName), -1)
}

func lastPublicMessagesExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetLastPublicMessages(roomName), -1)
}

func servernameExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetServername(), -1)
}

func listenPortExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetListenPort(roomName), -1)
}

func roomNameExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, roomName, -1)
}

func usersTotalExpander(p *Preprocessor, roomName, varName, data string) string {
	return strings.Replace(data, varName, p.rooms.GetUsersTotal(roomName), -1)
}
