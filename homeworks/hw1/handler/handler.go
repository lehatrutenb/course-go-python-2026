package handler

import (
	"context"
	"errors"

	errs "github.com/lehatrutenb/course-go-python-2026/homeworks/hw1/handler/errors"
)

type Handler struct {
	db UsersDB
}

func NewHandler(db UsersDB) *Handler {
	return &Handler{db: db}
}

// Believe that README writes error checks in expected order
func (h *Handler) UpdateUserBalance(ctx context.Context, userID, balance int64) error {
	err := h.db.UpdateBalance(ctx, userID, balance)
	i := 0
	for retryErr, ok := errors.AsType[*errs.RetryableError](err); ok && i < retryErr.RetryCount(); i++ {
		errInner := h.db.UpdateBalance(ctx, userID, balance)
		if !errors.As(errInner, new(*errs.RetryableError)) {
			err = errInner
			break
		}
	}
	if errors.As(err, new(*errs.RetryableError)) { // TODO some misdread with README and no tests, but i suppose that variant is expected
		return err
	}
	if errors.As(err, new(*errs.NotFoundError)) {
		return errs.NewAdditionalMessageError(err, "user not found")
	}
	if errors.As(err, new(*errs.AdditionalMessageError)) {
		return err
	}
	if err != nil {
		panic(err)
	}

	return err
}
