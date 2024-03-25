package task

type Yum struct {
	Title   string
	Names   []string
	State   string
	Timeout int
	Env     map[string]interface{}
}

func (y *Yum) Run(option *ModuleOptions) (result TaskResult, err error) {

	result = TaskResult{
		Res: Result{
			Title: y.Title,
			State: StateOK,
		},
	}

	for _, name := range y.Names {
		option.SshClient.Config.Timeout = y.Timeout
		err = option.SshClient.Yum(name)
		if err != nil {
			result.Res.State = StateFailed
			result.Res.Message = err.Error()
			return
		}
	}

	return
}

func (y *Yum) GetEnv() map[string]interface{} {
	return y.Env
}
