package services_test

import (
	"log"
	"microsservice-encoder/application/services"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}
}

func TestVideoServiceUpload(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("encodervideo")
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = "encodervideo"
	videoUpload.VideoPath = os.Getenv("localstoragePath" + "/" + video.ID)

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(10, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "uploaded completed successfully")
}
