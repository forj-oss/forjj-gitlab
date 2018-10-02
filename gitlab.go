package main

import(
	"os"
	"fmt"
	"path"
	"log"

	"github.com/forj-oss/goforjj"
	"github.com/xanzy/go-gitlab"
)

func (gls *GitlabPlugin) gitlab_connect(server string, ret *goforjj.PluginData) *gitlab.Client {
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
		//g.user = ...
	}
	return gls.Client
}

func (req *CreateReq) InitGroup(gls *GitlabPlugin) (ret bool) {
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found{
		gls.SetGroup(app)
		ret = true
	}
	return
}

func (gls *GitlabPlugin) SetGroup(fromApp AppInstanceStruct) {
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

func (gls *GitlabPlugin) ensure_group_exists(ret *goforjj.PluginData) (s bool){
	//TODO
	return																   
}

func (gls *GitlabPlugin) IsNewForge(ret *goforjj.PluginData) (_ bool){

	ClientProjects := gls.Client.Projects

	for name, project  := range gls.gitlabDeploy.Projects{
		if !project.Infra{
			continue
		}
		client, _, _ := gls.Client.Users.CurrentUser() //!\\ manage err
		URLEncPathProject := client.Username + "/" + name
		if _, resp, e := ClientProjects.GetProject(URLEncPathProject); e!= nil && resp == nil {
			ret.Errorf("Unable to identify the infra project. Unknown issue: %s",e)
			return
		} else {
			gls.new_forge = (resp.StatusCode != 200)
		}
		return true
	}

	ret.Errorf("Unable to identify the infra repository. At least, one repo must be identified with "+"`%s` in %s. You can use Forjj update to fix this.","Infra: true", "github")
	return
}

func (gls *GitlabPlugin) gitlab_set_url(server string) (err error) {
	//gl_url := ""
	if gls.gitlab_source.Urls == nil {
		gls.gitlab_source.Urls = make(map[string]string)
	}
	//...
	return //TODO
}

func (r *ProjectStruct) ensure_exists(gls *GitlabPlugin, ret *goforjj.PluginData) /*error*/ {
	//test existence
	//TODO
}

func (gls *GitlabPlugin) projects_exists(ret *goforjj.PluginData) (err error) {
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
}

func (gls *GitlabPlugin) checkSourcesExistence(when string) (err error){
	log.Print("Checking Infrastructure code existence.")
	sourceProject := gls.source_path // ! \\
	sourcePath := path.Join(sourceProject, gls.instance)
	gls.sourceFile = path.Join(sourcePath, gitlab_file)

	deployProject := path.Join(gls.deployMount, gls.deployTo)
	deployBase := path.Join(deployProject, gls.instance)

	gls.deployFile = path.Join(deployBase, gitlab_file)

	gls.gitFile = path.Join(gls.instance, gitlab_file)

	switch when {
	case "create":
		if _, err := os.Stat(sourcePath); err != nil{
			if err = os.MkdirAll(sourcePath, 0755); err != nil{
				return fmt.Errorf("Unable to create '%s'. %s", sourcePath, err)
			}
		}

		if _, err := os.Stat(deployProject); err != nil{
			return fmt.Errorf("Unable to create '%s'. Forjj must create it. %s", deployProject, err)
		}

		if _, err := os.Stat(gls.sourceFile); err == nil{
			return fmt.Errorf("Unable to create the gitlab configuration which already exist.\nUse 'update' to update it "+"(or update %s), and 'maintain' to update your github service according to his configuration.",path.Join(gls.instance, gitlab_file))
		}

		if _, err := os.Stat(deployBase); err != nil{
			if err = os.Mkdir(deployBase, 0755); err != nil{
				return fmt.Errorf("Unable to create '%s'. %s", deployBase, err)
			}
		}
		return
	
	case "update":
		log.Printf("TODO UPDATE")
	}
	return
}