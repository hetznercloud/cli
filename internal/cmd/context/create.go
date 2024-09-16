package context

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/cli/internal/state/config"
)

func NewCreateCommand(s state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "create [--token-from-env] <name>",
		Short:                 "Create a new context",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runCreate),
	}
	cmd.Flags().Bool("token-from-env", false, "If true, the HCLOUD_TOKEN from the environment will be used without asking")
	return cmd
}

func runCreate(s state.State, cmd *cobra.Command, args []string) error {
	tokenFromEnv, _ := cmd.Flags().GetBool("token-from-env")

	cfg := s.Config()
	if !s.Terminal().StdoutIsTerminal() && !tokenFromEnv {
		return errors.New("non-interactive tty detected. Use --token-from-env to use HCLOUD_TOKEN from the environment")
	}

	name := strings.TrimSpace(args[0])
	if name == "" {
		return errors.New("invalid name")
	}
	if config.ContextByName(cfg, name) != nil {
		return errors.New("name already used")
	}

	var token string

	envToken := os.Getenv("HCLOUD_TOKEN")
	if envToken != "" {
		if len(envToken) != 64 {
			if tokenFromEnv {
				return errors.New("invalid token: must be 64 characters in length")
			}
			cmd.Println("Warning: HCLOUD_TOKEN is set, but token is invalid (must be exactly 64 characters long)")
		} else if tokenFromEnv {
			token = envToken
		} else {
			cmd.Print("The HCLOUD_TOKEN environment variable is set. Do you want to use the token from HCLOUD_TOKEN for the new context? (Y/n): ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if s := strings.ToLower(scanner.Text()); s == "" || s == "y" || s == "yes" {
				token = envToken
			}
		}
	}

	if token == "" {
		if tokenFromEnv {
			return errors.New("no token provided")
		}
		for {
			cmd.Printf("Token: ")
			// Conversion needed for compilation on Windows
			//                                       vvv
			btoken, err := s.Terminal().ReadPassword(int(syscall.Stdin))
			cmd.Print("\n")
			if err != nil {
				return err
			}
			token = string(bytes.TrimSpace(btoken))
			if token == "" {
				continue
			}
			if len(token) != 64 {
				cmd.Print("Entered token is invalid (must be exactly 64 characters long)\n")
				continue
			}
			break
		}
	}

	context := config.NewContext(name, token)

	cfg.SetContexts(append(cfg.Contexts(), context))
	cfg.SetActiveContext(context)

	if err := cfg.Write(nil); err != nil {
		return err
	}

	cmd.Printf("Context %s created and activated\n", name)

	return nil
}
