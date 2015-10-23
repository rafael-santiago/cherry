package parser

import (
    "../../config"
    "fmt"
    "strings"
)

type CherryFileError struct {
    src string
    line int
    msg string
}

func (c *CherryFileError) Error() string {
    return fmt.Sprintf("ERROR: %s: at line %d: %s.\n", c.src, c.line, c.msg);
}

func NewCherryFileError(src string, line int, msg string) *CherryFileError {
    return &CherryFileError{src, line, msg}
}

func GetDataFromSection(section, config_data string, curr_line int) (string, int, int) {
    var s int
    var temp string
    for s = 0; s < len(config_data); s++ {
        switch config_data[s] {
            case  ' ', '\t':
                continue
            case '#':
                for config_data[s] != '\n' && s < len(config_data) {
                    s++
                }
                if s < len(config_data) {
                    curr_line++;
                }
                continue
            case '(', '\n':
                if config_data[s] == '\n' {
                    curr_line++;
                }
                if temp == section {
                    if config_data[s] == '\n' {
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
                    return data, s, curr_line
                }
                temp = ""
                break
            default:
                temp += string(config_data[s])
                break
        }
    }
    return "", s, curr_line
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
    var next_data string
    if s < len(data) {
        next_data = data[s:]
    }
    return set, curr_line, next_data
}

func ParseCherryFile(filepath string) (*config.CherryRooms, error) {
    var cherry_rooms *config.CherryRooms = nil
    return cherry_rooms, nil
}
