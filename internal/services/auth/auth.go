package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/q2rd/gRPC_sso_go/internal/custom_logger/sl"
	"github.com/q2rd/gRPC_sso_go/internal/domain/models"
	"github.com/q2rd/gRPC_sso_go/internal/lib/customjwt"
	"github.com/q2rd/gRPC_sso_go/internal/storage"

	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uuid string, err error)
}

type UserProvider interface {
	User(ctc context.Context, email string) (models.UserDatabase, error)
	IsAdmin(ctx context.Context, userId string) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

func NewAuth(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appId int) (string, error) {
	const op = "services.auth.Login"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("Attempt login" + email)
	user, err := a.usrProvider.User(ctx, email)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("app not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to fetch user ", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			lo
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user logged in successfully")
	token, err := customjwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email, password string) (string, error) {
	// mb confirmPassword also needs
	const op = "services.auth.RegisterNewUSer"
	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("user registration")
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {

		log.Error("generation failed", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	uuid, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}
		log.Error("failed to save user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User successfully registered.")
	return uuid, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userId string) (bool, error) {
	const op = "services.auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.String("userId", userId),
	)
	isAdmin, err := a.usrProvider.IsAdmin(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User admin status checked: ", slog.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}
