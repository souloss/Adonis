package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type devNull int

func (devNull) Write(p []byte) (int, error) {
	return len(p), nil
}

func CreateDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return cli
}

func PullImages(cli *client.Client, images []string) {
	ctx := context.Background()

	for _, image := range images {
		if IsImageExist(cli, image) {
			fmt.Printf("Image < %s > already exists\n", image)
			continue
		} else {
			fmt.Printf("Pull image < %s >\n", image)
		}

		reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
		io.Copy(devNull(0), reader)

		if err != nil {
			fmt.Printf("  Image < %s > failed to pull\n", image)
			panic(err)
		}

		fmt.Printf("  Image < %s > is sucessfully pulled\n", image)
	}
}

func SaveImages(cli *client.Client, images []string, fileDir string) {
	ctx := context.Background()

	for _, image := range images {
		imageFilePath := EnsureFilePathAvailable(filepath.Join(fileDir, GenerateImageFileName(image)))
		fmt.Printf("Save image < %s >\n", image)

		reader, err := cli.ImageSave(ctx, []string{image})
		if err != nil {
			fmt.Printf("  Image < %s > failed to save\n", image)
			panic(err)
		}

		SaveFileFromReadCloser(reader, imageFilePath)

		fmt.Printf("  Image < %s > is sucessfully saved as < %s >\n", image, imageFilePath)
	}
}

func SaveFileFromReadCloser(reader io.ReadCloser, imageFilePath string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	ioutil.WriteFile(imageFilePath, buf.Bytes(), 0666)
}

func GetImageLastName(image string) string {
	splited := strings.Split(image, "/")

	return splited[len(splited)-1]
}

func GenerateImageFileName(image string) string {
	return strings.Replace(GetImageLastName(image), ":", "_", 1) + ".tar"
}

func IsFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func GetTimetamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func EnsureFilePathAvailable(filePath string) string {
	if IsFileExists(filePath) {
		dir, name := filepath.Split(filePath)
		return filepath.Join(dir, GetTimetamp()+"_"+name)
	}
	return filePath
}

func TagImages(cli *client.Client, images []string, repoPath string) []string {
	ctx := context.Background()

	newImages := make([]string, 0)
	for _, image := range images {
		newImage := filepath.Join(repoPath, GetImageLastName(image))
		fmt.Printf("Tag image < %s >\n", image)

		if image == newImage {
			fmt.Printf("  The same tag < %s >, skip\n", image)
			continue
		}

		err := cli.ImageTag(ctx, image, newImage)
		if err != nil {
			fmt.Printf("  Image < %s > tag failed\n", image)
			panic(err)
		}

		fmt.Printf("  Image < %s > is sucessfully taged as < %s >\n", image, newImage)

		newImages = append(newImages, newImage)
	}

	return newImages
}

func IsImageExist(cli *client.Client, image string) bool {
	ctx := context.Background()

	imageSummaries, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, imageSummary := range imageSummaries {
		for _, repoTag := range imageSummary.RepoTags {
			if repoTag == image {
				return true
			}
		}
	}

	return false
}

func RemoveImages(cli *client.Client, images []string) {
	ctx := context.Background()

	for _, image := range images {
		fmt.Printf("Untag image < %s >\n", image)

		_, err := cli.ImageRemove(ctx, image, types.ImageRemoveOptions{})
		if err != nil {
			fmt.Printf("  Image < %s > untag failed\n", image)
			panic(err)
		}

		fmt.Printf("  Image < %s > is sucessfully untagged\n", image)
	}
}

func PushImages(cli *client.Client, images []string) {
	ctx := context.Background()

	for _, image := range images {
		fmt.Printf("Push image < %s >\n", image)

		_, err := cli.ImagePush(ctx, image, types.ImagePushOptions{})
		if err != nil {
			fmt.Printf("  Image < %s > push failed\n", image)
			panic(err)
		}

		fmt.Printf("  Image < %s > is sucessfully pushed\n", image)
	}
}