package bot

var Registry = map[string]Command{}

func Register(cmd Command) {
	Registry[cmd.Name] = cmd
}

func Get(name string) (Command, bool) {
	cmd, ok := Registry[name]
	return cmd, ok
}

func Execute(name string, ctx *Context) (string, bool) {
	cmd, ok := Get(name)
	if !ok {
		return "", false
	}
	return cmd.Execute(ctx), true
}
