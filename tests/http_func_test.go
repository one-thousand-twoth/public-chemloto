package tests

import (
	"net/http"
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
	jresp := resp.JSON().Object()
	jresp.Value("error").IsNull()
	checkUserInfo(t, username, jresp.Value("token").String().Raw(), "", "Player_Role")

	// test create admin
	username = getRandomUsername()

	resp = Createuser(e, username, "test_code")
	jresp = resp.JSON().Object()
	jresp.Value("error").IsNull()
	checkUserInfo(t, username, jresp.Value("token").String().Raw(), "", "Admin_Role")

	// test create  with same name is not null
	username = username

	resp = Createuser(e, username, "test_code")
	resp.Status(http.StatusBadRequest)
	jresp = resp.JSON().Object()
	jresp.Value("error").Object().Value("kind").IsEqual("item already exists")
}

func Createuser(e *httpexpect.Expect, username string, code string) *httpexpect.Response {
	resp := e.POST("/users").
		WithJSON(map[string]interface{}{"name": username, "code": code}).
		Expect()
	return resp
}
