package checkmk

import (
	"fmt"
	"io"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// AuditLogEntry represents a row in the audit log
type AuditLogEntry struct {
	Host        string
	Date        time.Time
	Username    string
	Description string
}

// AuditLog is a slice of audit log entries
type AuditLog []*AuditLogEntry

// FindEntriesByDescription returns audit log entries based on a matching description
func (auditLog *AuditLog) FindEntriesByDescription(description string) []*AuditLogEntry {
	matchingEntries := []*AuditLogEntry{}
	for _, entry := range *auditLog {
		if entry.Description == description {
			matchingEntries = append(matchingEntries, entry)
		}
	}

	return matchingEntries
}

// FindEntriesByUsername returns audit log entries based on a matching username
func (auditLog *AuditLog) FindEntriesByUsername(username string) []*AuditLogEntry {
	matchingEntries := []*AuditLogEntry{}
	for _, entry := range *auditLog {
		if entry.Username == username {
			matchingEntries = append(matchingEntries, entry)
		}
	}

	return matchingEntries
}

// ParseFromReader populates the audit log with entries found in the reader
func (auditLog *AuditLog) ParseFromReader(r io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}

	tableCells := doc.Find("table.auditlog tr td")
	if tableCells.Length() <= 0 {
		return nil
	}

	counter := 1
	auditLogEntry := &AuditLogEntry{}
	dateString := ""
	for i := range tableCells.Nodes {
		selection := tableCells.Eq(i)
		value := ""
		if selection != nil {
			value = selection.Text()
		}

		switch counter {
		case 1:
			auditLogEntry.Host = value
			break
		case 2:
			dateString = value
			break
		case 3:
			timeString := value
			date, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT%sZ", dateString, timeString))
			if err != nil {
				break
			}

			auditLogEntry.Date = date
			break
		case 4:
			auditLogEntry.Username = value
			break
		case 5:
			auditLogEntry.Description = value
			break
		}

		counter++
		if counter == 6 {
			*auditLog = append(*auditLog, auditLogEntry)
			counter = 1
			auditLogEntry = &AuditLogEntry{}
			dateString = ""
		}
	}

	return nil
}
