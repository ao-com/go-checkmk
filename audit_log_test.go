package checkmk

import (
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuditLog(t *testing.T) {
	auditLog := AuditLog{}

	Convey("Given an audit log", t, func() {
		Convey("When parsing entries from reader", func() {
			bytes, _ := ioutil.ReadFile("./samples/audit_log_sample.txt")
			auditLogSample := string(bytes)
			r := strings.NewReader(auditLogSample)

			Convey("They should be parsed correctly", func() {
				err := auditLog.ParseFromReader(r)

				So(err, ShouldBeNil)
				So(len(auditLog), ShouldEqual, 2)
				So(auditLog[0].Host, ShouldEqual, "some-host-1")
				So(auditLog[0].Description, ShouldEqual, "Created new host some-host-1.")
				So(auditLog[1].Host, ShouldEqual, "some-host-2")
				So(auditLog[1].Description, ShouldEqual, "")
			})

			Convey("They should be returned when finding by description", func() {
				entries := auditLog.FindEntriesByDescription("Created new host some-host-1.")

				So(len(entries), ShouldEqual, 1)
			})

			Convey("They should be returned when finding by username", func() {
				entries := auditLog.FindEntriesByUsername("autoscale")

				So(len(entries), ShouldEqual, 2)
			})
		})
	})
}
