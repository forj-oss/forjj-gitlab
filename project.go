package main

import(
	"github.com/forj-oss/goforjj"
)

type ProjectStruct struct {
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

}

func (r *RepoInstanceStruct) IsValid(repo_name string, ret *goforjj.PluginData) (valid bool){
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

func (r *ProjectStruct) set(project *RepoInstanceStruct) *ProjectStruct{
	if r == nil {
		r = new(ProjectStruct)
	}
	r.Name = project.Name
	return r
}