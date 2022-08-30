package usecases_test

import (
	"log"
	"microsservice-encoder/application/repositories"
	"microsservice-encoder/application/usecases"
	"microsservice-encoder/domain"
	"microsservice-encoder/framework/database"
	"testing"
	"time"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading environment: %v", err)
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDbTest()
	defer db.Close()
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "video_teste.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}

func TestVideoUseCaseDownload(t *testing.T) {
	video, repo := prepare()
	videoUseCase := usecases.NewVideoUseCase()
	videoUseCase.Video = video
	videoUseCase.VideoRepository = repo

	err := videoUseCase.Download("encodervideo")
	require.Nil(t, err)

	err = videoUseCase.Fragment()
	require.Nil(t, err)
}
