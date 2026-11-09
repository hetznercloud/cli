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
		Use:                   "create [options] <name>",
		Short:                 "Create a new context",
		Args:                  util.Validate,
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE:                  state.Wrap(s, runCreate),
	}
	cmd.Flags().Bool("token-from-env", false, "If true, the HCLOUD_TOKEN from the environment will be used without asking")
	cmd.Flags().String("credential-helper", "", "Credential helper executable (uses the Git credential protocol)")
	cmd.Flags().StringArray("credential-helper-arg", nil, "Argument to pass to the credential helper (repeatable)")
	cmd.Flags().String("credential-url", "", "URL identifying the credential to retrieve")
	return cmd
}

func runCreate(s state.State, cmd *cobra.Command, args []string) error {
	tokenFromEnv, _ := cmd.Flags().GetBool("token-from-env")
	helperConfigured := cmd.Flags().Changed("credential-helper") || cmd.Flags().Changed("credential-helper-arg") || cmd.Flags().Changed("credential-url")
	if helperConfigured && tokenFromEnv {
		return errors.New("credential helper settings and --token-from-env are mutually exclusive")
	}

	cfg := s.Config()
	if !helperConfigured && !s.Terminal().StdoutIsTerminal() && !tokenFromEnv {
		return errors.New("non-interactive tty detected. Use --token-from-env to use HCLOUD_TOKEN from the environment")
	}

	name := strings.TrimSpace(args[0])
	if name == "" {
		return errors.New("invalid name")
	}
	if config.ContextByName(cfg, name) != nil {
		return errors.New("name already used")
	}
	if helperConfigured {
		helper, _ := cmd.Flags().GetString("credential-helper")
		helperArgs, _ := cmd.Flags().GetStringArray("credential-helper-arg")
		credentialURL, _ := cmd.Flags().GetString("credential-url")
		ctx, err := config.NewCredentialContext(name, append([]string{helper}, helperArgs...), credentialURL)
		if err != nil {
			return err
		}
		return saveCreatedContext(s, cmd, ctx)
	}

	var token string

	envToken := os.Getenv("HCLOUD_TOKEN")
	if envToken != "" {
		switch {
		case len(envToken) != 64:
			if tokenFromEnv {
				return errors.New("invalid token: must be 64 characters in length")
			}
			cmd.Println("Warning: HCLOUD_TOKEN is set, but token is invalid (must be exactly 64 characters long)")
		case tokenFromEnv:
			token = envToken
		default:
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

	return saveCreatedContext(s, cmd, config.NewContext(name, token))
}

func saveCreatedContext(s state.State, cmd *cobra.Command, ctx config.Context) error {
	cfg := s.Config()
	cfg.SetContexts(append(cfg.Contexts(), ctx))
	cfg.SetActiveContext(ctx)

	if err := cfg.Write(nil); err != nil {
		return err
	}

	cmd.Printf("Context %s created and activated\n", ctx.Name())

	return nil
}
