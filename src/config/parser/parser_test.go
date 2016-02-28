/*
 *                          Copyright (C) 2015, 2016 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package parser

import (
    "os"
    ".."
    "testing"
    "fmt"
)

func TestGetDataFromSection(t *testing.T) {
    var config_data string = "# The commentary\n\n\n # More commentaries\ncherry.rooms ( aliens-on-earth:1024\nbackyard-science:911\n )\n"
    sec_data, _, _, err := GetDataFromSection("cherry.room", config_data, 1, "foobar.cherry")
    if err == nil || len(sec_data) > 0 {
        t.Fail();
    }
    sec_data, _, _, err = GetDataFromSection("cherry.rooms", config_data, 1, "foobar.cherry")
    if err != nil || sec_data != " aliens-on-earth:1024\nbackyard-science:911\n " {
        t.Fail();
    }
}

func TestStripBlanks(t *testing.T) {
    var data string = " \t\t\t   \t\t\t   o u t e r  s p a c e   \t\t\t \t \t  \t"
    fmt.Println(StripBlanks(data))
    if StripBlanks(data) != "o u t e r  s p a c e" {
        t.Fail()
    }
}

func TestGetNextSetFromData(t *testing.T) {
    var data string = "aliens-on-earth:1024\nbackyard-science : 911\n"
    set, _, data := GetNextSetFromData(data, 1, ":")
    if len(set) != 2 || len(data) == 0 || set[0] != "aliens-on-earth" || set[1] != "1024" {
        t.Fail()
    }
    set, _, data = GetNextSetFromData(data, 1, ":")
    if len(set) != 2 || len(data) == 0 || set[0] != "backyard-science" || set[1] != "911" {
        t.Fail()
    }
    set, _, data = GetNextSetFromData(data, 1, ":")
    if len(set) != 1 || len(data) != 0 {
        t.Fail()
    }
    data = "i01 = \"http://www.nowhere.com/images/i01.gif\"\ni02 = \"http://www.nowhere.com/images/i02.gif\"\ni03 = \"http://www.nowhere.com/images/i03.gif\"\n"
    set, _, data = GetNextSetFromData(data, 1, "=")
    if len(set) != 2 || len(data) == 0 || set[0] != "i01" || set[1] != "\"http://www.nowhere.com/images/i01.gif\"" {
        t.Fail()
    }
    set, _, data = GetNextSetFromData(data, 1, "=")
    if len(set) != 2 || len(data) == 0 || set[0] != "i02" || set[1] != "\"http://www.nowhere.com/images/i02.gif\"" {
        t.Fail()
    }
    set, _, data = GetNextSetFromData(data, 1, "=")
    if len(set) != 2 || len(data) == 0 || set[0] != "i03" || set[1] != "\"http://www.nowhere.com/images/i03.gif\"" {
        t.Fail()
    }
}

func TestRealCherryFileParsing(t *testing.T) {
    //  INFO(Santiago): This Homeric test should be splitted in the future.
    var cherry_rooms *config.CherryRooms
    cwd, _ := os.Getwd()
    os.Chdir("../../../sample")
    var error *CherryFileError
    cherry_rooms, error = ParseCherryFile("conf/sample.cherry")
    os.Chdir(cwd)
    if error != nil {
        fmt.Println(error)
       t.Fail()
    }
    if cherry_rooms == nil {
        t.Fail()
    }
    var rooms []string
    rooms = cherry_rooms.GetRooms()
    if len(rooms) != 1 {
        t.Fail()
    }
    if rooms[0] != "aliens-on-earth" {
        t.Fail()
    }
    if cherry_rooms.GetListenPort(rooms[0]) != "1024" {
        t.Fail()
    }
    if cherry_rooms.GetUsersTotal(rooms[0]) != "0" {
        t.Fail()
    }
    if cherry_rooms.GetGreetingMessage(rooms[0]) != "Take meeeeee to your leader!!!" {
        t.Fail()
    }
    if cherry_rooms.GetJoinMessage(rooms[0]) != "joined...<script>scrollIt();</script>" {
        t.Fail()
    }
    if cherry_rooms.GetExitMessage(rooms[0]) != "has left...<script>scrollIt();</script>" {
        t.Fail()
    }
    if cherry_rooms.GetOnIgnoreMessage(rooms[0]) != "(only you can see it) IGNORING " {
        t.Fail()
    }
    if cherry_rooms.GetOnDeIgnoreMessage(rooms[0]) != "(only you can see it) is NOT IGNORING " {
        t.Fail()
    }
    if cherry_rooms.GetPrivateMessageMarker(rooms[0]) != "(private)" {
        t.Fail()
    }
    if cherry_rooms.GetAllUsersAlias(rooms[0]) != "EVERYBODY" {
        t.Fail()
    }
    if cherry_rooms.GetMaxUsers(rooms[0]) != "10" {
        t.Fail()
    }
    if ! cherry_rooms.IsAllowingBriefs(rooms[0]) {
        t.Fail()
    }
    var exp_action_labels map[string]string
    exp_action_labels = make(map[string]string)
    exp_action_labels["a01"] = "talks to"
    exp_action_labels["a02"] = "screams with"
    exp_action_labels["a03"] = "IGNORE"
    exp_action_labels["a04"] = "NOT IGNORE"
    for a, l := range exp_action_labels {
        if cherry_rooms.GetRoomActionLabel(rooms[0], a) != l {
            t.Fail()
        }
        if len(cherry_rooms.GetRoomActionTemplate(rooms[0], a)) == 0 {
            t.Fail()
        }
    }
    if cherry_rooms.GetUsersTotal(rooms[0]) != "0" {
        t.Fail()
    }
    cherry_rooms.AddUser(rooms[0], "dunha", "0", false)
    if cherry_rooms.GetUsersTotal(rooms[0]) != "1" {
        t.Fail()
    }
    if len(cherry_rooms.GetSessionId("dunha", rooms[0])) == 0 {
        t.Fail()
    }
    if cherry_rooms.GetColor("dunha", rooms[0]) != "0" {
        t.Fail()
    }
    if len(cherry_rooms.GetIgnoreList("dunha", rooms[0])) != 0 {
        t.Fail()
    }
    cherry_rooms.AddToIgnoreList("dunha", "quiet", rooms[0])
    if cherry_rooms.GetIgnoreList("dunha", rooms[0]) == "\"quiet\"" {
        t.Fail()
    }
    cherry_rooms.DelFromIgnoreList("dunha", "quiet", rooms[0])
    if len(cherry_rooms.GetIgnoreList("dunha", rooms[0])) != 0 {
        t.Fail()
    }
    cherry_rooms.RemoveUser(rooms[0], "donha")
    if cherry_rooms.GetUsersTotal(rooms[0]) != "1" {
        t.Fail()
    }
    cherry_rooms.RemoveUser(rooms[0], "dunha")
    if cherry_rooms.GetUsersTotal(rooms[0]) != "0" {
        t.Fail()
    }
    if len(cherry_rooms.GetSessionId(rooms[0], "dunha")) != 0 {
        t.Fail()
    }
    message := cherry_rooms.GetNextMessage(rooms[0])
    if len(message.From)   != 0 ||
       len(message.To)     != 0 ||
       len(message.Action) != 0 ||
       len(message.Image)  != 0 ||
       len(message.Say)    != 0 ||
       len(message.Priv)   != 0 {
        t.Fail()
    }
    cherry_rooms.EnqueueMessage(rooms[0], "(null)", "(anyone)", "a01", "i01", "boo!", "1")
    message = cherry_rooms.GetNextMessage(rooms[0])
    if message.From   != "(null)"   ||
       message.To     != "(anyone)" ||
       message.Action != "a01"      ||
       message.Image  != "i01"      ||
       message.Say    != "boo!"     ||
       message.Priv   != "1" {
        t.Fail()
    }
    for i := 0; i < 2; i++ {
        cherry_rooms.DequeueMessage(rooms[0])
        message = cherry_rooms.GetNextMessage(rooms[0])
        if len(message.From)   != 0 ||
           len(message.To)     != 0 ||
           len(message.Action) != 0 ||
           len(message.Image)  != 0 ||
           len(message.Say)    != 0 ||
           len(message.Priv)   != 0 {
                t.Fail()
        }
    }
}
