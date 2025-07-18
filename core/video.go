package core

import (
	"go-ffmpeg/binds"
	"log"
	"strconv"
	"strings"
)

var ffprobe = "/usr/bin/ffprobe"

type Video struct {
	Path string
}

func (v Video) Length() uint16 {
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

func (v Video) Resolution() (uint16,uint16){
	params := []string{
		"-v", "error",
		"-select_streams",
		"v:0",
		"-show_entries",
		"stream=width,height",
		"-of", "csv=s=x:p=0",
		"-i" , v.Path,
	}
	out := binds.RunCommandWithOutput(ffprobe, params)

	parts := strings.Split(strings.TrimSpace(out),"x")

	width , err := strconv.Atoi(parts[0]) 
	if err != nil {
		log.Fatalf("Error Converting width to Integer:\n Error %d", err)
	}
	height , err := strconv.Atoi(parts[1])

	if err != nil{
		log.Fatalf("Error Converting height to Integer:\n Error %d", err)
	}
	
	return uint16(width),uint16(height)
}