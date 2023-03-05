package rest

import (
	"context"
	"crud_app/internal/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type Handler struct {
	usersService User
}

func NewHandler(users User) *Handler {
	return &Handler{
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodGet)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodPost)
	}

	return r
}
