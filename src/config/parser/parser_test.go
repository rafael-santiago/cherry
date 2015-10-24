package parser

import (
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
