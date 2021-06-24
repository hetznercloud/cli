// Code generated by interfacer; DO NOT EDIT

package hcapi2

import (
	"context"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

// ServerClientBase is an interface generated for "github.com/hetznercloud/hcloud-go/hcloud.ServerClient".
type ServerClientBase interface {
	All(context.Context) ([]*hcloud.Server, error)
	AllWithOpts(context.Context, hcloud.ServerListOpts) ([]*hcloud.Server, error)
	AttachISO(context.Context, *hcloud.Server, *hcloud.ISO) (*hcloud.Action, *hcloud.Response, error)
	AttachToNetwork(context.Context, *hcloud.Server, hcloud.ServerAttachToNetworkOpts) (*hcloud.Action, *hcloud.Response, error)
	ChangeAliasIPs(context.Context, *hcloud.Server, hcloud.ServerChangeAliasIPsOpts) (*hcloud.Action, *hcloud.Response, error)
	ChangeDNSPtr(context.Context, *hcloud.Server, string, *string) (*hcloud.Action, *hcloud.Response, error)
	ChangeProtection(context.Context, *hcloud.Server, hcloud.ServerChangeProtectionOpts) (*hcloud.Action, *hcloud.Response, error)
	ChangeType(context.Context, *hcloud.Server, hcloud.ServerChangeTypeOpts) (*hcloud.Action, *hcloud.Response, error)
	Create(context.Context, hcloud.ServerCreateOpts) (hcloud.ServerCreateResult, *hcloud.Response, error)
	CreateImage(context.Context, *hcloud.Server, *hcloud.ServerCreateImageOpts) (hcloud.ServerCreateImageResult, *hcloud.Response, error)
	Delete(context.Context, *hcloud.Server) (*hcloud.Response, error)
	DetachFromNetwork(context.Context, *hcloud.Server, hcloud.ServerDetachFromNetworkOpts) (*hcloud.Action, *hcloud.Response, error)
	DetachISO(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	DisableBackup(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	DisableRescue(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	EnableBackup(context.Context, *hcloud.Server, string) (*hcloud.Action, *hcloud.Response, error)
	EnableRescue(context.Context, *hcloud.Server, hcloud.ServerEnableRescueOpts) (hcloud.ServerEnableRescueResult, *hcloud.Response, error)
	Get(context.Context, string) (*hcloud.Server, *hcloud.Response, error)
	GetByID(context.Context, int) (*hcloud.Server, *hcloud.Response, error)
	GetByName(context.Context, string) (*hcloud.Server, *hcloud.Response, error)
	GetMetrics(context.Context, *hcloud.Server, hcloud.ServerGetMetricsOpts) (*hcloud.ServerMetrics, *hcloud.Response, error)
	List(context.Context, hcloud.ServerListOpts) ([]*hcloud.Server, *hcloud.Response, error)
	Poweroff(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Poweron(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Reboot(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Rebuild(context.Context, *hcloud.Server, hcloud.ServerRebuildOpts) (*hcloud.Action, *hcloud.Response, error)
	RequestConsole(context.Context, *hcloud.Server) (hcloud.ServerRequestConsoleResult, *hcloud.Response, error)
	Reset(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	ResetPassword(context.Context, *hcloud.Server) (hcloud.ServerResetPasswordResult, *hcloud.Response, error)
	Shutdown(context.Context, *hcloud.Server) (*hcloud.Action, *hcloud.Response, error)
	Update(context.Context, *hcloud.Server, hcloud.ServerUpdateOpts) (*hcloud.Server, *hcloud.Response, error)
}
