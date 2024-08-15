package main

import "fmt"

type ImageVersion struct {
	name    string
	version string
}

func NewImageVersion(name, version string) ImageVersion {
	if name == "" {
		name = "pandoc/core"
	}
	if version == "" {
		version = "latest"
	}
	return ImageVersion{name, version}
}

func (iv ImageVersion) String() string {
	return fmt.Sprintf("%s:%s", iv.name, iv.version)
}

func (iv ImageVersion) ImageName() string {
	return iv.String()
}
