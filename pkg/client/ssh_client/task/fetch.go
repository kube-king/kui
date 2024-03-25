package task

import (
	"path/filepath"
)

type Fetch struct {
	Title      string
	RemotePath string
	LocalDir   string
	DelegateTo string
	Force      bool
	Env        map[string]interface{}
}

func (f *Fetch) Run(option *ModuleOptions) (result TaskResult, err error) {
	result = TaskResult{
		Res: Result{
			Title: f.Title,
			State: StateOK,
		},
	}

	if f.DelegateTo != "" && option.SshClient.Config.Ip != f.DelegateTo {
		result.Res.State = StateNoExec
		return
	}

	f.LocalDir = filepath.Join(f.LocalDir, option.SshClient.Config.Ip)

	err = option.SshClient.FetchFile(f.RemotePath, f.LocalDir, f.Force)
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = err.Error()
		return
	}

	return
}

func (f *Fetch) GetEnv() map[string]interface{} {
	return f.Env
}
