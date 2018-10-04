// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"io/ioutil"
	"fmt"
	
	"github.com/forj-oss/goforjj"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type GitlabPlugin struct{
	yaml			YamlGitlab
	sourcePath		string

	deployMount		string
	instance		string
	deployTo		string
	token			string
	group			string

	app				*AppInstanceStruct		//forjfile access
	Client			*gitlab.Client			//gitlab client ~ api gitlab
	gitlabSource	GitlabSourceStruct		//urls...
	gitlabDeploy	GitlabDeployStruct		//

	gitFile			string
	deployFile		string
	sourceFile		string

	//maintain
	workspaceMount	string
	maintainCtxt	bool
	force			bool

	newForge		bool
}

type GitlabSourceStruct struct{
	goforjj.PluginService	`,inline`							//base url
	ProdGroup string		`yaml:"production-group-name"`		//`yaml:"production-group-name, omitempty"`
}

type GitlabDeployStruct struct{
	goforjj.PluginService								`yaml:",inline"`			//urls
	Projects				map[string]ProjectStruct								// projects managed in gitlab
	NoProjects				bool						`yaml:",omitempty"`
	ProdGroup				string
	Group					string
	GroupDisplayName		string
	//...
}

const gitlabFile = "forjj-gitlab.yaml"

type YamlGitlab struct {
}

func newPlugin(src string) (gls *GitlabPlugin) {
	gls = new(GitlabPlugin)

	gls.sourcePath = src
	return
}

func (gls *GitlabPlugin) initializeFrom(r *CreateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (gls *GitlabPlugin) loadFrom(ret *goforjj.PluginData) (status bool) {
	return true
}

func (gls *GitlabPlugin) updateFrom(r *UpdateReq, ret *goforjj.PluginData) (status bool) {
	return true
}

func (gls *GitlabPlugin) saveYaml(in interface{}, file string) (Updated bool, _ error) {
	d, err := yaml.Marshal(in)
	if err != nil {
		return false, fmt.Errorf("Unable to encode gitlab data in yaml. %s", err)
	}

	if dBefore, err := ioutil.ReadFile(file); err != nil{
		Updated = true
	} else {
		Updated = (string(d) != string(dBefore))
	}
	
	if !Updated {
		return
	}
	if err = ioutil.WriteFile(file, d, 0644); err != nil {
		return false, fmt.Errorf("Unable to save '%s'. %s", file, err)
	}
	return
}

func (gls *GitlabPlugin) loadYaml(file string) error {
	d, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Unable to load '%s'. %s", file, err)
	}

	err = yaml.Unmarshal(d, &gls.gitlabDeploy)

	if err != nil {
		return fmt.Errorf("Unable to decode gitlab data in yaml. %s", err)
	}

	return nil
}
