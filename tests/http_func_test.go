package tests

import (
	"net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

var u = url.URL{
	Scheme: "http",
	Host:   host,
	Path:   baseUrl,
}

func TestCreateUser(t *testing.T) {

	e := httpexpect.Default(t, u.String())

	// test create user
	username := getRandomUsername()

	resp := Createuser(e, username, "")
	resp.Value("error").IsNull()
	checkUserInfo(t, username, resp.Value("token").String().Raw(), "", "Player_Role")

	// test create admin
	username = getRandomUsername()

	resp = Createuser(e, username, "test_code")
	resp.Value("error").IsNull()
	checkUserInfo(t, username, resp.Value("token").String().Raw(), "", "Admin_Role")

	// test create  with same name is not null
	username = username

	resp = Createuser(e, username, "test_code")
	resp.Value("error").Array().Length().IsEqual(1)
}

func Createuser(e *httpexpect.Expect, username string, code string) *httpexpect.Object {
	resp := e.POST("/users").
		WithJSON(map[string]interface{}{"name": username, "code": code}).
		Expect().
		JSON().Object().ContainsKey("token").ContainsKey("error")
	return resp
}
