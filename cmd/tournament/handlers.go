package main

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"

// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/anrew1002/Tournament-ChemLoto/models"
// 	"github.com/anrew1002/Tournament-ChemLoto/sqlite"
// 	"github.com/go-chi/chi/v5"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/google/uuid"
// 	"github.com/gorilla/sessions"
// )

// func (app *App) HubHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userSession := r.Context().Value("user").(*sessions.Session)
// 		username, ok := userSession.Values["username"].(string)
// 		if !ok {
// 			log.Println("Fail to type assertion")
// 		}
// 		admin, ok := userSession.Values["admin"].(bool)
// 		if !ok {
// 			log.Println("Fail to type assertion")
// 		}
// 		// log.Println("admin", admin)

// 		data := struct {
// 			Username      string
// 			Admin         bool
// 			ErrMaxPlayers string
// 			ErrRoomName   string
// 			Error         string
// 		}{
// 			Username:      username,
// 			Admin:         admin,
// 			ErrMaxPlayers: "",
// 			ErrRoomName:   "",
// 		}

// 		app.render(w, http.StatusOK, "room_list", data)
// 	}
// }
// func (app *App) GetRooms() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rooms := sqlite.NewStorage().GetRooms()
// 		app.writeJSON(w, http.StatusOK, envelope{"rooms": rooms}, nil)
// 	}
// }
// func (app *App) CreateRoomHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// if err := r.ParseForm(); err != nil {
// 		// 	fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		// 	return
// 		// }
// 		// formErrors := make(map[string]string)
// 		// if r.FormValue("roomName") == "" {
// 		// 	formErrors["ErrRoomName"] = "Имя должно быть заполнено"
// 		// }
// 		// if r.FormValue("maxPlayers") == "" {
// 		// 	formErrors["ErrMaxPlayers"] = "Кол-во игроков не должно быть пустым"
// 		// }
// 		data := new(models.Room)
// 		err := json.NewDecoder(r.Body).Decode(data)
// 		if err != nil {
// 			log.Print(err)
// 			app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"errors": []string{"Fail to decode JSON"}, "success": false}, nil)
// 			return
// 		}
// 		// log.Printf("%+v", data)
// 		err = app.validate.Struct(data)
// 		errorsList := make([]string, 2)
// 		if err != nil {
// 			validationErrors := err.(validator.ValidationErrors)
// 			// log.Printf("%+v", validationErrors)
// 			for _, v := range validationErrors {
// 				if v.Field() == "Name" {
// 					errorsList = append(errorsList, "Имя комнаты указано неправильно")
// 				} else if v.Field() == "Max_partic" {
// 					errorsList = append(errorsList, "Количество игроков указано неправильно")
// 				} else if v.Field() == "Time" {
// 					errorsList = append(errorsList, "Время указано неправильно")
// 				} else {
// 					errorsList = append(errorsList, v.Error())
// 				}
// 			}
// 			app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"errors": errorsList, "success": false}, nil)
// 			return
// 		}

// 		if !data.IsAuto {
// 			// time, err := strconv.Atoi(r.FormValue("time"))
// 			// if err != nil {
// 			// 	time = 0
// 			// }
// 			// data.Time = time
// 			data.Time = 0
// 		}
// 		// data.Elements = map[string]int{
// 		// 	"H":     52,
// 		// 	"C":     40,
// 		// 	"CH":    24,
// 		// 	"CH2":   24,
// 		// 	"CH3":   28,
// 		// 	"O":     28,
// 		// 	"CL":    16,
// 		// 	"N":     16,
// 		// 	"C6H4":  16,
// 		// 	"TRADE": 4,
// 		// 	// "C6H4": 16,
// 		// }
// 		// log.Println(data)
// 		err = app.database.CreateRoom(*data)
// 		// fmt.Printf("err: %T\n", err)
// 		// log.Print(err)
// 		if err != nil {
// 			if errors.Is(err, sqlite.ErrDup) {
// 				log.Print(err)
// 				// formErrors["ErrRoomName"] = "Такая комната уже существует!"
// 				err = app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"errors": []string{"Такая комната уже существует!"}, "success": false}, nil)
// 				if err != nil {
// 					log.Println(err)
// 				}
// 			}
// 			return
// 		}
// 		app.clientManager.addRoom(*data)
// 		app.writeJSON(w, http.StatusCreated, envelope{"errors": nil, "success": true}, nil)
// 	}
// }

// func (app *App) RoomHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if roomID := chi.URLParam(r, "room_id"); roomID != "" {
// 			// log.Println(roomID)
// 			if room, err := app.database.GetRoom(roomID); err == nil {
// 				userSession := r.Context().Value("user").(*sessions.Session)
// 				username, ok := userSession.Values["username"].(string)
// 				if !ok {
// 					log.Println("Fail to type assertion")
// 				}

// 				admin, ok := userSession.Values["admin"].(bool)
// 				if !ok {
// 					log.Println("Fail to type assertion")
// 				}
// 				// log.Println("admin", admin)
// 				roomusers := app.database.GetUsers()
// 				count := 0
// 				for _, v := range roomusers {
// 					if !v.Admin && v.Room == roomID {
// 						count += 1
// 					}
// 				}
// 				data := struct {
// 					Room     string
// 					Username string
// 					Admin    bool
// 					Error    string
// 				}{
// 					Room:     roomID,
// 					Username: username,
// 					Admin:    admin,
// 				}
// 				user, err := app.database.GetUser(username)
// 				if err != nil {
// 					log.Println("roomhandler: failed to get user")
// 				}

// 				if count >= room.Max_partic && !admin && user.Room != roomID {
// 					data.Error = "В комнате больше нет мест"
// 					app.render(w, http.StatusTemporaryRedirect, "room_list", data)
// 					return
// 				}
// 				app.database.UpdateUserRoom(username, room.Name)

// 				app.render(w, http.StatusOK, "room", data)
// 			} else {
// 				w.WriteHeader(404)
// 				w.Write([]byte("404"))

// 			}
// 		}

// 	}
// }

// func (app *App) RoomDeleteHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if roomID := chi.URLParam(r, "room_id"); roomID != "" {
// 			log.Println("Удаление комнаты: ", roomID)
// 			err := app.database.DeleteRoom(roomID)
// 			if err != nil {
// 				log.Println("Не получилось удалить комнаты", err)
// 			}
// 			app.clientManager.removeRoom(roomID)
// 			app.writeJSON(w, http.StatusOK, envelope{"success": true}, nil)
// 		}
// 	}
// }
// func (app *App) RoomsClearHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		err := app.database.DeleteAllRooms()
// 		if err != nil {
// 			app.writeJSON(w, http.StatusInternalServerError, envelope{"success": false}, nil)
// 			log.Println(err)
// 		}
// 		app.writeJSON(w, http.StatusOK, envelope{"success": true}, nil)
// 	}
// }
// func (app *App) UserHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(r.RequestURI)
// 		if roomID := chi.URLParam(r, "user_id"); roomID != "" {
// 			if err := r.ParseForm(); err != nil {
// 				fmt.Fprintf(w, "ParseForm() err: %v", err)
// 				return
// 			}
// 			log.Println("update score on: ", r.FormValue("score"), roomID)
// 			// formErrors := make(map[string]string)
// 			if r.FormValue("score") != "" {
// 				score := r.FormValue("score")
// 				// userSession := r.Context().Value("user").(*sessions.Session)
// 				// username, ok := userSession.Values["username"].(string)
// 				// if !ok {
// 				// 	log.Println("Fail to type assertion")
// 				// }
// 				scoreValue, err := strconv.Atoi(score)
// 				if err != nil {
// 					log.Println("Failed to atoi", err, score)
// 				}
// 				err = app.database.UpdateUserScore(roomID, scoreValue)
// 				if err != nil {
// 					log.Println("user score update:", err)
// 					app.writeJSON(w, http.StatusBadRequest, envelope{"success": false}, nil)
// 				}
// 				app.writeJSON(w, http.StatusOK, envelope{"success": true}, nil)
// 				return
// 			}
// 			app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"success": false}, nil)
// 			return
// 		} else {
// 			app.writeJSON(w, http.StatusBadRequest, nil, nil)
// 		}
// 	}
// }
// func (app *App) GetUsers() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		users := app.database.GetUsers()
// 		app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
// 	}
// }
// func (app *App) UsersClearHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		err := app.database.DeleteAllRooms()
// 		if err != nil {
// 			app.writeJSON(w, http.StatusInternalServerError, envelope{"success": false}, nil)
// 			log.Println(err)
// 		}
// 		app.writeJSON(w, http.StatusOK, envelope{"success": true}, nil)
// 	}
// }

// func (app *App) AdminPanelHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userSession := r.Context().Value("user").(*sessions.Session)
// 		username, ok := userSession.Values["username"].(string)
// 		if !ok {
// 			log.Println("Fail to type assertion")
// 		}

// 		data := struct {
// 			Username string
// 		}{
// 			Username: username,
// 		}
// 		app.render(w, http.StatusOK, "admin", data)

// 	}
// }
// func (app *App) AdminLoginHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		app.render(w, http.StatusOK, "admin_login", nil)
// 	}
// }
// func (app *App) PostAdminLoginHandler(AdminCode string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		data := new(models.User)
// 		if err := r.ParseForm(); err != nil {
// 			fmt.Fprintf(w, "ParseForm() err: %v", err)
// 			return
// 		}
// 		formErrors := make(map[string]string)
// 		if r.FormValue("name") == "" {
// 			formErrors["name"] = "Имя должно быть заполнено"
// 		}
// 		if r.FormValue("code") == "" {
// 			formErrors["code"] = "Код должнен быть указан"
// 		}

// 		if len(formErrors) != 0 {
// 			app.render(w, http.StatusUnprocessableEntity, "admin_login", formErrors)
// 			return
// 		}

// 		seed := strconv.Itoa(rand.Intn(1000))
// 		data.Username = r.FormValue("name") + "#" + seed
// 		code := r.FormValue("code")
// 		if code == AdminCode {
// 			data.Admin = true
// 		}
// 		// log.Print(data)
// 		app.database.AddUser(data)
// 		err := SetCookie(w, r, data, app)
// 		if err != nil {
// 			log.Println("failed to set cookie")
// 		}
// 		http.Redirect(w, r, "/admin", http.StatusSeeOther)
// 	}
// }

// func (app *App) LoginHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		app.render(w, http.StatusOK, "login", nil)
// 	}
// }

// func (app *App) PostLoginHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		data := new(models.User)
// 		if err := r.ParseForm(); err != nil {
// 			fmt.Fprintf(w, "ParseForm() err: %v", err)
// 			return
// 		}
// 		formErrors := make(map[string]string)
// 		if r.FormValue("name") == "" {
// 			formErrors["name"] = "Имя должно быть заполнено"
// 		}

// 		if len(formErrors) != 0 {
// 			app.render(w, http.StatusUnprocessableEntity, "login", formErrors)
// 			return
// 		}

// 		seed := strconv.Itoa(rand.Intn(1000))
// 		data.Username = r.FormValue("name") + "#" + seed
// 		// log.Print(data)
// 		app.database.AddUser(data)
// 		err := SetCookie(w, r, data, app)
// 		if err != nil {
// 			log.Println("failed to set cookie")
// 		}
// 		http.Redirect(w, r, "/hub", http.StatusSeeOther)
// 	}
// }

// func SetCookie(w http.ResponseWriter, r *http.Request, data *models.User, app *App) error {
// 	sessionToken := uuid.NewString()
// 	expiresAt := time.Now().Add(120 * time.Hour)

// 	session, err := app.CS.Get(r, sessionToken)
// 	if err != nil {
// 		return err
// 	}
// 	// session.Values["id"] = data.Id
// 	session.Values["username"] = data.Username
// 	session.Values["admin"] = data.Admin

// 	err = session.Save(r, w)
// 	if err != nil {
// 		return err
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "session_token",
// 		Value:   sessionToken,
// 		Expires: expiresAt,
// 	})
// 	return nil
// }
