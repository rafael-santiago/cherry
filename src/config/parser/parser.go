/*
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package parser

import (
    "../../config"
    "fmt"
    "strings"
    "io/ioutil"
    "strconv"
)

type CherryFileError struct {
    src string
    line int
    msg string
}

func (c *CherryFileError) Error() string {
    if c.line > -1 {
        return fmt.Sprintf("ERROR: %s: at line %d: %s\n", c.src, c.line, c.msg);
    }
    return fmt.Sprintf("ERROR: %s: %s\n", c.src, c.msg);
}

func NewCherryFileError(src string, line int, msg string) *CherryFileError {
    return &CherryFileError{src, line, msg}
}

func GetDataFromSection(section, config_data string, curr_line int, curr_file string) (string, int, int, *CherryFileError) {
    var s int
    var temp string
    for s = 0; s < len(config_data); s++ {
        switch config_data[s] {
            case '#':
                for config_data[s] != '\n' && s < len(config_data) {
                    s++
                }
                if s < len(config_data) {
                    curr_line++;
                }
                continue
            case '(', '\n', ' ', '\t':
                if config_data[s] == '\n' {
                    curr_line++;
                }
                if temp == section {
                    if config_data[s] == '\n' || config_data[s] == ' ' || config_data[s] == '\t' {
                        for s < len(config_data) && config_data[s] != '(' {
                            s++
                            if s < len(config_data) && config_data[s] == '\n' {
                                curr_line++
                            }
                        }
                    }
                    if s < len(config_data) && config_data[s] == '(' {
                        s++
                    }
                    var data string
                    for s < len(config_data) {
                        if config_data[s] == '"' {
                            data += string(config_data[s])
                            s++
                            for s < len(config_data) && config_data[s] != '"' {
                                if config_data[s] != '\\' {
                                    data += string(config_data[s])
                                } else {
                                    data += string(config_data[s + 1]);
                                    s++
                                }
                                s++
                            }
                            if s < len(config_data) {
                                data += string(config_data[s])
                            }
                        } else if config_data[s] != ')' {
                            data += string(config_data[s])
                        } else {
                            break
                        }
                        s++
                    }
                    return data, s, curr_line, nil
                } else if temp == "cherry.branch" {
                    for s < len(config_data) && (config_data[s] == ' ' || config_data[s] == '\t') {
                        s++
                    }
                    if s < len(config_data) {
                        var branch_filepath string
                        for s < len(config_data) && config_data[s] != '\n' {
                            branch_filepath += string(config_data[s])
                            s++
                        }
                        if s < len(config_data) {
                            curr_line++
                        }
                        branch_buffer, err := ioutil.ReadFile(branch_filepath)
                        if err != nil {
                            fmt.Println(fmt.Sprintf("WARNING: %s: at line %d: %s. Be tidy... removing or commenting this dry branch from your cherry.", curr_file, curr_line - 1, err.Error()))
                            //return "", s, curr_line, NewCherryFileError(curr_file,
                            //                                            curr_line - 1,
                            //                                            "unable to read cherry.branch from \"" + branch_filepath + "\" [ details: " + err.Error() + " ]")
                        } else {
                            branch_data, branch_offset, branch_line, _ := GetDataFromSection(section, string(branch_buffer), 1, branch_filepath)
                            if len(branch_data) > 0 {
                                return branch_data, branch_offset, branch_line, nil
                            }
                        }
                    }
                }
                temp = ""
                break
            default:
                temp += string(config_data[s])
                break
        }
    }
    return "", s, curr_line, NewCherryFileError(curr_file, -1, "section \"" + section + "\" not found.")
}

func GetNextSetFromData(data string, curr_line int, tok string) ([]string, int, string) {
    if len(data) == 0 {
        return make([]string, 0), curr_line, ""
    }
    var s int
    for s = 0; s < len(data) && (data[s] == ' ' || data[s] == '\t' || data[s] == '\n'); s++ {
        if data[s] == '\n' {
            curr_line++
        }
    }
    var line string
    for s < len(data) && data[s] != '\n' {
        if data[s] == '"' {
            line += string(data[s])
            s++
            for s < len(data) && data[s] != '"' {
                if data[s] == '\\' {
                    s++
                }
                line += string(data[s])
                if data[s] == '\n' {
                    curr_line++
                }
                s++
            }
            if s < len(data) {
                line += string(data[s])
            }
        } else if data[s] == '#' {
            for s < len(data) && data[s] != '\n' {
                s++
            }
            if s < len(data) {
                curr_line++
            }
        } else {
            line += string(data[s])
        }
        s++
    }
    if len(line) == 0 {
        if len(data) == 0 {
            return make([]string, 0), curr_line, ""
        }
    }
    set := strings.Split(line, tok)
    if len(set) == 2 {
        set[0] = StripBlanks(set[0])
        set[1] = StripBlanks(set[1])
    }
    var next_data string
    if s < len(data) {
        next_data = data[s:]
    }
    return set, curr_line, next_data
}

func StripBlanks(data string) string {
    var retval string
    var d_start int = 0
    for d_start < len(data) && (data[d_start] == ' ' || data[d_start] == '\t') {
        d_start++
    }
    var d_end int = len(data) - 1
    for d_end > 0 && (data[d_end] == ' ' || data[d_end] == '\t') {
        d_end--
    }
    retval = data[d_start:d_end+1]
    return retval
}

func ParseCherryFile(filepath string) (*config.CherryRooms, *CherryFileError) {
    var cherry_rooms *config.CherryRooms = nil
    var cherry_file_data []byte
    var data string
    var err *CherryFileError
    var line int
    cherry_file_data, io_err := ioutil.ReadFile(filepath)
    if io_err != nil {
        return nil, NewCherryFileError("(no file)", -1, fmt.Sprintf("unable to read from \"%s\" [more details: %s].", filepath, io_err.Error()))
    }
    data, _, line,  err = GetDataFromSection("cherry.rooms", string(cherry_file_data), 1, filepath)
    if err != nil {
        return nil, err
    }
    //  INFO(Santiago): Adding all scanned rooms from the first cherry.rooms section found
    //                  [cherry branches were scanned too at this point].
    var set []string
    cherry_rooms = config.NewCherryRooms()
    set, line, data = GetNextSetFromData(data, line, ":")
    for len(set) == 2 {
        if cherry_rooms.HasRoom(set[0]) {
            return nil, NewCherryFileError(filepath, line, fmt.Sprintf("room \"%s\" redeclared.", set[0]))
        }
        var value int64
        var conv_err error
        value, conv_err = strconv.ParseInt(set[1], 10, 16)
        if conv_err != nil {
            return nil, NewCherryFileError(filepath, line, fmt.Sprintf("invalid port value \"%s\" [more details: %s].", set[1], conv_err))
        }
        var port int16
        port = int16(value)
        if cherry_rooms.PortBusyByAnotherRoom(port) {
            return nil, NewCherryFileError(filepath, line, fmt.Sprintf("the port \"%s\" is already busy by another room.", set[1]))
        }

        cherry_rooms.AddRoom(set[0], port)

        err_room_config := GetRoomTemplates(set[0], cherry_rooms, string(cherry_file_data), filepath)
        if err_room_config != nil {
            return nil, err_room_config
        }

        err_room_config = GetRoomActions(set[0], cherry_rooms, string(cherry_file_data), filepath)
        if err_room_config != nil {
            return nil, err_room_config
        }

        //  INFO(Santiago): until now these two following sections are non-mandatory.

        _ = GetRoomImages(set[0], cherry_rooms, string(cherry_file_data), filepath)

        _ = GetRoomSounds(set[0], cherry_rooms, string(cherry_file_data), filepath)

        err_room_config = GetRoomMisc(set[0], cherry_rooms, string(cherry_file_data), filepath)
        if err_room_config != nil {
            return nil, err_room_config
        }

        //  INFO(Santiago): Let's transfer the next room from file to the memory.
        set, line, data = GetNextSetFromData(data, line, ":")
    }
    return cherry_rooms, nil
}

func GetRoomTemplates(room_name string, cherry_rooms *config.CherryRooms, config_data, filepath string) *CherryFileError {
    var data string
    var line int
    var err *CherryFileError
    data, _, line, err = GetDataFromSection("cherry." + room_name + ".templates",
                                             config_data, 1, filepath)
    if err != nil {
        return err
    }
    var set []string
    set, line, data = GetNextSetFromData(data, line, "=")
    for len(set) == 2 {
        if cherry_rooms.HasTemplate(room_name, set[0]) {
            return NewCherryFileError(filepath, line, "room template \"" + set[0] + "\" redeclared.")
        }
        if len(set[1]) == 0 {
            return NewCherryFileError(filepath, line, "room template with no value.")
        }
        if set[1][0] != '"' || set[1][len(set[1])-1] != '"' {
            return NewCherryFileError(filepath, line, "room template must be set with a valid string.")
        }
        var template_data []byte
        var template_data_err error
        template_data, template_data_err = ioutil.ReadFile(set[1][1:len(set[1])-1])
        if template_data_err != nil {
            return NewCherryFileError(filepath, line, "unable to access room template file [more details: " + template_data_err.Error() + "].")
        }
        cherry_rooms.AddTemplate(room_name, set[0], string(template_data))
        set, line, data = GetNextSetFromData(data, line, "=")
    }
    return nil
}

func GetRoomActions(room_name string, cherry_rooms *config.CherryRooms, config_data, filepath string) *CherryFileError {
    return get_indirect_config("cherry." + room_name + ".actions",
                               "cherry." + room_name + ".actions.templates",
                                room_action_main_verifier, room_action_sub_verifier, room_action_setter,
                                room_name, cherry_rooms, config_data, filepath)
}

func GetRoomImages(room_name string, cherry_rooms *config.CherryRooms, config_data, filepath string) *CherryFileError {
    return get_indirect_config("cherry." + room_name + ".images",
                               "cherry." + room_name + ".images.url",
                               room_image_main_verifier, room_image_sub_verifier, room_image_setter,
                               room_name, cherry_rooms, config_data, filepath)
}

func GetRoomSounds(room_name string, cherry_rooms *config.CherryRooms, config_data, filepath string) *CherryFileError {
    return get_indirect_config("cherry." + room_name + ".sounds",
                               "cherry." + room_name + ".sounds.url",
                               room_sound_main_verifier, room_sound_sub_verifier, room_sound_setter,
                               room_name, cherry_rooms, config_data, filepath)
}

func GetRoomMisc(room_name string, cherry_rooms *config.CherryRooms, config_data, filepath string) *CherryFileError {
    var m_data string
    var m_line int
    var m_err *CherryFileError
    m_data, _, m_line, m_err = GetDataFromSection("cherry." + room_name + ".misc", config_data, 1, filepath)
    if m_err != nil {
        return m_err
    }

    var verifier map[string]func(string) bool
    verifier = make(map[string]func(string) bool)
    verifier["join-message"]                  = verify_string
    verifier["exit-message"]                  = verify_string
    verifier["on-ignore-message"]             = verify_string
    verifier["on-deignore-message"]           = verify_string
    verifier["greeting-message"]              = verify_string
    verifier["private-message-marker"]        = verify_string
    verifier["max-users"]                     = verify_number
    verifier["allow-brief"]                   = verify_bool
    verifier["flooding-police"]               = verify_bool
    verifier["max-flood-allowed-before-kick"] = verify_number
    verifier["all-users-alias"]               = verify_string

    var setter map[string]func(*config.CherryRooms, string, string)
    setter = make(map[string]func(*config.CherryRooms, string, string))
    setter["join-message"]                  = set_join_message
    setter["exit-message"]                  = set_exit_message
    setter["on-ignore-message"]             = set_on_ignore_message
    setter["on-deignore-message"]           = set_on_deignore_message
    setter["greeting-message"]              = set_greeting_message
    setter["private-message-marker"]        = set_private_message_marker
    setter["max-users"]                     = set_max_users
    setter["allow-brief"]                   = set_allow_brief
    setter["flooding-police"]               = set_flooding_police
    setter["max-flood-allowed-before-kick"] = set_max_flood_allowed_before_kick
    setter["all-users-alias"]               = set_all_users_alias

    var already_set map[string]bool
    already_set = make(map[string]bool)
    already_set["join-message"]                  = false
    already_set["exit-message"]                  = false
    already_set["on-ignore-message"]             = false
    already_set["on-deignore-message"]           = false
    already_set["greeting-message"]              = false
    already_set["private-message-marker"]        = false
    already_set["max-users"]                     = false
    already_set["flooding-police"]               = false
    already_set["max-flood-allowed-before-kick"] = false
    already_set["all-users-alias"]               = false

    var m_set []string
    m_set, m_line, m_data = GetNextSetFromData(m_data, m_line, "=")
    for len(m_set) == 2 {
        _, exists := verifier[m_set[0]]
        if !exists {
            return NewCherryFileError(filepath, m_line, "misc configuration named as \"" + m_set[0] + "\" is unrecognized.")
        }
        if already_set[m_set[0]] {
            return NewCherryFileError(filepath, m_line, "misc configuration \"" + m_set[0] + "\" re-configured.")
        }
        if !verifier[m_set[0]](m_set[1]) {
            return NewCherryFileError(filepath, m_line, "misc configuration \"" + m_set[0] + "\" has invalid value : " + m_set[1])
        }
        setter[m_set[0]](cherry_rooms, room_name, m_set[1])
        already_set[m_set[0]] = true
        m_set, m_line, m_data = GetNextSetFromData(m_data, m_line, "=")
    }

    return nil
}

func set_join_message(cherry_rooms *config.CherryRooms, room_name, message string) {
    cherry_rooms.SetJoinMessage(room_name, message)
}

func set_exit_message(cherry_rooms *config.CherryRooms, room_name, message string) {
    cherry_rooms.SetExitMessage(room_name, message)
}

func set_on_ignore_message(cherry_rooms *config.CherryRooms, room_name, message string) {
    cherry_rooms.SetOnIgnoreMessage(room_name, message)
}

func set_on_deignore_message(cherry_rooms *config.CherryRooms, room_name, message string) {
    cherry_rooms.SetOnDeIgnoreMessage(room_name, message)
}

func set_greeting_message(cherry_rooms *config.CherryRooms, room_name, message string) {
    cherry_rooms.SetGreetingMessage(room_name, message)
}

func set_private_message_marker(cherry_rooms *config.CherryRooms, room_name, marker string) {
    cherry_rooms.SetPrivateMessageMarker(room_name, marker)
}

func set_max_users(cherry_rooms *config.CherryRooms, room_name, value string) {
    var int_value int64
    int_value, _ = strconv.ParseInt(value, 10, 64)
    cherry_rooms.SetMaxUsers(room_name, int(int_value))
}

func set_allow_brief(cherry_rooms *config.CherryRooms, room_name, value string) {
    var allow bool
    allow = (value == "yes")
    cherry_rooms.SetAllowBrief(room_name, allow)
}

func set_flooding_police(cherry_rooms *config.CherryRooms, room_name, value string) {
    var impose bool
    impose = (value == "yes")
    cherry_rooms.SetFloodingPolice(room_name, impose)
}

func set_all_users_alias(cherry_rooms *config.CherryRooms, room_name, value string) {
    cherry_rooms.SetAllUsersAlias(room_name, value)
}

func set_max_flood_allowed_before_kick(cherry_rooms *config.CherryRooms, room_name, value string) {
    var int_value int64
    int_value, _ = strconv.ParseInt(value, 10, 64)
    cherry_rooms.SetMaxFloodAllowedBeforeKick(room_name, int(int_value))
}

func verify_number(buffer string) bool {
    if len(buffer) == 0 {
        return false
    }
    for _, b := range buffer {
        if b != '0' &&
           b != '1' &&
           b != '2' &&
           b != '3' &&
           b != '4' &&
           b != '5' &&
           b != '6' &&
           b != '7' &&
           b != '8' &&
           b != '9' {
            return false
        }
    }
    return true
}

func verify_string(buffer string) bool {
    if len(buffer) <= 1 {
        return false
    }
    return (buffer[0] == '"' && buffer[len(buffer)-1] == '"')
}

func verify_bool(buffer string) bool {
    if len(buffer) == 0 {
        return false
    }
    return (buffer == "yes" || buffer == "no")
}

//  WARN(Santiago): The following codes are a brain damage. I am sorry.

func get_indirect_config(main_section,
                         sub_section string,
                         main_verifier,
                         sub_verifier func([]string, []string, int, int, string, string, *config.CherryRooms) *CherryFileError,
                         setter func(*config.CherryRooms, string, []string, []string),
                         room_name string,
                         cherry_rooms *config.CherryRooms,
                         config_data,
                         filepath string) *CherryFileError {
    var m_data string
    var m_line int
    var m_err *CherryFileError
    m_data, _, m_line, m_err = GetDataFromSection(main_section, config_data, 1, filepath)
    if m_err != nil {
        return m_err
    }

    var s_data string
    var s_line int
    var s_err *CherryFileError
    s_data, _, s_line, s_err = GetDataFromSection(sub_section, config_data, 1, filepath)

    if s_err != nil {
        return s_err
    }

    var m_set []string
    m_set, m_line, m_data = GetNextSetFromData(m_data, m_line, "=")
    for len(m_set) == 2 {
        var s_set []string
        m_err = main_verifier(m_set, s_set, m_line, s_line, room_name, filepath, cherry_rooms)
        if m_err != nil {
            return m_err
        }

        //  INFO(Santiago): Getting the template for the current action label from a section to another.
        var temp string = s_data
        var temp_line int = s_line
        s_set, temp_line, temp = GetNextSetFromData(temp, temp_line, "=")
        for len(s_set) == 2 && s_set[0] != m_set[0] {
            s_set, temp_line, temp = GetNextSetFromData(temp, temp_line, "=")
        }

        s_err = sub_verifier(m_set, s_set, m_line, s_line, room_name, filepath, cherry_rooms)

        setter(cherry_rooms, room_name, m_set, s_set)

        m_set, m_line, m_data = GetNextSetFromData(m_data, m_line, "=")
    }
    return nil
}

func room_action_main_verifier(m_set, s_set []string, m_line, s_line int, room_name, filepath string, cherry_rooms *config.CherryRooms) *CherryFileError {
    if cherry_rooms.HasAction(room_name, m_set[0]) {
        return NewCherryFileError(filepath, m_line, "room action \"" + m_set[0] + "\" redeclared.")
    }
    if len(m_set[1]) == 0 {
        return NewCherryFileError(filepath, m_line, "unlabeled room action.")
    }
    if m_set[1][0] != '"' || m_set[1][len(m_set[1])-1] != '"' {
        return NewCherryFileError(filepath, m_line, "room action must be set with a valid string.")
    }
    return nil
}

func room_action_sub_verifier(m_set, s_set []string, m_line, s_line int, room_name, filepath string, cherry_rooms *config.CherryRooms) *CherryFileError {
    if s_set[0] != m_set[0] {
        return NewCherryFileError(filepath, s_line, "there is no template for action \"" + m_set[0] + "\".")
    }
    if len(s_set[1]) == 0 {
        return NewCherryFileError(filepath, s_line, "empty room action template.")
    }
    if s_set[1][0] != '"' || s_set[1][len(s_set[1])-1] != '"' {
        return NewCherryFileError(filepath, s_line, "room action template must be set with a valid string.")
    }
    return nil
}

func room_action_setter(cherry_rooms *config.CherryRooms, room_name string, m_set, s_set []string) {
    cherry_rooms.AddAction(room_name, m_set[0], m_set[1][1:len(m_set[1])-1], s_set[1][1:len(s_set[1])-1])
}

func room_image_main_verifier(m_set, s_set []string, m_line, s_line int, room_name, filepath string, cherry_rooms *config.CherryRooms) *CherryFileError {
    if cherry_rooms.HasImage(room_name, m_set[0]) {
        return NewCherryFileError(filepath, m_line, "room image \"" + m_set[0] + "\" redeclared.")
    }
    if len(m_set[1]) == 0 {
        return NewCherryFileError(filepath, m_line, "unlabeled room image.")
    }
    if m_set[1][0] != '"' || m_set[1][len(m_set[1])-1] != '"' {
        return NewCherryFileError(filepath, m_line, "room image must be set with a valid string.")
    }
    return nil
}

func room_image_sub_verifier(m_set, s_set []string, m_line, s_line int, room_name, filepath string, cherry_rooms *config.CherryRooms) *CherryFileError {
    if s_set[0] != m_set[0] {
        return NewCherryFileError(filepath, s_line, "there is no url for image \"" + m_set[0] + "\".")
    }
    if len(s_set[1]) == 0 {
        return NewCherryFileError(filepath, s_line, "empty room image url.")
    }
    if s_set[1][0] != '"' || s_set[1][len(s_set[1])-1] != '"' {
        return NewCherryFileError(filepath, s_line, "room image url must be set with a valid string.")
    }
    return nil
}

func room_image_setter(cherry_rooms *config.CherryRooms, room_name string, m_set, s_set []string) {
    //  WARN(Santiago): by now we will pass the image template as empty.
    cherry_rooms.AddImage(room_name, m_set[0], m_set[1][1:len(m_set[1])-1], "", s_set[1][1:len(s_set[1])-1])
}

func room_sound_main_verifier(m_set, s_set []string, m_line, s_line int, room_name, filepath string, cherry_rooms *config.CherryRooms) *CherryFileError {
    if cherry_rooms.HasImage(room_name, m_set[0]) {
        return NewCherryFileError(filepath, m_line, "room sound \"" + m_set[0] + "\" redeclared.")
    }
    if len(m_set[1]) == 0 {
        return NewCherryFileError(filepath, m_line, "unlabeled room sound.")
    }
    if m_set[1][0] != '"' || m_set[1][len(m_set[1])-1] != '"' {
        return NewCherryFileError(filepath, m_line, "room sound must be set with a valid string.")
    }
    return nil
}

func room_sound_sub_verifier(m_set, s_set []string, m_line, s_line int, room_name, filepath string, cherry_rooms *config.CherryRooms) *CherryFileError {
    if s_set[0] != m_set[0] {
        return NewCherryFileError(filepath, s_line, "there is no url for sound \"" + m_set[0] + "\".")
    }
    if len(s_set[1]) == 0 {
        return NewCherryFileError(filepath, s_line, "empty room sound url.")
    }
    if s_set[1][0] != '"' || s_set[1][len(s_set[1])-1] != '"' {
        return NewCherryFileError(filepath, s_line, "room sound url must be set with a valid string.")
    }
    return nil
}

func room_sound_setter(cherry_rooms *config.CherryRooms, room_name string, m_set, s_set []string) {
    //  WARN(Santiago): by now we will pass the sound template as empty.
    cherry_rooms.AddSound(room_name, m_set[0], m_set[1][1:len(m_set[1])-1], "", s_set[1][1:len(s_set[1])-1])
}


