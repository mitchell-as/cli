package captain

import (
	"fmt"

	"github.com/ActiveState/cli/internal/analytics"
	"github.com/ActiveState/cli/internal/locale"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cobraCommander interface {
	GetCobraCmd() *cobra.Command
}

type Executor func(cmd *Command, args []string) error

type Command struct {
	cobra *cobra.Command

	name string

	flags     []*Flag
	arguments []*Argument

	execute func(cmd *Command, args []string) error
}

func NewCommand(name string, flags []*Flag, args []*Argument, executor Executor) *Command {
	// Validate args
	for idx, arg := range args {
		if idx > 0 && arg.Required && !args[idx-1].Required {
			msg := fmt.Sprintf(
				"Cannot have a non-required argument followed by a required argument.\n\n%v\n\n%v",
				arg, args[len(args)-1],
			)
			panic(msg)
		}
	}

	cmd := &Command{
		execute:   executor,
		arguments: args,
		flags:     flags,
	}

	cmd.cobra = &cobra.Command{
		Use:  name,
		RunE: cmd.runner,
	}

	cmd.setFlags(flags)

	return cmd
}

func (c *Command) Usage() error {
	return c.cobra.Usage()
}

func (c *Command) Execute(args []string) error {
	c.cobra.SetArgs(args)
	err := c.cobra.Execute()
	c.cobra.SetArgs(nil)
	return err
}

func (c *Command) SetAliases(aliases []string) {
	c.cobra.Aliases = aliases
}

func (c *Command) SetDescription(description string) {
	c.cobra.Use = description
}

func (c *Command) SetUsageTemplate(usageTemplate string) {
	localizedArgs := []map[string]string{}
	for _, arg := range c.Arguments() {
		req := ""
		if arg.Required {
			req = "1"
		}
		localizedArgs = append(localizedArgs, map[string]string{
			"Name":        locale.T(arg.Name),
			"Description": locale.T(arg.Description),
			"Required":    req,
		})
	}
	c.cobra.SetUsageTemplate(locale.Tt(usageTemplate, map[string]interface{}{
		"Arguments": localizedArgs,
	}))
}

func (c *Command) Arguments() []*Argument {
	return c.arguments
}

func (c *Command) AddChildren(children ...*Command) {
	for _, child := range children {
		c.cobra.AddCommand(child.cobra)
	}
}

func (c *Command) AddLegacyChildren(children ...cobraCommander) {
	for _, child := range children {
		c.cobra.AddCommand(child.GetCobraCmd())
	}
}

func (c *Command) flagByName(name string, persistOnly bool) *Flag {
	for _, flag := range c.flags {
		if flag.Name == name && (!persistOnly || flag.Persist) {
			return flag
		}
	}
	return nil
}

func (c *Command) runner(cobraCmd *cobra.Command, args []string) error {
	analytics.Event(analytics.CatRunCmd, c.cobra.Name())

	// Run OnUse functions for flags
	if !cobraCmd.DisableFlagParsing {
		cobraCmd.Flags().VisitAll(func(cobraFlag *pflag.Flag) {
			if !cobraFlag.Changed {
				return
			}

			flag := c.flagByName(cobraFlag.Name, false)
			if flag == nil || flag.OnUse == nil {
				return
			}

			flag.OnUse()
		})
	}

	for idx, arg := range c.arguments {
		if len(args) > idx {
			(*arg.Variable) = args[idx]
		}
	}
	return c.execute(c, args)
}

func (c *Command) argValidator(cobraCmd *cobra.Command, args []string) error {
	return nil
}