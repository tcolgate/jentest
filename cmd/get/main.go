package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bndr/gojenkins"
)

func main() {
	jenkins := gojenkins.CreateJenkins(nil, os.Getenv("JENKINS_URL"), os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASSWORD"))

	job, err := jenkins.GetJob("deploy-staging", "hubot")
	if err != nil {
		log.Fatalf("couldn't get job, %v", err)
	}

	config, err := job.GetConfig()
	if err != nil {
		log.Fatalf("couldn't get config, %v", err)
	}

	fmt.Println(config)
}
