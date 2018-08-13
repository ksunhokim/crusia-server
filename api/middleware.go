package api

import (
	"context"
	"net/http"

	"github.com/sdbx/crusia-server/store"
	"github.com/sdbx/crusia-server/utils"
)

var (
	userCtxKey = contextKey{"user"}
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "sdbx/crusia-server/api context value " + k.name
}

func getUser(r *http.Request) *store.User {
	entry, _ := r.Context().Value(userCtxKey).(*store.User)
	return entry
}

func withUser(r *http.Request, u *store.User) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), userCtxKey, u))
	return r
}

func (a *Api) UserMiddleWare(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		tok := r.Header.Get("Authorization")
		id, err := a.in.GetToken(tok)
		if err != nil {
			utils.HttpError(w, err, 403)
			return
		}

		u, err := a.in.GetStore().GetUser(id)
		if err != nil {
			utils.HttpError(w, err, 500)
			return
		}

		next.ServeHTTP(w, withUser(r, u))
	}
	return http.HandlerFunc(fn)
}
