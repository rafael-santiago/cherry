/*
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package main

import (
//   "./config"
   "./config/parser"
    "fmt"
)

func main() {
    //var cherry_rooms *config.CherryRooms
    var err *parser.CherryFileError
    _, err = parser.ParseCherryFile("config.cherry")
    if err != nil {
        fmt.Println(err.Error())
    } else {
        fmt.Println("*** Configuration loaded!")
    }
}
