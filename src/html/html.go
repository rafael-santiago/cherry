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
)

type Preprocessor struct {
    rooms *config.CherryRooms
    data_expander map[string]func(*Preprocessor, string, string, string) string
}

func NewHtmlPreprocessor(rooms *config.CherryRooms) *Preprocessor {
    var preprocessor *Preprocessor
    preprocessor = new(Preprocessor)
    preprocessor.Init(rooms)
    return preprocessor
}



func (p *Preprocessor) Init(rooms *config.CherryRooms) {
    p.rooms = rooms
    p.data_expander = make(map[string]func(*Preprocessor, string, string, string) string)
    p.data_expander["{{.nickname}}"] = nick_name_expander
    p.data_expander["{{.session-id}}"] = session_id_expander
    p.data_expander["{{.color}}"] = color_expander
    p.data_expander["{{.ignore-list}}"] = ignore_list_expander
    p.data_expander["{{.hour}}"] = hour_expander
    p.data_expander["{{.minute}}"] = minute_expander
    p.data_expander["{{.second}}"] = second_expander
    p.data_expander["{{.month}}"] = month_expander
    p.data_expander["{{.day}}"] = day_expander
    p.data_expander["{{.year}}"] = year_expander
    p.data_expander["{{.greeting-message}}"] = greeting_message_expander
    p.data_expander["{{.join-message}}"] = join_message_expander
    p.data_expander["{{.exit-message}}"] = exit_message_expander
    p.data_expander["{{.on-ignore-message}}"] = on_ignore_message_expander
    p.data_expander["{{.on-deignore-message}}"] = on_deignore_message_expander
    p.data_expander["{{.private-message-marker}}"] = private_message_marker_expander
    p.data_expander["{{.max-users}}"] = max_users_expander
    p.data_expander["{{.all-users-alias}}"] = all_users_alias_expander
    p.data_expander["{{.action-list}}"] = action_list_expander
    p.data_expander["{{.image-list}}"] = image_list_expander
    p.data_expander["{{.sound-list}}"] = sound_list_expander
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
}

func (p *Preprocessor) ExpandData(room_name, data string) string {
    if p.rooms.HasRoom(room_name) {
        for var_name, expander := range p.data_expander {
            data = expander(p, room_name, var_name, data)
        }
    }
    return data
}

func nick_name_expander(p *Preprocessor, room_name, var_name, data string) string {
    return strings.Replace(data, var_name, p.rooms.GetNextMessage(room_name).From, -1)
}

func session_id_expander(p *Preprocessor, room_name, var_name, data string) string {
//  TODO(Santiago):    return strings.Replace(data, var_name, p.rooms.GetSessionId(from, room_name).Id, -1)
    return ""
}

func color_expander(p *Preprocessor, room_name, var_name, data string) string {
//  TODO(Santiago):    return strings.Replace(data, var_name, p.rooms.GetColor(from, room_name).Color, -1)
    return ""
}

func ignore_list_expander(p *Preprocessor, room_name, var_name, data string) string {
//  TODO(Santiago):    return strings.Replace(data, var_name, p.rooms.GetIgnoreList(from, room_name), -1)
    return ""
}

func hour_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func minute_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func second_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func month_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func day_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func year_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func greeting_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func join_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func exit_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func on_ignore_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func on_deignore_message_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func private_message_marker_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func max_users_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func all_users_alias_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func action_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func image_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func sound_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func users_list_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func top_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func body_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func banner_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func highlight_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func entrance_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func exit_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func nickclash_template_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func last_public_messages_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func servername_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func listen_port_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}

func room_name_expander(p *Preprocessor, room_name, var_name, data string) string {
    return room_name
}

func users_total_expander(p *Preprocessor, room_name, var_name, data string) string {
    return ""
}
