package task

type Replace struct {
	Title          string
	RemoteFilePath string
	Regexp         string
	Replace        string
	Env            map[string]interface{}
}

func (r *Replace) Run(option *ModuleOptions) (result TaskResult, err error) {
	result = TaskResult{
		Res: Result{
			Title: r.Title,
			State: StateOK,
		},
	}

	err = option.SshClient.Replace(r.RemoteFilePath, r.Regexp, r.Replace)
	if err != nil {
		result.Res.State = StateFailed
		result.Res.Message = err.Error()
		return
	}

	return
}

func (r *Replace) GetEnv() map[string]interface{} {
	return r.Env
}
