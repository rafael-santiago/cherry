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
    "strings"
    "fmt"
    "net/url"
)

func cherryDefaultHTTPReplyHeader(statusCode int, closeConnection bool) string {
    var header string = "HTTP/1.1 "
    switch (statusCode) {
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

func MakeReplyBuffer(buffer string, statusCode int, closeConnection bool) []byte {
    return []byte(strings.Replace(cherryDefaultHTTPReplyHeader(statusCode, closeConnection) + buffer,
                                  "{{.content-length}}",
                                  fmt.Sprintf("%d", len(buffer)),
                                  -1))
}

func GetHTTPFieldFromBuffer(field, buffer string) string {
    index := strings.Index(buffer, field + ":")
    if index == -1 {
        return ""
    }
    retval := ""
    if len(buffer) > index + len(field) {
        for _, b := range buffer[index+len(field + ":"):] {
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

func GetFieldsFromPost(buffer string) map[string]string {
    if !strings.HasPrefix(buffer, "POST /") {
        return make(map[string]string)
    }
    index := strings.Index(buffer, "\r\n\r\n")
    if index == -1 || len(buffer) <= index + 4 {
        return make(map[string]string)
    }
    return splitFields(buffer[index+4:])
}

func GetFieldsFromGet(buffer string) map[string]string {
    if !strings.HasPrefix(buffer, "GET /") {
        return make(map[string]string)
    }
    var startIndex int = 0
    var endIndex int = 0
    for buffer[startIndex] != '&' {
        startIndex++
    }
    for buffer[endIndex] != '\n' && buffer[endIndex] != '\r' && endIndex < len(buffer) {
        endIndex++
    }
    return splitFields(buffer[startIndex+1:endIndex])
}
