package main

import (
	"os"

	"strconv"

	"log"
	"os/exec"
	"time"

	ocp4cli "bitbucket.org/veldrane/golibs/ocp4cli"
)

type ConfigT struct {
	logger     log.Logger
	session    ocp4cli.SessionT
	namespace  string
	hostname   string
	replicaSet string
	path       string
}

var (
	p ocp4cli.PodsT
)

func main() {

	config := GetConfig()

	config.logger.Println("Replicaset is:", config.replicaSet)
	config.logger.Println("Hostname of this pod:", config.hostname)

	time.Sleep(5 * time.Second)

	podlist, err := p.GetList(&config.session, &config.namespace, &config.replicaSet)

	if err != nil {
		config.logger.Fatalln("Pod list has not been found, please check your access rights for the service account!")
	}

	idx, err := ocp4cli.GetIndex(podlist, &config.hostname, &config.path)

	if err != nil {
		config.logger.Fatalln(err)
	}

	logcmd := "cloudflared tunnel --config /etc/cloudflared/config/config-" + strconv.Itoa(idx) + ".yaml run"

	for _, i := range podlist {
		config.logger.Println("Found pod:", i)
	}

	config.logger.Println("Index of the pod is:", idx)
	config.logger.Println(logcmd)

	cmd := exec.Command("cloudflared", "tunnel", "--config", "/etc/cloudflared/config/config-"+strconv.Itoa(idx)+".yaml", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

}

func GetConfig() *ConfigT {

	config := ConfigT{
		logger:     *log.New(os.Stderr, "[initcfd] -> ", log.Ltime),
		hostname:   os.Getenv("HOSTNAME"),
		session:    *ocp4cli.Session(),
		namespace:  ocp4cli.GetNamespace(),
		replicaSet: ocp4cli.GetReplicationSet(),
		path:       "/etc/cloudflared/config",
	}

	if config.namespace == "" {
		config.logger.Fatalln("No current namespace found!")
	}

	return &config

}
