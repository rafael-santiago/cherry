package parser

import (
    "../../config"
    "fmt"
)

type CherryFileError struct {
    src string
    line int
    msg string
}

func (c *CherryFileError) Error() string {
    return fmt.Sprintf("ERROR: %s: LINE %d: %s.\n", c.src, c.line, c.msg);
}

func NewCherryFileError(src string, line int, msg string) *CherryFileError {
    return &CherryFileError{src, line, msg}
}

func ParseCherryFile(filepath string) (*config.CherryRooms, error) {
    var cherry_rooms *config.CherryRooms = nil
    return cherry_rooms, nil
}
