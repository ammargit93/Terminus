package agents

type ToolExecutor func(str1, str2 string) string

func ExecuteTool(action string, args []string) {
	if action == "WriteFile" && len(args) == 2 {
		WriteFile(args[0], args[1])
	}
}
