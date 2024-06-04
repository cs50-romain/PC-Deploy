package workspace 

type Workspace struct {
	Prompt			string
	availableCommands	map[string]func(...string) // Maybe not needed if there is a personal shell
}

// Should be initiated only when we are using a workspace
func InitWorkspace(clientName string) *Workspace {
	// Init available commands and its options.
	workspace := &Workspace{
		availableCommands: make(map[string]func(s ...string)),	
	}

	return workspace
}

func (w *Workspace) AddCommand(command string, handler func(...string)) {
	w.availableCommands[command] = handler
}

func (w *Workspace) HandleCommands(command string, options ...string) {
	if _, ok := w.availableCommands[command]; !ok {
		return
	}

	executeCommand := w.availableCommands[command]
	executeCommand()
	return 
}

func StopWorkspace() {
	
}
