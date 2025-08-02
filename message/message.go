package message

import (
	"go-ffmpeg/minio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type AudioItem struct {
	TTSAudioName string `json:"tts_audio_name"`
	TTSDirectory string `json:"tts_audio_directory"`
	FileGetter   string `json:"file_getter"`
	Pitch        int    `json:"pitch"`
	TTSVoice     string `json:"tts_voice"`
	TTSRate      int    `json:"tts_rate"`
	PTHVoice     string `json:"pth_voice"`
}

type SubtitleItem struct {
	SubtitlesName      string `json:"subtitles_name"`
	FileGetter         string `json:"file_getter"`
	SubtitlesDirectory string `json:"subtitles_directory"`
}

type BackgroundMusic struct {
	AudioName  string `json:"audio_name"`
	FileGetter string `json:"file_getter"`
	StartTime  int `json:"start_time"`
	Duration   int `json:"duration"`
}

type ImageItem struct {
	ImageName      string `json:"image_name"`
	ImageModifier  string `json:"image_modifier"`
	FileGetter     string `json:"file_getter"`
	ImageDirectory string `json:"image_directory"`
	TimeStamp      int `json:"timestamp"`
	Duration       int `json:"duration"`
}

type Message struct {
	Tema                string            `json:"tema"`
	Usuario             string            `json:"usuario"`
	Idioma              string            `json:"idioma"`
	Personaje           string            `json:"personaje"`
	Script              string            `json:"script"`
	AudioItem           []AudioItem       `json:"audio_item"`
	SubtitleItem        []SubtitleItem    `json:"subtitle_item"`
	Author              string            `json:"author"`
	GameplayName        string            `json:"gameplay_name"`
	BackgroundMusic     []BackgroundMusic `json:"background_music"`
	Images              []ImageItem       `json:"images"`
	RandomImages        bool             `json:"random_images"`
	RandomAmountImages  int               `json:"random_amount_images"`
	GptModel            string            `json:"gpt_model"`
}

func (m Message) DownloadAudio(fileGetter minio.FileGetter, destinationFolder string) []string{
	
	destinationPathList := make([]string, len(m.AudioItem))
	
	for i :=0; i < len(m.AudioItem);i++{
		object := m.AudioItem[i].TTSAudioName
		directory := m.AudioItem[i].TTSDirectory
		destinationPath := destinationFolder + "/" + object
		destinationPathList[i] = destinationPath
		fileGetter.GetFile(directory,object,destinationPath)
	}

	return destinationPathList
}

func (m Message) DownloadSubtitles(fileGetter minio.FileGetter, destinationFolder string) []string{
	
	destinationPathList := make([]string, len(m.SubtitleItem))

	for i :=0; i < len(m.SubtitleItem);i++{
		object := m.SubtitleItem[i].SubtitlesName
		directory := m.SubtitleItem[i].SubtitlesDirectory
		destinationPath := destinationFolder + "/" + object
		destinationPathList[i] = destinationPath
		fileGetter.GetFile(directory,object,destinationPath)
	}
	return destinationPathList
}

func (m Message) DownloadGameplay(fileGetter minio.FileGetter, destinationFolder string)string{
	destinationPath := destinationFolder + "/" + m.GameplayName
	fileGetter.GetFile("gameplays", m.GameplayName,destinationPath)
	return destinationPath
}

func (m Message) DownloadAssets(fileGetter minio.FileGetter, destinationFolder string){
	m.DownloadAudio(fileGetter,destinationFolder)
	m.DownloadSubtitles(fileGetter,destinationFolder)
	m.DownloadGameplay(fileGetter, destinationFolder)
}

type randomImageResponse struct {
	ObjectName string `json:"object_name"`
	ObjectURL  string `json:"object_url"`
}

// Descarga una imagen aleatoria desde el API Gateway, usando HTTP puro
func (m *Message) DownloadRandomImage(destinationFolder string, adminApi string) (string, error) {
	personaje := m.Personaje
	endpoint := fmt.Sprintf("%s/random-image/%s", adminApi, personaje)

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("error en la petición al API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API respondió con error (%d): %s", resp.StatusCode, string(body))
	}

	var imageResp randomImageResponse
	if err := json.NewDecoder(resp.Body).Decode(&imageResp); err != nil {
		return "", fmt.Errorf("error parseando JSON del API: %w", err)
	}

	// Descargar la imagen desde la URL pública del MinIO
	imageRespHttp, err := http.Get(imageResp.ObjectURL)
	if err != nil {
		return "", fmt.Errorf("error descargando imagen desde MinIO (URL): %w", err)
	}
	defer imageRespHttp.Body.Close()

	if imageRespHttp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fallo al descargar imagen (status %d)", imageRespHttp.StatusCode)
	}

	// Crear archivo local
	localPath := filepath.Join(destinationFolder, imageResp.ObjectName)
	outFile, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("error creando archivo local: %w", err)
	}
	defer outFile.Close()

	// Copiar el contenido de la imagen al archivo
	_, err = io.Copy(outFile, imageRespHttp.Body)
	if err != nil {
		return "", fmt.Errorf("error guardando imagen localmente: %w", err)
	}

	return localPath, nil
}
