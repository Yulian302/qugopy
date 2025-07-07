package shell

var (
	tokenGroups = [][]string{
		{"send", "email"},
		{"send", "file"},
		{"set", "name"},
		{"set", "password"},
		{"run", "script"},
		{"run", "executable"},
		{"start", "command", "workers"},
		{"start", "command", "redis"},
	}
)
