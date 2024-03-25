package task

import (
	"io"
	"time"
)

type Command struct {
	Title       string
	CmdList     []string
	Timeout     int
	IgnoreError bool
	Stdout      io.Writer
	Callback    ResultFunc
	Retries     int
	Delay       time.Duration
	DelegateTo  string
	When        func() bool
	Until       func(output string) bool
	Env         map[string]interface{}
}

func (c *Command) Run(option *ModuleOptions) (result TaskResult, err error) {
	result = TaskResult{
		Res: Result{
			Title: c.Title,
			State: StateOK,
		},
	}

	if c.When != nil && !c.When() {
		result.Res.State = StateSkip
		return
	}

	if c.DelegateTo != "" && option.SshClient.Config.Ip != c.DelegateTo {
		result.Res.State = StateNoExec
		return
	}

	option.SetEnv(c.Env)
	var render string
	for _, cmd := range c.CmdList {
		render, err = TextRender(option.GetEnv(), cmd)
		if err != nil {
			result.Res.State = StateFailed
			result.Res.Message = result.Res.Output
			break
		}

		//global.Log.Info("Command:" + render)
		if c.Retries > 0 && c.Delay > 0 {
			result.Res.State = StateFailed
			for i := 0; i < c.Retries; i++ {
				time.Sleep(c.Delay * time.Second)
				result.Res.Output, err = option.SshClient.Command(render)
				if c.IgnoreError {
					result.Res.State = StateSkip
					break
				}
				if c.Until(result.Res.Output) && err == nil {
					result.Res.State = StateOK
					break
				}
			}
		} else {
			result.Res.Output, err = option.SshClient.Command(render)
			if !c.IgnoreError && err != nil {
				result.Res.State = StateFailed
				result.Res.Message = result.Res.Output
				break
			}
		}
	}

	if result.Res.State != StateFailed && c.Callback != nil {
		result.Res = c.Callback(c, result.Res)
	}

	return
}

func (c *Command) GetEnv() map[string]interface{} {
	return c.Env
}
