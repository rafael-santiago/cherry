/*
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package html

import (
    "../config"
    "strings"
    "time"
    "fmt"
)

type Preprocessor struct {
    rooms *config.CherryRooms
    data_expander map[string]func(*Preprocessor, string, string, string) string
    data_value map[string]string
}

func NewHtmlPreprocessor(rooms *config.CherryRooms) *Preprocessor {
    var preprocessor *Preprocessor
    preprocessor = new(Preprocessor)
    preprocessor.Init(rooms)
    return preprocessor
}

func (p *Preprocessor) SetDataValue(field, data string) {
    p.data_value[field] = data
}

func (p *Preprocessor) UnsetDataValue(field string) {
    p.data_value[field] = ""
}

func (p *Preprocessor) Init(rooms *config.CherryRooms) {
    p.rooms = rooms
    p.data_value = make(map[string]string)
    p.data_expander = make(map[string]func(*Preprocessor, string, string, string) string)
    p.data_expander["{{.nickname}}"] = nick_name_expander
    p.data_expander["{{.session-id}}"] = session_id_expander
    p.data_expander["{{.color}}"] = color_expander
    p.data_expander["{{.ignore-list}}"] = ignore_list_expander
    p.data_expander["{{.hour}}"] = hour_expander
    p.data_expander["{{.minute}}"] = minute_expander
    p.data_expander["{{.second}}"] = second_expander
//    p.data_expander["{{.month}}"] = month_expander
//    p.data_expander["{{.day}}"] = day_expander
//    p.data_expander["{{.year}}"] = year_expander
    p.data_expander["{{.greeting-message}}"] = greeting_message_expander
    p.data_expander["{{.join-message}}"] = join_message_expander
    p.data_expander["{{.exit-message}}"] = exit_message_expander
    p.data_expander["{{.on-ignore-message}}"] = on_ignore_message_expander
    p.data_expander["{{.on-deignore-message}}"] = on_deignore_message_expander
    p.data_expander["{{.max-users}}"] = max_users_expander
    p.data_expander["{{.all-users-alias}}"] = all_users_alias_expander
    p.data_expander["{{.action-list}}"] = action_list_expander
    p.data_expander["{{.image-list}}"] = image_list_expander
//    p.data_expander["{{.sound-list}}"] = sound_list_expander
    p.data_expander["{{.users-list}}"] = users_list_expander
    p.data_expander["{{.top-template}}"] = top_template_expander
    p.data_expander["{{.body-template}}"] = body_template_expander
    p.data_expander["{{.banner-template}}"] = banner_template_expander
    p.data_expander["{{.highlight-template}}"] = highlight_template_expander
    p.data_expander["{{.entrance-template}}"] = entrance_template_expander
    p.data_expander["{{.exit-template}}"] = exit_template_expander
    p.data_expander["{{.nickclash-template}}"] = nickclash_template_expander
    p.data_expander["{{.last-public-messages}}"] = last_public_messages_expander
    p.data_expander["{{.servername}}"] = servername_expander
    p.data_expander["{{.listen-port}}"] = listen_port_expander
    p.data_expander["{{.room-name}}"] = room_name_expander
    p.data_expander["{{.users-total}}"] = users_total_expander
    p.data_expander["{{.message-action-label}}"] = message_action_label_expander
    p.data_expander["{{.message-whoto}}"] = message_whoto_expander
    p.data_expander["{{.message-user}}"] = nick_name_expander
    p.data_expander["{{.message-colored-user}}"] = colored_nick_name_expander
    p.data_expander["{{.message-says}}"] = message_says_expander
//    p.data_expander["{{.message-sound}}"] = message_sound_expander
    p.data_expander["{{.message-image}}"] = message_image_expander
    p.data_expander["{{.message-private-marker}}"] = message_private_marker_expander
    p.data_expander["{{.current-formatted-message}}"] = nil
    p.data_expander["{{.priv}}"] = nil
    p.data_expander["{{.brief-last-public-messages}}"] = brief_last_public_messages
    p.data_expander["{{.brief-who-are-talking}}"] = brief_who_are_talking
    p.data_expander["{{.brief-users-total}}"] = brief_users_total
    p.data_expander["{{.find-result-user}}"] = nil
    p.data_expander["{{.find-result-room-name}}"] = nil
    p.data_expander["{{.find-result-users-total}}"] = nil
}

func (p *Preprocessor) ExpandData(room_name, data string) string {
    if p.rooms.HasRoom(room_name) {
        for var_name, expander := range p.data_expander {
            local_value, exists := p.data_value[var_name]
            if exists && len(local_value) > 0 {
                data = strings.Replace(data, var_name, local_value, -1)
            } else {
                if expander == nil {
                    continue
                }
                data = expander(p, room_name, var_name, data)
            }
        }
    }
    return data
}

func brief_users_total(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetUsersTotal(room_name), -1)
}

func brief_who_are_talking(p *Preprocessor, room_name, var_name, data string) string {
    var users []string = p.rooms.GetRoomUsers(room_name)
    var table_data string
    table_data = "<table border = 0>"
    for _, u := range users {
        table_data += "\n\t<tr><td>" + u + "</td></tr>"
    }
    table_data += "\n</table>"
    return strings.Replace(data, var_name, table_data, -1)
}

func brief_last_public_messages(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetLastPublicMessages(room_name), -1)
}

func message_action_label_expander(p *Preprocessor, room_name, var_name, data string) string {
    action := p.rooms.GetNextMessage(room_name).Action
    if !p.rooms.HasAction(room_name, action) {
        return data
    }
    return strings.Replace(data, var_name, p.rooms.GetRoomActionLabel(room_name, action), -1)
}

func message_whoto_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetNextMessage(room_name).To, -1)
}

func message_says_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetNextMessage(room_name).Say, -1)
}

//func message_sound_expander(p *Preprocessor, room_name, var_name, data string) string {
//    sound := p.rooms.GetNextMessage(room_name).Sound
//    if len(sound) > 0 {
//    }
//    return strings.Replace(data, var_name, sound, -1)
//}

func message_image_expander(p *Preprocessor, room_name, var_name, data string) string {
    image := p.rooms.GetNextMessage(room_name).Image
    if len(image) > 0 {
        image = "<br><img src = \"" + image + "\">"
    }
    return strings.Replace(data, var_name, image, -1)
}

func nick_name_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetNextMessage(room_name).From, -1)
}

func get_hex_color(cl_key string) string {
    var hex_colors map[string]string = make(map[string]string)
    hex_colors["0"] = "#000000"
    hex_colors["1"] = "#d10019"
    hex_colors["2"] = "#0d7000"
    hex_colors["3"] = "#c0c0c0"
    hex_colors["4"] = "#b533ff"
    hex_colors["5"] = "#ff3db5"
    hex_colors["6"] = "#0019d1"
    hex_colors["7"] = "#3de5ff"
    return hex_colors[cl_key]
}

func colored_nick_name_expander(p *Preprocessor, room_name, var_name, data string) string {
    color := p.rooms.GetColor(p.rooms.GetNextMessage(room_name).From, room_name)
    colored_nick_name := "<font color = \"" + get_hex_color(color) + "\">" + p.rooms.GetNextMessage(room_name).From + "</font>"
    return strings.Replace(data, var_name, colored_nick_name, -1)
}

func session_id_expander(p *Preprocessor, room_name, var_name, data string) string {
    from := p.rooms.GetNextMessage(room_name).From
    return strings.Replace(data, var_name, p.rooms.GetSessionId(from, room_name), -1)
}

func color_expander(p *Preprocessor, room_name, var_name, data string) string {
    from := p.rooms.GetNextMessage(room_name).From
    return strings.Replace(data, var_name, p.rooms.GetColor(from, room_name), -1)
}

func ignore_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    from := p.rooms.GetNextMessage(room_name).From
    return strings.Replace(data, var_name, p.rooms.GetIgnoreList(from, room_name), -1)
}

func hour_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, fmt.Sprintf("%.2d", time.Now().Hour()), -1)
}

func minute_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, fmt.Sprintf("%.2d", time.Now().Minute()), -1)
}

func second_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, fmt.Sprintf("%.2d", time.Now().Second()), -1)
}

func greeting_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetGreetingMessage(room_name), -1)
}

func join_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetJoinMessage(room_name), -1)
}

func exit_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetExitMessage(room_name), -1)
}

func on_ignore_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetOnIgnoreMessage(room_name), -1)
}

func on_deignore_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetOnDeIgnoreMessage(room_name), -1)
}

func message_private_marker_expander(p *Preprocessor, room_name, var_name, data string) string {
    var private_marker string = ""
    if p.rooms.GetNextMessage(room_name).Priv == "1" {
        private_marker = p.rooms.GetPrivateMessageMarker(room_name)
    }
    return strings.Replace(data, var_name, private_marker, -1)
}

func max_users_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetMaxUsers(room_name), -1)
}

func all_users_alias_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetAllUsersAlias(room_name), -1)
}

func action_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetActionList(room_name), -1)
}

func image_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetImageList(room_name), -1)
}

//func sound_list_expander(p *Preprocessor, room_name, var_name, data string) string {
//    return strings.Replace(data, var_name, p.rooms.GetSoundList(room_name), -1)
//}

func users_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetUsersList(room_name), -1)
}

func top_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetTopTemplate(room_name), -1)
}

func body_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetBodyTemplate(room_name), -1)
}

func banner_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetBannerTemplate(room_name), -1)
}

func highlight_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetHighlightTemplate(room_name), -1)
}

func entrance_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetEntranceTemplate(room_name), -1)
}

func exit_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetExitTemplate(room_name), -1)
}

func nickclash_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetNickclashTemplate(room_name), -1)
}

func last_public_messages_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetLastPublicMessages(room_name), -1)
}

func servername_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetServername(), -1)
}

func listen_port_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetListenPort(room_name), -1)
}

func room_name_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, room_name, -1)
}

func users_total_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetUsersTotal(room_name), -1)
}

func GetBadAssErrorData() string {
    return "<html><h1>404 Bad ass error</h1><h3>No cherry for you!</h3></html>"
}
