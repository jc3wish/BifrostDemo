package src

import (
	"github.com/garyburd/redigo/redis"
	"strings"
	"strconv"
)


func getRedisUriParam(uri string)(pwd string, network string, url string, database int){
	i := strings.IndexAny(uri, "@")
	pwd = ""
	if i > 0{
		pwd = uri[0:i]
		url = uri[i+1:]
	}else{
		url = uri
	}
	i = strings.IndexAny(url, "/")
	if i > 0 {
		databaseString := url[i+1:]
		intv,err:=strconv.Atoi(databaseString)
		if err != nil{
			database = -1
		}
		database = intv
		url = url[0:i]
	}else{
		database = 0
	}
	i = strings.IndexAny(url, "(")
	if i > 0{
		network = url[0:i]
		url = url[i+1:len(url)-1]
	}else{
		network = "tcp"
	}
	return
}
func NewRedisClient(redisUri string) (*redisClient,error){
	var err error
	var db redis.Conn
	pwd,network,Uri,_ := getRedisUriParam(redisUri)
	if pwd != ""{
		db, err = redis.Dial(network, Uri,redis.DialPassword(pwd))
	}else{
		db, err = redis.Dial(network, Uri)
	}
	if err != nil{
		return nil,err
	}
	return &redisClient{db:db},nil
}
type redisClient struct {
	db redis.Conn
}

func(This *redisClient) Get(key string) string{
	data, err := redis.String(This.db.Do("GET", key))
	if err != nil{
		return "not found data"
	}
	return data
}

func(This *redisClient) LPop(key string) string{
	data, err := redis.String(This.db.Do("lpop", key))
	if err != nil{
		return ""
	}
	return data
}

func(This *redisClient) Close(){
	This.db.Close()
}
