package test

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigProfiles struct {
	Default   ConfigBody `yaml:"default"`
	Secondary ConfigBody `yaml:"secondary"`
}

type ConfigBody struct {
	AuthorARN    string   `yaml:"author-arn"`
	Repositories []string `yaml:"repositories"`
}

func CreateTestConfigFile() {
	p := ConfigProfiles{
		Default: ConfigBody{
			AuthorARN: "arn:aws:sts::596944721148:role/test-role",
			Repositories: []string{
				"first-repo",
				"second-repo",
			},
		},
		Secondary: ConfigBody{
			AuthorARN: "arn:aws:sts::596944721148:role/test-role",
			Repositories: []string{
				"second-repo",
			},
		},
	}
	yamlData, _ := yaml.Marshal(&p)
	os.WriteFile("cprl.yaml", yamlData, 0777)
}

func DeleteTestConfigFile() error {
	return os.Remove("cprl.yaml")
}
