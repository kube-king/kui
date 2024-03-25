package task

import (
	"os"
)

type Unarchive struct {
	Title         string
	LocalFilePath string
	RemoteDir     string
	Mode          os.FileMode
	Env           map[string]interface{}
	Force         bool
}

func (u *Unarchive) Run(option *ModuleOptions) (result TaskResult, err error) {

	result = TaskResult{
		Res: Result{
			Title: u.Title,
			State: StateOK,
		},
	}

	err = option.SshClient.Unarchive(u.LocalFilePath, u.RemoteDir, u.Mode, u.Force)
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = err.Error()
		return
	}

	return
}

func (u *Unarchive) GetEnv() map[string]interface{} {
	return u.Env
}
