package task

import (
	"errors"
	"fmt"
	"io"
	ssh "kube-invention/pkg/client/ssh_client"
	"kube-invention/pkg/installer/global"
	"kube-invention/pkg/utils/pool"
	"os"
	"sync"
)

type Module interface {
	Run(option *ModuleOptions) (message TaskResult, err error)
	GetEnv() map[string]interface{}
}

const (
	MaxLimit = 20
)

type Task struct {
	name           string
	modules        []Module
	output         io.Writer
	envs           map[string]interface{}
	hosts          []ssh.Config
	maxLimit       int
	resultCallback func(string)
}

func New(name string, hosts ...ssh.Config) *Task {
	t := &Task{
		name:     name,
		hosts:    hosts,
		maxLimit: MaxLimit,
		envs:     make(map[string]interface{}, 0),
	}
	return t
}

func (t *Task) SetPoolLimit(maxLimit int) *Task {
	t.maxLimit = maxLimit
	return t
}

func (t *Task) SetResultCallback(callback func(string)) *Task {
	t.resultCallback = callback
	return t
}

func (t *Task) SetOutput(output io.Writer) *Task {
	t.output = io.MultiWriter(output)
	return t
}

func (t *Task) SetEnv(envs map[string]interface{}) *Task {
	t.envs = envs
	return t
}

func (t *Task) writeStatus(result Result) {
	if t.output == nil || result.State == StateNoExec {
		return
	}

	msg := fmt.Sprintf("=> Task:[%v] | Status:[%v] | IP:[%v]\n", result.Title, result.State, result.Ip)

	//if t.resultCallback != nil {
	//	t.resultCallback(msg)
	//}

	//_, err := t.output.Write([]byte(msg))
	global.Log.Info(msg)
	var err error
	if result.Output != "" && result.State == StateFailed {
		//global.Log.Info("Command Exec Output:" + result.Output)
		global.Log.Error(result.Output)
		//successMsg := fmt.Sprintf("*** [ERROR] Result: %v \n", result.Output)
		//_, err = t.output.Write([]byte(successMsg))
		//if t.resultCallback != nil {
		//	t.resultCallback(successMsg)
		//}
	}
	if result.Message != "" && result.State == StateFailed {
		//global.Log.Info("Command Exec Message:" + result.Message)
		errMsg := fmt.Sprintf("*** [ERROR] Result: %v \n", result.Message)
		_, err = t.output.Write([]byte(errMsg))
		if t.resultCallback != nil {
			t.resultCallback(errMsg)
		}
	}

	if err != nil {
		return
	}
}

func (t *Task) Run(modules ...Module) (statusList []Result, err error) {

	global.Log.Info(fmt.Sprintf("=> ########## Run Task【%v】########## \n", t.name))

	var client *ssh.Client
	t.SetOutput(os.Stdout)

	clients := make([]*ssh.Client, 0)
	for _, host := range t.hosts {
		client, err = ssh.NewClient(host)
		if err != nil {
			return
		}
		client.Stdout = t.output
		clients = append(clients, client)
	}

	var limit int
	var wg pool.SizedWaitGroup

	var isSuccess bool

	hostLength := len(clients)
	if hostLength < MaxLimit {
		limit = hostLength
	} else {
		limit = MaxLimit
	}

	statusList = make([]Result, 0)

	for _, module := range modules {
		isSuccess = true
		wg = pool.New(limit)
		for i := 0; i < hostLength; i++ {
			wg.Add()
			go func(client *ssh.Client) {

				op := ModuleOptions{
					SshClient: client,
					Env:       sync.Map{},
				}
				op.SetEnv(t.envs)

				op.Env.Store("hostname", client.Config.Hostname)
				op.Env.Store("ip", client.Config.Ip)
				op.Env.Store("password", client.Config.Password)
				op.Env.Store("port", client.Config.Port)
				op.Env.Store("username", client.Config.Username)

				result, _ := module.Run(&op)
				result.Config = op.SshClient.Config
				result.Res.Ip = op.SshClient.Config.Ip

				statusList = append(statusList, result.Res)
				if result.Res.State == StateFailed {
					isSuccess = false
					err = errors.New(fmt.Sprintf("[%v] %v : %v", result.Res.Ip, result.Res.Title, result.Res.Message))
				}

				t.writeStatus(result.Res)
				wg.Done()

			}(clients[i])
		}
		wg.Wait()
		if !isSuccess {
			break
		}
	}

	for _, client := range clients {
		client.Close()
	}

	if err != nil {
		return
	}

	global.Log.Info(fmt.Sprintf("=> ########## End Task【%v】########## \n", t.name))

	return statusList, nil
}
