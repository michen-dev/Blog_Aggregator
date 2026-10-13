package commands

import(
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Commands map[string]func(*State, Command) error
}


func (c *Commands) Run(s *State, cmd Command) error {
	f, exists := c.Commands[cmd.Name]
	if !exists {
		return fmt.Errorf("Command not found: %s", cmd.Name)
	}

	return f(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.Commands == nil {
        c.Commands = make(map[string]func(*State, Command) error)
    }
	c.Commands[name] = f
}