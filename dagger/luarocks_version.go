package main

import "fmt"

type LuarocksVersion struct {
	version string
}

func NewLuarocksVersion(version string) LuarocksVersion {
	if version == "" {
		return LuarocksVersion{"3.11.1"}
	}
	return LuarocksVersion{version}
}

func (lv LuarocksVersion) String() string {
	return lv.version
}

func (lv LuarocksVersion) DownloadURL() string {
	return fmt.Sprintf("https://luarocks.org/releases/luarocks-%s.tar.gz", lv.version)
}

func (lv LuarocksVersion) ArchiveName() string {
	return fmt.Sprintf("/src/luarocks-%s.tar.gz", lv.version)
}

func (lv LuarocksVersion) ExtractedDirPath() string {
	return fmt.Sprintf("/src/luarocks-%s", lv.version)
}
