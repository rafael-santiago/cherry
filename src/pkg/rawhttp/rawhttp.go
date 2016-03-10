/*
Package rawhttp implements conveniences related with HTTP request/reply operations.
 --
 *
 *                               Copyright (C) 2015 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
*/
package rawhttp

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

var charLT map[string]string
var charLTInitialized bool

func initCharLT() {
	charLT = make(map[string]string)
	charLT["%E2%82%AC"] = "%80"
	charLT["%E2%80%9A"] = "%82"
	charLT["%C6%92"] = "%83"
	charLT["%E2%80%9E"] = "%84"
	charLT["%E2%80%A6"] = "%85"
	charLT["%E2%80%A0"] = "%86"
	charLT["%E2%80%A1"] = "%87"
	charLT["%CB%86"] = "%88"
	charLT["%E2%80%B0"] = "%89"
	charLT["%C5%A0"] = "%8A"
	charLT["%E2%80%B9"] = "%8B"
	charLT["%C5%92"] = "%8C"
	charLT["%C5%8D"] = "%8D"
	charLT["%C5%BD"] = "%8E"
	charLT["%C2%90"] = "%90"
	charLT["%E2%%80%98"] = "%91"
	charLT["%E2%80%99"] = "%92"
	charLT["%E2%80%9C"] = "%93"
	charLT["%E2%80%9D"] = "%94"
	charLT["%E2%80%A2"] = "%95"
	charLT["%E2%80%93"] = "%96"
	charLT["%E2%80%94"] = "%97"
	charLT["%CB%9C"] = "%98"
	charLT["%E2%84"] = "%99"
	charLT["%C5%A1"] = "%9A"
	charLT["%E2%80"] = "%9B"
	charLT["%C5%93"] = "%9C"
	charLT["%C5%BE"] = "%9E"
	charLT["%C5%B8"] = "%9F"
	charLT["%C2%A0"] = "%A0"
	charLT["%C2%A1"] = "%A1"
	charLT["%C2%A2"] = "%A2"
	charLT["%C2%A3"] = "%A3"
	charLT["%C2%A4"] = "%A4"
	charLT["%C2%A5"] = "%A5"
	charLT["%C2%A6"] = "%A6"
	charLT["%C2%A7"] = "%A7"
	charLT["%C2%A8"] = "%A8"
	charLT["%C2%A9"] = "%A9"
	charLT["%C2%AA"] = "%AA"
	charLT["%C2%AB"] = "%AB"
	charLT["%C2%AC"] = "%AC"
	charLT["%C2%AD"] = "%AD"
	charLT["%C2%AE"] = "%AE"
	charLT["%C2%AF"] = "%AF"
	charLT["%C2%B0"] = "%B0"
	charLT["%C2%B1"] = "%B1"
	charLT["%C2%B2"] = "%B2"
	charLT["%C2%B3"] = "%B3"
	charLT["%C2%B4"] = "%B4"
	charLT["%C2%B5"] = "%B5"
	charLT["%C2%B6"] = "%B6"
	charLT["%C2%B7"] = "%B7"
	charLT["%C2%B8"] = "%B8"
	charLT["%C2%B9"] = "%B9"
	charLT["%C2%BA"] = "%BA"
	charLT["%C2%BB"] = "%BB"
	charLT["%C2%BC"] = "%BC"
	charLT["%C2%BD"] = "%BD"
	charLT["%C2%BE"] = "%BE"
	charLT["%C2%BF"] = "%BF"
	charLT["%C3%80"] = "%C0"
	charLT["%C3%81"] = "%C1"
	charLT["%C3%82"] = "%C2"
	charLT["%C3%83"] = "%C3"
	charLT["%C3%84"] = "%C4"
	charLT["%C3%85"] = "%C5"
	charLT["%C3%86"] = "%C6"
	charLT["%C3%87"] = "%C7"
	charLT["%C3%88"] = "%C8"
	charLT["%C3%89"] = "%C9"
	charLT["%C3%8A"] = "%CA"
	charLT["%C3%8B"] = "%CB"
	charLT["%C3%8C"] = "%CC"
	charLT["%C3%8D"] = "%CD"
	charLT["%C3%8E"] = "%CE"
	charLT["%C3%8F"] = "%CF"
	charLT["%C3%90"] = "%D0"
	charLT["%C3%91"] = "%D1"
	charLT["%C3%92"] = "%D2"
	charLT["%C3%93"] = "%D3"
	charLT["%C3%94"] = "%D4"
	charLT["%C3%95"] = "%D5"
	charLT["%C3%96"] = "%D6"
	charLT["%C3%97"] = "%D7"
	charLT["%C3%98"] = "%D8"
	charLT["%C3%99"] = "%D9"
	charLT["%C3%9A"] = "%DA"
	charLT["%C3%9B"] = "%DB"
	charLT["%C3%9C"] = "%DC"
	charLT["%C3%9D"] = "%DD"
	charLT["%C3%9E"] = "%DE"
	charLT["%C3%9F"] = "%DF"
	charLT["%C3%A0"] = "%E0"
	charLT["%C3%A1"] = "%E1"
	charLT["%C3%A2"] = "%E2"
	charLT["%C3%A3"] = "%E3"
	charLT["%C3%A4"] = "%E4"
	charLT["%C3%A5"] = "%E5"
	charLT["%C3%A6"] = "%E6"
	charLT["%C3%A7"] = "%E7"
	charLT["%C3%A8"] = "%E8"
	charLT["%C3%A9"] = "%E9"
	charLT["%C3%AA"] = "%EA"
	charLT["%C3%AB"] = "%EB"
	charLT["%C3%AC"] = "%EC"
	charLT["%C3%AD"] = "%ED"
	charLT["%C3%AE"] = "%EE"
	charLT["%C3%AF"] = "%EF"
	charLT["%C3%B0"] = "%F0"
	charLT["%C3%B1"] = "%F1"
	charLT["%C3%B2"] = "%F2"
	charLT["%C3%B3"] = "%F3"
	charLT["%C3%B4"] = "%F4"
	charLT["%C3%B5"] = "%F5"
	charLT["%C3%B6"] = "%F6"
	charLT["%C3%B7"] = "%F7"
	charLT["%C3%B8"] = "%F8"
	charLT["%C3%B9"] = "%F9"
	charLT["%C3%BA"] = "%FA"
	charLT["%C3%BB"] = "%FB"
	charLT["%C3%BC"] = "%FC"
	charLT["%C3%BD"] = "%FD"
	charLT["%C3%BE"] = "%FE"
	charLT["%C3%BF"] = "%FF"
	charLTInitialized = true
}

func cherryDefaultHTTPReplyHeader(statusCode int, closeConnection bool) string {
	var header = "HTTP/1.1 "
	switch statusCode {
	case 200:
		header += "200 OK"
		break

	case 404:
		header += "404 NOT FOUND"
		break

	case 403:
		header += "403 FORBIDDEN"
		break

	default:
		header += "501 NOT IMPLEMENTED"
		break
	}
	if closeConnection {
		header += "\n\rConnection: close\n" +
			"Server: Cherry/0.1\n" +
			"Accept-ranges: bytes\n" +
			"Content-type: text/html\n" +
			"Content-length: {{.content-length}}\n\n"
	} else {
		header += "200 Document follows\nContent-type: text/html\n\n"
	}
	return header
}

// MakeReplyBuffer assembles the reply buffer based on the statusCode.
func MakeReplyBuffer(buffer string, statusCode int, closeConnection bool) []byte {
	return []byte(strings.Replace(cherryDefaultHTTPReplyHeader(statusCode, closeConnection)+buffer,
		"{{.content-length}}",
		fmt.Sprintf("%d", len(buffer)),
		-1))
}

// MakeReplyBufferByFilePath assembles the reply buffer base on the file data and the statusCode.
func MakeReplyBufferByFilePath(filePath string, statusCode int, closeConnection bool) []byte {
	buffer, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte("")
	}
	tempReply := strings.Replace(cherryDefaultHTTPReplyHeader(statusCode, closeConnection)+string(buffer),
		"{{.content-length}}",
		fmt.Sprintf("%d", len(buffer)),
		-1)
	tempReply = strings.Replace(tempReply, "Content-type: text/html", getContentTypeFromFilePath(filePath), -1)
	return []byte(tempReply)
}

func getContentTypeFromFilePath(filePath string) string {
	if strings.HasSuffix(filePath, ".gif") {
		return "Content-type: image/gif"
	}
	if strings.HasSuffix(filePath, ".jpeg") || strings.HasSuffix(filePath, ".jpg") {
		return "Content-type: image/jpeg"
	}
	if strings.HasSuffix(filePath, ".png") {
		return "Content-type: image/png"
	}
	if strings.HasSuffix(filePath, ".bmp") {
		return "Content-type: image/bmp"
	}
	return "Content-type: text/plain"
}

// GetHTTPFieldFromBuffer returns data carried by some HTTP field inside a HTTP request.
func GetHTTPFieldFromBuffer(field, buffer string) string {
	index := strings.Index(buffer, field+":")
	if index == -1 {
		return ""
	}
	retval := ""
	if len(buffer) > index+len(field) {
		for _, b := range buffer[index+len(field+":"):] {
			if b == '\n' || b == '\r' {
				break
			}
			if len(retval) == 0 && (b == ' ' || b == '\t') {
				continue
			}
			retval += string(b)
		}
	}
	return retval
}

func utf8Unescape(data string) string {
	if !charLTInitialized {
		initCharLT()
	}
	var buffer string = data
	for k, v := range charLT {
		buffer = strings.Replace(buffer, k, v, -1)
	}
	finalData, _ := url.QueryUnescape(buffer)
	return finalData
}

func splitFields(buffer string) map[string]string {
	dataValue := strings.Split(buffer, "&")
	retval := make(map[string]string)
	for _, dv := range dataValue {
		set := strings.Split(dv, "=")
		if len(set) == 2 {
			retval[set[0]] = utf8Unescape(set[1])
		}
	}
	return retval
}

// GetFieldsFromPost returns a map containing all fields in a HTTP post.
func GetFieldsFromPost(buffer string) map[string]string {
	if !strings.HasPrefix(buffer, "POST /") {
		return make(map[string]string)
	}
	index := strings.Index(buffer, "\r\n\r\n")
	if index == -1 || len(buffer) <= index+4 {
		return make(map[string]string)
	}
	return splitFields(buffer[index+4:])
}

// GetFieldsFromGet returns a map containing all fields in a HTTP get.
func GetFieldsFromGet(buffer string) map[string]string {
	if !strings.HasPrefix(buffer, "GET /") {
		return make(map[string]string)
	}
	var startIndex int
	var endIndex int
	for buffer[startIndex] != '&' {
		startIndex++
	}
	for buffer[endIndex] != '\n' && buffer[endIndex] != '\r' && endIndex < len(buffer) {
		endIndex++
	}
	return splitFields(buffer[startIndex+1 : endIndex])
}
