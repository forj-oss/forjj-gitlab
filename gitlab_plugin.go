// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

type GitlabPlugin struct {
	yaml        YamlGitlab
	source_path string
}

const gitlab_file = "forjj-gitlab.yaml"

type YamlGitlab struct {
}

func new_plugin(src string) (p *GitlabPlugin) {
	p = new(GitlabPlugin)

	p.source_path = src
	return
}

func (p *GitlabPlugin) initialize_from(r *CreateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *GitlabPlugin) load_from(ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *GitlabPlugin) update_from(r *UpdateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (p *GitlabPlugin) save_yaml(ret *goforjj.PluginData, instance string) (status bool) {
	file := path.Join(instance, gitlab_file)

	d, err := yaml.Marshal(&p.yaml)
	if err != nil {
		ret.Errorf("Unable to encode forjj gitlab configuration data in yaml. %s", err)
		return
	}

	if err := ioutil.WriteFile(file, d, 0644); err != nil {
		ret.Errorf("Unable to save '%s'. %s", file, err)
		return
	}
	return true
}

func (p *GitlabPlugin) load_yaml(ret *goforjj.PluginData, instance string) (status bool) {
	file := path.Join(instance, gitlab_file)

	d, err := ioutil.ReadFile(file)
	if err != nil {
		ret.Errorf("Unable to load '%s'. %s", file, err)
		return
	}

	err = yaml.Unmarshal(d, &p.yaml)
	if err != nil {
		ret.Errorf("Unable to decode forjj gitlab data in yaml. %s", err)
		return
	}
	return true
}
