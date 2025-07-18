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

	// video_width , video_height :=input_video.Resolution()

	input_image_path := core.Image{Path:"assets/homer.jpg", PosX: 0 , PosY: 0}
	input_audio_path := core.Audio{Path: "assets/audio.mp3"}
	output_video_path := "temp/output.mp4"

	cmd := "/usr/bin/ffmpeg"

	// Esto sería "crear video con una sola imágen"
	
	/*"ffmpeg -i output.mp4 -vf subtitles=subtitles.srt mysubtitledmovie.mp4"*/	

	params:= []string{
		"-i", input_video.Path,
		"-i", input_image_path.Path,
		"-filter_complex",fmt.Sprintf("[0:v][1:v] overlay=%d:%d:enable='between(t,0,%d)'",input_image_path.PosX, input_image_path.PosY, 10),
		"-pix_fmt", "yuv420p",
		"-c:a", "copy",
		output_video_path,
	}
	binds.RunCommand(cmd, params)

	output_video_audio := "temp/output_audio.mp4"
	paramsAudio := []string{
		"-i", output_video_path,
		"-i", input_audio_path.Path,
		"-c", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		output_video_audio,
	}
	binds.RunCommand(cmd, paramsAudio)


}
