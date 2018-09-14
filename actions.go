// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"net/http"
	"log"
	"os"
	"fmt"

	"github.com/forj-oss/goforjj"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sys/unix"
)

type GitlabStruct struct{
	sourceMount 	string
	deployMount 	string
	instance 		string
	deployTo 		string
	token 			string
	group			string

	app *AppInstanceStruct  			//permet l'accès au forjfile
	Client *gitlab.Client 				//client gitlab cf api gitlab ~ 
	gitlab_source GitlabSourceStruct 	//urls...
	gitlabDeploy GitlabDeployStruct     //

}

type GitlabSourceStruct struct{
	goforjj.PluginService `,inline` //base url
	ProdGroup string `yaml:"production-group-name, omitempty"`
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

type ProjectStruct struct {
	Name 			string
	Flow 			string
	Description 	string
	Disabled 		bool
	IssueTracker 	bool 
	Users 			map[string]string
	//...
}

// Linux support only
func IsWritable(path string) (res bool) {
	return unix.Access(path, unix.W_OK) == nil
}

// check path is writable.
// return false if something is wrong.
func reqCheckPath(name, path string, ret *goforjj.PluginData) bool {

	if path == "" {
		ret.ErrorMessage = name + " is empty."
		return true
	}

	if _, err := os.Stat(path); err != nil {
		ret.ErrorMessage = fmt.Sprintf(name+" mounted '%s' is inexistent.", path)
		return true
	}

	if !IsWritable(path) {
		ret.ErrorMessage = fmt.Sprintf(name+" mounted '%s' is NOT writable", path)
		return true
	}

	return false
}

func (gls *GitlabStruct) verify_req_fails(ret *goforjj.PluginData, check map[string]bool) bool{
	if v, ok := check["source"]; ok && v {
		if reqCheckPath("source (forjj-source-mount)", gls.sourceMount, ret){
			return true
		}
	}

	if v, ok := check["token"]; ok && v {
		if gls.token == ""{
			ret.ErrorMessage = fmt.Sprint("gitlab token is empty - Required")
			return true
		}
	}

	return false
}

func (gls *GitlabStruct) gitlab_set_url(server string) (err error) {
	//gl_url := ""
	if gls.gitlab_source.Urls == nil {
		gls.gitlab_source.Urls = make(map[string]string)
	}
	//...
	return //TODO
}

func (gls *GitlabStruct) gitlab_connect(server string, ret *goforjj.PluginData) *gitlab.Client {
	//
	gls.Client = gitlab.NewClient(nil, gls.token)

	//Set url
	if err := gls.gitlab_set_url(server); err != nil{
		ret.Errorf("Invalid url. %s", err)
		return nil
	}

	//check connected
	_, _, err := gls.Client.Users.CurrentUser()
	if err != nil {
		ret.Errorf("Unable to get the owner of the token given.", err)
		return nil
	} else {
		ret.StatusAdd("Connection successful.")
		/*g.user =  */
	}


	return gls.Client
}

func (gls *GitlabStruct) projects_exists (ret *goforjj.PluginData) (err error) {
	client := gls.Client.Projects
	
	//loop
	for name, project_data := range gls.gitlabDeploy.Projects{
		log.Print(client)
		log.Print(name)
		log.Print(project_data)
	}

	return
}

func (r *RepoInstanceStruct) IsValid (repo_name string, ret *goforjj.PluginData) (valid bool){ // /!\ change struct (project)
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
}

func (gls *GitlabStruct) DefineRepoUrls(name string) (upstream goforjj.PluginRepoRemoteUrl){
	upstream = goforjj.PluginRepoRemoteUrl{
		Ssh: gls.gitlab_source.Urls["gitlab-ssh"] + gls.gitlabDeploy.Group + "/" + name + ".git",
		Url: gls.gitlab_source.Urls["gitlab-url"] + "/" + gls.gitlabDeploy.Group + "/" + name,
	}
	return
}

func (r *ProjectStruct) set( project *RepoInstanceStruct) *ProjectStruct{
	if r == nil {
		r = new(ProjectStruct)
	}
	r.Name = project.Name
	return r
}

func (gls *GitlabStruct) SetProject(project *RepoInstanceStruct, isInfra, isDeployable bool) { //SetRepo
	//upstream := gls.DefineRepoUrls(project.Name)

	/*owner := gls.gitlabDeploy.Group
	if isInfra {
		owner = gls.gitlabDeploy.ProdGroup
	}*/

	//set it, found or not
	r := ProjectStruct{}
	r.set(project)
	gls.gitlabDeploy.Projects[project.Name] = r
}

func (gls *GitlabStruct) create_yaml_data (req *CreateReq, ret *goforjj.PluginData) error{
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

	//more ...

	return nil
}

func (req *CreateReq) InitGroup (gls *GitlabStruct) (ret bool) {
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found{
		gls.SetGroup(app)
		ret = true
	}
	return
}

func (gls *GitlabStruct) SetGroup (fromApp AppInstanceStruct) {
	if group := fromApp.Group; group == ""{
		gls.gitlabDeploy.Group =fromApp.ForjjGroup
	} else {
		gls.gitlabDeploy.Group = group
	}
	if group := fromApp.ProductionGroup; group == ""{
		gls.gitlabDeploy.ProdGroup = fromApp.ForjjGroup
	} else {
		gls.gitlabDeploy.ProdGroup = group
	}
	gls.gitlab_source.ProdGroup = gls.gitlabDeploy.ProdGroup
}

// Do creating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoCreate(r *http.Request, req *CreateReq, ret *goforjj.PluginData) (httpCode int) {

	//log.Printf(ret.StatusAdd(""))

	//Get instance name
	instance := req.Forj.ForjjInstanceName

	log.Printf(ret.StatusAdd("instance: " + instance))
	log.Println(ret.StatusAdd("token: " + req.Objects.App[instance].Token))
	log.Println(ret.StatusAdd("Source Mount: " + req.Forj.ForjjSourceMount))
	log.Println(ret.StatusAdd("Deploy Mount: " + req.Forj.ForjjDeployMount))
	log.Println(ret.StatusAdd("Deployto: " + req.Forj.ForjjDeploymentEnv))
	

	//init GitlabStruct
	gls := GitlabStruct{
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

	//ensure source path is writeable, token isn't empty
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

	// InitGroup 
	if !req.InitGroup(&gls){
		ret.Errorf("Internal Error. Unable to define the group.")
	}

	//Create yaml data
	if err := gls.create_yaml_data(req, ret); err != nil {
		ret.Errorf("Unable to create. %s",err)
		return
	}

	//repos exist ?
	if err := gls.projects_exists(ret); err != nil {
		ret.Errorf("%s\nUnable to 'create' your forge when gitlab already has an infra project created. Clone it and use 'update' instead.", err)
		return 419
	}

	//finalisation


	
	log.Println(ret.StatusAdd("end"))
	return
}

// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {
	var p *GitlabPlugin

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
	}
	return
}

// Do maintaining plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoMaintain(r *http.Request, req *MaintainReq, ret *goforjj.PluginData) (httpCode int) {
	// This is where you shoud write your Maintain code. Following lines are typical code for a basic plugin.
	if !req.check_source_existence(ret) {
		return
	}

	if !req.instantiate(ret) {
		return
	}
	return
}
