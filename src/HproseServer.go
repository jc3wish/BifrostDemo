package src

import (
	"github.com/hprose/hprose-golang/rpc"
	"log"
	"encoding/json"
)

func init()  {
	go HproseServer(hproseUri)
}

func sendToWS(SchemaName,TableName string,data interface{}){
	switch data.(type) {
	case string:
		sendMsgToWsClient("SchemaName:"+SchemaName+ "TableName:"+TableName+ "data: "+data.(string))
		break
	default:
		b,_ := json.Marshal(data)
		sendMsgToWsClient("SchemaName:"+SchemaName+ "TableName:"+TableName+ "data: "+string(b))
		break
	}
}

func Check(context *rpc.HTTPContext) (e error) {
	log.Println("Check success")
	return nil
}

func Insert(SchemaName string,TableName string, data map[string]interface{}) (e error) {
	log.Println("Insert","SchemaName:",SchemaName,"TableName:",TableName,"data:",data)
	sendToWS(SchemaName,TableName,data)
	return nil
}

func Update(SchemaName string,TableName string, data []map[string]interface{}) (e error){
	for k,v := range data[0]{
		log.Println(k,"before:",v, "after:",data[1][k])
	}
	sendToWS(SchemaName,TableName,data)
	return nil
}

func Delete(SchemaName string,TableName string,data map[string]interface{}) (e error) {
	log.Println("Delete","SchemaName:",SchemaName,"TableName:",TableName,"data:",data)
	sendToWS(SchemaName,TableName,data)
	return nil
}

func Query(SchemaName string,TableName string,data interface{}) (e error) {
	log.Println("Query","SchemaName:",SchemaName,"TableName:",TableName,"data:",data)
	sendToWS(SchemaName,TableName,data)
	return nil
}
func HproseServer(pluginServer string){
	service := rpc.NewTCPServer(pluginServer)
	service.Debug = true
	service.AddFunction("Insert", Insert)
	service.AddFunction("Update", Update)
	service.AddFunction("Delete", Delete)
	service.AddFunction("Query", Query)
	service.AddFunction("Check", Check)
	service.Start()
}
