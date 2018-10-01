// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"io/ioutil"
	//"path"
	"fmt"
	
	"github.com/forj-oss/goforjj"
	"github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v2"
)

type GitlabPlugin struct {
	yaml        YamlGitlab
	source_path string

	sourceMount 	string
	deployMount 	string
	instance 		string
	deployTo 		string
	token 			string
	group			string

	app 			*AppInstanceStruct  	//forjfile access
	Client 			*gitlab.Client 			//gitlab client ~ api gitlab
	gitlab_source 	GitlabSourceStruct 		//urls...
	gitlabDeploy 	GitlabDeployStruct     	//

	gitFile 		string
	deployFile 		string
	sourceFile  	string

	//maintain
	workspace_mount string
	maintain_ctxt 	bool
	force 			bool

	new_forge 		bool
}

type GitlabSourceStruct struct{
	goforjj.PluginService `,inline` //base url
	ProdGroup string `yaml:"production-group-name"`//`yaml:"production-group-name, omitempty"`
}

type GitlabDeployStruct struct{
	goforjj.PluginService 	`yaml:",inline"` //urls
	Projects 				map[string]ProjectStruct // projects managed in gitlab
	NoProjects				bool `yaml:",omitempty"`
	ProdGroup 				string
	Group 					string
	GroupDisplayName 		string
	//...

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

func (p *GitlabPlugin) save_yaml(/*ret *goforjj.PluginData, instance string*/in interface{}, file string) (Updated bool, _ error) {
	/*file := path.Join(instance, gitlab_file)

	d, err := yaml.Marshal(&p.yaml)
	if err != nil {
		ret.Errorf("Unable to encode forjj gitlab configuration data in yaml. %s", err)
		return
	}

	if err := ioutil.WriteFile(file, d, 0644); err != nil {
		ret.Errorf("Unable to save '%s'. %s", file, err)
		return
	}
	return true*/
	d, err := yaml.Marshal(in)
	if err != nil {
		return false, fmt.Errorf("Unable to encode gitlab data in yaml. %s", err)
	}

	if d_before, err := ioutil.ReadFile(file); err != nil{
		Updated = true
	} else {
		Updated = (string(d) != string(d_before))
	}
	
	if !Updated {
		return
	}
	if err = ioutil.WriteFile(file, d, 0644); err != nil {
		return false, fmt.Errorf("Unable to save '%s'. %s", file, err)
	}
	return
}

func (gls *GitlabPlugin) load_yaml(/*ret *goforjj.PluginData, instance string*/file string) error {
	/*file := path.Join(instance, gitlab_file)

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
	return true*/
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
