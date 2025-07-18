package tests

import (
	"fmt"
	"go-ffmpeg/core"
	"testing"
)

func TestResolution(t *testing.T) {
	video := core.Video{Path:"../assets/subway1.mp4"}
	width, height := video.Resolution()
    
	fmt.Println("Width: ", width)
	fmt.Println("Height: ", height)
}

func TestLength(t *testing.T) {
	video := core.Video{Path:"../assets/subway1.mp4"}
	length := video.Length()
    
	fmt.Println("Length: ", length)
}

