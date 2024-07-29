package Minecraft

import (
	"context"
	"net"
)

type Resolver struct {
	internalResolver *net.Resolver
}

func NewResolver() (resolver *Resolver) {
	return &Resolver{
		&net.Resolver{
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, network, "1.1.1.1:53")
			},
		},
	}
}

func (resolver *Resolver) SetInternalResolver(internalResolver *net.Resolver) {
	resolver.internalResolver = internalResolver
}

func (resolver *Resolver) SRVLookup(host string) (res bool, tHost string, tPort uint16) {
	_, srvs, err := resolver.internalResolver.LookupSRV(context.Background(), "minecraft", "tcp", host)
	if err != nil || len(srvs) == 0 {
		return
	}
	return true, srvs[0].Target[:len(srvs[0].Target)-1], srvs[0].Port
}
