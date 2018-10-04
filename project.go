package main

import(
	"github.com/forj-oss/goforjj"
)

//ProjectStruct (TODO)
type ProjectStruct struct {
	Name 			string
	Flow 			string 				`yaml:",omitempty"`
	Description		string 				`yaml:",omitempty"`
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

	Owner 			string 				`yaml:",omitempty"`

}

//isValid verify name project
func (r *RepoInstanceStruct) isValid(repo_name string, ret *goforjj.PluginData) (valid bool){
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

//set (TODO)
func (r *ProjectStruct) set(project *RepoInstanceStruct, remotes map[string]goforjj.PluginRepoRemoteUrl, branchConnect map[string]string, is_infra, IsDeployable bool, owner string) *ProjectStruct{
	if r == nil {
		r = new(ProjectStruct)
	}
	r.Name = project.Name
	r.Description = project.Title

	//issueTracker

	r.Flow = project.Flow
	r.Infra = is_infra

	//r.addUsers(project.Users)
	//Groups

	r.remotes = remotes
	r.branchConnect = branchConnect
	
	//WebHooks

	r.Role = project.Role
	r.Owner = owner
	r.IsDeployable = IsDeployable

	return r
}

//addUsers (TODO)
func (r *ProjectStruct) addUsers(users string) {
	//
}