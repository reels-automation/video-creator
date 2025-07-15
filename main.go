package main

import (
	"go-ffmpeg/core"
	"go-ffmpeg/binds"
	"fmt"
)

func main(){
	
	/*
		var command1 string = "date"
	out  := runCommandWithOutput(command1)
	fmt.Println(out)
	command2 := "dater"
	out2  := runCommandWithOutput(command2)
	fmt.Println(out2)
	*/

	
	input_video := core.Video{Path:"assets/subway1.mp4"}
	input_image_path := "assets/homer.jpg"
	output_video_path := "assets/output.mp4"

	cmd := "/usr/bin/ffmpeg"

	// Esto sería "crear video con una sola imágen"
	
	/*"ffmpeg -i output.mp4 -vf subtitles=subtitles.srt mysubtitledmovie.mp4"*/	

	params:= []string{
		"-i", input_video.Path,
		"-i", input_image_path,
		"-filter_complex",fmt.Sprintf("[0:v][1:v] overlay=25:25:enable='between(t,0,%d)'",input_video.Length()),
		"-pix_fmt", "yuv420p",
		"-c:a", "copy",
		output_video_path,
	}
	binds.RunCommand(cmd, params)

}
