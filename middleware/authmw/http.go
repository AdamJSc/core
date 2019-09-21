package authmw

import (
	"net/http"
	"strings"

	"github.com/LUSHDigital/core/auth"
	"github.com/LUSHDigital/core/rest"
)

const (
	authHeader               = "Authorization"
	authHeaderPrefix         = "Bearer "
	msgMissingRequiredGrants = "missing required grants"
)

// HandlerValidateJWT takes a JWT from the request headers, attempts validation and returns a http handler.
func HandlerValidateJWT(brk auth.RSAPublicKeyCopierRenewer, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimPrefix(r.Header.Get(authHeader), authHeaderPrefix)
		if raw == "" {
			rest.Response{
				Code:    http.StatusUnauthorized,
				Message: "missing token",
			}.WriteTo(w)
			return
		}
		pk := brk.Copy()
		parser := auth.NewParser(&pk)
		claims, err := parser.Claims(raw)
		if err != nil {
			switch err.(type) {
			case auth.TokenSignatureError:
				brk.Renew() // Renew the public key if there's an error validating the token signature
			}
			rest.Response{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			}.WriteTo(w)
			return
		}
		ctx := auth.ContextWithConsumer(r.Context(), claims.Consumer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// HandlerGrants is an HTTP handler to check that the consumer in the request context has the required grants.
func HandlerGrants(grants []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumer := auth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyGrant(grants...) {
			rest.Response{
				Code:    http.StatusUnauthorized,
				Message: msgMissingRequiredGrants,
			}.WriteTo(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// HandlerRoles is an HTTP handler to check that the consumer in the request context has the required roles.
func HandlerRoles(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		consumer := auth.ConsumerFromContext(r.Context())
		if !consumer.HasAnyRole(roles...) {
			rest.Response{
				Code:    http.StatusUnauthorized,
				Message: msgMissingRequiredGrants,
			}.WriteTo(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
