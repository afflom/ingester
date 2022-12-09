package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/util/templates"
)

func main() {
	rootCmd := NewRootCmd()
	cobra.CheckErr(rootCmd.Execute())
}

var clientLong = templates.LongDesc(
	`
	Ingester helps label content by ingesting an Emporous workspace along with an
	Emporous schema address and returning a populated dataset-config for the 
	referenced workspace.
	`,
)

// Common describes global configuration options that can be set.
type Ingester struct {
	IOStreams genericclioptions.IOStreams
	LogLevel  string
	CacheDir  string
	Workspace string
	Schema    string
}

// Init initializes default values for Common options.
func (o *Ingester) Init() error {

	o.CacheDir = filepath.Join(xdg.CacheHome, "uor")

	return nil
}

// NewRootCmd creates a new cobra.Command for the command root.
func NewRootCmd() *cobra.Command {
	o := Ingester{}

	o.IOStreams = genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	cmd := &cobra.Command{
		Use:           filepath.Base(os.Args[0]),
		Short:         "Ingester",
		Long:          clientLong,
		SilenceErrors: false,
		SilenceUsage:  false,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if err := o.Init(); err != nil {
				return err
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Complete(args))
			cobra.CheckErr(o.Validate())
			cobra.CheckErr(o.Ingest())
		},
	}

	cmd.Flags().StringVarP(&o.Workspace, "workspace", "", o.Workspace, "workspace to be ingested")
	cmd.Flags().StringVar(&o.Schema, "schema", o.Schema, "schema address")

	return cmd
}

func (o *Ingester) Complete(args []string) error {
	if len(args) < 0 {
		fmt.Printf("ERROR: %v", errors.New("bug: ingester uses flags"))
	}

	return nil
}

func (o *Ingester) Validate() error {
	//if _, err := os.IsExist(os.o.Workspace); err != nil {
	//	fmt.Printf("ERROR: %v", err)
	//}
	return nil
}
