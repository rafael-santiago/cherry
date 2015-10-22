package main

import (
   "./config"
   //"./config/parser"
)

func main() {
    cherry_rooms := config.NewCherryRooms()
    cherry_rooms.AddRoom("aliens-on-earth", 8811)
    //e, rooms := parser.ParseCherryFile("config.cherry")
}
