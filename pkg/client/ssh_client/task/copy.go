package task

import "os"

type Copy struct {
	Title          string
	LocalFilePath  string
	RemoteFilePath string
	Force          bool
	Mode           os.FileMode
	Env            map[string]interface{}
}

func (c *Copy) Run(option *ModuleOptions) (result TaskResult, err error) {

	result = TaskResult{
		Res: Result{
			Title: c.Title,
			State: StateOK,
		},
	}

	err = option.SshClient.CopyFile(c.LocalFilePath, c.RemoteFilePath, c.Force, c.Mode)
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = err.Error()
		return
	}

	return
}

func (c *Copy) GetEnv() map[string]interface{} {
	return c.Env
}
