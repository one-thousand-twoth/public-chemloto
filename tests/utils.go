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

func checkUserInfo(t *testing.T, username string, token string, room string, role string) {
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
func getRandomUsername() string {
	gofakeit.Seed(0)
	return gofakeit.Name() + strconv.Itoa(gofakeit.Number(0, 10))
}
