package symbols

import "reflect"

// This file is required for yaegi to extract the symbols from the shared package
// and make them available to the interpreter. This is necessary for the plugins
// to be able to access the shared package.

var Symbols = map[string]map[string]reflect.Value{}

var MapTypes = map[reflect.Value][]reflect.Type{}

func init() {
	Symbols["."] = map[string]reflect.Value{
		"MapTypes": reflect.ValueOf(MapTypes),
	}
}

//go:generate $GOPATH/bin/yaegi extract git.dbyte.xyz/distro/gerry/bot
//go:generate $GOPATH/bin/yaegi extract git.dbyte.xyz/distro/gerry/shared
