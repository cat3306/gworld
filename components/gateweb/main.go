package main

import (
	"flag"
	"github.com/cat3306/goworld/components/gateweb/conf"
	"github.com/cat3306/goworld/components/gateweb/handler"
	"github.com/cat3306/goworld/glog"
)

func ParseFlag() string {
	var file string
	flag.StringVar(&file, "c", "", "use -c to bind conf file")
	flag.Parse()
	return file
}
func main() {
	filePath := ParseFlag()
	glog.Init()
	err := conf.Load(filePath)
	if err != nil {
		panic(err)
	}

	handler.StartHttpServer(conf.Global.Port)
}
