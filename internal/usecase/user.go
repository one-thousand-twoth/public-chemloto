package usecase

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log/slog"
	"regexp"

	"github.com/anrew1002/Tournament-ChemLoto/internal/common"
	"github.com/anrew1002/Tournament-ChemLoto/internal/common/enerr"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	enmodels "github.com/anrew1002/Tournament-ChemLoto/internal/engines/models"
	"github.com/anrew1002/Tournament-ChemLoto/internal/entities"
	"github.com/invopop/validation"
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

func (r LoginRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(1, 25).Error("Должно быть меньше 26 символов"),
			validation.Match(regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9- ]+$`)).
				Error("Должно содержать только буквы, цифры и пробел"),
		))
}
func stringEquals(str string) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if s != str {
			return errors.New("Неправильный код администратора")
		}
		return nil
	}
}

func (uc *Usecases) Login(log *slog.Logger, req LoginRequest, code string) (*LoginResponse, error) {
	const op enerr.Op = "usecase.user/Login"
	err := req.Validate()
	if err != nil {
		return nil, enerr.E(op, err, enerr.Validation)
	}

	var role common.Role
	role = common.Player_Role
	if req.Code != "" {
		err = validation.Errors{
			"code": validation.Validate(req.Code, validation.By(stringEquals(code))),
		}.Filter()
		if err != nil {
			return nil, enerr.E(op, err, enerr.Validation)
		}
		role = common.Admin_Role
	}
	token, err := generateRandomStringURLSafe(32)
	if err != nil {
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

	_, err = uc.UserRepo.CreateUser(params)

	if err != nil {
		return nil, enerr.E(op, err)
	}
	log.Info("user registred", "name", req.Name, "role", role)
	resp := &LoginResponse{
		Token: token,
		Role:  user.Role,
		Error: []string{},
	}
	return resp, nil

}

func (uc *Usecases) DeleteUser(ctx context.Context, username string) error {
	const op enerr.Op = "usecase.user/DeleteUser"

	user, err := uc.UserRepo.GetUserByName(username)
	if err != nil {
		return enerr.E(op, err, enerr.Database)
	}
	if user.IsInRoom() {
		return enerr.E(op, "Пользователь сейчас в комнате, его невозможно удалить")
	}

	err = uc.queries.DeleteUser(ctx, int64(user.ID))
	if err != nil {
		return enerr.E(op, err, enerr.Database)
	}

	return nil
}
func generateRandomStringURLSafe(n int) (string, error) {
	b, err := generateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b), err
}
func generateRandomBytes(n int) ([]byte, error) {
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
func (uc *Usecases) PatchUserRole(ctx context.Context, req PatchRequest) (*entities.User, error) {
	const op enerr.Op = "usecase.user/PatchUser"
	// log = log.With(slog.String("op", op))

	params := database.PatchUserRoleParams{
		Role: int64(req.Role),
		Name: req.Name,
	}
	row, err := uc.queries.PatchUserRole(context.TODO(), params)
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
func (uc *Usecases) RouteActionToUserRoom(ctx context.Context, userID entities.ID, msg map[string]any) error {
	const op enerr.Op = "usecase.user/RouteActionToUserRoom"
	user, err := uc.UserRepo.GetUserByID(userID)
	if err != nil {
		return enerr.E(err)
	}
	room, err := uc.RoomRepo.GetRoom(user.Room)
	if err != nil {
		return err
	}
	go room.Engine.Input(enmodels.Action{
		Player:   user.Name,
		Envelope: msg,
	})

	return nil

}
