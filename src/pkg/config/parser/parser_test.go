/*
 *                          Copyright (C) 2015, 2016 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package parser

import (
	"fmt"
	"os"
	"pkg/config"
	"testing"
)

func TestGetDataFromSection(t *testing.T) {
	var configData = "# The commentary\n\n\n # More commentaries\ncherry.rooms ( aliens-on-earth:1024\nbackyard-science:911\n )\n"
	secData, _, _, err := GetDataFromSection("cherry.room", configData, 1, "foobar.cherry")
	if err == nil || len(secData) > 0 {
		t.Fail()
	}
	secData, _, _, err = GetDataFromSection("cherry.rooms", configData, 1, "foobar.cherry")
	if err != nil || secData != " aliens-on-earth:1024\nbackyard-science:911\n " {
		t.Fail()
	}
}

func TestStripBlanks(t *testing.T) {
	var data = " \t\t\t   \t\t\t   o u t e r  s p a c e   \t\t\t \t \t  \t"
	fmt.Println(StripBlanks(data))
	if StripBlanks(data) != "o u t e r  s p a c e" {
		t.Fail()
	}
}

func TestGetNextSetFromData(t *testing.T) {
	var data = "aliens-on-earth:1024\nbackyard-science : 911\n"
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
	var cherryRooms *config.CherryRooms
	cwd, _ := os.Getwd()
	os.Chdir("../../../../sample")
	var error *CherryFileError
	cherryRooms, error = ParseCherryFile("conf/sample.cherry")
	os.Chdir(cwd)
	if error != nil {
		fmt.Println(error)
		t.Fail()
	}
	if cherryRooms == nil {
		t.Fail()
	}
	var rooms []string
	rooms = cherryRooms.GetRooms()
	if len(rooms) != 1 {
		t.Fail()
	}
	if rooms[0] != "aliens-on-earth" {
		t.Fail()
	}
	if cherryRooms.GetListenPort(rooms[0]) != "1024" {
		t.Fail()
	}
	if cherryRooms.GetUsersTotal(rooms[0]) != "0" {
		t.Fail()
	}
	if cherryRooms.GetGreetingMessage(rooms[0]) != "Take meeeeee to your leader!!!" {
		t.Fail()
	}
	if cherryRooms.GetJoinMessage(rooms[0]) != "joined...<script>scrollIt();</script>" {
		t.Fail()
	}
	if cherryRooms.GetExitMessage(rooms[0]) != "has left...<script>scrollIt();</script>" {
		t.Fail()
	}
	if cherryRooms.GetOnIgnoreMessage(rooms[0]) != "(only you can see it) IGNORING " {
		t.Fail()
	}
	if cherryRooms.GetOnDeIgnoreMessage(rooms[0]) != "(only you can see it) is NOT IGNORING " {
		t.Fail()
	}
	if cherryRooms.GetPrivateMessageMarker(rooms[0]) != "(private)" {
		t.Fail()
	}
	if cherryRooms.GetAllUsersAlias(rooms[0]) != "EVERYBODY" {
		t.Fail()
	}
	if cherryRooms.GetMaxUsers(rooms[0]) != "10" {
		t.Fail()
	}
	if !cherryRooms.IsAllowingBriefs(rooms[0]) {
		t.Fail()
	}
	var expActionLabels map[string]string
	expActionLabels = make(map[string]string)
	expActionLabels["a01"] = "talks to"
	expActionLabels["a02"] = "screams with"
	expActionLabels["a03"] = "IGNORE"
	expActionLabels["a04"] = "NOT IGNORE"
	for a, l := range expActionLabels {
		if cherryRooms.GetRoomActionLabel(rooms[0], a) != l {
			t.Fail()
		}
		if len(cherryRooms.GetRoomActionTemplate(rooms[0], a)) == 0 {
			t.Fail()
		}
	}
	if cherryRooms.GetUsersTotal(rooms[0]) != "0" {
		t.Fail()
	}
	cherryRooms.AddUser(rooms[0], "dunha", "0", false)
	if cherryRooms.GetUsersTotal(rooms[0]) != "1" {
		t.Fail()
	}
	if len(cherryRooms.GetSessionID("dunha", rooms[0])) == 0 {
		t.Fail()
	}
	if cherryRooms.GetColor("dunha", rooms[0]) != "0" {
		t.Fail()
	}
	if len(cherryRooms.GetIgnoreList("dunha", rooms[0])) != 0 {
		t.Fail()
	}
	cherryRooms.AddToIgnoreList("dunha", "quiet", rooms[0])
	if cherryRooms.GetIgnoreList("dunha", rooms[0]) == "\"quiet\"" {
		t.Fail()
	}
	cherryRooms.DelFromIgnoreList("dunha", "quiet", rooms[0])
	if len(cherryRooms.GetIgnoreList("dunha", rooms[0])) != 0 {
		t.Fail()
	}
	cherryRooms.RemoveUser(rooms[0], "donha")
	if cherryRooms.GetUsersTotal(rooms[0]) != "1" {
		t.Fail()
	}
	cherryRooms.RemoveUser(rooms[0], "dunha")
	if cherryRooms.GetUsersTotal(rooms[0]) != "0" {
		t.Fail()
	}
	if len(cherryRooms.GetSessionID(rooms[0], "dunha")) != 0 {
		t.Fail()
	}
	message := cherryRooms.GetNextMessage(rooms[0])
	if len(message.From) != 0 ||
		len(message.To) != 0 ||
		len(message.Action) != 0 ||
		len(message.Image) != 0 ||
		len(message.Say) != 0 ||
		len(message.Priv) != 0 {
		t.Fail()
	}
	cherryRooms.EnqueueMessage(rooms[0], "(null)", "(anyone)", "a01", "i01", "boo!", "1")
	message = cherryRooms.GetNextMessage(rooms[0])
	if message.From != "(null)" ||
		message.To != "(anyone)" ||
		message.Action != "a01" ||
		message.Image != "i01" ||
		message.Say != "boo!" ||
		message.Priv != "1" {
		t.Fail()
	}
	for i := 0; i < 2; i++ {
		cherryRooms.DequeueMessage(rooms[0])
		message = cherryRooms.GetNextMessage(rooms[0])
		if len(message.From) != 0 ||
			len(message.To) != 0 ||
			len(message.Action) != 0 ||
			len(message.Image) != 0 ||
			len(message.Say) != 0 ||
			len(message.Priv) != 0 {
			t.Fail()
		}
	}
}
