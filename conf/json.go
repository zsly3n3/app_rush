package conf

import (
	"app/log"
	"encoding/json"
	"io/ioutil"
)

const (
	Debug   = "debug"   //开发
	Test    = "test"    //测试
	Release = "release" //发布
)

var Server struct {
	Domain        string
	HttpServer    string
	DB_IP         string
	DB_Name       string
	DB_UserName   string
	DB_Pwd        string
	Redis_IP      string
	Redis_Name    string
	Redis_Pwd     string
	OutHttpServer string
}

var Common struct {
	Mode string
}

func init() {
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		log.Fatal("Read server.json err:%v", err)
	}
	err = json.Unmarshal(data, &Common)
	if err != nil {
		log.Fatal("Unmarshal server.json err:%v", err)
	}

	var file_str string
	switch Common.Mode {
	case Debug:
		file_str = "conf/server_dev.json"
	case Test:
		file_str = "conf/server_test.json"
	case Release:
		file_str = "conf/server_release.json"
	}
	data, err = ioutil.ReadFile(file_str)
	if err != nil {
		log.Fatal("Read json err:%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
