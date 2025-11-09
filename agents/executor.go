package agents

type ToolExecutor func(args []string) string

var ActionMapper = map[string]ToolExecutor{
	"WriteFile": WriteFile,
	"MakeDirs":  MakeDirs,
}

func ExecuteTool(actions []string, args []Argument) {
	for _, arg := range args {
		argNames := arg.ArgNames
		actionName := arg.ActionName
		ActionMapper[actionName](argNames)
	}
}
