// This file has been created by "go generate" as initial code. go generate will never update it, EXCEPT if you remove it.

// So, update it for your need.
package main

import (
	"github.com/forj-oss/goforjj"
	"log"
	"os"
)

// Return ok if the jenkins instance exist
func (r *MaintainReq) checkSourceExistence(ret *goforjj.PluginData) (status bool) {
	log.Print("Checking gitlab source code path existence.")

	if _, err := os.Stat(r.Forj.ForjjSourceMount); err == nil {
		ret.Errorf("Unable to maintain gitlab instances. '%s' is inexistent or innacessible.\n",
			r.Forj.ForjjSourceMount)
		return
	}
	ret.StatusAdd("environment checked.")
	return true
}

func (r *MaintainReq) instantiate(ret *goforjj.PluginData) (status bool) {

	return true
}
