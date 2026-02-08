package commands

import (
	"context"

	"github.com/maycon-jesus/mj-cli/internal/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
	"github.com/maycon-jesus/mj-cli/pkg/logger"
	"github.com/maycon-jesus/mj-cli/pkg/mjterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Command representa um comando customizado da CLI
type Command struct {
	// Propriedades básicas
	Name                string
	Args                []Arg
	ShortDescriptionKey string
	LongDescriptionKey  string
	ExampleKey          string
	Aliases             []string

	// Flags personalizadas
	Flags []Flag

	// Handler para execução do comando
	Handler CommandHandler

	// Subcomandos (opcional)
	SubCommands []*Command

	// Hooks (opcional)
	BeforeRun CommandHandler
	AfterRun  CommandHandler

	// Traduções (opcional)
	Translations intl.Translations

	// Se true, exibe ajuda quando executado sem argumentos/subcomandos
	RunHelpOnNoArgs bool
}

// CommandHandler é a função que executa a lógica do comando
type CommandHandler func(ctx context.Context, execData *ExecData) error

// Flag representa uma flag customizada
type Flag struct {
	Name           string
	Shorthand      string
	DescriptionKey string
	DefaultValue   interface{}
	Required       bool
	Global         bool
	ConfigRegistry FlagConfigRegistry
}

type FlagConfigRegistry struct {
	RegistryName string
	Key          string
}

type Arg struct {
	Name           string
	DescriptionKey string
	Required       bool
}

type ExecData struct {
	Args       []string
	Flags      *pflag.FlagSet
	Config     *config.ConfigRegistry
	Translator *intl.Translator
	Terminal   *mjterm.Terminal
	Logger     *logger.Logger
}

func (c *Command) ToCobraCommand(app *App) *cobra.Command {
	app.Logger.Debug("Converting custom command to Cobra command", logger.Metadata{"command": c.Name, "event": "convert_command_to_cobra"})
	// Adiciona traduções ao tradutor
	app.Translator.AddMessagesBulk(c.Translations)

	// Cria o comando cobra
	cmd := &cobra.Command{
		Use:     makeUseString(c.Name, c.Args),
		Short:   app.Translator.T(c.ShortDescriptionKey, nil),
		Long:    app.Translator.T(c.LongDescriptionKey, nil),
		Example: app.Translator.T(c.ExampleKey, nil),
		Aliases: c.Aliases,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if c.BeforeRun != nil {
				data := &ExecData{
					Args:       args,
					Flags:      cmd.Flags(),
					Config:     app.Config,
					Translator: app.Translator,
					Terminal:   mjterm.New(),
					Logger:     app.Logger,
				}

				err := c.BeforeRun(ctx, data)
				data.Terminal.Close()
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Se RunHelpOnNoArgs está ativo e não há handler, exibe ajuda
			if c.RunHelpOnNoArgs && c.Handler == nil {
				return cmd.Help()
			}

			// Executa o handler principal
			if c.Handler != nil {
				data := &ExecData{
					Args:       args,
					Flags:      cmd.Flags(),
					Config:     app.Config,
					Translator: app.Translator,
					Terminal:   mjterm.New(),
					Logger:     app.Logger,
				}
				err := c.Handler(ctx, data)
				data.Terminal.Close()
				return err
			}

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if c.AfterRun != nil {
				data := &ExecData{
					Args:       args,
					Flags:      cmd.Flags(),
					Config:     app.Config,
					Translator: app.Translator,
					Terminal:   mjterm.New(),
					Logger:     app.Logger,
				}
				err := c.AfterRun(ctx, data)
				data.Terminal.Close()
				return err
			}
			return nil
		},
	}

	// Adiciona as flags
	for _, flag := range c.Flags {
		c.addFlag(app, cmd, flag)
	}

	// Adiciona subcomandos recursivamente
	for _, subCmd := range c.SubCommands {
		cmd.AddCommand(subCmd.ToCobraCommand(app))
	}

	return cmd
}

func makeUseString(name string, args []Arg) string {
	use := name
	for _, arg := range args {
		if arg.Required {
			use += " <" + arg.Name + ">"
		} else {
			use += " [" + arg.Name + "]"
		}
	}
	return use
}

// addFlag adiciona uma flag ao comando cobra baseado no tipo
func (c *Command) addFlag(app *App, cmd *cobra.Command, flag Flag) {
	flags := cmd.Flags()
	if flag.Global {
		flags = cmd.PersistentFlags()
	}
	switch v := flag.DefaultValue.(type) {
	case string:
		flags.StringP(flag.Name, flag.Shorthand, v, app.Translator.T(flag.DescriptionKey, nil))
	case int:
		flags.IntP(flag.Name, flag.Shorthand, v, app.Translator.T(flag.DescriptionKey, nil))
	case bool:
		flags.BoolP(flag.Name, flag.Shorthand, v, app.Translator.T(flag.DescriptionKey, nil))
	case []string:
		flags.StringSliceP(flag.Name, flag.Shorthand, v, app.Translator.T(flag.DescriptionKey, nil))
	}
	if flag.ConfigRegistry != (FlagConfigRegistry{}) {
		app.Config.GetModule(flag.ConfigRegistry.RegistryName).BindPFlag(flag.ConfigRegistry.Key, flags.Lookup(flag.Name))
	}

	if flag.Required {
		err := cmd.MarkFlagRequired(flag.Name)
		if err != nil {
			panic(err)
		}
	}
}
