//The API supports idempotency for safely retrying requests without accidentally performing
// the same operation twice. This is useful when an API call is disrupted in transit and you
// do not receive a response. For example, if a request to create something and it does not respond
// due to a network connection error, you can retry the request with the same idempotency key
// to guarantee that no more than one charge is created.
package util

import (
	"context"
)

//To perform an idempotent request, provide an additional Idempotency-Key to the context request.
const IdempotencyKey = "Idempotency-Key"

// Check to see if contains a IdempotencyKey value in the context and define if it is a Idempotency Request
func isIdempotencyRequest(ctx context.Context) bool {
	if v := ctx.Value(IdempotencyKey); v != nil {
		return true
	}
	return false
}

func WithIdempotencyRequest(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, IdempotencyKey, key)
}
