package loadbalancer

import (
	"fmt"

	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/cmd/util"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func newAddServiceCommand(cli *state.State) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add-service LOADBALANCER FLAGS",
		Short:                 "Add a service from a Load Balancer",
		Args:                  cobra.ExactArgs(1),
		ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(cli.LoadBalancerNames)),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               util.ChainRunE(validateAddService, cli.EnsureToken),
		RunE:                  cli.Wrap(runAddService),
	}
	cmd.Flags().String("protocol", "", "Protocol of the service (required)")
	cmd.MarkFlagRequired("protocol")

	cmd.Flags().Int("listen-port", 0, "Listen port of the service")
	cmd.Flags().Int("destination-port", 0, "Destination port of the service on the targets")
	cmd.Flags().Bool("proxy-protocol", false, "Enable proxyprotocol")

	cmd.Flags().Bool("http-sticky-sessions", false, "Enable Sticky Sessions")
	cmd.Flags().String("http-cookie-name", "", "Sticky Sessions: Cookie Name we set")
	cmd.Flags().Duration("http-cookie-lifetime", 0, "Sticky Sessions: Lifetime of the cookie")
	cmd.Flags().IntSlice("http-certificates", []int{}, "ID of Certificates which are attached to this Load Balancer")
	cmd.Flags().Bool("http-redirect-https", false, "Redirect all traffic on port 80 to port 443")

	return cmd
}

func validateAddService(cmd *cobra.Command, args []string) error {
	protocol, _ := cmd.Flags().GetString("protocol")
	listenPort, _ := cmd.Flags().GetInt("listen-port")
	destinationPort, _ := cmd.Flags().GetInt("destination-port")
	httpCertificates, _ := cmd.Flags().GetIntSlice("http-certificates")

	if protocol == "" {
		return fmt.Errorf("required flag protocol not set")
	}

	switch hcloud.LoadBalancerServiceProtocol(protocol) {
	case hcloud.LoadBalancerServiceProtocolHTTP:
		break
	case hcloud.LoadBalancerServiceProtocolTCP:
		if listenPort == 0 {
			return fmt.Errorf("please specify a listen port")
		}

		if destinationPort == 0 {
			return fmt.Errorf("please specify a destination port")
		}
		break
	case hcloud.LoadBalancerServiceProtocolHTTPS:
		if len(httpCertificates) == 0 {
			return fmt.Errorf("no certificate specified")
		}
	default:
		return fmt.Errorf("%s is not a valid protocol", protocol)
	}
	if listenPort > 65535 {
		return fmt.Errorf("%d is not a valid listen port", listenPort)
	}

	if destinationPort > 65535 {
		return fmt.Errorf("%d is not a valid destination port", destinationPort)
	}
	return nil
}

func runAddService(cli *state.State, cmd *cobra.Command, args []string) error {
	idOrName := args[0]
	protocol, _ := cmd.Flags().GetString("protocol")
	listenPort, _ := cmd.Flags().GetInt("listen-port")
	destinationPort, _ := cmd.Flags().GetInt("destination-port")
	proxyProtocol, _ := cmd.Flags().GetBool("proxy-protocol")

	httpStickySessions, _ := cmd.Flags().GetBool("http-sticky-sessions")
	httpCookieName, _ := cmd.Flags().GetString("http-cookie-name")
	httpCookieLifetime, _ := cmd.Flags().GetDuration("http-cookie-lifetime")
	httpCertificates, _ := cmd.Flags().GetIntSlice("http-certificates")
	httpRedirect, _ := cmd.Flags().GetBool("http-redirect-https")

	loadBalancer, _, err := cli.Client().LoadBalancer.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if loadBalancer == nil {
		return fmt.Errorf("Load Balancer not found: %s", idOrName)
	}

	opts := hcloud.LoadBalancerAddServiceOpts{
		Protocol:      hcloud.LoadBalancerServiceProtocol(protocol),
		Proxyprotocol: hcloud.Bool(proxyProtocol),
	}

	if listenPort != 0 {
		opts.ListenPort = hcloud.Int(listenPort)
	}
	if destinationPort != 0 {
		opts.DestinationPort = hcloud.Int(destinationPort)
	}

	if protocol != string(hcloud.LoadBalancerServiceProtocolTCP) {
		opts.HTTP = &hcloud.LoadBalancerAddServiceOptsHTTP{
			StickySessions: hcloud.Bool(httpStickySessions),
			RedirectHTTP:   hcloud.Bool(httpRedirect),
		}
		if httpCookieName != "" {
			opts.HTTP.CookieName = hcloud.String(httpCookieName)
		}
		if httpCookieLifetime != 0 {
			opts.HTTP.CookieLifetime = hcloud.Duration(httpCookieLifetime)
		}
		for _, certificateID := range httpCertificates {
			opts.HTTP.Certificates = append(opts.HTTP.Certificates, &hcloud.Certificate{ID: certificateID})
		}
	}
	action, _, err := cli.Client().LoadBalancer.AddService(cli.Context, loadBalancer, opts)
	if err != nil {
		return err
	}
	if err := cli.ActionProgress(cli.Context, action); err != nil {
		return err
	}
	fmt.Printf("Service was added to Load Balancer %d\n", loadBalancer.ID)

	return nil
}
