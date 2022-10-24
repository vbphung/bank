package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type metadataRes struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *metadataRes {
	mdRes := &metadataRes{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if usrAgs := md.Get("grpcgateway-user-agent"); len(usrAgs) > 0 {
			mdRes.UserAgent = usrAgs[0]
		} else if usrAgs = md.Get("user-agent"); len(usrAgs) > 0 {
			mdRes.UserAgent = usrAgs[0]
		}

		if ips := md.Get("x-forwarded-for"); len(ips) > 0 {
			mdRes.ClientIp = ips[0]
		} else if p, ok := peer.FromContext(ctx); ok {
			mdRes.ClientIp = p.Addr.String()
		}
	}

	return mdRes
}
