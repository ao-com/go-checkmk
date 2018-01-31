package checkmk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client struct for check mk
type Client struct {
	URL                string
	Username           string
	Password           string
	requestCredentials string
	httpClient         *http.Client
}

// NewClient will create a new check mk client
func NewClient(url string, username string, password string) Client {
	return Client{
		URL:                url,
		Username:           username,
		Password:           password,
		requestCredentials: fmt.Sprintf("_username=%s&_secret=%s", username, password),
		httpClient:         &http.Client{},
	}
}

// ActivateChanges activates any pending changes in check_mk
func (client Client) ActivateChanges() error {
	url := fmt.Sprintf("%s/webapi.py?action=activate_changes&mode=dirty&allow_foreign_changes=1&%s", client.URL, client.requestCredentials)
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

// AuditLog returns the default audit log which consists of today's entries
func (client Client) AuditLog() (AuditLog, error) {
	url := fmt.Sprintf("%s/wato.py?folder=&mode=auditlog&%s", client.URL, client.requestCredentials)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return AuditLog{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return AuditLog{}, err
	}

	defer resp.Body.Close()
	auditLog := AuditLog{}
	err = auditLog.ParseFromReader(resp.Body)
	return auditLog, err
}

// IsAuthenticated checks if the client is successfully authenticated
func (client Client) IsAuthenticated() (bool, error) {
	url := fmt.Sprintf("%s/dashboard.py?%s", client.URL, client.requestCredentials)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	body := string(bodyBytes)
	isAuthenticated := !strings.Contains(body, "Permission denied") &&
		!strings.Contains(resp.Request.URL.String(), "login.py") &&
		resp.StatusCode != http.StatusUnauthorized

	return isAuthenticated, nil
}
