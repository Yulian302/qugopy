package shell

var (
	tokenGroups = [][]string{
		{"add", "task", "--type", "download_file", "--payload", "*", "--priority", "*"},
		{"add", "task", "--type", "send_email", "--payload", "*", "--priority", "*"},
		{"add", "task", "--type", "process_image", "--payload", "*", "--priority", "*"},
	}
)
