package main

import(
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