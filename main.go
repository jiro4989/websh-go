package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type CLIConfig struct {
	UseJSONOutput bool
}

type Config struct {
	UseJSONOutput bool
}

const (
	appName = "websh"
	version = "dev"
)

var (
	config CLIConfig
)

func init() {
	cobra.OnInitialize()
	rootCommand.Flags().SortFlags = false
	rootCommand.Flags().BoolVarP(&config.UseJSONOutput, "json-output", "j", false, "set output format to json")
}


func main() {
	os.Exit(Main())
}

func Main() int {
	if err := rootCommand.Execute(); err != nil {
		return 1
	}
	return 0
}

var rootCommand = &cobra.Command{
	Use:	appName,
	Short:   appName + " is websh cli",
	Example: appName + ` "echo hello world"`,
	Version: version,
	RunE:    runRootCommand,
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	conf := Config {
		UseJSONOutput: config.UseJSONOutput,
	}

	// 引数がない時は標準入力を受け取る
	if len(args) < 1 {
		b , err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		code := string(b)
		if err := runShellgei(code, conf); err != nil {
			return err
		}
		return nil
	}

	code := args[0]
	if err := runShellgei(code, conf); err != nil {
		return err
	}

	return nil
}

func runShellgei(code string, conf Config) error {
	req := &RequestParamPostShellgei{
		Code: code,
		Images: []string{},
	}
	c := NewClient(WebshHost)
	resp, err := c.PostShellgei(req)
	if err != nil {
		return err
	}

	out := resp.Stdout
	if conf.UseJSONOutput{
		b, err := json.Marshal(&resp)
		if err != nil {
			return err
		}
		out = string(b)
	}
	fmt.Println(out)
	return nil
}
