package potfile

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/vknabel/lithia/world"
)

type PotfileState struct {
	Cmds map[string]PotfileCmd
}

type PotfileCmd struct {
	Name    string
	Summary string
	Flags   map[string]PotfileFlag
	Envs    map[string]string
	Bin     string
	Args    []string
}

type PotfileFlag struct {
	Name         string
	Short        string
	Summary      string
	DefaultValue string
	Required     bool
}

func (potCmd PotfileCmd) RunCmd(args []string) {
	bin := potCmd.Bin
	if bin == "lithia" {
		bin = os.Args[0]
	}
	cmd := exec.CommandContext(context.Background(), bin, (append(potCmd.Args, args...))...)
	cmd.Env = append(cmd.Env, os.Environ()...)

	for key, val := range potCmd.Envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%q", key, val))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Fprint(world.Current.Stderr, err)
		world.Current.Env.Exit(1)
	}
}
