/*
 *                                Copyright (C) 2016 by Rafael Santiago
 *
 * This is a free software. You can redistribute it and/or modify under
 * the terms of the GNU General Public License version 2.
 *
 */
package cherry_test

import (
	"pkg/config"
	"pkg/config/parser"
	"pkg/html"
	"os"
	"testing"
)

func PreprocessorBasicTest(t *testing.T) {
	preprocessor := html.NewHTMLPreprocessor(nil)
	if preprocessor.ExpandData("land-of-competition", "{{.FoD}} Zzz...") != "{{.FoD}} Zzz..." {
		t.Fail()
	}
	var cherry_rooms *config.CherryRooms
	cwd, _ := os.Getwd()
	os.Chdir("../../sample")
	var error *parser.CherryFileError
	cherry_rooms, error = parser.ParseCherryFile("conf/sample.cherry")
	os.Chdir(cwd)
	if error != nil {
		t.Fail()
	}
	preprocessor.Init(cherry_rooms)
	preprocessor.SetDataValue("{{.foo}}", "bar")
	preprocessor.SetDataValue("{{.bar}}", "foo")
	if preprocessor.ExpandData("aliens-on-earth", "{{.foo}}{{.bar}}") != "barfoo" {
		t.Fail()
	}
	if preprocessor.ExpandData("aliens-on-earth", "{{.greeting-message}}") != "Take meeeeee to your leader!!!" {
		t.Fail()
	}
}
