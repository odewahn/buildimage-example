package main

import (
	"archive/tar"
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/samalba/dockerclient"
)

func buildDockerFileTar(directoryPath string) (string, error) {
	dockerFileTar, err := os.Create(os.TempDir() + "/Dockerfile.tar")
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	directories, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	tarWriter := tar.NewWriter(dockerFileTar)
	defer tarWriter.Close()

	for _, v := range directories {
		fileContent, _ := ioutil.ReadFile(directoryPath + "/" + v.Name())
		header, err := tar.FileInfoHeader(v, "")
		if err != nil {
			log.Fatal(err)
			return "", err
		}
		tarWriter.WriteHeader(header)
		tarWriter.Write(fileContent)
	}

	tarWriter.Flush()
	return dockerFileTar.Name(), nil
}

func buildDockerImage() (string, error) {
	log.Println("Building docker test image")
	buildContextDir := "/Users/odewahn/Desktop/tensorflow-for-poets"
	dockerFileTar, err := buildDockerFileTar(buildContextDir)
	if err != nil {
		return "", err
	}
	dockerFile, err := os.Open(dockerFileTar)
	if err != nil {
		return "", err
	}
	defer dockerFile.Close()

	// Init the client
	docker, _ := dockerclient.NewDockerClient("unix:///var/run/docker.sock", nil)

	imageName := "salamander:1.7"

	labels := make(map[string]string)
	labels["metadata"] = "this is a f*cking miracle!"

	reader, err := docker.BuildImage(&dockerclient.BuildImage{
		Context:        dockerFile,
		RepoName:       imageName,
		SuppressOutput: false,
		Remove:         true,
		Labels:         labels,
	})
	defer reader.Close()

	if err != nil {
		log.Fatal("error building image", err)
		return "", err
	}
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() == true {
		log.Println(scanner.Text())
	}

	return imageName, nil
}

func main() {
	buildDockerImage()
}
