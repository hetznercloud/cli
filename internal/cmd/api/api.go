package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
)

var APICmd = base.Cmd{
	BaseCobraCommand: func(_ hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:   "api [options] <path>",
			Short: "Make API call",
		}

		cmd.Flags().StringP("method", "X", "GET", "HTTP method to use for API call")
		_ = cmd.RegisterFlagCompletionFunc("method", cmpl.SuggestCandidates("GET", "POST", "PUT", "DELETE"))

		cmd.Flags().StringP("data", "d", "", "HTTP request body content (use - to read from stdin)")
		cmd.Flags().StringToStringP("value", "v", nil, "key=value pairs to pass as query parameters")
		return cmd
	},
	Run: func(s state.State, cmd *cobra.Command, args []string) error {
		path := args[0]
		method, _ := cmd.Flags().GetString("method")
		data, _ := cmd.Flags().GetString("data")
		values, _ := cmd.Flags().GetStringToString("value")

		method = strings.ToUpper(method)
		switch method {
		case http.MethodGet, http.MethodHead, http.MethodPost,
			http.MethodPut, http.MethodPatch, http.MethodDelete,
			http.MethodConnect, http.MethodOptions, http.MethodTrace:
			break
		default:
			return fmt.Errorf("unknown HTTP method: %s", method)
		}

		var body io.Reader
		switch data {
		case "":
			body = nil
		case "-":
			body = os.Stdin
		default:
			body = strings.NewReader(data)
		}

		if !strings.HasPrefix(path, "/") {
			path = "/" + path
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

		var respBuf bytes.Buffer
		_, err = s.Client().Do(req, &respBuf)
		if err != nil {
			return err
		}
		respBody := respBuf.Bytes()

		if json.Valid(respBody) {
			var formatted bytes.Buffer
			if err := json.Indent(&formatted, respBody, "", "  "); err != nil {
				return err
			}
			cmd.Println(formatted.String())
		} else if len(respBody) > 0 {
			cmd.Println(string(respBody))
		}
		return nil
	},
}
