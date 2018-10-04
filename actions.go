// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"net/http"
	"log"
	"fmt"
	"path"

	"github.com/forj-oss/goforjj"

)

//DoCreate --> forjj create
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
		sourcePath: 	req.Forj.ForjjSourceMount, // ! \\
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
	if gls.verifyReqFails(ret, check){
		return
	}

	//verify if instance existe
	if a, found := req.Objects.App[instance]; !found{
		ret.Errorf("Internal issue. Forjj has not given the Application information for '%s'. Aborted.", instance)
		return
	} else {
		gls.app = &a
	}

	//Check Gitlab connection
	log.Println("Checking gitlab connection.")
	ret.StatusAdd("Connect to gitlab...")

	if git := gls.gitlabConnect("X"/*"myserv"*/, ret); git == nil{
		return
	}

	/*log.Println(ret.StatusAdd("Org: " + gls.app.ForjjOrganization))
	log.Println(ret.StatusAdd("Org: " + gls.app.Organization))*/

	// Init Group of project
	if !req.InitGroup(&gls){
		ret.Errorf("Internal Error. Unable to define the group.")
	}

	//Create yaml data for maintain function
	if err := gls.createYamlData(req, ret); err != nil {
		ret.Errorf("Unable to create. %s",err)
		return
	}

	//repos exist ?
	if err := gls.projectsExists(ret); err != nil {
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
	gitFile := path.Join(gls.instance, gitlabFile)

	//Save gitlab source
	if _, err := gls.saveYaml(&gls.gitlabSource, gls.sourceFile/*ret, instance*/); err != nil {
		ret.Errorf("%s", err)
		return
	}

	log.Printf(ret.StatusAdd("Configuration saved in source project '%s' (%s).", gitFile, gls.sourcePath)) // ! \\

	//Save gitlab deploy
	if _, err := gls.saveYaml(&gls.gitlabDeploy, gls.deployFile); err != nil{
		ret.Errorf("%s", err)
		return
	}

	log.Printf(ret.StatusAdd("Configuration saved in deploy project '%s' (%s).", gitFile, path.Join(gls.deployMount, gls.deployTo)))

	//Build final post answer
	for k, v := range gls.gitlabSource.Urls{
		ret.Services.Urls[k] = v
	}

	//API by forjj
	ret.Services.Urls["api_url"] = gls.gitlabSource.Urls["gitlab-base-url"]

	ret.CommitMessage = fmt.Sprint("Gitlab configuration created.")

	ret.AddFile(goforjj.FilesSource, gitFile)
	ret.AddFile(goforjj.FilesDeploy, gitFile)

	
	log.Println(ret.StatusAdd("end"))
	return
}

//DoUpdate --> forjj update
// Do updating plugin task
// req_data contains the request data posted by forjj. Structure generated from 'gitlab.yaml'.
// ret_data contains the response structure to return back to forjj.
//
// By default, if httpCode is not set (ie equal to 0), the function caller will set it to 422 in case of errors (error_message != "") or 200
func DoUpdate(r *http.Request, req *UpdateReq, ret *goforjj.PluginData) (httpCode int) {
	/*var p *GitlabPlugin

	// This is where you shoud write your Update code. Following lines are typical code for a basic plugin.
	if pr, ok := req.checkSourceExistence(ret); !ok {
		return
	} else {
		p = pr
	}

	instance := req.Forj.ForjjInstanceName
	if !p.loadYaml(ret, instance) {
		return
	}

	if !p.updateFrom(req, ret) {
		return
	}

	// Example of the core task
	//if ! p.update_jenkins_sources(ret) {
	//    return
	//}

	if !p.saveYaml(ret, instance) {
		return
	}*/
	return
}

//DoMaintain --> forjj maintain
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
	} else {
		gls = GitlabPlugin{
			deployMount: 		req.Forj.ForjjDeployMount,
			workspaceMount: 	req.Forj.ForjjWorkspaceMount,
			token: 				a.Token,
			maintainCtxt: 		true,
			force: 				req.Forj.Force == "true",
		}
	}
	
	check := make(map[string]bool)
	check["token"] = true
	check["workspace"] = true
	check["deploy"] = true

	if gls.verifyReqFails(ret, check) {
		return
	}

	confFile := path.Join(gls.deployMount, req.Forj.ForjjDeploymentEnv, instance, gitlabFile)

	//read yaml file
	if err := gls.loadYaml(confFile/*ret, instance*/); err != nil{
		ret.Errorf("%s"/*, err*/)
		return
	}
	
	if gls.gitlabConnect("", ret) == nil{
		return
	}

	if !gls.ensureGroupExists(ret){
		return
	}

	if !gls.IsNewForge(ret){
		return
	}

	if gls.gitlabDeploy.NoProjects{
		log.Printf(ret.StatusAdd("Projects maintained limited to your infra project"))
	}

	//loop verif
	for name, projectData := range gls.gitlabDeploy.Projects{
		if !projectData.Infra && gls.gitlabDeploy.NoProjects{
			log.Printf(ret.StatusAdd("Project ignored: %s", name))
			continue
		}
		if projectData.Role == "infra" && !projectData.IsDeployable{
			log.Printf(ret.StatusAdd("Project ignored: %s - Infra project owned by '%s'", name, gls.gitlabDeploy.ProdGroup))
			continue
		}
		/*if err := projectData.ensureExists(&gls, ret); err != nil{
			return
		}*/
		//...
		log.Printf(ret.StatusAdd("Project maintained: %s", name))
	}
	

	return
}
