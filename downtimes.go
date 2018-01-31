package checkmk

import (
	"io"

	"github.com/PuerkitoBio/goquery"
)

// DowntimeEntry represents a row in the downtimes view
type DowntimeEntry struct {
	Origin           string
	Author           string
	Entry            string
	Start            string
	End              string
	Mode             string
	FlexibleDuration string
	Recurring        string
	Comment          string
}

// Downtimes is a slice of downtime entries
type Downtimes []*DowntimeEntry

// ParseFromReader populates downtimes with entries found in the reader
func (downtimes *Downtimes) ParseFromReader(r io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	tableCells := doc.Find("#data_container > table tr td")
	if tableCells.Length() <= 0 {
		return nil
	}

	counter := 1
	downtimeEntry := &DowntimeEntry{}
	for i := range tableCells.Nodes {
		selection := tableCells.Eq(i)
		value := ""
		if selection != nil {
			value = selection.Text()
		}

		switch counter {
		case 1:
			downtimeEntry.Origin = value
			break
		case 2:
			downtimeEntry.Author = value
			break
		case 3:
			downtimeEntry.Entry = value
			break
		case 4:
			downtimeEntry.Start = value
			break
		case 5:
			downtimeEntry.End = value
			break
		case 6:
			downtimeEntry.Mode = value
			break
		case 7:
			downtimeEntry.FlexibleDuration = value
			break
		case 8:
			downtimeEntry.Recurring = value
			break
		case 9:
			downtimeEntry.Comment = value
			break
		}

		counter++
		if counter == 10 {
			*downtimes = append(*downtimes, downtimeEntry)
			counter = 1
			downtimeEntry = &DowntimeEntry{}
		}
	}

	return nil
}
