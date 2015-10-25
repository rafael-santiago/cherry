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
        return fmt.Sprintf("ERROR: %s: at line %d: %s.\n", c.src, c.line, c.msg);
    }
    return fmt.Sprintf("ERROR: %s: %s.\n", c.src, c.msg);
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
    return "", s, curr_line, NewCherryFileError(curr_file, -1, "section not found")
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

        //  TODO(Santiago): Load the images & sounds & misc configurations.

        //  INFO(Santiago): until now these following section are non-mandatory

        //_ = GetRoomImages(set[0], cherry_rooms, string(cherry_file_data), filepath)

        //_ = GetRoomSounds(set[0], cherry_rooms, string(cherry_file_data), filepath)

        //err_room_config = GetRoomMisc(set[0], cherry_rooms, string(cherry_file_data), filepath)

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
        cherry_rooms.AddTemplate(room_name, set[0], set[1][1:len(set[1])-1]);
    }
    return nil
}

func GetRoomActions(room_name string, cherry_rooms *config.CherryRooms, config_data, filepath string) *CherryFileError {
    var labels string
    var labels_line int
    var labels_err *CherryFileError
    labels, _, labels_line, labels_err = GetDataFromSection("cherry." + room_name + ".actions",
                                                            config_data, 1, filepath)
    if labels_err != nil {
        return labels_err
    }

    var templates string
    var templates_line int
    var templates_err *CherryFileError
    templates, _, templates_line, templates_err = GetDataFromSection("cherry." + room_name + ".actions.templates",
                                                                     config_data, 1, filepath)

    if templates_err != nil {
        return templates_err
    }

    var action_labels []string
    action_labels, labels_line, labels = GetNextSetFromData(labels, labels_line, "=")
    for len(action_labels) == 2 {
        if cherry_rooms.HasAction(room_name, action_labels[0]) {
            return NewCherryFileError(filepath, labels_line, "room action \"" + action_labels[0] + "\" redeclared.")
        }
        if len(action_labels[1]) == 0 {
            return NewCherryFileError(filepath, labels_line, "room action with no value.")
        }
        if action_labels[1][0] != '"' || action_labels[1][len(action_labels[1])-1] != '"' {
            return NewCherryFileError(filepath, labels_line, "room action must be set with a valid string.")
        }

        //  INFO(Santiago): Getting the template for the current action label from a section to another.
        var action_templates []string
        var temp string = templates
        action_templates, templates_line, temp = GetNextSetFromData(temp, templates_line, "=")
        for len(action_templates) == 2 && action_templates[0] != action_labels[0] {
            action_templates, templates_line, temp = GetNextSetFromData(temp, templates_line, "=")
        }

        if action_templates[0] != action_labels[0] {
            return NewCherryFileError(filepath, templates_line, "there is no template for action \"" + action_labels[0] + "\".")
        }
        if len(action_templates[1]) == 0 {
            return NewCherryFileError(filepath, templates_line, "room action template with no value.")
        }
        if action_templates[1][0] != '"' || action_templates[1][len(action_templates[1])-1] != '"' {
            return NewCherryFileError(filepath, templates_line, "room action template must be set with a valid string.")
        }

        cherry_rooms.AddAction(room_name, action_labels[0], action_labels[1][1:len(action_labels[1])-1], action_templates[1][1:len(action_templates[1])-1]);
    }
    return nil
}
