package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"auth/custom/service"
	"auth/custom/util"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var (
	validate = validator.New()
)

type Service interface {
	CreateSession(context.Context, *service.User) (*service.Session, error)
	RefreshSession(context.Context, string) (*service.Session, error)
	DestroySession(context.Context, string)

	CreateUser(context.Context, string, string) (*service.User, error)
	ReadUser(context.Context, uuid.UUID) (*service.User, error)
	ReadUserByCredentials(context.Context, string, string) (*service.User, error)
	ReadUsers(context.Context) []service.User
}

type Controller struct {
	*chi.Mux

	service     Service
	tokenParser TokenParser
	logger      *zap.Logger
}

func NewController(service *service.Service, tokenParser TokenParser, logger *zap.Logger) *Controller {
	c := &Controller{
		service:     service,
		tokenParser: tokenParser,
		logger:      logger,
	}
	c.initRouter()
	return c
}

func (c *Controller) initRouter() {
	r := chi.NewRouter()

	authenticate := Authenticate(c.tokenParser)
	authorizeAdmin := Authorize(service.RoleAdmin)
	authorizeUser := Authorize(service.RoleUser)

	r.Route("/sessions", func(r chi.Router) {
		r.Post("/native", c.CreateSession)
		r.Post("/refresh", c.RefreshSession)
		r.Post("/destroy", c.DestroySession)
	})
	r.Route("/users", func(r chi.Router) {
		r.Post("/", c.CreateUser)
		r.With(authenticate, authorizeAdmin).Get("/", c.ReadUsers)
		r.With(authenticate, authorizeUser).Get("/me", c.ReadLoggedUser)
	})

	c.Mux = r
}

func (c *Controller) CreateSession(w http.ResponseWriter, r *http.Request) {
	createSessionInput := CreateSessionInput{}
	if err := parseRequestBody(r, &createSessionInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.service.ReadUserByCredentials(
		r.Context(),
		createSessionInput.Email,
		createSessionInput.Password,
	)
	if err != nil {
		c.logger.Error("reading user by email", zap.Error(err))
		if errors.Is(err, service.ErrUserNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, err := c.service.CreateSession(r.Context(), user)
	if err != nil {
		c.logger.Error("creating session", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &CreateUserResp{
		User:    FromUser(user),
		Session: FromSession(session),
	}
	writeResponse(c.logger, w, http.StatusCreated, resp)
}

func (c *Controller) RefreshSession(w http.ResponseWriter, r *http.Request) {
	refreshSessionInput := RefreshSessionInput{}
	if err := parseRequestBody(r, &refreshSessionInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := c.service.RefreshSession(r.Context(), refreshSessionInput.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := FromSession(session)
	writeResponse(c.logger, w, http.StatusOK, resp)
}

func (c *Controller) DestroySession(w http.ResponseWriter, r *http.Request) {
	refreshSessionInput := RefreshSessionInput{}
	if err := parseRequestBody(r, &refreshSessionInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.service.DestroySession(r.Context(), refreshSessionInput.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

func (c *Controller) CreateUser(w http.ResponseWriter, r *http.Request) {
	createUserInput := CreateUserInput{}
	if err := parseRequestBody(r, &createUserInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.service.CreateUser(
		r.Context(),
		createUserInput.Email,
		createUserInput.Password,
	)
	if err != nil {
		c.logger.Error("creating user", zap.Error(err))
		if errors.Is(err, service.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	session, err := c.service.CreateSession(r.Context(), user)
	if err != nil {
		c.logger.Error("creating session", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := &CreateUserResp{
		User:    FromUser(user),
		Session: FromSession(session),
	}
	writeResponse(c.logger, w, http.StatusCreated, resp)
}

func (c *Controller) ReadLoggedUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := util.UserIDFromCtx(r.Context())
	user, err := c.service.ReadUser(r.Context(), userID)
	if err != nil {
		c.logger.Error("reading logged user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := FromUser(user)
	writeResponse(c.logger, w, http.StatusOK, resp)
}

func (c *Controller) ReadUsers(w http.ResponseWriter, r *http.Request) {
	users := c.service.ReadUsers(r.Context())
	var resp []User
	for _, u := range users {
		resp = append(resp, FromUser(&u))
	}
	writeResponse(c.logger, w, http.StatusOK, resp)
}

func parseRequestBody(r *http.Request, target any) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return err
	}
	if err := validate.Struct(target); err != nil {
		return err
	}
	return nil
}

func writeResponse(l *zap.Logger, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		l.Error("writing http response", zap.Error(err))
	}
}
