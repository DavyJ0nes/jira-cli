package config

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/andygrunwald/go-jira.v1"
)

func LoadConfigFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.jli")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetEnvPrefix("jli")
	err = viper.BindEnv("token")
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal Binding Envar: %s \n", err))
	}
}

type Jira struct {
	url      string
	username string
	token    string
	boardID  string
	Client   *jira.Client
}

func NewJIRAClient() *Jira {
	j := &Jira{}

	j.url = viper.GetString("endpoint")
	j.username = viper.GetString("username")
	j.token = viper.GetString("token")
	j.boardID = viper.GetString("boardid")

	return j
}

func (j *Jira) Connect() *Jira {
	tp := jira.BasicAuthTransport{
		Username: j.username,
		Password: j.token,
	}

	jiraClient, err := jira.NewClient(tp.Client(), j.url)
	if err != nil {
		panic(err)
	}

	j.Client = jiraClient

	return j
}

func (j *Jira) GetSprints() []jira.Sprint {
	sprints, _, err := j.Client.Board.GetAllSprints(j.boardID)
	if err != nil {
		panic(err)
	}

	return sprints
}
