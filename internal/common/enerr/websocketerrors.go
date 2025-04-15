package enerr

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"

	"github.com/anrew1002/Tournament-ChemLoto/internal/sl"
)

func ErrorResponse(unicast models.UnicastFunction, user string, log *slog.Logger, err error) {
	if err == nil {
		// nilErrorResponse(w, lgr)
		return
	}
	var e *EngineError
	if errors.As(err, &e) {
		switch e.Kind {
		default:
			engineErrorResponse(unicast, user, log, e)
			return
		}
	}

	// unknownErrorResponse(unicast, user, log, err)
}

// ErrorResponse unicast message with ENGINE_ACTION type to user and log error.
func engineErrorResponse(unicast models.UnicastFunction, user string, log *slog.Logger, e *EngineError) {
	const op Op = "errs/typicalErrorResponse"

	// Error should not be empty, but it's
	// theoretically possible, so this is just in case...
	if e.isZero() {
		log.Error(fmt.Sprintf("error sent to %s, but empty - very strange, investigate", op))
		// return
	}

	// typical errors
	const errMsg = "error response sent to client"

	ops := OpStack(e)
	if len(ops) > 0 {
		// j, _ := json.Marshal(ops)
		log.Error(errMsg, slog.Any("stack", ops),
			slog.String("Kind", e.Kind.String()),
			slog.String("Parameter", string(e.Param)),
			slog.String("User", string(e.User)),
			sl.Err(e.Err))
	} else {
		log.Error(errMsg,
			slog.String("Kind", e.Kind.String()),
			slog.String("Parameter", string(e.Param)),
			slog.String("User", string(e.User)),
			sl.Err(e.Err))
	}

	// get ErrResponse
	er := newErrResponse(e)
	// unicast it to user
	unicast(user, er)

}

func newErrResponse(err *EngineError) common.Message {
	var msg = []string{"Ошибка сервера", "Cообщите о ней организатору"}

	switch err.Kind {
	case Internal, Database:
		return common.Message{
			Type:   common.ENGINE_ACTION,
			Ok:     false,
			Errors: msg,
		}
	default:
		return common.Message{
			Type:   common.ENGINE_ACTION,
			Ok:     false,
			Errors: []string{err.Error()},
		}
	}
}
