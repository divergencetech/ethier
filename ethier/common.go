package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// getStringFlag gets the named string flag from the Command's FlagSet, checks
// that it isn't empty, and returns the value.
func getStringFlag(cmd *cobra.Command, name string) (string, error) {
	f, err := cmd.Flags().GetString(name)
	if err != nil {
		return "", err
	}
	if f == "" {
		return "", fmt.Errorf("flag --%s not provided", name)
	}
	return f, nil
}

// boolPrompt runs a confirmation Prompt with the label, returning true iff the
// Prompt's returned value begins with y or Y.
func boolPrompt(format string, a ...interface{}) (bool, error) {
	p := &promptui.Prompt{
		Label:     fmt.Sprintf(format, a...),
		IsConfirm: true,
	}
	s, err := p.Run()
	if err != nil && err != promptui.ErrAbort {
		return false, err
	}
	return strings.HasPrefix(strings.ToLower(s), "y"), nil
}
