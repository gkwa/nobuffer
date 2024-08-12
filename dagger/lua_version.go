package main

import "fmt"

type LuaVersion struct {
	version string
}

func NewLuaVersion(version string) LuaVersion {
	if version == "" {
		return LuaVersion{"5.4"}
	}
	return LuaVersion{version}
}

func (lv LuaVersion) String() string {
	return lv.version
}

func (lv LuaVersion) PackageName() string {
	return fmt.Sprintf("lua%s", lv.version)
}

func (lv LuaVersion) DevPackageName() string {
	return fmt.Sprintf("lua%s-dev", lv.version)
}

func (lv LuaVersion) Executable() string {
	return fmt.Sprintf("lua%s", lv.version)
}

func (lv LuaVersion) LuaIncludePath() string {
	return fmt.Sprintf("--with-lua-include=/usr/include/lua%s", lv.version)
}

func (lv LuaVersion) InterpreterFlag() string {
	return fmt.Sprintf("--with-lua-interpreter=lua%s", lv.version)
}
