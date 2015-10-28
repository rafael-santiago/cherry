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
)

type Preprocessor struct {
    room_configs map[string]*config.RoomConfig
    data_expander map[string]func(string, string, string) string
}

func NewHtmlPreprocessor(room_configs map[string]*config.RoomConfig) *Preprocessor {
    var preprocessor *Preprocessor
    preprocessor = new(Preprocessor)
    preprocessor.Init(room_configs)
    return preprocessor
}

func (p *Preprocessor) Init(room_configs map[string]*config.RoomConfig) {
    p.room_configs = room_configs
    p.data_expander = make(map[string]func(string, string, string) string)
    p.data_expander["{{.nickname}}"] = nil
    p.data_expander["{{.session-id}}"] = nil
    p.data_expander["{{.color}}"] = nil
    p.data_expander["{{.ignore-list}}"] = nil
    p.data_expander["{{.hour}}"] = nil
    p.data_expander["{{.minute}}"] = nil
    p.data_expander["{{.second}}"] = nil
    p.data_expander["{{.month}}"] = nil
    p.data_expander["{{.day}}"] = nil
    p.data_expander["{{.year}}"] = nil
    p.data_expander["{{.greeting-message}}"] = nil
    p.data_expander["{{.join-message}}"] = nil
    p.data_expander["{{.exit-message}}"] = nil
    p.data_expander["{{.on-ignore-message}}"] = nil
    p.data_expander["{{.on-deignore-message}}"] = nil
    p.data_expander["{{.private-message-marker}}"] = nil
    p.data_expander["{{.max-users}}"] = nil
    p.data_expander["{{.all-users-alias}}"] = nil
    p.data_expander["{{.room-action-list}}"] = nil
    p.data_expander["{{.room-image-list}}"] = nil
    p.data_expander["{{.room-sound-list}}"] = nil
    p.data_expander["{{.users-list}}"] = nil
    p.data_expander["{{.top-template}}"] = nil
    p.data_expander["{{.body-template}}"] = nil
    p.data_expander["{{.banner-template}}"] = nil
    p.data_expander["{{.highlight-template}}"] = nil
    p.data_expander["{{.entrance-template}}"] = nil
    p.data_expander["{{.exit-template}}"] = nil
    p.data_expander["{{.nickclash-template}}"] = nil
    p.data_expander["{{.last-public-messages}}"] = nil
    p.data_expander["{{.servername}}"] = nil
    p.data_expander["{{.listen-port}}"] = nil
    p.data_expander["{{.room-name}}"] = nil
    p.data_expander["{{.users-total}}"] = nil
}

func (p *Preprocessor) ExpandData(room_name, data string) string {
    _, exists := p.room_configs[room_name]
    if ! exists {
        return data
    }
    for var_name, expander := range p.data_expander {
        data = expander(room_name, var_name, data)
    }
    return data
}