package main

import (
	"encoding/json"
	"fmt"
	"go-ffmpeg/binds"
	"go-ffmpeg/core"
	"go-ffmpeg/message"
	"go-ffmpeg/minio"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
)



func main(){
	os.RemoveAll("/temp")
	err := godotenv.Load()
	tempFolder := "temp_assets"
	minioUrl := os.Getenv("PUBLIC_MINIO_URL")
	publicMinioAccessKey := os.Getenv("PUBLIC_MINIO_ACCESS_KEY")
	publicMinioSecretKey := os.Getenv("PUBLIC_MINIO_SECRET_KEY")
	var useSSL bool 

	if os.Getenv("USESSL") == "true"{
		useSSL = true
	}else{
		useSSL = false
	}
	fmt.Println(useSSL)
		
	if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
	}

	currentFileGetter := minio.NewMinioFileGetter(minioUrl, publicMinioAccessKey, publicMinioSecretKey,useSSL)

	topic := "subtitles-audios"
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BROKER"),
		"group.id":"go-ffmpeg",
		"auto.offset.reset": "earliest",
	})

	if err != nil{
		panic(err)
	}
	
	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		panic(err)
	}

	run := true
	for run {

		msg, err := c.ReadMessage(time.Second*5)
		if err == nil{
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			
			var m message.Message

			re := regexp.MustCompile(`'([^']*?)'`)
			fixedJSON := re.ReplaceAllString(string(msg.Value), `"$1"`)

			err = json.Unmarshal([]byte(fixedJSON), &m)
			if err != nil {
				log.Fatalf("Error parsing message JSON: %s", err)
			}


			fmt.Println("This is the message", m)
			// Aca elegir el fileGetter según el que diga en el mensaje del producer.
			
			audioDestinationPaths := m.DownloadAudio(currentFileGetter, tempFolder)
			subtitleDestinationPaths := m.DownloadSubtitles(currentFileGetter,tempFolder)
			gameplayDestinationPath := m.DownloadGameplay(currentFileGetter, tempFolder)
			
			// Video Creation
			os.RemoveAll("temp")
			os.Mkdir("temp",0755)
			start := time.Now()
			input_video := core.Video{Path:gameplayDestinationPath}

			w_video , h_video := input_video.Resolution()

			input_image_path := core.Image{Path:"assets/homer.png", PosX: 0 , PosY: uint16(float32(h_video) * 0.30)}
			input_audio_path := core.Audio{Path: audioDestinationPaths[0]}
			input_subtitles_path := core.Subtitles{Path: subtitleDestinationPaths[0]}
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
			
		}else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}else {
			fmt.Printf("No new message. Waiting... (%s)\n", os.Getenv("KAFKA_BROKER"))

		}
		}
	c.Close()

	

}
