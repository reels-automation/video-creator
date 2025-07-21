package core

import (
	"fmt"
	"go-ffmpeg/binds"
)

type VideoBuilder interface{
	CreateVideo(cmd string, outputFile string)	
}

type NormalVideoBuilder struct{
	Video Video
	Subtitles Subtitles
	Audio Audio
	Image Image
}

func (n NormalVideoBuilder) CreateVideo(cmd string, outputFile string){

	width_video, _ := n.Video.Resolution()
	params := []string{
				"-i", n.Video.Path,
				"-i", n.Audio.Path,
				"-i", n.Image.Path,
				"-filter_complex", fmt.Sprintf(
				"[2:v]scale=%d:-1[img];" +
				"[0:v][img]overlay=%d:%d:enable='between(t,0,%d)'[overlaid];" +
		        "[overlaid]ass=%s[v]",
				width_video/2,
				n.Image.PosX,
				n.Image.PosY,
				n.Video.Length(),
				n.Subtitles.Path,
			),
				"-map", "[v]",
				"-map", "1:a:0",
				"-c:v", "libx264",
				"-c:a", "copy",
				"-pix_fmt", "yuv420p",
				"-shortest",
				outputFile,
	}
	binds.RunCommand(cmd, params)

}	
