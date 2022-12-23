package cmd

import (
	"github.com/spf13/cobra"

	"github.com/vknabel/lithia/langsrv"
)

func init() {
	rootCmd.AddCommand(lspCmd)
	lspCmd.AddCommand(lspStdioCmd)
	lspCmd.AddCommand(lspSocketCmd)
	lspCmd.AddCommand(lspIPCCmd)
	lspCmd.AddCommand(lspTCPCmd)

	lspSocketCmd.Flags().StringVarP(
		&lspSocketAddress,
		"listen",
		"l",
		"127.0.0.1:7998",
		"Address and port on which to listen for LSP connections",
	)
}

var lspCmd = &cobra.Command{
	Use:   "lsp",
	Short: "Language Server",
	Long:  `Runs the language server for the use inside an editor.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		lspStdioCmd.Run(lspStdioCmd, args)
	},
}

var lspStdioCmd = &cobra.Command{
	Use:     "stdio",
	Aliases: []string{"stdin", "-"},
	Short:   "stdio mode. Supported by most editors.",
	Run: func(cmd *cobra.Command, args []string) {
		err := langsrv.RunStdio()
		if err != nil {
			panic(err)
		}
	},
}

var lspIPCCmd = &cobra.Command{
	Use:   "ipc",
	Short: `opens a nodejs ipc connection.`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		err := langsrv.RunIPC()
		if err != nil {
			panic(err)
		}
	},
}

var lspSocketAddress string = "127.0.0.1:7998"
var lspSocketCmd = &cobra.Command{
	Use:   "socket",
	Short: `opens a socket on the specified address. Make sure the port is free.`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		err := langsrv.RunSocket(lspSocketAddress)
		if err != nil {
			panic(err)
		}
	},
}

var lspTCPAddress string = "127.0.0.1:7998"
var lspTCPCmd = &cobra.Command{
	Use:   "tcp",
	Short: `opens a tcp connection on the specified address. Make sure the port is free.`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		err := langsrv.RunTCP(lspTCPAddress)
		if err != nil {
			panic(err)
		}
	},
}
