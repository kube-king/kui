package task

import (
	"kube-invention/pkg/client/ssh_client"
	"sync"
)

const (
	StateOK     = "ok"
	StateFailed = "failed"
	StateSkip   = "skip"
	StateNoExec = "no_exec"
)

type TaskResult struct {
	Res    Result
	Config ssh_client.Config `json:"-"`
}

type ResultFunc func(module Module, status Result) Result

type Result struct {
	Title   string
	Ip      string
	State   string
	Output  string
	Message string
}

type ModuleOptions struct {
	SshClient *ssh_client.Client
	Env       sync.Map
}

func (m *ModuleOptions) SetEnv(envs map[string]interface{}) {
	for k, v := range envs {
		m.Env.Store(k, v)
	}
}

func (m *ModuleOptions) GetEnv() map[string]interface{} {
	envs := make(map[string]interface{}, 0)
	m.Env.Range(func(key, value any) bool {
		envs[key.(string)] = value
		return true
	})
	return envs
}
