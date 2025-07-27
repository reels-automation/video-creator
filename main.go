package main

import (
	"encoding/json"
	"fmt"
	"go-ffmpeg/core"
	"go-ffmpeg/message"
	"go-ffmpeg/minio"
	"os"
	"regexp"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
)

func init() {
    log.SetOutput(os.Stdout)
    log.SetLevel(log.InfoLevel)
    log.SetFormatter(&log.TextFormatter{
        DisableColors: true,
        FullTimestamp: true,
    })
    log.Info("Starting application inside container")
}

func main(){
	log.Info("Starting application inside container")
	os.RemoveAll("/temp")
	err := godotenv.Load()
	tempFolder := "temp_assets"
	minioUrl := os.Getenv("PUBLIC_MINIO_URL")
	publicMinioAccessKey := os.Getenv("PUBLIC_MINIO_ACCESS_KEY")
	publicMinioSecretKey := os.Getenv("PUBLIC_MINIO_SECRET_KEY")
	apiGatewayUrl := os.Getenv("API_GATEWAY_URL")
	var useSSL bool 

	if os.Getenv("USESSL") == "true"{
		useSSL = true
	}else{
		useSSL = false
	}
		
	if err != nil {
		log.Errorf("Error loading .env file: %v", err)
	}

	currentFileGetter := minio.NewMinioFileGetter(minioUrl, publicMinioAccessKey, publicMinioSecretKey,useSSL)

	topic := "subtitles-audios"
	
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
    "bootstrap.servers": os.Getenv("KAFKA_BROKER"),
    "group.id": "go-ffmpeg",
    "auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer c.Close()
		
	err = c.SubscribeTopics([]string{topic}, nil)

	if err != nil {
		log.Fatal(err)
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
				log.Errorf("Error parsing message JSON: %s", err)
			}

			fmt.Println("This is the message", m)
			os.RemoveAll("temp_assets")
			os.Mkdir("temp_assets",0755)
			audioDestinationPaths := m.DownloadAudio(currentFileGetter, tempFolder)
			subtitleDestinationPaths := m.DownloadSubtitles(currentFileGetter,tempFolder)
			gameplayDestinationPath := m.DownloadGameplay(currentFileGetter, tempFolder)
			
			// Video Creation
			os.RemoveAll("temp")
			os.Mkdir("temp",0755)
			start := time.Now()
			input_video := core.Video{Path:gameplayDestinationPath}

			_ , h_video := input_video.Resolution()

			input_image_path := core.Image{Path:"assets/homer.png", PosX: 0 , PosY: uint16(float32(h_video) * 0.30)}
			input_audio_path := core.Audio{Path: audioDestinationPaths[0]}
			input_subtitles_path := core.Subtitles{Path: subtitleDestinationPaths[0]}
			output_video_path := "temp/output.mp4"
			cmd := "/usr/bin/ffmpeg"

			video_builder := core.NormalVideoBuilder{
				Video: input_video,
				Audio: input_audio_path,
				Image: input_image_path,
				Subtitles: input_subtitles_path,
			}
			video_builder.CreateVideo(cmd,output_video_path)
			bucket := "videos-homero"
			fileName := m.Tema + ".mp4"
			video_uploader := core.VideoUploader{FileGetter: currentFileGetter}
			video_uploader.UploadVideo(bucket,fileName,output_video_path,apiGatewayUrl, &m)
			elapsed := time.Since(start)
			log.Infof("Video creation took %s", elapsed)
			
		}else if !err.(kafka.Error).IsTimeout() {
			log.Errorf("Consumer error: %v (%v)\n", err, msg)
		}else {
			log.Infof("No new message. Waiting... (%s)\n", os.Getenv("KAFKA_BROKER"))
		}
		}
	c.Close()
}
