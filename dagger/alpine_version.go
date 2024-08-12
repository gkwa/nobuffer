package main

import "fmt"

type AlpineVersion struct {
	version string
}

func NewAlpineVersion(version string) AlpineVersion {
	if version == "" {
		return AlpineVersion{"latest"}
	}
	return AlpineVersion{version}
}

func (av AlpineVersion) String() string {
	return av.version
}

func (av AlpineVersion) ImageName() string {
	return fmt.Sprintf("alpine:%s", av.version)
}
