package usecase

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"log/slog"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type LoginRequest struct {
	Name string `json:"name" validate:"required,min=1,safeinput"`
	Code string `json:"code,omitempty"`
}
type LoginResponse struct {
	Token string      `json:"token"`
	Role  common.Role `json:"role"`
	Error []string    `json:"error"`
}

func Login(log *slog.Logger, req LoginRequest, db *sql.DB, code string) (*LoginResponse, error) {
	const op enerr.Op = "usecase.user/Login"
	// TODO: add validation
	// validate := appvalidation.Ins
	// if err := validate.Struct(req); err != nil {
	// 	validateErr := err.(validator.ValidationErrors)
	// 	encode(w, r, http.StatusBadRequest, Response{Error: appvalidation.ValidationError(validateErr)})
	// 	return
	// }

	var role common.Role
	role = common.Player_Role
	if req.Code != "" {
		if req.Code == code {
			role = common.Admin_Role
		} else {
			// encode(w, r, http.StatusBadRequest, Response{Error: []string{"Неправильный код администратора"}})
			return nil, enerr.E(op, "Неправильный код администратора", enerr.InvalidRequest)
		}
	}
	token, err := GenerateRandomStringURLSafe(32)
	if err != nil {
		// encode(w, r, http.StatusInternalServerError, Response{Error: []string{"Ошибка сервера"}})
		return nil, enerr.E(op, err)
	}
	channels := []string{"default"}
	user := entities.NewUser(req.Name, token, "", role, channels)
	params := database.InsertUserParams{
		Name:   user.Name,
		Apikey: user.Apikey,
		Room:   sql.NullString{},
		Role:   int64(user.Role),
	}
	_, err = database.New(db).InsertUser(context.TODO(), params)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return nil, enerr.E(op, err, enerr.Exist)
			}
		}
		// encode(w, r, http.StatusConflict, Response{Error: []string{"Пользователь с таким именем уже существует"}})
		return nil, enerr.E(op, err, enerr.Internal)
	}
	log.Info("user registred", "name", req.Name, "role", role)
	// encode(w, r, http.StatusCreated, Response{Token: token, Role: role})
	resp := &LoginResponse{
		Token: token,
		Role:  user.Role,
		Error: []string{},
	}
	return resp, nil

}

func GenerateRandomStringURLSafe(n int) (string, error) {
	b, err := GenerateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

type PatchRequest struct {
	Name string
	Role common.Role
}

//	type Response struct {
//		Error []string `json:"error"`
//	}
func PatchUser(req PatchRequest, db *sql.DB) (*entities.User, error) {
	const op = "usecase.user/PatchUser"
	// log = log.With(slog.String("op", op))

	params := database.PatchUserRoleParams{
		Role: int64(req.Role),
		Name: req.Name,
	}
	row, err := database.New(db).PatchUserRole(context.TODO(), params)
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok {
			if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
				return nil, enerr.E(op, err, enerr.Exist)
			}
		}
		return nil, enerr.E(op, err, enerr.Internal)
	}
	user := &entities.User{
		Name:   row.Name,
		Apikey: row.Apikey,
		Room:   row.Room.String,
		Role:   common.Role(row.Role),
	}
	return user, nil

}
