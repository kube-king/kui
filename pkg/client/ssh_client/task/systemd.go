package task

type Systemd struct {
	Title       string
	Names       []string
	State       string
	IgnoreError bool
	Enable      bool
	Env         map[string]interface{}
}

func (s *Systemd) Run(option *ModuleOptions) (result TaskResult, err error) {
	result = TaskResult{
		Res: Result{
			Title: s.Title,
			State: StateOK,
		},
	}

	for _, name := range s.Names {
		err = option.SshClient.Systemd(name, s.State, s.Enable)
		if err != nil {
			result.Res.State = StateFailed
			result.Res.Message = err.Error()
			return
		}
	}

	if s.IgnoreError {
		result.Res.State = StateOK
		result.Res.Message = ""
	}

	return
}

func (s *Systemd) GetEnv() map[string]interface{} {
	return s.Env
}
