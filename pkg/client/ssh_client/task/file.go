package task

import "os"

type File struct {
	Title       string
	Paths       []string
	Type        string
	Mode        os.FileMode
	Owner       string
	Group       string
	IgnoreError bool
	Env         map[string]interface{}
}

func (f *File) Run(option *ModuleOptions) (result TaskResult, err error) {
	result = TaskResult{
		Res: Result{
			Title: f.Title,
			State: StateOK,
		},
	}

	for _, p := range f.Paths {
		err = option.SshClient.CreatePath(p, f.Type, f.Mode, f.Owner, f.Group)
		if err != nil {
			result.Res.State = StateFailed
			result.Res.Message = err.Error()
			break
		}
	}

	return
}

func (f *File) GetEnv() map[string]interface{} {
	return f.Env
}
