package hcloud

import "github.com/hetznercloud/hcloud-go/hcloud/schema"

// This file provides converter functions to convert models in the
// schema package to models in the hcloud package.

// ActionFromSchema converts a schema.Action to an Action.
func ActionFromSchema(s schema.Action) *Action {
	action := &Action{
		ID:       s.ID,
		Status:   s.Status,
		Command:  s.Command,
		Progress: s.Progress,
		Started:  s.Started,
		Finished: s.Finished,
	}
	if s.Error != nil {
		action.ErrorCode = s.Error.Code
		action.ErrorMessage = s.Error.Message
	}
	return action
}

// FloatingIPFromSchema converts a schema.FloatingIP to a FloatingIP.
func FloatingIPFromSchema(s schema.FloatingIP) *FloatingIP {
	f := &FloatingIP{
		ID:           s.ID,
		IP:           s.IP,
		Type:         FloatingIPType(s.Type),
		HomeLocation: LocationFromSchema(s.HomeLocation),
	}
	if s.Description != nil {
		f.Description = *s.Description
	}
	if s.Server != nil {
		f.Server = &Server{ID: *s.Server}
	}
	if s.DNSPtr.IPv4 != nil {
		f.DNSPtr = map[string]string{
			s.IP: *s.DNSPtr.IPv4,
		}
	} else if s.DNSPtr.IPv6 != nil {
		f.DNSPtr = map[string]string{}
		for _, entry := range s.DNSPtr.IPv6 {
			f.DNSPtr[entry.IP] = entry.DNSPtr
		}
	}
	return f
}

// ISOFromSchema converts a schema.ISO to an ISO.
func ISOFromSchema(s schema.ISO) *ISO {
	return &ISO{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Type:        ISOType(s.Type),
	}
}

// LocationFromSchema converts a schema.Location to a Location.
func LocationFromSchema(s schema.Location) *Location {
	return &Location{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Country:     s.Country,
		City:        s.City,
		Latitude:    s.Latitude,
		Longitude:   s.Longitude,
	}
}

// ServerFromSchema converts a schema.Server to a Server.
func ServerFromSchema(s schema.Server) *Server {
	server := &Server{
		ID:              s.ID,
		Name:            s.Name,
		Status:          ServerStatus(s.Status),
		Created:         s.Created,
		PublicNet:       ServerPublicNetFromSchema(s.PublicNet),
		ServerType:      ServerTypeFromSchema(s.ServerType),
		IncludedTraffic: s.IncludedTraffic,
		OutgoingTraffic: s.OutgoingTraffic,
		IngoingTraffic:  s.IngoingTraffic,
		BackupWindow:    s.BackupWindow,
		RescueEnabled:   s.RescueEnabled,
	}
	if s.ISO != nil {
		server.ISO = ISOFromSchema(*s.ISO)
	}
	return server
}

// ServerPublicNetFromSchema converts a schema.ServerPublicNet to a ServerPublicNet.
func ServerPublicNetFromSchema(s schema.ServerPublicNet) ServerPublicNet {
	publicNet := ServerPublicNet{
		IPv4: ServerPublicNetIPv4FromSchema(s.IPv4),
		IPv6: ServerPublicNetIPv6FromSchema(s.IPv6),
	}
	for _, id := range s.FloatingIPs {
		publicNet.FloatingIPs = append(publicNet.FloatingIPs, &FloatingIP{ID: id})
	}
	return publicNet
}

// ServerPublicNetIPv4FromSchema converts a schema.ServerPublicNetIPv4 to
// a ServerPublicNetIPv4.
func ServerPublicNetIPv4FromSchema(s schema.ServerPublicNetIPv4) ServerPublicNetIPv4 {
	return ServerPublicNetIPv4{
		IP:      s.IP,
		Blocked: s.Blocked,
		DNSPtr:  s.DNSPtr,
	}
}

// ServerPublicNetIPv6FromSchema converts a schema.ServerPublicNetIPv6 to
// a ServerPublicNetIPv6.
func ServerPublicNetIPv6FromSchema(s schema.ServerPublicNetIPv6) ServerPublicNetIPv6 {
	ipv6 := ServerPublicNetIPv6{
		IP:      s.IP,
		Blocked: s.Blocked,
	}
	for _, dnsPtr := range s.DNSPtr {
		ipv6.DNSPtr = append(ipv6.DNSPtr, ServerPublicNetIPv6DNSPtrFromSchema(dnsPtr))
	}
	return ipv6
}

// ServerPublicNetIPv6DNSPtrFromSchema converts a schema.ServerPublicNetIPv6DNSPtr
// to a ServerPublicNetIPv6DNSPtr.
func ServerPublicNetIPv6DNSPtrFromSchema(s schema.ServerPublicNetIPv6DNSPtr) ServerPublicNetIPv6DNSPtr {
	return ServerPublicNetIPv6DNSPtr{
		IP:     s.IP,
		DNSPtr: s.DNSPtr,
	}
}

// ServerTypeFromSchema converts a schema.ServerType to a ServerType.
func ServerTypeFromSchema(s schema.ServerType) *ServerType {
	return &ServerType{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Cores:       s.Cores,
		Memory:      s.Memory,
		Disk:        s.Disk,
		StorageType: StorageType(s.StorageType),
	}
}

// SSHKeyFromSchema converts a schema.SSHKey to a SSHKey.
func SSHKeyFromSchema(s schema.SSHKey) *SSHKey {
	return &SSHKey{
		ID:          s.ID,
		Name:        s.Name,
		Fingerprint: s.Fingerprint,
		PublicKey:   s.PublicKey,
	}
}

// PaginationFromSchema converts a schema.MetaPagination to a Pagination.
func PaginationFromSchema(s schema.MetaPagination) Pagination {
	return Pagination{
		Page:         s.Page,
		PerPage:      s.PerPage,
		PreviousPage: s.PreviousPage,
		NextPage:     s.NextPage,
		LastPage:     s.LastPage,
		TotalEntries: s.TotalEntries,
	}
}
