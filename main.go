package main

import (
	"encoding/json"
	"fmt"
	"go-ffmpeg/core"
	"go-ffmpeg/message"
	"go-ffmpeg/minio"
	"os"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/joho/godotenv"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel) // full verbosity for debugging
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	log.Info("🔧 Starting application initialization")
}

func LoadEnv() {
	// Try loading .env only if it exists
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Warnf("⚠️  Could not load .env file: %v", err)
		} else {
			log.Info("✅ .env file successfully loaded")
		}
	} else {
		log.Info("ℹ️  No .env file found — using Docker or system environment variables")
	}

	// Dump environment variables for debugging
	log.Debug("========= 🔍 ENVIRONMENT VARIABLES (BEGIN) =========")

	envKeys := []string{
		"ENVIRONMENT",
		"LOG_LEVEL",
		"DEBUG_MODE",
		"PUBLIC_MINIO_URL",
		"PUBLIC_MINIO_ACCESS_KEY",
		"PUBLIC_MINIO_SECRET_KEY",
		"API_GATEWAY_URL",
		"ADMIN_API",
		"ADMIN_API_NO_BACKSLASH",
		"KAFKA_URL",
		"KAFKA_BROKER",
		"OLLAMA_IP",
		"OLLAMA_IP_API",
		"MONGO_URL",
		"TTS_RVC_IMAGE",
		"USESSL",
		"SECURE",
		"JWT_KEY",
	}

	for _, key := range envKeys {
		val := os.Getenv(key)
		if val == "" {
			log.Warnf("🚨 %s is not set!", key)
		} else {
			// Mask sensitive keys
			if strings.Contains(strings.ToLower(key), "key") || strings.Contains(strings.ToLower(key), "secret") {
				log.Debugf("%s = *****", key)
			} else {
				log.Debugf("%s = %s", key, val)
			}
		}
	}

	log.Debug("========= 🔍 ENVIRONMENT VARIABLES (END) =========")

	// Example check for USESSL
	useSSL := strings.ToLower(os.Getenv("USESSL")) == "true"
	log.Infof("🔐 USESSL = %v", useSSL)
}

func main() {
	log.Info("🚀 Starting main application")
	LoadEnv()

	// Remove old temp folders
	os.RemoveAll("/temp")
	os.RemoveAll("temp_assets")
	os.Mkdir("temp_assets", 0755)

	// Retrieve environment variables
	minioUrl := os.Getenv("PUBLIC_MINIO_URL")
	publicMinioAccessKey := os.Getenv("PUBLIC_MINIO_ACCESS_KEY")
	publicMinioSecretKey := os.Getenv("PUBLIC_MINIO_SECRET_KEY")
	apiGatewayUrl := os.Getenv("API_GATEWAY_URL")
	adminApi := os.Getenv("ADMIN_API")
	kafkaBroker := os.Getenv("KAFKA_BROKER")

	useSSL := strings.ToLower(os.Getenv("USESSL")) == "true"
	log.Infof("🔗 Connecting to MinIO at %s (useSSL=%v)", minioUrl, useSSL)

	currentFileGetter := minio.NewMinioFileGetter(minioUrl, publicMinioAccessKey, publicMinioSecretKey, useSSL)
	topic := "subtitles-audios"

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          "go-ffmpeg",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("❌ Failed to create Kafka consumer: %v", err)
	}
	defer c.Close()

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		log.Fatalf("❌ Failed to subscribe to topic: %v", err)
	}
	log.Infof("✅ Subscribed to Kafka topic: %s", topic)

	for {
		msg, err := c.ReadMessage(time.Second * 5)
		if err == nil {
			log.Infof("📩 Message on %s: %s", msg.TopicPartition, string(msg.Value))
			
			var m message.Message
			re := regexp.MustCompile(`'([^']*?)'`)
			fixedJSON := re.ReplaceAllString(string(msg.Value), `"$1"`)
			
			if err := json.Unmarshal([]byte(fixedJSON), &m); err != nil {
				log.Errorf("❌ Error parsing message JSON: %s", err)
				continue
			}

			log.Infof("📦 Parsed message: %+v", m)

			audioPaths := m.DownloadAudio(currentFileGetter, "temp_assets")
			subPaths := m.DownloadSubtitles(currentFileGetter, "temp_assets")
			gameplayPath := m.DownloadGameplay(currentFileGetter, "temp_assets")

			log.Debugf("🎵 Audio paths: %v", audioPaths)
			log.Debugf("💬 Subtitle paths: %v", subPaths)
			log.Debugf("🎮 Gameplay path: %s", gameplayPath)

			os.RemoveAll("temp")
			os.Mkdir("temp", 0755)

			start := time.Now()
			inputVideo := core.Video{Path: gameplayPath}
			_, hVideo := inputVideo.Resolution()

			imagePath, err := m.DownloadRandomImage("temp_assets", adminApi)
			if err != nil {
				log.Errorf("⚠️  Error obtaining random image: %v", err)
				imagePath = "assets/404.png"
			}

			inputImage := core.Image{
				Path: imagePath,
				PosX: 0,
				PosY: uint16(float32(hVideo) * 0.30),
			}

			inputAudio := core.Audio{Path: audioPaths[0]}
			inputSubs := core.Subtitles{Path: subPaths[0]}
			outputVideo := "temp/output.mp4"
			cmd := "/usr/bin/ffmpeg"

			videoBuilder := core.NormalVideoBuilder{
				Video:     inputVideo,
				Audio:     inputAudio,
				Image:     inputImage,
				Subtitles: inputSubs,
			}

			log.Infof("🎬 Starting FFmpeg video creation...")
			videoBuilder.CreateVideo(cmd, outputVideo)

			bucket := "videos-homero"
			fileName := fmt.Sprintf("%s.mp4", m.Tema)
			videoUploader := core.VideoUploader{FileGetter: currentFileGetter}

			log.Infof("⬆️  Uploading video %s to bucket %s", fileName, bucket)
			videoUploader.UploadVideo(bucket, fileName, outputVideo, apiGatewayUrl, &m)

			elapsed := time.Since(start)
			log.Infof("✅ Video creation took %s", elapsed)

		} else if kafkaErr, ok := err.(kafka.Error); ok && !kafkaErr.IsTimeout() {
			log.Errorf("❌ Kafka consumer error: %v", kafkaErr)
		} else {
			log.Debugf("⌛ No new message. Waiting... (%s)", kafkaBroker)
		}
	}
}
