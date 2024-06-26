package cmd

import (
	"os"
	"runtime"
	"time"

	"github.com/fmotalleb/watch2do/cli"
	"github.com/fsnotify/fsnotify"
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
  watch2do --execute "echo files changed" --watch "*" --debounce 2500ms
`,

	Example: "  watch2do --execute 'echo files changed' --watch '*' --debounce 2500ms",

	Run: func(cmd *cobra.Command, args []string) {
		var level logrus.Level
		if getBool(cmd.Flags(), "verbose") {
			level = logrus.DebugLevel
		} else {
			level = logrus.InfoLevel
		}

		Params = cli.Params{
			Shell:             getString(cmd.Flags(), "shell"),
			WatchList:         getArray(cmd.Flags(), "watch"),
			ExcludeWatchList:  getArray(cmd.Flags(), "exclude"),
			MatchList:         getArray(cmd.Flags(), "match"),
			Commands:          getArray(cmd.Flags(), "execute"),
			Debounce:          getDuration(cmd.Flags(), "debounce"),
			LogLevel:          level,
			Operations:        getTriggerFlags(cmd.Flags()),
			JsonOutput:        getBool(cmd.Flags(), "log-json"),
			KillBeforeExecute: !getBool(cmd.Flags(), "no-kill"),
			Recursive:         getBool(cmd.Flags(), "recursive"),
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
func getTriggerFlags(flags *pflag.FlagSet) []fsnotify.Op {
	mapper := map[string]fsnotify.Op{
		"no-create": fsnotify.Create,
		"no-write":  fsnotify.Write,
		"no-rename": fsnotify.Rename,
		"no-remove": fsnotify.Remove,
		"no-chmod":  fsnotify.Chmod,
	}

	result := make([]fsnotify.Op, 0)

	for k, v := range mapper {
		r, err := flags.GetBool(k)
		if err != nil {
			os.Exit(1)
		}
		if r == false {
			result = append(result, v)
		}
	}

	return result
}

func init() {
	rootCmd.Flags().StringSliceP("execute", "x", []string{}, "Commands to execute after receiving a change event")
	rootCmd.Flags().Bool("no-kill", false, "Don't kill old processes from last trigger")
	rootCmd.Flags().StringSliceP("watch", "w", []string{"."}, "Directories to watch")
	rootCmd.Flags().StringSlice("exclude", []string{".git"}, "Directories to ignore")
	rootCmd.Flags().BoolP("recursive", "r", false, "recursively watch subdirectories")
	rootCmd.Flags().StringSliceP("match", "m", []string{"*", "*.*", "**/*"}, "Match with given globs (supports glob pattern), by default matches everything")
	rootCmd.Flags().DurationP("debounce", "d", time.Microsecond, "Debounce time (wait time before executing command or receiving another event)")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose logging")
	if runtime.GOOS == "windows" {
		rootCmd.Flags().StringP("shell", "s", "cmd /c", "Shell executable for windows by default uses `cmd /c`")
	} else {
		rootCmd.Flags().StringP("shell", "s", "sh -c", "Shell executable for linux by default uses `sh -c`")
	}
	rootCmd.Flags().Bool("no-write", false, "Trigger on write")
	rootCmd.Flags().Bool("no-create", false, "Trigger on create")
	rootCmd.Flags().Bool("no-rename", false, "Trigger on rename")
	rootCmd.Flags().Bool("no-remove", false, "Trigger on remove")
	rootCmd.Flags().Bool("no-chmod", false, "Trigger on chmod")
	rootCmd.Flags().Bool("log-json", false, "Use json logger instead of text logger")
}
