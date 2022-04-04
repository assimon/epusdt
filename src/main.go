package main

import (
	"github.com/assimon/luuu/bootstrap"
	"github.com/assimon/luuu/config"
	"github.com/gookit/color"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			color.Error.Println("[Start Server Err!!!] ", err)
		}
	}()
	color.Green.Printf("%s\n", "  _____                     _ _   \n | ____|_ __  _   _ ___  __| | |_ \n |  _| | '_ \\| | | / __|/ _` | __|\n | |___| |_) | |_| \\__ \\ (_| | |_ \n |_____| .__/ \\__,_|___/\\__,_|\\__|\n       |_|                        ")
	color.Infof("Epusdt version(%s) Powered by %s %s \n", config.GetAppVersion(), "assimon", "https://github.com/assimon/epusdt")
	bootstrap.Start()
}
