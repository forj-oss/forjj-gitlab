// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"net/http"
	"log"
	//"os"
	"fmt"
	"path"
	//"io/ioutil"

	"github.com/forj-oss/goforjj"
	//"github.com/xanzy/go-gitlab"
	//"golang.org/x/sys/unix"
	//"gopkg.in/yaml.v2"

)

/*type ProjectStruct struct {
	Name 			string
	Flow 			string 				`yaml:",omitempty"`
	Description 	string 				`yaml:",omitempty"`
	Disabled 		bool				`yaml:",omitempty"`
	IssueTracker 	bool 				`yaml:"issue_tracker,omitempty"`
	Users 			map[string]string 	`yaml:",omitempty"`
	//Groups

	exist 			bool 				`yaml:",omitempty"`
	remotes 		map[string]goforjj.PluginRepoRemoteUrl
	branchConnect 	map[string]string 	
	//...

	//maintain
	Infra 			bool 				`yaml:",omitempty"`
	Role 			string 				`yaml:",omitempty"`
	IsDeployable 	bool

}*/

/*func (gls *GitlabPlugin) gitlab_set_url(server string) (err error) {
	//gl_url := ""
	if gls.gitlab_source.Urls == nil {
		gls.gitlab_source.Urls = make(map[string]string)
	}
	//...
	return //TODO
}*/

/*func (gls *GitlabPlugin) projects_exists(ret *goforjj.PluginData) (err error) {
	clientProjects := gls.Client.Projects // Projects of user
	client, _, err := gls.Client.Users.CurrentUser() // Get current user
	
	//loop
	for name, project_data := range gls.gitlabDeploy.Projects{

		URLEncPathProject := client.Username + "/" + name // UserName/ProjectName
		//Get X repo, if find --> err
		if found_project, _, e := clientProjects.GetProject(URLEncPathProject); e == nil{
			if err == nil && name == gls.app.ForjjInfra {
				err = fmt.Errorf("Infra projects '%s' already exist in gitlab server.", name)
			}
			project_data.exist = true
			if project_data.remotes == nil{
				project_data.remotes = make(map[string]goforjj.PluginRepoRemoteUrl)
				project_data.branchConnect = make(map[string]string)
			}
			project_data.remotes["origin"] = goforjj.PluginRepoRemoteUrl{
				Ssh: found_project.SSHURLToRepo,
				Url: found_project.HTTPURLToRepo,
			}
			project_data.branchConnect["master"] = "origin/master"
		}
		ret.Repos[name] = goforjj.PluginRepo{ //Project
			Name: 		project_data.Name,
			Exist: 		project_data.exist,
			Remotes: 		project_data.remotes,
			BranchConnect: 		project_data.branchConnect,
			//Owner: 		project_data.Owner,
		}

	}

	return
}*/

/*func (r *RepoInstanceStruct) IsValid(repo_name string, ret *goforjj.PluginData) (valid bool){
	if r.Name == "" {
		ret.Errorf("Invalid project '%s'. Name is empty.", repo_name)
		return
	}
	if r.Name != repo_name {
		ret.Errorf("Invalid project '%s'. Name must be equal to '%s'. But the project name is set to '%s'.", repo_name, repo_name, r.Name)
		return
	}
	valid = true
	return
}*/

/*func (gls *GitlabPlugin) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl){
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: gls.gitlab_source.Urls["gitlab-ssh"] + gls.gitlabDeploy.Group + "/" + name + ".git",
		Url: gls.gitlab_source.Urls["gitlab-url"] + "/" + gls.gitlabDeploy.Group + "/" + name,
	}
	return
}*/

/*func (r *ProjectStruct) set(project *RepoInstanceStruct) *ProjectStruct{
	if r == nil {
		r = new(ProjectStruct)
	}
	r.Name = project.Name
	return r
}*/

/*func (gls *GitlabPlugin) SetProject(project *RepoInstanceStruct, isInfra, isDeployable bool) { //SetRepo
	//upstream := gls.DefineRepoUrls(project.Name)

	/*owner := gls.gitlabDeploy.Group
	if isInfra {
		owner = gls.gitlabDeploy.ProdGroup
	}*//*

	//set it, found or not
	r := ProjectStruct{}
	r.set(project)
	gls.gitlabDeploy.Projects[project.Name] = r
}*/

/*func (gls *GitlabPlugin) create_yaml_data(req *CreateReq, ret *goforjj.PluginData) error{
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
}*/




// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int) {

	//Get instance name
	instance := req.Forj.ForjjInstanceName

	/*log.Printf(ret.StatusAdd("instance: " + instance))
	log.Println(ret.StatusAdd("token: " + req.Objects.App[instance].Token))
	log.Println(ret.StatusAdd("Source Mount: " + req.Forj.ForjjSourceMount))
	log.Println(ret.StatusAdd("Deploy Mount: " + req.Forj.ForjjDeployMount))
	log.Println(ret.StatusAdd("Deployto: " + req.Forj.ForjjDeploymentEnv))*/
	

	//init GitlabPlugin
	gls := GitlabPlugin{
		sourceMount: 	req.Forj.ForjjSourceMount,
		deployMount: 	req.Forj.ForjjDeployMount,
		instance: 		req.Forj.ForjjInstanceName,
		deployTo: 		req.Forj.ForjjDeploymentEnv,
		token: 			req.Objects.App[instance].Token,
		group:			req.Objects.App[instance].Group,
	}

	log.Printf("Checking parameters : %#v", gls)

	//Check token and source
	check := make(map[string]bool)
	check["token"] = true
	check["source"] = true

	//ensure source path is writeable && token isn't empty
	if gls.verify_req_fails(ret, check){
		return
	}

	//verify if instance existe
	if a, found := req.Objects.App[instance]; !found{
		ret.Errorf("Internal issue. Forjj has not given the Application information for '%s'. Aborted.")
		return
	} else {
		gls.app = &a
	}

	//Check Gitlab connection
	log.Println("Checking gitlab connection.")
	ret.StatusAdd("Connect to gitlab...")

	if git := gls.gitlab_connect("myserv", ret); git == nil{
		return
	}

	/*log.Println(ret.StatusAdd("Org: " + gls.app.ForjjOrganization))
	log.Println(ret.StatusAdd("Org: " + gls.app.Organization))*/

	// Init Group of project
	if !req.InitGroup(&gls){
		ret.Errorf("Internal Error. Unable to define the group.")
	}

	//Create yaml data for maintain function
	if err := gls.create_yaml_data(req, ret); err != nil {
		ret.Errorf("Unable to create. %s",err)
		return
	}

	//repos exist ?
	if err := gls.projects_exists(ret); err != nil {
		ret.Errorf("%s\nUnable to 'create' your forge when gitlab already has an infra project created. Clone it and use 'update' instead.", err)
		return 419
	}

	//CheckSourceExistence
	if err := gls.checkSourcesExistence("create"); err != nil{
		ret.Errorf("%s\nUnable to 'create' your forge", err)
		return
	}
	ret.StatusAdd("Environment checked. Ready to be created.")

	//Path in ctxt git
	gitFile := path.Join(gls.instance, gitlab_file)

	//Save gitlab source
	if _, err := gls.save_yaml(&gls.gitlab_source, gls.sourceFile/*ret, instance*/); err != nil {
		ret.Errorf("%s", err)
		return
	}

	log.Printf(ret.StatusAdd("Configuration saved in source project '%s' (%s).", gitFile, gls.sourceMount))

	//Save gitlab deploy
	if _, err := gls.save_yaml(&gls.gitlabDeploy, gls.deployFile); err != nil{
		ret.Errorf("%s", err)
		return
	}

	log.Printf(ret.StatusAdd("Configuration saved in deploy project '%s' (%s).", gitFile, path.Join(gls.deployMount, gls.deployTo)))

	//Build final post answer
	for k, v := range gls.gitlab_source.Urls{
		ret.Services.Urls[k] = v
	}

	//API by forjj
	ret.Services.Urls["api_url"] = gls.gitlab_source.Urls["gitlab-base-url"]

	ret.CommitMessage = fmt.Sprint("Gitlab configuration created.")

	ret.AddFile(goforjj.FilesSource, gitFile)
	ret.AddFile(goforjj.FilesDeploy, gitFile)

	
	log.Println(ret.StatusAdd("end"))
	return
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {
	/*var p *GitlabPlugin

	// This is where you shoud write your Update code. Following lines are typical code for a basic plugin.
	if pr, ok := req.check_source_existence(ret); !ok {
		return
	} else {
		p = pr
	}

	instance := req.Forj.ForjjInstanceName
	if !p.load_yaml(ret, instance) {
		return
	}

	if !p.update_from(req, ret) {
		return
	}

	// Example of the core task
	//if ! p.update_jenkins_sources(ret) {
	//    return
	//}

	if !p.save_yaml(ret, instance) {
		return
	}*/
	return
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoMaintain(r *http.Request, req *MaintainReq, ret *goforjj.PluginData) (httpCode int) {
	instance := req.Forj.ForjjInstanceName

	var gls GitlabPlugin
	if a, found := req.Objects.App[instance]; !found {
		ret.Errorf("Invalid request. Missing Objects/App/%s", instance)
		return
	}else{
		gls = GitlabPlugin{
			deployMount: 		req.Forj.ForjjDeployMount,
			workspace_mount: 	req.Forj.ForjjWorkspaceMount,
			token: 				a.Token,
			maintain_ctxt: 		true,
			force: 				req.Forj.Force == "true",
		}
	}
	
	check := make(map[string]bool)
	check["token"] = true
	check["workspace"] = true
	check["deploy"] = true

	if gls.verify_req_fails(ret, check) {
		return
	}

	confFile := path.Join(gls.deployMount, req.Forj.ForjjDeploymentEnv, instance, gitlab_file)

	//read yaml file
	if err := gls.load_yaml(confFile/*ret, instance*/); err != nil{
		ret.Errorf("%s"/*, err*/)
		return
	}
	
	if gls.gitlab_connect("", ret) == nil{
		return
	}

	if !gls.ensure_group_exists(ret){
		return
	}

	if !gls.IsNewForge(ret){
		return
	}

	if gls.gitlabDeploy.NoProjects{
		log.Printf(ret.StatusAdd("Projects maintained limited to your infra project"))
	}

	//loop verif
	for name, project_data := range gls.gitlabDeploy.Projects{
		if !project_data.Infra && gls.gitlabDeploy.NoProjects{
			log.Printf(ret.StatusAdd("Project ignored: %s", name))
			continue
		}
		if project_data.Role == "infra" && !project_data.IsDeployable{
			log.Printf(ret.StatusAdd("Project ignored: %s - Infra project owned by '%s'", name, gls.gitlabDeploy.ProdGroup))
			continue
		}
		/*if err := project_data.ensure_exists(&gls, ret); err != nil{
			return
		}*/
		//...
		log.Printf(ret.StatusAdd("Project maintained: %s", name))
	}
	

	return
}
