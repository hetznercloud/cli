package loadbalancer

import (
	"context"
	"fmt"
	"time"

	"github.com/hetznercloud/cli/internal/cmd/base"
	"github.com/hetznercloud/cli/internal/cmd/cmpl"
	"github.com/hetznercloud/cli/internal/hcapi2"
	"github.com/hetznercloud/cli/internal/state"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

var UpdateServiceCommand = base.Cmd{
	BaseCobraCommand: func(client hcapi2.Client) *cobra.Command {
		cmd := &cobra.Command{
			Use:                   "update-service LOADBALANCER FLAGS",
			Short:                 "Updates a service from a Load Balancer",
			Args:                  cobra.ExactArgs(1),
			ValidArgsFunction:     cmpl.SuggestArgs(cmpl.SuggestCandidatesF(client.LoadBalancer().Names)),
			TraverseChildren:      true,
			DisableFlagsInUseLine: true,
		}

		cmd.Flags().Int("listen-port", 0, "The listen port of the service that you want to update (required)")
		cmd.MarkFlagRequired("listen-port")

		cmd.Flags().Int("destination-port", 0, "Destination port of the service on the targets")

		cmd.Flags().String("protocol", "", "The protocol the health check is performed over")
		cmd.Flags().Bool("proxy-protocol", false, "Enable or disable (with --proxy-protocol=false) Proxy Protocol")
		cmd.Flags().Bool("http-redirect-http", false, "Enable or disable redirect all traffic on port 80 to port 443")

		cmd.Flags().Bool("http-sticky-sessions", false, "Enable or disable (with --http-sticky-sessions=false) Sticky Sessions")
		cmd.Flags().String("http-cookie-name", "", "Sticky Sessions: Cookie Name which will be set")
		cmd.Flags().Duration("http-cookie-lifetime", 0, "Sticky Sessions: Lifetime of the cookie")
		cmd.Flags().IntSlice("http-certificates", []int{}, "ID of Certificates which are attached to this Load Balancer")

		cmd.Flags().String("health-check-protocol", "", "The protocol the health check is performed over")
		cmd.Flags().Int("health-check-port", 0, "The port the health check is performed over")
		cmd.Flags().Duration("health-check-interval", 15*time.Second, "The interval the health check is performed")
		cmd.Flags().Duration("health-check-timeout", 10*time.Second, "The timeout after a health check is marked as failed")
		cmd.Flags().Int("health-check-retries", 3, "Number of retries after a health check is marked as failed")

		cmd.Flags().String("health-check-http-domain", "", "The domain we request when performing a http health check")
		cmd.Flags().String("health-check-http-path", "", "The path we request when performing a http health check")

		cmd.Flags().StringSlice("health-check-http-status-codes", []string{}, "List of status codes we expect to determine a target as healthy")
		cmd.Flags().String("health-check-http-response", "", "The response we expect to determine a target as healthy")
		cmd.Flags().Bool("health-check-http-tls", false, "Determine if the health check should verify if the target answers with a valid TLS certificate")
		return cmd
	},
	Run: func(ctx context.Context, client hcapi2.Client, waiter state.ActionWaiter, cmd *cobra.Command, args []string) error {
		idOrName := args[0]
		listenPort, _ := cmd.Flags().GetInt("listen-port")

		loadBalancer, _, err := client.LoadBalancer().Get(ctx, idOrName)
		if err != nil {
			return err
		}
		if loadBalancer == nil {
			return fmt.Errorf("Load Balancer not found: %s", idOrName)
		}
		var service hcloud.LoadBalancerService
		for _, _service := range loadBalancer.Services {
			if _service.ListenPort == listenPort {
				if _service.HealthCheck.HTTP != nil {
					service = _service
				}
			}
		}
		opts := hcloud.LoadBalancerUpdateServiceOpts{
			HTTP:        &hcloud.LoadBalancerUpdateServiceOptsHTTP{},
			HealthCheck: &hcloud.LoadBalancerUpdateServiceOptsHealthCheck{},
		}
		if cmd.Flag("protocol").Changed {
			protocol, _ := cmd.Flags().GetString("protocol")
			opts.Protocol = hcloud.LoadBalancerServiceProtocol(protocol)
		}
		if cmd.Flag("destination-port").Changed {
			destinationPort, _ := cmd.Flags().GetInt("destination-port")
			opts.DestinationPort = &destinationPort
		}
		if cmd.Flag("proxy-protocol").Changed {
			proxyProtocol, _ := cmd.Flags().GetBool("proxy-protocol")
			opts.Proxyprotocol = hcloud.Bool(proxyProtocol)
		}
		// HTTP
		if cmd.Flag("http-redirect-http").Changed {
			redirectHTTP, _ := cmd.Flags().GetBool("http-redirect-http")
			opts.HTTP.RedirectHTTP = hcloud.Bool(redirectHTTP)
		}
		if cmd.Flag("http-sticky-sessions").Changed {
			stickySessions, _ := cmd.Flags().GetBool("http-sticky-sessions")
			opts.HTTP.StickySessions = hcloud.Bool(stickySessions)
		}
		if cmd.Flag("http-cookie-name").Changed {
			cookieName, _ := cmd.Flags().GetString("http-cookie-name")
			opts.HTTP.CookieName = hcloud.String(cookieName)
		}
		if cmd.Flag("http-cookie-lifetime").Changed {
			cookieLifetime, _ := cmd.Flags().GetDuration("http-cookie-lifetime")
			opts.HTTP.CookieLifetime = hcloud.Duration(cookieLifetime)
		}
		if cmd.Flag("http-certificates").Changed {
			certificates, _ := cmd.Flags().GetIntSlice("http-certificates")
			for _, certificateID := range certificates {
				opts.HTTP.Certificates = append(opts.HTTP.Certificates, &hcloud.Certificate{ID: certificateID})
			}
		}
		// Health Check
		if cmd.Flag("health-check-protocol").Changed {
			healthCheckProtocol, _ := cmd.Flags().GetString("health-check-protocol")
			opts.HealthCheck.Protocol = hcloud.LoadBalancerServiceProtocol(healthCheckProtocol)
		}
		if cmd.Flag("health-check-port").Changed {
			healthCheckPort, _ := cmd.Flags().GetInt("health-check-port")
			opts.HealthCheck.Port = hcloud.Int(healthCheckPort)
		}
		if cmd.Flag("health-check-interval").Changed {
			healthCheckInterval, _ := cmd.Flags().GetDuration("health-check-interval")
			opts.HealthCheck.Interval = hcloud.Duration(healthCheckInterval)
		}
		if cmd.Flag("health-check-timeout").Changed {
			healthCheckTimeout, _ := cmd.Flags().GetDuration("health-check-timeout")
			opts.HealthCheck.Timeout = hcloud.Duration(healthCheckTimeout)
		}
		if cmd.Flag("health-check-retries").Changed {
			healthCheckRetries, _ := cmd.Flags().GetInt("health-check-retries")
			opts.HealthCheck.Retries = hcloud.Int(healthCheckRetries)
		}

		// Health Check HTTP
		healthCheckProtocol, _ := cmd.Flags().GetString("health-check-protocol")
		if healthCheckProtocol != string(hcloud.LoadBalancerServiceProtocolTCP) || service.HealthCheck.Protocol != hcloud.LoadBalancerServiceProtocolTCP {
			opts.HealthCheck.HTTP = &hcloud.LoadBalancerUpdateServiceOptsHealthCheckHTTP{}

			if cmd.Flag("health-check-http-domain").Changed {
				healthCheckHTTPDomain, _ := cmd.Flags().GetString("health-check-http-domain")
				opts.HealthCheck.HTTP.Domain = hcloud.String(healthCheckHTTPDomain)
			}
			if cmd.Flag("health-check-http-path").Changed {
				healthCheckHTTPPath, _ := cmd.Flags().GetString("health-check-http-path")
				opts.HealthCheck.HTTP.Path = hcloud.String(healthCheckHTTPPath)
			}
			if cmd.Flag("health-check-http-response").Changed {
				healthCheckHTTPResponse, _ := cmd.Flags().GetString("health-check-http-response")
				opts.HealthCheck.HTTP.Response = hcloud.String(healthCheckHTTPResponse)
			}
			if cmd.Flag("health-check-http-status-codes").Changed {
				healthCheckHTTPStatusCodes, _ := cmd.Flags().GetStringSlice("health-check-http-status-codes")
				opts.HealthCheck.HTTP.StatusCodes = healthCheckHTTPStatusCodes
			}
			if cmd.Flag("health-check-http-tls").Changed {
				healthCheckHTTPTLS, _ := cmd.Flags().GetBool("health-check-http-tls")
				opts.HealthCheck.HTTP.TLS = hcloud.Bool(healthCheckHTTPTLS)
			}
		}

		action, _, err := client.LoadBalancer().UpdateService(ctx, loadBalancer, listenPort, opts)
		if err != nil {
			return err
		}
		if err := waiter.ActionProgress(ctx, action); err != nil {
			return err
		}
		fmt.Printf("Service %d on Load Balancer %d was updated\n", listenPort, loadBalancer.ID)

		return nil
	},
}
