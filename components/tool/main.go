package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexeyco/simpletable"
	"github.com/cat3306/gocommon/confutil"
	"github.com/cat3306/goworld/conf"
	"github.com/cat3306/goworld/util"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	bash = "/bin/bash"
)

var (
	appName  = os.Args[0]
	confFile = ".gamecluster.json"
)

func main() {

	app := cli.App{
		Name:        "gametool",
		Usage:       "game cluster manage tool",
		Description: "",
		Commands:    commands(),
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
func commands() cli.Commands {
	tmp := cli.Commands{
		&cli.Command{
			Name:        "init",
			Usage:       fmt.Sprintf("%s init", appName),
			Description: "set config path and init template warning:this action will restore default settings ",
			Action:      initConfig,
		},
		&cli.Command{
			Name:        "get",
			Usage:       fmt.Sprintf("%s get", appName),
			Description: "get game cluster",
			Action:      GetCluster,
		},
		&cli.Command{
			Name:        "start",
			Usage:       fmt.Sprintf("%s get", appName),
			Description: "start game cluster",
			Action:      StartCluster,
		},

		&cli.Command{
			Name:        "stop",
			Usage:       fmt.Sprintf("%s get", appName),
			Description: "stop game cluster",
			Action:      StopCluster,
		},
		&cli.Command{
			Name:        "restart",
			Usage:       fmt.Sprintf("%s get", appName),
			Description: "get game cluster",
			Action:      GetCluster,
		},
	}
	return tmp
}
func StopCluster(ctx *cli.Context) error {
	target, _, err := getTargetClusterConf(ctx)
	name := ctx.Args().First()
	if err != nil || target == nil {
		if target == nil {
			err = errors.New(fmt.Sprintf("not found %s config", name))
		}
		return err
	}
	for _, servers := range target.Servers {
		for _, server := range servers {
			raw, err := ioutil.ReadFile(server.Deploy.PidFile)
			if err != nil {
				return err
			}
			pid := util.BytesToString(raw)
			args := []string{pid}
			fmt.Println(strings.Join([]string{"kill", pid, server.Deploy.BinPath}, " "))
			str, err := Cmd("kill", args, false, time.Hour)
			if err != nil {
				fmt.Println(err)
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	return nil
}
func getTargetClusterConf(ctx *cli.Context) (*conf.ClusterConf, string, error) {
	if ctx.NArg() != 1 {
		return nil, "", errors.New(fmt.Sprintf("invalid args example:%s stop local", appName))
	}
	clusterList, p, err := loadConf()
	if err != nil {
		return nil, p, err
	}
	name := ctx.Args().First()
	for _, v := range clusterList {
		if v.Name == name {
			return &v, p, nil
		}
	}
	return nil, "", errors.New(fmt.Sprintf("not found %s config", name))
}
func StartCluster(ctx *cli.Context) error {
	target, configPath, err := getTargetClusterConf(ctx)
	name := ctx.Args().First()
	if err != nil || target == nil {
		if target == nil {
			err = errors.New(fmt.Sprintf("not found %s config", name))
		}
		return err
	}
	for _, servers := range target.Servers {
		for _, server := range servers {
			args := []string{"-c", configPath, "-name", name, "-idx", strconv.Itoa(server.Idx)}

			str, err := Cmd(server.Deploy.BinPath, args, true, time.Hour)
			if err != nil {
				fmt.Println(err)
			}
			if str != "" {
				fmt.Println(str)
			}
		}
	}
	return err
}

func Cmd(name string, arg []string, showArgs bool, timeoutArgs ...time.Duration, ) (string, error) {
	if showArgs {
		fmt.Println(name + " " + strings.Join(arg, " "))
	}
	timeout := 3 * time.Second
	if len(arg) == 0 {
		return "", errors.New("arg empty")
	}
	if len(timeoutArgs) > 0 {
		timeout = timeoutArgs[0]
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, arg...)
	raw, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(raw)), err
}

func getConfPath() (string, error) {
	u, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	u = path.Join(u, confFile)
	return u, nil
}
func initConfig(ctx *cli.Context) error {
	confPath, err := getConfPath()
	if err != nil {
		return err
	}
	config := confutil.Config{}
	clusterList := []conf.ClusterConf{
		conf.ClusterConf{
			Name: "local",
			AuthConfig: conf.AuthConfig{
				IsAuth:         true,
				PrivateKeyPath: "/Users/joker/code/go/src/github.com/cat3306/goworld/cert/private_key.pem",
			},
			Servers: map[util.ClusterType][]conf.ServerConf{
				util.ClusterTypeGate: []conf.ServerConf{
					{
						Idx:             0,
						Host:            "0.0.0.0",
						Port:            8888,
						MaxConn:         1000,
						ConnWriteBuffer: 1048576,
						ConnReadBuffer:  1048576,
						KV:              map[string]interface{}{},
						Logic:           "gate",
						OuterIp:         "",
						Deploy: conf.DeployConf{
							IsDaemon: true,
							BinPath:  "/Users/joker/code/go/bin/gate",
							PidFile:  "/Users/joker/Documents/game/pid/gate0.pid",
							LogFile:  "/Users/joker/Documents/game/log/gate0.out",
						},
					},
				},
				util.ClusterTypeGame: []conf.ServerConf{
					{
						Idx:             0,
						Logic:           "base",
						Host:            "127.0.0.1",
						Port:            8890,
						MaxConn:         1000,
						ConnWriteBuffer: 1048576,
						ConnReadBuffer:  1048576,
						KV:              map[string]interface{}{},
						OuterIp:         "127.0.0.1",
						Deploy: conf.DeployConf{
							IsDaemon: true,
							BinPath:  "/Users/joker/code/go/bin/game",
							PidFile:  "/Users/joker/Documents/game/pid/game0.pid",
							LogFile:  "/Users/joker/Documents/game/log/game0.out",
						},
					},
				},
			},
		},
	}

	err = config.Save(confPath, clusterList)
	if err != nil {
		return err
	}
	return errors.New(confPath)
}
func loadConf() ([]conf.ClusterConf, string, error) {
	p, err := getConfPath()
	if err != nil {
		return nil, "", err
	}
	config := confutil.Config{}
	var clusterList []conf.ClusterConf
	err = config.Load(p, &clusterList)
	return clusterList, p, err
}
func GetCluster(cxt *cli.Context) error {
	clusterList, _, err := loadConf()
	if err != nil {
		return err
	}
	tmpTable := simpletable.New()
	if cxt.NArg() == 0 {

		tmpTable.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "#"},
				{Align: simpletable.AlignCenter, Text: "Name"},
			},
		}
		for i, row := range clusterList {
			r := []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%d", i)},
				{Align: simpletable.AlignCenter, Text: fmt.Sprintf("%s", row.Name)},
				//{Text: row.Servers},
				//{Align: simpletable.AlignRight, Text: fmt.Sprintf("$ %.2f", row[2].(float64))},
			}

			tmpTable.Body.Cells = append(tmpTable.Body.Cells, r)
		}

	} else {
		tmpTable.Header = &simpletable.Header{
			Cells: []*simpletable.Cell{
				{Align: simpletable.AlignCenter, Text: "id"},
				{Align: simpletable.AlignCenter, Text: "Name"},
			},
		}
	}
	tmpTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(tmpTable.String())
	return nil
}
