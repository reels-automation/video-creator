package core

import (
	"go-ffmpeg/binds"
	"log"
	"strconv"
	"strings"
)

type Video struct {
	Path string
}

func (v Video) Length() int16 {
	cmd := "/usr/bin/ffprobe"
	params:= []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"-i", v.Path,
	}
	out := binds.RunCommandWithOutput(cmd, params)
	parsedOut , err := strconv.ParseFloat(strings.TrimSpace(out), 32)

	if err!= nil{
		log.Fatal("Error al convertir la duración del video a string", err)
	}
	return int16(parsedOut)
}
