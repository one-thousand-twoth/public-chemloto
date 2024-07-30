package tests

import (
	"net/url"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const (
	host    = "localhost:1090"
	baseUrl = "/api/v1"
)

func TestCreateUser(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   baseUrl,
	}
	e := httpexpect.Default(t, u.String())

	// test create user
	username := gofakeit.Name() + strconv.Itoa(gofakeit.Number(0, 10))

	resp := e.POST("/users").
		WithForm(map[string]interface{}{"name": username}).
		Expect().
		Status(201).Text()
	resp.NotEmpty()
	testUserInfo(t, username, resp.Raw(), "", "Player_Role")

	// test crate admin
	username = gofakeit.Name() + strconv.Itoa(gofakeit.Number(0, 10))

	resp = e.POST("/admin").
		WithForm(map[string]interface{}{"name": username, "code": "test_code"}).
		Expect().
		Status(201).Text()
	resp.NotEmpty()
	testUserInfo(t, username, resp.Raw(), "", "Admin_Role")
}

func testUserInfo(t *testing.T, username string, token string, room string, role string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   baseUrl,
	}
	e := httpexpect.Default(t, u.String())
	e.GET("/users/"+token).Expect().JSON().Object().
		HasValue("username", username).
		HasValue("room", room).
		HasValue("role", role)
}
