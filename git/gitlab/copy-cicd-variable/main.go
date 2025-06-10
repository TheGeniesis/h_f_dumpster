package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/xanzy/go-gitlab"
)

type config struct {
	AccessToken string `env:"GITLAB_ACCESS_TOKEN,required"`
	ProjectId   int    `env:"GITLAB_PROJECT_ID,required"`
	Host        string `env:"GITLAB_HOST,required"`
}

func loadConfig() (c *config, err error) {
	c = &config{}
	err = env.Parse(c)

	return
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	config, err := loadConfig()
	check(err)

	oldVarName := flag.String("old", "", "Original variable name")
	newVarName := flag.String("new", "", "New variable name")

	flag.Parse()
	flag.Usage = func() {
		log.Println("Usage:")
		flag.PrintDefaults()
	}

	git, err := gitlab.NewClient(config.AccessToken, gitlab.WithBaseURL(config.Host))
	check(err)
	log.Printf("Connection to %s established\n", config.Host)

	log.Printf("Fetching variable '%s'...\n", *oldVarName)
	oldGitVar, _, err := git.ProjectVariables.GetVariable(config.ProjectId, *oldVarName, &gitlab.GetProjectVariableOptions{})
	check(err)

	log.Printf("Creating new variable '%s' with type '%s'...\n", *newVarName, oldGitVar.VariableType)
	_, _, err = git.ProjectVariables.CreateVariable(config.ProjectId, &gitlab.CreateProjectVariableOptions{
		Key:          gitlab.String(*newVarName),
		Value:        gitlab.String(oldGitVar.Value),
		Protected:        gitlab.Bool(oldGitVar.Protected),
		Masked:        gitlab.Bool(oldGitVar.Masked),
		EnvironmentScope:        gitlab.String(oldGitVar.EnvironmentScope),
		VariableType: gitlab.VariableType(oldGitVar.VariableType),
	})
	check(err)

	log.Println("New variable created successfully.")
}
