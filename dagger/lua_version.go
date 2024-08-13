package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LuaVersion struct {
	version string
}

func NewLuaVersion(version string) (LuaVersion, error) {
	if version == "" {
		v, err := FetchLatestLuaVersion()
		if err != nil {
			return LuaVersion{}, fmt.Errorf("failed to fetch latest Lua version: %w", err)
		}
		return LuaVersion{v}, nil
	}
	return LuaVersion{version}, nil
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

func (lv LuaVersion) AssertSingleLuaH() string {
	return `test $(find /usr/include -type f -name lua.h | wc -l) -eq 1`
}

func (lv LuaVersion) InterpreterFlag() string {
	return fmt.Sprintf("--with-lua-interpreter=lua%s", lv.version)
}

func (lv LuaVersion) GetConfigureArgs() []string {
	return strings.Split(fmt.Sprintf("./configure --prefix=/usr --with-lua=/usr %s %s", lv.LuaIncludePath(), lv.InterpreterFlag()), " ")
}

func FetchLatestLuaVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/lua/lua/releases")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var releases []struct {
		TagName string `json:"tag_name"`
	}
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return "", err
	}
	if len(releases) == 0 {
		return "", fmt.Errorf("no releases found")
	}

	version := strings.TrimPrefix(releases[0].TagName, "v")
	parts := strings.Split(version, ".")
	if len(parts) > 2 {
		version = strings.Join(parts[:2], ".")
	}
	return version, nil
}
