package checkmk

import (
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDowntimes(t *testing.T) {
	downtimes := Downtimes{}

	Convey("Given downtimes", t, func() {
		Convey("When parsing entries from reader", func() {
			bytes, _ := ioutil.ReadFile("./samples/downtimes_example.txt")
			downtimesExample := string(bytes)
			r := strings.NewReader(downtimesExample)

			Convey("They should be parsed correctly", func() {
				err := downtimes.ParseFromReader(r)

				So(err, ShouldBeNil)
				So(len(downtimes), ShouldEqual, 19)
				So(downtimes[0].Origin, ShouldEqual, "command")
				So(downtimes[0].Author, ShouldEqual, "some-user-1")
				So(downtimes[0].Entry, ShouldEqual, "21 min")
				So(downtimes[0].Start, ShouldEqual, "21 min ago")
				So(downtimes[0].End, ShouldEqual, "in 9 hrs")
				So(downtimes[0].Mode, ShouldEqual, "fixed")
				So(downtimes[0].FlexibleDuration, ShouldEqual, "")
				So(downtimes[0].Recurring, ShouldEqual, "(not supported)")
				So(downtimes[0].Comment, ShouldEqual, "Testing.")
			})
		})
	})
}
