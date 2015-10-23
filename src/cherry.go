package main

import (
//   "./config"
   "./config/parser"
    "fmt"
)

func main() {
    var s string = "# This is a cherry config script.\ncherry.root (\n# System configuration as whole)\ncherry.rooms (\nfoobar:8810\naliens-on-earth:8811\nbackyard-science:8812\n)"
    var data string
    var line int
    var set []string
    //cherry_rooms := config.NewCherryRooms()
    //cherry_rooms.AddRoom("aliens-on-earth", 8811)
    data, _, line = parser.GetDataFromSection("cherry.rooms", s, 1)
    set, line, data = parser.GetNextSetFromData(data, line, ":")
    for len(set) == 2 {
        fmt.Println(set)
        set, line, data = parser.GetNextSetFromData(data, line, ":")
    }
}
