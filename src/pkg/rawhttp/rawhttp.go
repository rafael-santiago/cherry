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

func splitFields(buffer string) map[string]string {
	dataValue := strings.Split(buffer, "&")
	retval := make(map[string]string)
	for _, dv := range dataValue {
		set := strings.Split(dv, "=")
		if len(set) == 2 {
			data, _ := url.QueryUnescape(set[1])
			retval[set[0]] = data
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
