package api_rpc

import (
	"context"

	"google.golang.org/grpc/peer"

	"google.golang.org/grpc/metadata"
)

const (
	userAgentHeader            = "user-agent"
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	md := &Metadata{}

	// context metadata
	if ctxmd, ok := metadata.FromIncomingContext(ctx); ok {

		if userAgents := ctxmd.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}

		if userAgents := ctxmd.Get(userAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}

		if ClientIPs := ctxmd.Get(xForwardedForHeader); len(ClientIPs) > 0 {
			md.ClientIP = ClientIPs[0]
		}
	}

	// grpc sub package peer
	if p, ok := peer.FromContext(ctx); ok {
		md.ClientIP = p.Addr.String()
	}

	return md
}
