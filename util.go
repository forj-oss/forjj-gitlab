package main 

import(
	"os"
	"fmt"

	"github.com/forj-oss/goforjj"
	"golang.org/x/sys/unix"
)

//IsWritable Linux support only
func isWritable(path string) (res bool) {
	return unix.Access(path, unix.W_OK) == nil
}


// reqCheckPath return false if something is wrong.
// check path is writable.
func reqCheckPath(name, path string, ret *goforjj.PluginData) bool {

	if path == "" {
		ret.ErrorMessage = name + " is empty."
		return true
	}

	if _, err := os.Stat(path); err != nil {
		ret.ErrorMessage = fmt.Sprintf(name+" mounted '%s' is inexistent.", path)
		return true
	}

	if !isWritable(path) {
		ret.ErrorMessage = fmt.Sprintf(name+" mounted '%s' is NOT writable", path)
		return true
	}

	return false
}

//verifyReqFails verify check params (TODO)
func (gls *GitlabPlugin) verifyReqFails(ret *goforjj.PluginData, check map[string]bool) bool{
	if v, ok := check["source"]; ok && v {
		if reqCheckPath("source (forjj-source-mount)", gls.sourcePath, ret){
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
