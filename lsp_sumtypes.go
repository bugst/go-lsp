package lsp

import (
	"encoding/json"
	"errors"
)

type CommandOrCodeAction struct {
	command    *Command
	codeAction *CodeAction
}

func (c *CommandOrCodeAction) Set(value interface{}) {
	c.command = nil
	c.codeAction = nil
	switch v := value.(type) {
	case *Command:
		c.command = v
	case Command:
		c.command = &v
	case *CodeAction:
		c.codeAction = v
	case CodeAction:
		c.codeAction = &v
	default:
		panic("value must be a Command or a CodeAction")
	}
}

func (c *CommandOrCodeAction) Get() interface{} {
	if c.command != nil {
		return *(c.command)
	}
	if c.codeAction != nil {
		return *(c.codeAction)
	}
	panic("empty value")
}

func (c *CommandOrCodeAction) UnmarshalJSON(data []byte) error {
	c.command = nil
	c.codeAction = nil
	var co Command
	if err := json.Unmarshal(data, &co); err == nil {
		c.command = &co
		return nil
	}
	var ca CodeAction
	if err := json.Unmarshal(data, &ca); err == nil {
		c.codeAction = &ca
		return nil
	}
	return errors.New("expected Command or CodeAction")
}

func (c CommandOrCodeAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Get())
}
