package task

import (
	"os"
)

type Template struct {
	Title            string
	TemplateFilePath string
	RemoteFilePath   string
	Force            bool
	DelegateTo       string
	Mode             os.FileMode
	When             func() bool
	Env              map[string]interface{}
}

func (t *Template) Run(option *ModuleOptions) (result TaskResult, err error) {

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
	err = option.SshClient.RenderTemplate(t.TemplateFilePath, t.RemoteFilePath, option.GetEnv(), t.Force, t.Mode)
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = err.Error()
		return
	}

	return
}

func (t *Template) GetEnv() map[string]interface{} {
	return t.Env
}
