package core

import (
	"go-ffmpeg/binds"
	"log"
	"strconv"
	"strings"
)

type Audio struct {
	Path string
}

func (v Audio) Length() uint16 {
	params:= []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		"-i", v.Path,
	}
	out := binds.RunCommandWithOutput(ffprobe, params)
	parsedOut , err := strconv.ParseFloat(strings.TrimSpace(out), 32)

	if err!= nil{
		log.Fatal("Error al convertir la duración del video a string", err)
	}
	return uint16(parsedOut)
}