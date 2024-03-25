package task

type LineInFile struct {
	Title string
	Path  string
	Lines []Line
	Env   map[string]interface{}
}

type Line struct {
	Line        string
	Insertafter string
	State       string
	Pattern     string
}

func (l *LineInFile) Run(option *ModuleOptions) (result TaskResult, err error) {

	result = TaskResult{
		Res: Result{
			Title: l.Title,
			State: StateFailed,
		},
	}

	var lineRender, patternRender, insertafterRender string

	option.SetEnv(l.Env)

	for _, line := range l.Lines {

		env := option.GetEnv()

		lineRender, err = TextRender(env, line.Line)
		if err != nil {
			result.Res.Message = err.Error()
			return
		}
		patternRender, err = TextRender(env, line.Pattern)
		if err != nil {
			result.Res.Message = err.Error()
			return
		}

		insertafterRender, err = TextRender(env, line.Insertafter)
		if err != nil {
			result.Res.Message = err.Error()
			return
		}

		err = option.SshClient.LineInFile(l.Path, patternRender, insertafterRender, lineRender, line.State)
		if err != nil {
			result.Res.Message = err.Error()
			return
		}
	}

	result.Res.State = StateOK
	return
}

func (l *LineInFile) GetEnv() map[string]interface{} {
	return l.Env
}
