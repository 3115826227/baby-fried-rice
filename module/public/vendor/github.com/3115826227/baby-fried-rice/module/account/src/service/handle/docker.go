package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/3115826227/baby-fried-rice/module/account/src/config"
	"github.com/3115826227/baby-fried-rice/module/account/src/log"
	"github.com/3115826227/baby-fried-rice/module/account/src/service/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
}

func ImagesGet(c *gin.Context) {
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var rsp = make([]model.RspImages, 0)
	for _, image := range images {
		t := time.Unix(image.Created, 0).Format(config.TimeLayout)
		rsp = append(rsp, model.RspImages{
			Id:         strings.Split(image.ID, ":")[1][:12],
			Name:       image.RepoTags,
			Size:       image.Size,
			Timestamp:  image.Created,
			CreateTime: t,
		})
	}
	SuccessResp(c, "", rsp)
}

func ContainersGet(c *gin.Context) {
	containers, err := GetContainer()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var rsp = make([]model.RspContainers, 0)
	for _, container := range containers {
		var rspPorts = make([]model.ContainerPorts, 0)
		for _, p := range container.Ports {
			rspPorts = append(rspPorts, model.ContainerPorts{
				Ip:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
			})
		}
		rsp = append(rsp, model.RspContainers{
			Id:         container.ID[:12],
			Name:       container.Names[0][1:],
			ImageId:    strings.Split(container.ImageID, ":")[1][:12],
			ImageName:  container.Image,
			Timestamp:  container.Created,
			CreateTime: time.Unix(container.Created, 0).Format(config.TimeLayout),
			State:      container.State,
			Status:     container.Status,
			Ports:      rspPorts,
		})
	}
	SuccessResp(c, "", rsp)
}

func GetContainer() (containers []types.Container, err error) {
	containers, err = cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	return
}

func StatsGet(c *gin.Context) {
	rsp, err := ContainerStatsCmd()
	if err != nil {
		log.Logger.Warn(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	SuccessResp(c, "", rsp)
}

func ContainerStatsCmd() (stats []model.RspCmdContainerStats, err error) {
	name := "docker"
	args := []string{
		"stats",
		"--no-stream",
		"--format",
		fmt.Sprintf(`{"id":"{{ .Container }}","name": "{{ .Name }}","use_limit_memory": "{{ .MemUsage }}","memory_percent":"{{ .MemPerc }}","cpu_percent":"{{ .CPUPerc }}"}`),
	}
	var data []byte
	data, err = Cmd(name, args...)
	if err != nil {
		return
	}
	stats = make([]model.RspCmdContainerStats, 0)
	res := strings.Split(string(data), "\n")
	for _, d := range res {
		if d == "" {
			continue
		}
		var status = model.RspCmdContainerStats{}
		err = json.Unmarshal([]byte(d), &status)
		if err != nil {
			log.Logger.Warn(err.Error())
			continue
		}
		stats = append(stats, status)
	}
	return
}

func Cmd(name string, args ...string) (data []byte, err error) {
	cmd := exec.Command(name, args...)
	var stdout io.ReadCloser
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	defer stdout.Close()
	if err = cmd.Start(); err != nil {
		log.Logger.Warn(err.Error())
		return
	}
	data, err = ioutil.ReadAll(stdout)
	if err != nil {
		log.Logger.Warn(err.Error())
	}
	return
}
