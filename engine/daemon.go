package engine

import (
	"fmt"
	"github.com/cat3306/goworld/conf"
	"github.com/sevlyar/go-daemon"
	"os"
)

func DaemonMode(c *conf.DeployConf, name string) (*daemon.Context, error) {
	ctx := &daemon.Context{
		PidFileName: c.PidFile,
		PidFilePerm: 0644,
		LogFileName: c.LogFile,
		LogFilePerm: 0640,
		Umask:       027,
	}
	d, err := ctx.Reborn()
	if err != nil {

		return nil, err
	}
	if d != nil {
		fmt.Println(fmt.Sprintf("%s run in daemon mode", name))
		os.Exit(0)
	}

	return ctx, nil
}
