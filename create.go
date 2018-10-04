// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
	"path"
	"fmt"
)

// checkSourceExistence return true if instance doesn't exist.
func (r *CreateReq) checkSourceExistence(ret *goforjj.PluginData) (p *GitlabPlugin, status bool) {
	log.Print("Checking Gitlab source code existence.")
	srcPath := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(srcPath, gitlabFile)); err == nil {
		log.Printf(ret.Errorf("Unable to create the gitlab source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update gitlab according to his configuration. %s.", srcPath, srcPath, err))
		return
	}

	p = newPlugin(srcPath)

	log.Printf(ret.StatusAdd("environment checked."))
	return p, true
}

//SaveMaintainOptions save maintain options (TODO ?)
func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

//createYamlData prepare data for yaml file (TODO)
func (gls *GitlabPlugin) createYamlData(req *CreateReq, ret *goforjj.PluginData) error{
	if gls.gitlabSource.Urls == nil{
		return fmt.Errorf("Internal Error. Urls was not set")
	}

	gls.gitlabDeploy.Projects = make(map[string]ProjectStruct)
	//gls.gitlabDeploy.Users = ...

	//Norepo
	gls.gitlabDeploy.NoProjects = (gls.app.ReposDisabled == "true")
	if gls.gitlabDeploy.NoProjects {
		log.Print("Repositories_disabled is true. forjj_gitlab won't manage repositories except the infra repository.")
	}

	//SetOrgHooks

	for name, project := range req.Objects.Repo{
		isInfra := (name == gls.app.ForjjInfra)
		if gls.gitlabDeploy.NoProjects && !isInfra {
			continue
		}
		if !project.isValid(name, ret){
			ret.StatusAdd("Warning!!! Invalid project '%s' requested. Ignored.")
			continue
		}
		gls.SetProject(&project, isInfra, project.Deployable == "true")
		//gls.SetHooks(...)

	}

	log.Printf("forjj-gitlab manages %d project(s).", len(gls.gitlabDeploy.Projects))

	//more todo...

	return nil
}
// DefineRepoUrls return default repo url for the repo name given
func (gls *GitlabPlugin) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl){
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: gls.gitlabSource.Urls["gitlab-ssh"] + gls.gitlabDeploy.Group + "/" + name + ".git",
		Url: gls.gitlabSource.Urls["gitlab-url"] + "/" + gls.gitlabDeploy.Group + "/" + name,
	}
	return
}