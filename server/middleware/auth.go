// Copyright (c) The go-grpc-middleware Authors.
// Licensed under the Apache License 2.0.

package middleware

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"

	"google.golang.org/grpc"
)

// AuthFn is the pluggable function that performs authentication.
//
// The passed in `Context` will contain the gRPC metadata.MD object (for header-based authentication) and
// the peer.Peer information that can contain transport-based credentials (e.g. `credentials.AuthInfo`).
//
// The returned context will be propagated to handlers, allowing user changes to `Context`. However,
// please make sure that the `Context` returned is a child `Context` of the one passed in.
//
// If error is returned, its `grpc.Code()` will be returned to the user as well as the verbatim message.
// Please make sure you use `codes.Unauthenticated` (lacking auth) and `codes.PermissionDenied`
// (authed, but lacking perms) appropriately.
type AuthFn func(ctx context.Context) (context.Context, error)

// ServiceAuthFuncOverride allows a given gRPC service implementation to override the global `AuthFunc`.
//
// If a service implements the AuthFuncOverride method, it takes precedence over the `AuthFunc` method,
// and will be called instead of AuthFunc for all method invocations within that service.
type ServiceAuthFuncOverride interface {
	AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error)
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func UnaryServerInterceptor(f AuthFn) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		var newCtx context.Context
		var err error
		if overrideSrv, ok := info.Server.(ServiceAuthFuncOverride); ok {
			newCtx, err = overrideSrv.AuthFuncOverride(ctx, info.FullMethod)
		} else {
			newCtx, err = f(ctx)
		}
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func AuthFromMetadata(ctx context.Context, expectedScheme string, header string) (string, error) {
	headerVal := metadata.ExtractIncoming(ctx).Get(header)
	if headerVal == "" {
		return "", status.Error(codes.Unauthenticated, "Request is not authenticated properly")
	}
	scheme, token, found := strings.Cut(headerVal, " ")
	if !found {
		return "", status.Errorf(codes.Unauthenticated, "The %s header was found, but it's value is malformed", header)
	}
	if !strings.EqualFold(expectedScheme, scheme) {
		return "", status.Errorf(codes.Unauthenticated, "The %s header was found, but it's value is malformed", header)
	}

	if token == "" {
		return "", status.Errorf(codes.Unauthenticated, "The token in the header is empty")
	}
	return token, nil
}
