package checkmk

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	url := "http://localhost:8001"

	Convey("Given a client", t, func() {
		username := "go-checkmk"
		password := "somepassword"
		client := NewClient(url, username, password)
		client.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		Convey("URL should be set", func() {
			So(client.URL, ShouldNotBeEmpty)
		})

		Convey("Username should be set", func() {
			So(client.Username, ShouldNotBeEmpty)
		})

		Convey("Password should be set", func() {
			So(client.Password, ShouldNotBeEmpty)
		})

		Convey("requestCredentials should be set correctly", func() {
			So(client.requestCredentials, ShouldEqual, fmt.Sprintf("_username=%s&_secret=%s", username, password))
		})

		Convey("HTTPClient should be set", func() {
			So(client.HTTPClient, ShouldNotBeNil)
		})

		Convey("IsAuthenticated should return true", func() {
			isAuthenticated, err := client.IsAuthenticated()

			So(err, ShouldBeNil)
			So(isAuthenticated, ShouldBeTrue)
		})

		Convey("AuditLog should return valid a valid audit log", func() {
			auditLog, err := client.AuditLog()

			So(err, ShouldBeNil)
			So(auditLog, ShouldNotBeEmpty)
		})

		Convey("ActivateChanges should activate pending changes", func() {
			err := client.ActivateChanges()

			So(err, ShouldBeNil)
		})
	})

	Convey("Given a client with incorrect credentials", t, func() {
		username := "thiswontwork"
		password := "password"
		client := NewClient(url, username, password)
		client.HTTPClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		Convey("IsAuthenticated should return false", func() {
			isAuthenticated, err := client.IsAuthenticated()

			So(err, ShouldBeNil)
			So(isAuthenticated, ShouldBeFalse)
		})
	})
}
