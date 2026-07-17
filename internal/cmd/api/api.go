package api

import (
	_ "embed"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var ApiCmd = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "api [options] <path>",
			Short: "Make API call",
		}

		cmd.Flags().StringP("method", "X", "GET", "HTTP method to use for API call")
		cmd.Flags().StringP("data", "d", "", "HTTP request body content (use - to read from stdin)")
		cmd.Flags().StringToStringP("value", "v", nil, "key=value pairs to pass as query parameters")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		path := args[0]
		method, _ := cmd.Flags().GetString("method")
		data, _ := cmd.Flags().GetString("data")
		values, _ := cmd.Flags().GetStringToString("value")

		var body io.Reader
		switch data {
		case "":
			body = nil
		case "-":
			body = os.Stdin
		default:
			body = strings.NewReader(data)
		}

		u, err := url.Parse(path)
		if err != nil {
			return err
		}

		params := u.Query()
		for k, v := range values {
			params.Set(k, v)
		}
		u.RawQuery = params.Encode()

		req, err := s.Client().NewRequest(s, method, u.String(), body)
		if err != nil {
			return err
		}

		var respBody json.RawMessage
		_, err = s.Client().Do(req, &respBody)
		if err != nil {
			return err
		}

		cmd.Print(string(respBody))
		return nil
	},
}
