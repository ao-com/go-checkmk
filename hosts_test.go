package checkmk

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHosts(t *testing.T) {
	url := "http://localhost:8001"

	Convey("Given a client", t, func() {
		username := "go-checkmk"
		password := "somepassword"
		client := NewClient(url, username, password)
		client.httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}

		Convey("AddHost should create a new host", func() {
			hostname := "go-check-mk-test-host"
			err := client.AddHost(hostname, "go-check-mk-test-folder")
			expectedDescription := fmt.Sprintf("Created new host %s.", hostname)
			auditLog, err := client.AuditLog()
			auditLogEntries := auditLog.FindEntriesByDescription(expectedDescription)
			client.ActivateChanges()

			So(err, ShouldBeNil)
			So(len(auditLogEntries), ShouldBeGreaterThan, 0)

			Convey("Given a host", func() {
				Convey("ScheduleDowntime should set the host with downtimes", func() {
					err := client.ScheduleHostDowntime(hostname, 45, "hosts_test_schedule_downtime")
					downtimes, err := client.GetDowntimesForHost(hostname)

					So(err, ShouldBeNil)
					So(len(downtimes), ShouldBeGreaterThan, 0)
				})
			})
		})
	})
}
