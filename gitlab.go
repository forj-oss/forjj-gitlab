package main

import(
	"os"
	"fmt"
	"path"
	"log"

	"github.com/forj-oss/goforjj"
	"github.com/xanzy/go-gitlab"
)

//gitlabConnect connect user to gitlab (TODO)
func (gls *GitlabPlugin) gitlabConnect(server string, ret *goforjj.PluginData) *gitlab.Client {
	//
	gls.Client = gitlab.NewClient(nil, gls.token)

	//Set url
	if err := gls.gitlabSetUrl(server); err != nil{
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

//InitGroup TODO production group
func (req *CreateReq) InitGroup(gls *GitlabPlugin) (ret bool) {
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found{
		if group := app.Group; group == ""{
			return false
		}
		/*if prodGroup := app.ProductionGroup; prodGroup == ""{
			return false
		}*/
		gls.SetGroup(app)
		ret = true
	}
	return
}

//InitGroup Same for now, TODO production group
func (req *UpdateReq) InitGroup(gls *GitlabPlugin) (ret bool) {
	if app, found := req.Objects.App[req.Forj.ForjjInstanceName]; found{
		if group := app.Group; group == ""{
			return false
		}
		/*if prodGroup := app.ProductionGroup; prodGroup == ""{
			return false
		}*/
		gls.SetGroup(app)
		ret = true
	}
	return
}

//SetGroup TODO production group
func (gls *GitlabPlugin) SetGroup(fromApp AppInstanceStruct) {
	
	if group := fromApp.ProductionGroup; group == ""{			//
		//gls.gitlabDeploy.ProdGroup = fromApp.ForjjGroup		//
	} else {													//
		gls.gitlabDeploy.ProdGroup = group						//
	}															//

	gls.gitlabDeploy.Group = fromApp.Group
	//Get+set id
	groups, _, err := gls.Client.Groups.SearchGroup(fromApp.Group)
	if err != nil{
		log.Printf("Group not exists in Gitlab server. ID not set.") //Change to stop, without ID maintain not found
	}
	for _, element := range groups{
		if element.Name == fromApp.Group {
			gls.gitlabDeploy.GroupId = element.ID
		}
	}
	if gls.gitlabDeploy.GroupId == 0 {
		log.Printf("Group not exists in Gitlab server. ID not set.") //Change to stop, without ID maintain not found
	}

	//gls.gitlabDeploy.ProdGroup = fromApp.ProductionGroup
	gls.gitlabSource.ProdGroup = gls.gitlabDeploy.ProdGroup
}

//ensureGroupExists (TODO)
func (gls *GitlabPlugin) ensureGroupExists(ret *goforjj.PluginData) (s bool){
	//Ensure Group exist, todo: if not it is created.
	//Ensure user is owner (or same).
	return																   
}

//IsNewForge ...
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
			gls.newForge = (resp.StatusCode != 200)
		}
		return true
	}

	ret.Errorf("Unable to identify the infra repository. At least, one repo must be identified with "+"`%s` in %s. You can use Forjj update to fix this.","Infra: true", "github")
	return
}

//gitlabSetUrl (TODO: server)
func (gls *GitlabPlugin) gitlabSetUrl(server string) (err error) {
	glUrl := ""

	if gls.gitlabSource.Urls == nil {
		gls.gitlabSource.Urls = make(map[string]string)
	}

	if !gls.maintainCtxt {
		if server == "" { // || ? 
			gls.gitlabSource.Urls["gitlab-base-url"] = "https://gitlab.com/"
			gls.gitlabSource.Urls["gitlab-url"] = "https://gitlab.com"
			gls.gitlabSource.Urls["gitlab-ssh"] = "git@gitlab.com:"
		} else {
			//set from serveur // ! \\ TODO
			server = "gitlab.com"
			glUrl = "https://" + server + "/api/v4/"
			gls.gitlabSource.Urls["gitlab-url"] = "https://gitlab.com"
			gls.gitlabSource.Urls["gitlab-ssh"] = "git@gitlab.com:"
		}
	} else {
		//maintain context
		gls.gitlabSource.Urls = gls.gitlabDeploy.Urls
		glUrl = gls.gitlabSource.Urls["gitlab-base-url"]
	}

	if glUrl == ""{
		return
	}

	err = gls.Client.SetBaseURL(glUrl)
	
	if err != nil{
		return
	}

	return
}

//ensureExists (TODO UPDATE and group management)
func (r *ProjectStruct) ensureExists(gls *GitlabPlugin, ret *goforjj.PluginData) error {
	//test existence
	clientProjects := gls.Client.Projects
	//client, _, err := gls.Client.Users.CurrentUser() // Get current user
	URLEncPathProject := gls.gitlabDeploy.Group + "/" + r.Name // UserName/ProjectName or Group/ProjectName

	_, _, err := clientProjects.GetProject(URLEncPathProject)
	
	if err != nil {
		//if does'nt exists --> Create
		ABM := 0
		projectOptions := &gitlab.CreateProjectOptions{
			Name: &r.Name,
			NamespaceID: &gls.gitlabDeploy.GroupId,
			ApprovalsBeforeMerge: &ABM, //without: request error because is set to null (restriction SQL: not null)
		}
		_, _, e := gls.Client.Projects.CreateProject(projectOptions)
		if e != nil{
			ret.Errorf("Unable to create '%s'. %s.", r.Name, e)
			return e
		}
		log.Printf(ret.StatusAdd("Repo '%s': created", r.Name))

	} else {
		//Update TODO
		
	}
	
	//...

	return nil
}

//projectExists (TODO)
func (gls *GitlabPlugin) projectsExists(ret *goforjj.PluginData) (err error) {
	clientProjects := gls.Client.Projects // Projects of user
	//client, _, err := gls.Client.Users.CurrentUser() // Get current user
	
	//loop
	for name, projectData := range gls.gitlabDeploy.Projects{

		URLEncPathProject := gls.gitlabDeploy.Group + "/" + name // client.Username = UserName/ProjectName or gls... = Group/ProjectName
		//Get X repo, if find --> err
		if foundProject, _, e := clientProjects.GetProject(URLEncPathProject); e == nil{
			if err == nil && name == foundProject.Name {
				err = fmt.Errorf("Infra projects '%s' already exist in gitlab server.", name)
			}
			projectData.exist = true
			if projectData.remotes == nil{
				projectData.remotes = make(map[string]goforjj.PluginRepoRemoteUrl)
				projectData.branchConnect = make(map[string]string)
			}
			projectData.remotes["origin"] = goforjj.PluginRepoRemoteUrl{
				Ssh: foundProject.SSHURLToRepo,
				Url: foundProject.HTTPURLToRepo,
			}
			projectData.branchConnect["master"] = "origin/master"
		}
		ret.Repos[name] = goforjj.PluginRepo{ //Project
			Name: 			projectData.Name,
			Exist: 			projectData.exist,
			Remotes: 		projectData.remotes,
			BranchConnect: 		projectData.branchConnect,
			//Owner: 		projectData.Owner,
		}

	}

	return
}

//checkSourcesExistence (TODO UPDATE)
func (gls *GitlabPlugin) checkSourcesExistence(when string) (err error){
	log.Print("Checking Infrastructure code existence.")
	sourceProject := gls.sourcePath
	sourcePath := path.Join(sourceProject, gls.instance)
	gls.sourceFile = path.Join(sourcePath, gitlabFile)

	deployProject := path.Join(gls.deployMount, gls.deployTo)
	deployBase := path.Join(deployProject, gls.instance)

	gls.deployFile = path.Join(deployBase, gitlabFile)

	gls.gitFile = path.Join(gls.instance, gitlabFile)

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
			return fmt.Errorf("Unable to create the gitlab configuration which already exist.\nUse 'update' to update it "+"(or update %s), and 'maintain' to update your github service according to his configuration.",path.Join(gls.instance, gitlabFile))
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
