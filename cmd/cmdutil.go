package cmd
var dbSubCmds = []string{"migrate", "create", "drop", "status", "import"}

func isDbCmd(cmd string) bool {
	for _, c := range dbSubCmds {
		if cmd == c {
			return true
		}
	}
	return false
}