package main

import _"github.com/jc3wish/BifrostDemo/src"
import (
	"github.com/jc3wish/Bifrost/manager/xgo"
	"github.com/jc3wish/BifrostDemo/config"
	"os/exec"
	"os"
	"path/filepath"
	"net/http"
	"flag"
	"html/template"
	"encoding/json"
	"log"
)


var sessionMgr *xgo.SessionMgr = nil //session管理器
var execDir string

func main()  {

	BifrostConfigFile := flag.String("config", "", "Bifrost config file path")
	Port := flag.String("port", "11521", "-port 11521")
	flag.Parse()

	if *BifrostConfigFile == ""{
		log.Println("err config == ''")
		os.Exit(1)
	}
	sessionMgr = xgo.NewSessionMgr("xgo_demo_cookie", 3600)

	config.LoadConf(*BifrostConfigFile)
	execPath, _ := exec.LookPath(os.Args[0])
	execDir = filepath.Dir(execPath)+"/"
	xgo.SetFirstCallBack(controller_FirstCallback)
	xgo.AddStaticRoute("/plugin/demo/www/public/",TemplatePath("/"))
	xgo.AddRoute("/",index_controller)
	xgo.AddRoute("/login",user_login)
	xgo.AddRoute("/dologin",user_do_login)
	xgo.Start("0.0.0.0:"+*Port)
}

func index_controller(w http.ResponseWriter,req *http.Request) {
	http.Redirect(w, req, "/bifrost/demo/index", http.StatusFound)
}

func TemplatePath(fileName string) string{
	return execDir+fileName
}

func controller_FirstCallback(w http.ResponseWriter,req *http.Request) bool {
	var sessionID= sessionMgr.CheckCookieValid(w, req)

	if sessionID != "" {
		if _,ok:=sessionMgr.GetSessionVal(sessionID,"UserName");ok{
			return true
		}else{
			goto toLogin
		}
	}else{
		goto toLogin
	}

toLogin:
	if req.RequestURI != "/login" && req.RequestURI != "/dologin" && req.RequestURI != "/logout"{
		http.Redirect(w, req, "/login", http.StatusFound)
		return false
	}
	return true
}


type resultStruct struct {
	Status bool `json:"status"`
	Msg string `json:"msg"`
}
func returnResult(r bool,msg string)[]byte{
	b,_:=json.Marshal(resultStruct{Status:r,Msg:msg})
	return  b
}

func user_login(w http.ResponseWriter,req *http.Request){
	type TemplateHeader struct {
		Title string
	}

	data := TemplateHeader{Title:"Login - BifrostDemo"}
	t, _ := template.ParseFiles(TemplatePath("plugin/demo/www/template/login.html"))
	t.Execute(w, data)
}

func user_do_login(w http.ResponseWriter,req *http.Request){
	req.ParseForm()
	var sessionID = sessionMgr.StartSession(w, req)
	UserName := req.Form.Get("user_name")
	UserPwd := req.Form.Get("password")

	if UserName == ""{
		w.Write(returnResult(false," user no exsit"))
		return
	}
	pwd := config.GetConfigVal("user",UserName)
	if pwd == UserPwd{
		sessionMgr.SetSessionVal(sessionID, "UserName", UserName)
		w.Write(returnResult(true," success"))
		return
	}
	w.Write(returnResult(false," user error"))
	return
}