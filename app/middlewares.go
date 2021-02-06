package app

import (
	"context"
	kitHttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

func validateMiddleware(next kitHttp.DecodeRequestFunc, validate func(interface{}) error) kitHttp.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (request interface{}, err error) {
		request, err = next(ctx, req)
		if err != nil {
			return nil, err
		}
		if err = validate(request); err != nil {
			return nil, err
		}
		return request, nil
	}
}
