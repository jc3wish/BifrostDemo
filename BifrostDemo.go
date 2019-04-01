package main

import _"github.com/jc3wish/BifrostDemo/src"
import (
	"github.com/jc3wish/Bifrost/manager/xgo"
	"os/exec"
	"os"
	"path/filepath"
	"net/http"
)


func controller_FirstCallback(w http.ResponseWriter,req *http.Request) bool {
	return true
}

func main()  {
	var execDir string

	var TemplatePath = func (fileName string) string{
		return execDir+fileName
	}
	execPath, _ := exec.LookPath(os.Args[0])
	execDir = filepath.Dir(execPath)+"/"
	xgo.SetFirstCallBack(controller_FirstCallback)
	xgo.AddStaticRoute("/plugin/demo/www/public/",TemplatePath("/"))
	xgo.Start("0.0.0.0:11521")
}
