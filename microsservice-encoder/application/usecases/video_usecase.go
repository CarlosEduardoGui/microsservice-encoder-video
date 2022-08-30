package usecases

import (
	"context"
	"io/ioutil"
	"log"
	"microsservice-encoder/application/repositories"
	"microsservice-encoder/domain"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoUseCase struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoUseCase() VideoUseCase {
	return VideoUseCase{}
}

func (v *VideoUseCase) Download(bucketName string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}

	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	f, err := os.Create(os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	defer f.Close()

	log.Printf("%s", "The video has been stored"+v.Video.ID)

	return nil
}

func (v *VideoUseCase) Fragment() error {
	videoPath := os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID

	err := os.Mkdir(videoPath, os.ModePerm)
	if err != nil {
		return err
	}

	source := videoPath + ".mp4"
	target := videoPath + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}
