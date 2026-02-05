package commands

import (
	"context"

	"github.com/maycon-jesus/mj-cli/internal/config"
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
}

func (c *Command) ToCobraCommand(config *config.ConfigRegistry, translator *intl.Translator) *cobra.Command {

	// Adiciona traduções ao tradutor
	translator.AddMessagesBulk(c.Translations)

	// Cria o comando cobra
	cmd := &cobra.Command{
		Use:     makeUseString(c.Name, c.Args),
		Short:   translator.T(c.ShortDescriptionKey, nil),
		Long:    translator.T(c.LongDescriptionKey, nil),
		Example: translator.T(c.ExampleKey, nil),
		Aliases: c.Aliases,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if c.BeforeRun != nil {
				data := &ExecData{
					Args:       args,
					Flags:      cmd.Flags(),
					Config:     config,
					Translator: translator,
					Terminal:   mjterm.New(),
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
					Config:     config,
					Translator: translator,
					Terminal:   mjterm.New(),
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
					Config:     config,
					Translator: translator,
					Terminal:   mjterm.New(),
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
		c.addFlag(translator, config, cmd, flag)
	}

	// Adiciona subcomandos recursivamente
	for _, subCmd := range c.SubCommands {
		cmd.AddCommand(subCmd.ToCobraCommand(config, translator))
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
func (c *Command) addFlag(translator *intl.Translator, config *config.ConfigRegistry, cmd *cobra.Command, flag Flag) {
	switch v := flag.DefaultValue.(type) {
	case string:
		cmd.Flags().StringP(flag.Name, flag.Shorthand, v, translator.T(flag.DescriptionKey, nil))
	case int:
		cmd.Flags().IntP(flag.Name, flag.Shorthand, v, translator.T(flag.DescriptionKey, nil))
	case bool:
		cmd.Flags().BoolP(flag.Name, flag.Shorthand, v, translator.T(flag.DescriptionKey, nil))
	case []string:
		cmd.Flags().StringSliceP(flag.Name, flag.Shorthand, v, translator.T(flag.DescriptionKey, nil))
	}
	if flag.ConfigRegistry != (FlagConfigRegistry{}) {
		config.GetModule(flag.ConfigRegistry.RegistryName).BindPFlag(flag.ConfigRegistry.Key, cmd.Flags().Lookup(flag.Name))
	}

	if flag.Required {
		err := cmd.MarkFlagRequired(flag.Name)
		if err != nil {
			panic(err)
		}
	}
}
