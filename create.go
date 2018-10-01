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

// return true if instance doesn't exist.
func (r *CreateReq) check_source_existence(ret *goforjj.PluginData) (p *GitlabPlugin, status bool) {
	log.Print("Checking Gitlab source code existence.")
	src_path := path.Join(r.Forj.ForjjSourceMount, r.Forj.ForjjInstanceName)
	if _, err := os.Stat(path.Join(src_path, gitlab_file)); err == nil {
		log.Printf(ret.Errorf("Unable to create the gitlab source code for instance name '%s' which already exist.\nUse update to update it (or update %s), and maintain to update gitlab according to his configuration. %s.", src_path, src_path, err))
		return
	}

	p = new_plugin(src_path)

	log.Printf(ret.StatusAdd("environment checked."))
	return p, true
}

func (r *CreateArgReq) SaveMaintainOptions(ret *goforjj.PluginData) {
	if ret.Options == nil {
		ret.Options = make(map[string]goforjj.PluginOption)
	}
}

func (gls *GitlabPlugin) create_yaml_data(req *CreateReq, ret *goforjj.PluginData) error{
	if gls.gitlab_source.Urls == nil{
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
		is_infra := (name == gls.app.ForjjInfra)
		if gls.gitlabDeploy.NoProjects && !is_infra {
			continue
		}
		if !project.IsValid(name, ret){
			ret.StatusAdd("Warning!!! Invalid project '%s' requested. Ignored.")
			continue
		}
		gls.SetProject(&project, is_infra, project.Deployable == "true")
		//gls.SetHooks(...)

	}

	log.Printf("forjj-gitlab manages %d project(s).", len(gls.gitlabDeploy.Projects))

	//more todo...

	return nil
}

func (gls *GitlabPlugin) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl){
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: gls.gitlab_source.Urls["gitlab-ssh"] + gls.gitlabDeploy.Group + "/" + name + ".git",
		Url: gls.gitlab_source.Urls["gitlab-url"] + "/" + gls.gitlabDeploy.Group + "/" + name,
	}
	return
}