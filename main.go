package main

import (
	"fmt"
	"go-ffmpeg/binds"
	"go-ffmpeg/core"
	"time"
	"log"
)

func main(){
	start := time.Now()

	input_video := core.Video{Path:"assets/undertale.mp4"}

	w_video , h_video := input_video.Resolution()

	input_image_path := core.Image{Path:"assets/homer.png", PosX: 0 , PosY: uint16(float32(h_video) * 0.30)}
	input_audio_path := core.Audio{Path: "assets/audio.mp3"}
	
	input_subtitles_path := core.Subtitles{Path: "assets/subtitles.ass"}
	output_video_path := "temp/output.mp4"

	cmd := "/usr/bin/ffmpeg"

	output_video_audio := "temp/output_audio.mp4"
	paramsAudio := []string{
		"-i", input_video.Path,
		"-i", input_audio_path.Path,
		"-c", "copy",
		"-map", "0:v:0",
		"-map", "1:a:0",
		"-shortest",
		output_video_audio,
	}
	binds.RunCommand(cmd, paramsAudio)

	params:= []string{
		"-i", output_video_audio,
		"-i", input_image_path.Path,
		"-filter_complex",fmt.Sprintf("[1:v] scale=%d:-1 [resized]; [0:v][resized] overlay=%d:%d:enable='between(t,0,%d)'",w_video/2,input_image_path.PosX, input_image_path.PosY, input_video.Length()),
		"-pix_fmt", "yuv420p",
		"-c:a", "copy",
		output_video_path,
	}
	binds.RunCommand(cmd, params)

	
	
	output_video_subtitle := "temp/output_subtitles.mp4"
	
	paramsSubtitle := []string{
		"-i", output_video_path,
		"-vf", fmt.Sprintf("ass=%s", input_subtitles_path.Path),
		"-c:a", "copy", 
		output_video_subtitle,

	}
	binds.RunCommand(cmd, paramsSubtitle)
	
	elapsed := time.Since(start)
    log.Printf("Video creation took %s", elapsed)

}
