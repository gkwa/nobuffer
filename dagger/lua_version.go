package main

import (
	"fmt"
	"strings"
)

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
	return `--with-lua-include=$(find /usr/include -type f -name lua.h | sed 's#/lua\.h##' | head -n 1)`
}

func (lv LuaVersion) InterpreterFlag() string {
	return fmt.Sprintf("--with-lua-interpreter=lua%s", lv.version)
}

func (lv LuaVersion) GetConfigureArgs() []string {
	return strings.Split(fmt.Sprintf("./configure --prefix=/usr --with-lua=/usr %s %s", lv.LuaIncludePath(), lv.InterpreterFlag()), " ")
}
