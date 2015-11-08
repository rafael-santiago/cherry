/*
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

func cherry_default_http_reply_header(status_code int, close_connection bool) string {
    var header string = "HTTP/1.1 "
    switch (status_code) {
        case 200:
            header += "200 OK\n\r"
            break

        case 404:
            header += "404 NOT FOUND\n\r"
            break

        case 403:
            header += "403 FORBIDDEN\n\r"
            break

        default:
            header += "501 NOT IMPLEMENTED\n\r"
            break
    }
    if close_connection {
        header += "Connection: close\n\r"
    }
    header += "Server: Cherry/0.1\n\r" +
              "Accept-ranges: bytes\n\r" +
              "Content-type: text/html\n\r" +
              "Content-length: {{.content-length}}\n\r\n\r"
    return header
}

func MakeReplyBuffer(buffer string, status_code int, close_connection bool) []byte {
    return []byte(strings.Replace(cherry_default_http_reply_header(status_code, close_connection) + buffer,
                                  "{{.content-length}}",
                                  fmt.Sprintf("%d", len(buffer)),
                                  -1))
}

func GetHttpFieldFromBuffer(field, buffer string) string {
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

func split_fields(buffer string) map[string]string {
    data_value := strings.Split(buffer, "&")
    retval := make(map[string]string)
    for _, dv := range data_value {
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
    return split_fields(buffer[index+4:])
}

func GetFieldsFromGet(buffer string) map[string]string {
    if !strings.HasPrefix(buffer, "GET /") {
        return make(map[string]string)
    }
    var start_index int = 0
    var end_index int = 0
    for buffer[start_index] != '&' {
        start_index++
    }
    for buffer[end_index] != '\n' && buffer[end_index] != '\r' && end_index < len(buffer) {
        end_index++
    }
    return split_fields(buffer[start_index+1:end_index])
}
