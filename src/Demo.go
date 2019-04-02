package src

import (
"net/http"
"html/template"
"github.com/jc3wish/Bifrost/manager/xgo"
"os/exec"
"os"
"path/filepath"
"encoding/json"
	"strings"
	"strconv"
)

var execDir string

type TemplateHeader struct {
	Title string
}

type resultDataStruct struct {
	Status bool `json:"status"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

var sessionMgr *xgo.SessionMgr = nil //session管理器

func returnDataResult(r bool,msg string,data interface{})[]byte{
	b,_:=json.Marshal(resultDataStruct{Status:r,Msg:msg,Data:data})
	return  b
}

func GetFormInt(req *http.Request,key string) int{
	v := strings.Trim(req.Form.Get(key),"")
	intv,err:=strconv.Atoi(v)
	if err != nil{
		return 0
	}else{
		return intv
	}
}

func init()  {
	xgo.AddRoute("/bifrost/demo/index",index_controller)
	xgo.AddRoute("/bifrost/demo/execsql",execsql_controller)
	xgo.AddRoute("/bifrost/demo/redis",redisget_controller)
	execPath, _ := exec.LookPath(os.Args[0])
	execDir = filepath.Dir(execPath)+"/"
}

func TemplatePath(fileName string) string{
	return execDir+fileName
}


func index_controller(w http.ResponseWriter,req *http.Request){
	req.ParseForm()
	type DemoIndex struct{
		TemplateHeader
		SqlList []string
	}
	data := DemoIndex{SqlList:sqlList}
	data.Title = "Demo index"
	t, _ := template.ParseFiles(TemplatePath("plugin/demo/www/template/index.html"))
	t.Execute(w, data)
}

func execsql_controller(w http.ResponseWriter,req *http.Request){
	req.ParseForm()
	index := GetFormInt(req,"index")
	if index >= len(sqlList){
		w.Write(returnDataResult(false,"param error",0))
		return
	}
	db := NewMySQLConn(mysqlUri)
	defer db.Close()
	err := db.ExecSQL(sqlList[index])
	if err != nil{
		w.Write(returnDataResult(false,err.Error(),0))
		return
	}else{
		w.Write(returnDataResult(true,"",0))
		return
	}
}

func redisget_controller(w http.ResponseWriter,req *http.Request){
	req.ParseForm()
	redisClient,err := NewRedisClient(redisUri)
	if err != nil{
		w.Write(returnDataResult(false,err.Error(),0))
		return
	}
	defer redisClient.Close()
	Type := req.Form.Get("type")
	Key := req.Form.Get("key")
	var data string
	switch Type {
	case	"list":
			data = redisClient.LPop(Key)
			break
	default:
		data = redisClient.Get(Key)
		break

	}
	w.Write(returnDataResult(true,"success",data))
	return
}