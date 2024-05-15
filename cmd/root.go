package cmd

import (
	"os"
	"runtime"
	"time"

	"github.com/fmotalleb/watch2do/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "watch2do",
	Short: "watch2do can listen to given list of files/directories and run a command when it detect changes",
	Long: `watch2do will listen to any changes in given list of files/directories and then run a command when it detect any changes

simple usage:
  # in this example, watch2do will listen for any change in current working directory and wait 2.5seconds before printing "files changed"
  watch2do --execute "echo files changed" --watch "*" --debounce 2500
`,

	Example: "  watch2do --execute 'echo files changed' --watch '*' --debounce 2500",

	Run: func(cmd *cobra.Command, args []string) {
		var level logrus.Level
		if getBool(cmd.Flags(), "verbose") {
			level = logrus.DebugLevel
		} else {
			level = logrus.InfoLevel
		}
		Params = cli.Params{
			Shell:     getString(cmd.Flags(), "shell"),
			WatchList: getArray(cmd.Flags(), "watch"),
			Commands:  getArray(cmd.Flags(), "execute"),
			Debounce:  getDuration(cmd.Flags(), "debounce"),
			LogLevel:  level,
		}
	},
}

// Params will be after Execute
var Params cli.Params

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		// rootCmd.Flags().Usage()
		os.Exit(1)
	}

}
func getString(flags *pflag.FlagSet, name string) string {
	r, err := flags.GetString(name)
	if err != nil {
		os.Exit(1)
	}
	return r
}
func getArray(flags *pflag.FlagSet, name string) []string {
	r, err := flags.GetStringSlice(name)
	if err != nil {
		os.Exit(1)
	}
	return r
}
func getDuration(flags *pflag.FlagSet, name string) time.Duration {
	r, err := flags.GetDuration(name)
	if err != nil {
		os.Exit(1)
	}
	return r
}
func getBool(flags *pflag.FlagSet, name string) bool {
	r, err := flags.GetBool(name)
	if err != nil {
		os.Exit(1)
	}
	return r
}
func init() {
	rootCmd.Flags().StringSliceP("execute", "x", []string{}, "Commands to execute after receiving a change event")
	rootCmd.Flags().StringSliceP("watch", "w", []string{}, "Files/Directories to watch (supports glob pattern)")
	rootCmd.Flags().DurationP("debounce", "d", time.Microsecond, "Debounce time (wait time before executing command or receiving another event)")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose logging")
	if runtime.GOOS == "linux" {
		rootCmd.Flags().StringP("shell", "s", "sh -c", "Shell executable for linux by default uses `sh -c`")
	} else {
		rootCmd.Flags().StringP("shell", "s", "cmd /c", "Shell executable for windows by default uses `cmd /c`")
	}
}
