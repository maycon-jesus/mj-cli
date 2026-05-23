package commands

import (
	"context"
	"log/slog"

	"github.com/maycon-jesus/mj-cli/pkg/config"
	"github.com/maycon-jesus/mj-cli/pkg/intl"
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

	// Chaves utilizadas para configurar o comando a partir do arquivo de configuração
	Configs map[string]Config

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

type Config struct {
	DescriptionKey string
	DefaultValue   interface{}
}

type ExecData struct {
	Args       []string
	Flags      *pflag.FlagSet
	Config     config.ReadOnlyConfig
	Translator *intl.Translator
	Terminal   *mjterm.Terminal
	Logger     *slog.Logger
	RootCmd    *Command
}

func (c *Command) ToCobraCommand(app *App) *cobra.Command {
	app.Logger.Log.Debug("Converting custom command to Cobra command", "command", c.Name)

	// Adiciona traduções ao tradutor
	app.Translator.AddMessagesBulk(c.Translations)

	// Configs
	for key, cfg := range c.Configs {
		description := ""
		if cfg.DescriptionKey != "" {
			description = app.Translator.T(cfg.DescriptionKey, nil)
		}
		if err := app.Config.AddEntry(key, description, cfg.DefaultValue); err != nil {
			app.Logger.Log.Warn("Falha ao registrar entrada de configuração", "key", key, "error", err.Error())
		}
	}

	// Cria o comando cobra
	cmd := &cobra.Command{
		Use:     makeUseString(c.Name, c.Args),
		Short:   app.Translator.T(c.ShortDescriptionKey, nil),
		Long:    app.Translator.T(c.LongDescriptionKey, nil),
		Example: app.Translator.T(c.ExampleKey, nil),
		Aliases: c.Aliases,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if daddy := cmd.Parent(); daddy != nil {
				if fn := daddy.PreRunE; fn != nil {
					app.Logger.Log.Debug("Running PreRunE hook", "command", c.Name, "parent_command", daddy.Name())
					err := fn(daddy, args)
					if err != nil {
						return err
					}
				}
			}

			if c.BeforeRun != nil {
				data := createExecData(app, cmd, args)

				err := c.BeforeRun(ctx, data)
				if err != nil {
					logPath := app.Logger.FileHandler.Name()
					app.Logger.Log.Error("Error in BeforeRun hook", "command", c.Name, "error", err.Error())
					app.Terminal.Printf("An error occurred in the BeforeRun hook. Please check the log file for details: %s\n", logPath)
					return err
				}
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
				data := createExecData(app, cmd, args)
				err := c.Handler(ctx, data)
				if err != nil {
					logPath := app.Logger.FileHandler.Name()
					app.Logger.Log.Error("Error in main handler", "command", c.Name, "error", err.Error())
					app.Terminal.Printf("An error occurred in the main handler. Please check the log file for details: %s\n", logPath)
				}
				return err
			}

			return nil
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if c.AfterRun != nil {
				data := createExecData(app, cmd, args)
				err := c.AfterRun(ctx, data)
				if err != nil {
					logPath := app.Logger.FileHandler.Name()
					app.Logger.Log.Error("Error in AfterRun hook", "command", c.Name, "error", err.Error())
					app.Terminal.Printf("An error occurred in the AfterRun hook. Please check the log file for details: %s\n", logPath)
					return err
				}
			}

			if daddy := cmd.Parent(); daddy != nil {
				if fn := daddy.PostRunE; fn != nil {
					app.Logger.Log.Debug("Running PostRunE hook", "command", c.Name, "parent_command", daddy.Name())
					err := fn(daddy, args)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	if cmd.Name() == app.Name {
		cmd.Version = app.Version
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

	if flag.Required {
		app.Logger.Log.Debug("Marking flag as required", "flag", flag.Name, "command", c.Name)
		err := cmd.MarkFlagRequired(flag.Name)
		if err != nil {
			app.Logger.Log.Error("Failed to mark flag as required", "flag", flag.Name, "command", c.Name, "error", err.Error())
		}
	}
}

func createExecData(app *App, cmd *cobra.Command, args []string) *ExecData {
	return &ExecData{
		Args:       args,
		Flags:      cmd.Flags(),
		Config:     app.Config,
		Translator: app.Translator,
		Terminal:   app.Terminal,
		Logger:     app.Logger.Log,
		RootCmd:    app.RootCmd,
	}
}
