package task

import (
	"fmt"
	"os"
	"path"
)

type Script struct {
	Title          string
	ScriptFilePath string
	DelegateTo     string
	When           func() bool
	Env            map[string]interface{}
}

func (t *Script) Run(option *ModuleOptions) (result TaskResult, err error) {

	result = TaskResult{
		Res: Result{
			Title: t.Title,
			State: StateOK,
		},
	}

	if t.When != nil && !t.When() {
		result.Res.State = StateSkip
		return
	}

	if t.DelegateTo != "" && option.SshClient.Config.Ip != t.DelegateTo {
		result.Res.State = StateNoExec
		return
	}

	option.SetEnv(t.Env)

	remoteFile := fmt.Sprintf("/tmp/%v", path.Base(t.ScriptFilePath))
	defer func() {
		err = option.SshClient.RemoveFile(remoteFile)
		if err != nil {
			return
		}
	}()
	err = option.SshClient.RenderTemplate(t.ScriptFilePath, remoteFile, option.GetEnv(), true, os.ModePerm)
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = err.Error()
		return
	}

	result.Res.Output, err = option.SshClient.Command(fmt.Sprintf("bash %v", remoteFile))
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = result.Res.Output
		return
	}

	return
}

func (t *Script) GetEnv() map[string]interface{} {
	return t.Env
}
