package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"text/template"

	"github.com/bndr/gojenkins"
	"github.com/pkg/errors"
)

type xmlString string

func (xstr xmlString) String() string {
	str := bytes.NewBufferString("")
	err := xml.EscapeText(str, []byte(xstr))
	if err != nil {
		panic(errors.Wrap(err, "failed converting string to xml"))
	}

	return str.String()
}

type GitHubSpec struct {
	ProjectURL string
}

type GitSpec struct {
	URL        string
	BranchSpec string
	Credential string
}

type JenkinsSpec struct {
	PipelinePath string
}

type JobSpec struct {
	Name    string
	GitHub  GitHubSpec
	Git     GitSpec
	Jenkins JenkinsSpec
}

func main() {
	jenkins := gojenkins.CreateJenkins(nil, "https://jenkins-aws-eu-stg.qutics.com/", "tcolgate", "ab17fda404f789c813758f3cd3c75d3a")

	jobspec := JobSpec{
		Name: "hugot",
		GitHub: GitHubSpec{
			ProjectURL: "https://github.com/qubitdigital/hugot/",
		},
		Git: GitSpec{
			URL:        "git@github.com:qubitdigital/hugot.git",
			Credential: "github-ssh",
		},
		Jenkins: JenkinsSpec{
			PipelinePath: "Jenkinsfile",
		},
	}

	buildTmpl := template.Must(template.ParseFiles("job-build.xml.tmpl"))

	conf := bytes.NewBufferString("")
	err := buildTmpl.Execute(conf, &jobspec)
	if err != nil {
		log.Fatalf("failed executing template, %v", err)
	}

	newfld, err := jenkins.CreateFolder(jobspec.Name)
	if err != nil {
		log.Fatalf("couldn't create folder, %v", err)
	}

	fmt.Printf("Trying to create job:\n%s", conf.String())
	newjob, err := jenkins.CreateJobInFolder(conf.String(), "build", newfld.GetName())
	if err != nil {
		log.Fatalf("couldn't create job, %v", err)
	}

	newconf, err := newjob.GetConfig()
	if err != nil {
		log.Fatalf("couldn't create job, %v", err)
	}
	fmt.Printf("New job:\n%s\n", newconf)
}
