package agents

type ToolExecutor func(args []string) string

var ActionMapper = map[string]ToolExecutor{
	"WriteFile": WriteFile,
}

func ExecuteTool(actions []string, args []Argument) {
	for _, action := range actions {
		for _, arg := range args {
			argNames := arg.ArgNames
			actionName := arg.ActionName
			if actionName == action {
				ActionMapper[action](argNames)
				continue
			}
		}
	}
}
