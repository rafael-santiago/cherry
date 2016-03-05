/*
 *                                Copyright (C) 2016 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package cherry_test

import (
	"pkg/rawhttp"
	"testing"
)

func TestGetFieldsFromPost(t *testing.T) {
	fields := rawhttp.GetFieldsFromPost("POST /path/script.cgi HTTP/1.0\r\n\r\nfoo=bar&bar=foo")
	if len(fields) != 2 {
		t.Fail()
	}
	if fields["foo"] != "bar" {
		t.Fail()
	}
	if fields["bar"] != "foo" {
		t.Fail()
	}
}

func TestGetHttpFieldFromBuffer(t *testing.T) {
	if rawhttp.GetHTTPFieldFromBuffer("Content-Length", "GET /abc/xy/z.log HTTP/1.0\r\nContent-Length: 255\r\n") != "255" {
		t.Fail()
	}
}

func TestGetFieldsFromGet(t *testing.T) {
	fields := rawhttp.GetFieldsFromGet("GET /abc/&foo=foo&bar=bar&\r\n\r\n")
	if len(fields) != 2 {
		t.Fail()
	}
	if fields["foo"] != "foo" {
		t.Fail()
	}
	if fields["bar"] != "bar" {
		t.Fail()
	}
}
