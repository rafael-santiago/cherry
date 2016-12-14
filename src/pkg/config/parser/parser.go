/*
Package parser parse and loads a cherry file to the memory.
--
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
*/
package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"pkg/config"
	"strconv"
	"strings"
)

// CherryFileError is returned by any parse function implemented in "parser.go".
type CherryFileError struct {
	src  string
	line int
	msg  string
}

// Error spits the error string.
func (c *CherryFileError) Error() string {
	if c.line > -1 {
		return fmt.Sprintf("ERROR: %s: at line %d: %s\n", c.src, c.line, c.msg)
	}
	return fmt.Sprintf("ERROR: %s: %s\n", c.src, c.msg)
}

// NewCherryFileError creates a *CherryFileError.
func NewCherryFileError(src string, line int, msg string) *CherryFileError {
	return &CherryFileError{src, line, msg}
}

// GetDataFromSection returns the raw section data from a cherry file section data.
func GetDataFromSection(section, configData string, currLine int, currFile string) (string, int, int, *CherryFileError) {
	var s int
	var temp string
	for s = 0; s < len(configData); s++ {
		switch configData[s] {
		case '#':
			for configData[s] != '\n' && configData[s] != '\r' && s < len(configData) {
				s++
			}
			if s < len(configData) {
				currLine++
			}
			continue
		case '(', '\n', ' ', '\t', '\r':
			if configData[s] == '\n' {
				currLine++
			}
			if temp == section {
				if configData[s] == '\n' || configData[s] == '\r' || configData[s] == ' ' || configData[s] == '\t' {
					for s < len(configData) && configData[s] != '(' {
						s++
						if s < len(configData) && configData[s] == '\n' {
							currLine++
						}
					}
				}
				if s < len(configData) && configData[s] == '(' {
					s++
				}
				var data string
				for s < len(configData) {
					if configData[s] == '"' {
						data += string(configData[s])
						s++
						for s < len(configData) && configData[s] != '"' {
							if configData[s] != '\\' {
								data += string(configData[s])
							} else {
								data += string(configData[s+1])
								s++
							}
							s++
						}
						if s < len(configData) {
							data += string(configData[s])
						}
					} else if configData[s] != ')' {
						data += string(configData[s])
					} else {
						break
					}
					s++
				}
				return data, s, currLine, nil
			} else if temp == "cherry.branch" {
				for s < len(configData) && (configData[s] == ' ' || configData[s] == '\t') {
					s++
				}
				if s < len(configData) {
					var branchFilepath string
					for s < len(configData) && configData[s] != '\n' && configData[s] != '\r' {
						branchFilepath += string(configData[s])
						s++
					}
					if s < len(configData) {
						currLine++
					}
					branchBuffer, err := ioutil.ReadFile(branchFilepath)
					if err != nil {
						fmt.Println(fmt.Sprintf("WARNING: %s: at line %d: %s. Be tidy... removing or commenting this dry branch from your cherry.", currFile, currLine-1, err.Error()))
						//return "", s, currLine, NewCherryFileError(currFile,
						//                                            currLine - 1,
						//                                            "unable to read cherry.branch from \"" + branchFilepath + "\" [ details: " + err.Error() + " ]")
					} else {
						branchData, branchOffset, branchLine, _ := GetDataFromSection(section, string(branchBuffer), 1, branchFilepath)
						if len(branchData) > 0 {
							return branchData, branchOffset, branchLine, nil
						}
					}
				}
			}
			temp = ""
			break
		default:
			temp += string(configData[s])
			break
		}
	}
	return "", s, currLine, NewCherryFileError(currFile, -1, "section \""+section+"\" not found.")
}

// GetNextSetFromData returns the next "field = value".
func GetNextSetFromData(data string, currLine int, tok string) ([]string, int, string) {
	if len(data) == 0 {
		return make([]string, 0), currLine, ""
	}
	var s int
	for s = 0; s < len(data) && (data[s] == ' ' || data[s] == '\t' || data[s] == '\n' || data[s] == '\r'); s++ {
		if data[s] == '\n' {
			currLine++
		}
	}
	var line string
	for s < len(data) && data[s] != '\n' && data[s] != '\r' {
		if data[s] == '"' {
			line += string(data[s])
			s++
			for s < len(data) && data[s] != '"' {
				if data[s] == '\\' {
					s++
				}
				line += string(data[s])
				if data[s] == '\n' {
					currLine++
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
				currLine++
			}
		} else {
			line += string(data[s])
		}
		s++
	}
	if len(line) == 0 {
		if len(data) == 0 {
			return make([]string, 0), currLine, ""
		}
	}
	set := strings.Split(line, tok)
	if len(set) == 2 {
		set[0] = StripBlanks(set[0])
		set[1] = StripBlanks(set[1])
	}
	var nextData string
	if s < len(data) {
		nextData = data[s:]
	}
	return set, currLine, nextData
}

// StripBlanks strips all blanks.
func StripBlanks(data string) string {
	var retval string
	var dStart int
	for dStart < len(data) && (data[dStart] == ' ' || data[dStart] == '\t') {
		dStart++
	}
	var dEnd = len(data) - 1
	for dEnd > 0 && (data[dEnd] == ' ' || data[dEnd] == '\t') {
		dEnd--
	}
	retval = data[dStart : dEnd+1]
	return retval
}

// ParseCherryFile parses a file at @filepath and returns a *config.CherryRooms or a *CherryFileError.
func ParseCherryFile(filepath string) (*config.CherryRooms, *CherryFileError) {
	var cherryRooms *config.CherryRooms
	var cherryFileData []byte
	var data string
	var err *CherryFileError
	var line int
	cherryFileData, ioErr := ioutil.ReadFile(filepath)
	if ioErr != nil {
		return nil, NewCherryFileError("(no file)", -1, fmt.Sprintf("unable to read from \"%s\" [more details: %s].", filepath, ioErr.Error()))
	}
	data, _, line, err = GetDataFromSection("cherry.root", string(cherryFileData), 1, filepath)
	if err != nil {
		return nil, err
	}
	var set []string
	cherryRooms = config.NewCherryRooms()
	set, line, data = GetNextSetFromData(data, line, "=")
	for len(set) == 2 {
		switch set[0] {
		case "servername", "certificate", "private-key":
			if set[1][0] != '"' || set[1][len(set[1])-1] != '"' {
				return nil, NewCherryFileError(filepath, line, fmt.Sprintf("invalid string."))
			}
			data := set[1][1 : len(set[1])-1]

			if set[0] == "certificate" || set[0] == "private-key" {
				if _, err := os.Stat(data); os.IsNotExist(err) {
					return nil, NewCherryFileError(filepath, line, fmt.Sprintf("\"%s\" must receive an accessible file path.", set[0]))
				}
			}

			if set[0] == "servername" {
				cherryRooms.SetServername(data)
			} else if set[0] == "certificate" {
				cherryRooms.SetCertificatePath(data)
			} else if set[0] == "private-key" {
				cherryRooms.SetPrivateKeyPath(data)
			}
			break

		default:
			return nil, NewCherryFileError(filepath, line, fmt.Sprintf("unknown config set \"%s\".", set[0]))
		}
		set, line, data = GetNextSetFromData(data, line, "=")
	}
	if cherryRooms.GetServername() == "localhost" {
		fmt.Println("WARN: cherry.root.servername is equals to \"localhost\". Things will not work outside this node.")
	}
	data, _, line, err = GetDataFromSection("cherry.rooms", string(cherryFileData), 1, filepath)
	if err != nil {
		return nil, err
	}
	//  INFO(Santiago): Adding all scanned rooms from the first cherry.rooms section found
	//                  [cherry branches were scanned too at this point].
	set, line, data = GetNextSetFromData(data, line, ":")
	for len(set) == 2 {
		if cherryRooms.HasRoom(set[0]) {
			return nil, NewCherryFileError(filepath, line, fmt.Sprintf("room \"%s\" redeclared.", set[0]))
		}
		var value int64
		var convErr error
		value, convErr = strconv.ParseInt(set[1], 10, 16)
		if convErr != nil {
			return nil, NewCherryFileError(filepath, line, fmt.Sprintf("invalid port value \"%s\" [more details: %s].", set[1], convErr))
		}
		var port int16
		port = int16(value)
		if cherryRooms.PortBusyByAnotherRoom(port) {
			return nil, NewCherryFileError(filepath, line, fmt.Sprintf("the port \"%s\" is already busy by another room.", set[1]))
		}

		cherryRooms.AddRoom(set[0], port)

		errRoomConfig := GetRoomTemplates(set[0], cherryRooms, string(cherryFileData), filepath)
		if errRoomConfig != nil {
			return nil, errRoomConfig
		}

		errRoomConfig = GetRoomActions(set[0], cherryRooms, string(cherryFileData), filepath)
		if errRoomConfig != nil {
			return nil, errRoomConfig
		}

		//  INFO(Santiago): until now these two following sections are non-mandatory.

		_ = GetRoomImages(set[0], cherryRooms, string(cherryFileData), filepath)

		//_ = GetRoomSounds(set[0], cherryRooms, string(cherryFileData), filepath)

		errRoomConfig = GetRoomMisc(set[0], cherryRooms, string(cherryFileData), filepath)
		if errRoomConfig != nil {
			return nil, errRoomConfig
		}

		//  INFO(Santiago): Let's transfer the next room from file to the memory.
		set, line, data = GetNextSetFromData(data, line, ":")
	}
	return cherryRooms, nil
}

// GetRoomTemplates parses "cherry.[roomName].templates" section.
func GetRoomTemplates(roomName string, cherryRooms *config.CherryRooms, configData, filepath string) *CherryFileError {
	var data string
	var line int
	var err *CherryFileError
	data, _, line, err = GetDataFromSection("cherry."+roomName+".templates",
		configData, 1, filepath)
	if err != nil {
		return err
	}
	var set []string
	set, line, data = GetNextSetFromData(data, line, "=")
	for len(set) == 2 {
		if cherryRooms.HasTemplate(roomName, set[0]) {
			return NewCherryFileError(filepath, line, "room template \""+set[0]+"\" redeclared.")
		}
		if len(set[1]) == 0 {
			return NewCherryFileError(filepath, line, "room template with no value.")
		}
		if set[1][0] != '"' || set[1][len(set[1])-1] != '"' {
			return NewCherryFileError(filepath, line, "room template must be set with a valid string.")
		}
		var templateData []byte
		var templateDataErr error
		templateData, templateDataErr = ioutil.ReadFile(set[1][1 : len(set[1])-1])
		if templateDataErr != nil {
			return NewCherryFileError(filepath, line, "unable to access room template file [more details: "+templateDataErr.Error()+"].")
		}
		cherryRooms.AddTemplate(roomName, set[0], string(templateData))
		set, line, data = GetNextSetFromData(data, line, "=")
	}
	return nil
}

// GetRoomActions parses "cherry.[roomName].actions" section.
func GetRoomActions(roomName string, cherryRooms *config.CherryRooms, configData, filepath string) *CherryFileError {
	return getIndirectConfig("cherry."+roomName+".actions",
		"cherry."+roomName+".actions.templates",
		roomActionMainVerifier, roomActionSubVerifier, roomActionSetter,
		roomName, cherryRooms, configData, filepath)
}

// GetRoomImages parses "cherry.[roomName].images" and "cherry.[roomName].images.url".
func GetRoomImages(roomName string, cherryRooms *config.CherryRooms, configData, filepath string) *CherryFileError {
	return getIndirectConfig("cherry."+roomName+".images",
		"cherry."+roomName+".images.url",
		roomImageMainVerifier, roomImageSubVerifier, roomImageSetter,
		roomName, cherryRooms, configData, filepath)
}

//func GetRoomSounds(roomName string, cherryRooms *config.CherryRooms, configData, filepath string) *CherryFileError {
//    return getIndirectConfig("cherry." + roomName + ".sounds",
//                               "cherry." + roomName + ".sounds.url",
//                               room_sound_main_verifier, room_sound_sub_verifier, room_sound_setter,
//                               roomName, cherryRooms, configData, filepath)
//}

// GetRoomMisc parses "cherry.[roomName].misc" section.
func GetRoomMisc(roomName string, cherryRooms *config.CherryRooms, configData, filepath string) *CherryFileError {
	var mData string
	var mLine int
	var mErr *CherryFileError
	mData, _, mLine, mErr = GetDataFromSection("cherry."+roomName+".misc", configData, 1, filepath)
	if mErr != nil {
		return mErr
	}

	var verifier map[string]func(string) bool
	verifier = make(map[string]func(string) bool)
	verifier["join-message"] = verifyString
	verifier["exit-message"] = verifyString
	verifier["on-ignore-message"] = verifyString
	verifier["on-deignore-message"] = verifyString
	verifier["greeting-message"] = verifyString
	verifier["private-message-marker"] = verifyString
	verifier["max-users"] = verifyNumber
	verifier["allow-brief"] = verifyBool
	//verifier["flooding-police"]               = verifyBool
	//verifier["max-flood-allowed-before-kick"] = verifyNumber
	verifier["all-users-alias"] = verifyString
	verifier["ignore-action"] = verifyString
	verifier["deignore-action"] = verifyString
	verifier["public-directory"] = verifyString

	var setter map[string]func(*config.CherryRooms, string, string)
	setter = make(map[string]func(*config.CherryRooms, string, string))
	setter["join-message"] = setJoinMessage
	setter["exit-message"] = setExitMessage
	setter["on-ignore-message"] = setOnIgnoreMessage
	setter["on-deignore-message"] = setOnDeIgnoreMessage
	setter["greeting-message"] = setGreetingMessage
	setter["private-message-marker"] = setPrivateMessageMarker
	setter["max-users"] = setMaxUsers
	setter["allow-brief"] = setAllowBrief
	//setter["flooding-police"]               = set_flooding_police
	//setter["max-flood-allowed-before-kick"] = set_max_flood_allowed_before_kick
	setter["all-users-alias"] = setAllUsersAlias
	setter["ignore-action"] = setIgnoreAction
	setter["deignore-action"] = setDeIgnoreAction
	setter["public-directory"] = setPublicDirectory

	var alreadySet map[string]bool
	alreadySet = make(map[string]bool)
	alreadySet["join-message"] = false
	alreadySet["exit-message"] = false
	alreadySet["on-ignore-message"] = false
	alreadySet["on-deignore-message"] = false
	alreadySet["greeting-message"] = false
	alreadySet["private-message-marker"] = false
	alreadySet["max-users"] = false
	//alreadySet["flooding-police"]               = false
	//alreadySet["max-flood-allowed-before-kick"] = false
	alreadySet["all-users-alias"] = false
	alreadySet["ignore-action"] = false
	alreadySet["deignore-action"] = false
	alreadySet["public-directory"] = false

	var mSet []string
	mSet, mLine, mData = GetNextSetFromData(mData, mLine, "=")
	for len(mSet) == 2 {
		_, exists := verifier[mSet[0]]
		if !exists {
			return NewCherryFileError(filepath, mLine, "misc configuration named as \""+mSet[0]+"\" is unrecognized.")
		}
		if alreadySet[mSet[0]] {
			return NewCherryFileError(filepath, mLine, "misc configuration \""+mSet[0]+"\" re-configured.")
		}
		if !verifier[mSet[0]](mSet[1]) {
			return NewCherryFileError(filepath, mLine, "misc configuration \""+mSet[0]+"\" has invalid value : "+mSet[1])
		}
		setter[mSet[0]](cherryRooms, roomName, mSet[1])
		alreadySet[mSet[0]] = true
		mSet, mLine, mData = GetNextSetFromData(mData, mLine, "=")
	}

	return nil
}

func setIgnoreAction(cherryRooms *config.CherryRooms, roomName, action string) {
	cherryRooms.SetIgnoreAction(roomName, action[1:len(action)-1])
}

func setDeIgnoreAction(cherryRooms *config.CherryRooms, roomName, action string) {
	cherryRooms.SetDeIgnoreAction(roomName, action[1:len(action)-1])
}

func setJoinMessage(cherryRooms *config.CherryRooms, roomName, message string) {
	cherryRooms.SetJoinMessage(roomName, message[1:len(message)-1])
}

func setExitMessage(cherryRooms *config.CherryRooms, roomName, message string) {
	cherryRooms.SetExitMessage(roomName, message[1:len(message)-1])
}

func setOnIgnoreMessage(cherryRooms *config.CherryRooms, roomName, message string) {
	cherryRooms.SetOnIgnoreMessage(roomName, message[1:len(message)-1])
}

func setOnDeIgnoreMessage(cherryRooms *config.CherryRooms, roomName, message string) {
	cherryRooms.SetOnDeIgnoreMessage(roomName, message[1:len(message)-1])
}

func setGreetingMessage(cherryRooms *config.CherryRooms, roomName, message string) {
	cherryRooms.SetGreetingMessage(roomName, message[1:len(message)-1])
}

func setPrivateMessageMarker(cherryRooms *config.CherryRooms, roomName, marker string) {
	cherryRooms.SetPrivateMessageMarker(roomName, marker[1:len(marker)-1])
}

func setMaxUsers(cherryRooms *config.CherryRooms, roomName, value string) {
	var intValue int64
	intValue, _ = strconv.ParseInt(value, 10, 64)
	cherryRooms.SetMaxUsers(roomName, int(intValue))
}

func setAllowBrief(cherryRooms *config.CherryRooms, roomName, value string) {
	var allow bool
	allow = (value == "yes" || value == "true")
	cherryRooms.SetAllowBrief(roomName, allow)
}

func setPublicDirectory(cherryRooms *config.CherryRooms, roomName, value string) {
	cherryRooms.SetPublicDirectory(roomName, value[1:len(value)-1])
}

//func set_flooding_police(cherryRooms *config.CherryRooms, roomName, value string) {
//    var impose bool
//    impose = (value == "yes")
//    cherryRooms.SetFloodingPolice(roomName, impose)
//}

func setAllUsersAlias(cherryRooms *config.CherryRooms, roomName, value string) {
	cherryRooms.SetAllUsersAlias(roomName, value[1:len(value)-1])
}

//func set_max_flood_allowed_before_kick(cherryRooms *config.CherryRooms, roomName, value string) {
//    var intValue int64
//    intValue, _ = strconv.ParseInt(value, 10, 64)
//    cherryRooms.SetMaxFloodAllowedBeforeKick(roomName, int(intValue))
//}

func verifyNumber(buffer string) bool {
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

func verifyString(buffer string) bool {
	if len(buffer) <= 1 {
		return false
	}
	return (buffer[0] == '"' && buffer[len(buffer)-1] == '"')
}

func verifyBool(buffer string) bool {
	if len(buffer) == 0 {
		return false
	}
	return (buffer == "yes" || buffer == "no" || buffer == "true" || buffer == "false")
}

//  WARN(Santiago): The following codes are a brain damage. I am sorry.

func getIndirectConfig(mainSection,
	subSection string,
	mainVerifier,
	subVerifier func([]string, []string, int, int, string, string, *config.CherryRooms) *CherryFileError,
	setter func(*config.CherryRooms, string, []string, []string),
	roomName string,
	cherryRooms *config.CherryRooms,
	configData,
	filepath string) *CherryFileError {
	var mData string
	var mLine int
	var mErr *CherryFileError
	mData, _, mLine, mErr = GetDataFromSection(mainSection, configData, 1, filepath)
	if mErr != nil {
		return mErr
	}

	var sData string
	var sLine int
	var sErr *CherryFileError
	sData, _, sLine, sErr = GetDataFromSection(subSection, configData, 1, filepath)

	if sErr != nil {
		return sErr
	}

	var mSet []string
	mSet, mLine, mData = GetNextSetFromData(mData, mLine, "=")
	for len(mSet) == 2 {
		var sSet []string
		mErr = mainVerifier(mSet, sSet, mLine, sLine, roomName, filepath, cherryRooms)
		if mErr != nil {
			return mErr
		}

		//  INFO(Santiago): Getting the template for the current action label from a section to another.
		var temp = sData

		var tempLine = sLine
		sSet, tempLine, temp = GetNextSetFromData(temp, tempLine, "=")
		for len(sSet) == 2 && sSet[0] != mSet[0] {
			sSet, tempLine, temp = GetNextSetFromData(temp, tempLine, "=")
		}

		sErr = subVerifier(mSet, sSet, mLine, sLine, roomName, filepath, cherryRooms)

		if sErr != nil {
			return sErr
		}

		setter(cherryRooms, roomName, mSet, sSet)

		mSet, mLine, mData = GetNextSetFromData(mData, mLine, "=")
	}
	return nil
}

func roomActionMainVerifier(mSet, sSet []string, mLine, sLine int, roomName, filepath string, cherryRooms *config.CherryRooms) *CherryFileError {
	if cherryRooms.HasAction(roomName, mSet[0]) {
		return NewCherryFileError(filepath, mLine, "room action \""+mSet[0]+"\" redeclared.")
	}
	if len(mSet[1]) == 0 {
		return NewCherryFileError(filepath, mLine, "unlabeled room action.")
	}
	if mSet[1][0] != '"' || mSet[1][len(mSet[1])-1] != '"' {
		return NewCherryFileError(filepath, mLine, "room action must be set with a valid string.")
	}
	return nil
}

func roomActionSubVerifier(mSet, sSet []string, mLine, sLine int, roomName, filepath string, cherryRooms *config.CherryRooms) *CherryFileError {
	if sSet[0] != mSet[0] {
		return NewCherryFileError(filepath, sLine, "there is no template for action \""+mSet[0]+"\".")
	}
	if len(sSet[1]) == 0 {
		return NewCherryFileError(filepath, sLine, "empty room action template.")
	}
	if sSet[1][0] != '"' || sSet[1][len(sSet[1])-1] != '"' {
		return NewCherryFileError(filepath, sLine, "room action template must be set with a valid string.")
	}
	var templatePath = sSet[1][1 : len(sSet[1])-1]
	_, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return NewCherryFileError(filepath, sLine, fmt.Sprintf("unable to access file, details: [ %s ]", err.Error()))
	}
	return nil
}

func roomActionSetter(cherryRooms *config.CherryRooms, roomName string, mSet, sSet []string) {
	data, _ := ioutil.ReadFile(sSet[1][1 : len(sSet[1])-1])
	cherryRooms.AddAction(roomName, mSet[0], mSet[1][1:len(mSet[1])-1], string(data))
}

func roomImageMainVerifier(mSet, sSet []string, mLine, sLine int, roomName, filepath string, cherryRooms *config.CherryRooms) *CherryFileError {
	if cherryRooms.HasImage(roomName, mSet[0]) {
		return NewCherryFileError(filepath, mLine, "room image \""+mSet[0]+"\" redeclared.")
	}
	if len(mSet[1]) == 0 {
		return NewCherryFileError(filepath, mLine, "unlabeled room image.")
	}
	if mSet[1][0] != '"' || mSet[1][len(mSet[1])-1] != '"' {
		return NewCherryFileError(filepath, mLine, "room image must be set with a valid string.")
	}
	return nil
}

func roomImageSubVerifier(mSet, sSet []string, mLine, sLine int, roomName, filepath string, cherryRooms *config.CherryRooms) *CherryFileError {
	if sSet[0] != mSet[0] {
		return NewCherryFileError(filepath, sLine, "there is no url for image \""+mSet[0]+"\".")
	}
	if len(sSet[1]) == 0 {
		return NewCherryFileError(filepath, sLine, "empty room image url.")
	}
	if sSet[1][0] != '"' || sSet[1][len(sSet[1])-1] != '"' {
		return NewCherryFileError(filepath, sLine, "room image url must be set with a valid string.")
	}
	return nil
}

func roomImageSetter(cherryRooms *config.CherryRooms, roomName string, mSet, sSet []string) {
	//  WARN(Santiago): by now we will pass the image template as empty.
	cherryRooms.AddImage(roomName, mSet[0], mSet[1][1:len(mSet[1])-1], "", sSet[1][1:len(sSet[1])-1])
}

//func room_sound_main_verifier(mSet, sSet []string, mLine, sLine int, roomName, filepath string, cherryRooms *config.CherryRooms) *CherryFileError {
//    if cherryRooms.HasImage(roomName, mSet[0]) {
//        return NewCherryFileError(filepath, mLine, "room sound \"" + mSet[0] + "\" redeclared.")
//    }
//    if len(mSet[1]) == 0 {
//        return NewCherryFileError(filepath, mLine, "unlabeled room sound.")
//    }
//    if mSet[1][0] != '"' || mSet[1][len(mSet[1])-1] != '"' {
//        return NewCherryFileError(filepath, mLine, "room sound must be set with a valid string.")
//    }
//    return nil
//}

//func room_sound_sub_verifier(mSet, sSet []string, mLine, sLine int, roomName, filepath string, cherryRooms *config.CherryRooms) *CherryFileError {
//    if sSet[0] != mSet[0] {
//        return NewCherryFileError(filepath, sLine, "there is no url for sound \"" + mSet[0] + "\".")
//    }
//    if len(sSet[1]) == 0 {
//        return NewCherryFileError(filepath, sLine, "empty room sound url.")
//    }
//    if sSet[1][0] != '"' || sSet[1][len(sSet[1])-1] != '"' {
//        return NewCherryFileError(filepath, sLine, "room sound url must be set with a valid string.")
//    }
//    return nil
//}

//func room_sound_setter(cherryRooms *config.CherryRooms, roomName string, mSet, sSet []string) {
//    //  WARN(Santiago): by now we will pass the sound template as empty.
//    cherryRooms.AddSound(roomName, mSet[0], mSet[1][1:len(mSet[1])-1], "", sSet[1][1:len(sSet[1])-1])
//}
