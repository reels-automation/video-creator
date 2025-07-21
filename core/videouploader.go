package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-ffmpeg/minio"
	"go-ffmpeg/message"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type VideoUploader struct {
	FileGetter minio.FileGetter
}

type Url struct{
	Url string  `json:"url"`
}

func PostEndpoint(apiGatewayUrl string, videoName string, message *message.Message) {
	getVideoEndpoint := apiGatewayUrl + "get-video"
	uploadMongoEndpoint := apiGatewayUrl + "add-video"
	body := []byte(fmt.Sprintf(`{"video_name":"%s"}`, videoName))

	request, err := http.NewRequest("POST", getVideoEndpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Errorf("Error creating request at Endpoint: %s \n Error: %s", getVideoEndpoint, err)
		return
	}
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		log.Errorf("Error posting at Endpoint: %s \n Error: %s", getVideoEndpoint, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		log.Errorf("Get-video endpoint returned status %d", res.StatusCode)
		return
	}

	urlResp := &Url{}
	if derr := json.NewDecoder(res.Body).Decode(urlResp); derr != nil {
		log.Errorf("Error decoding Response Body: %s", derr)
		return
	}
	log.Infof("Video data received: %+v", urlResp)

	// Map message to map[string]interface{}
	var dataMap map[string]interface{}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		log.Errorf("Error marshaling Message: %s", err)
		return
	}
	err = json.Unmarshal(msgBytes, &dataMap)
	if err != nil {
		log.Errorf("Error unmarshaling Message to map: %s", err)
		return
	}

	dataMap["url"] = urlResp.Url

	dataBytes, err := json.Marshal(dataMap)
	if err != nil {
		log.Errorf("Error marshaling dataMap: %s", err)
		return
	}

	mongoRequest, err := http.NewRequest("POST", uploadMongoEndpoint, bytes.NewBuffer(dataBytes))
	if err != nil {
		log.Errorf("Error creating Mongo request at Endpoint: %s \n Error: %s", uploadMongoEndpoint, err)
		return
	}
	mongoRequest.Header.Add("Content-Type", "application/json")

	mongoClient := &http.Client{}
	mongoResponse, err := mongoClient.Do(mongoRequest) // CORRECTO: usar mongoRequest
	if err != nil {
		log.Errorf("Error posting to Mongo Endpoint: %s \n Error: %s", uploadMongoEndpoint, err)
		return
	}
	defer mongoResponse.Body.Close()

	if mongoResponse.StatusCode >= 300 {
		log.Errorf("Mongo endpoint returned status %d", mongoResponse.StatusCode)
		return
	}

	log.Infof("Video Posted successfully!")
}


func (v VideoUploader) UploadVideo(bucket string, fileName string, filePath string, endpoint string, message *message.Message){
	v.FileGetter.UploadFile(bucket,fileName,filePath)
	PostEndpoint(endpoint, fileName, message)
}