package middleware

import (
	"context"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func Forward(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	w.Header().Set("test", "middleware")
	w.WriteHeader(200)
	return nil
}
