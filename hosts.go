package checkmk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vjeantet/jodaTime"
)

type addHostRequestAttributes struct {
	TagAgent  string `json:"tag_agent"`
	IPAddress string `json:"ipaddress"`
}

type addHostRequest struct {
	Hostname   string                   `json:"hostname"`
	Folder     string                   `json:"folder"`
	Attributes addHostRequestAttributes `json:"attributes"`
}

// AddHost adds a host to check mk
func (client Client) AddHost(hostname string, folder string) error {
	reqBody := addHostRequest{
		Hostname: hostname,
		Folder:   folder,
		Attributes: addHostRequestAttributes{
			TagAgent:  "cmk-agent",
			IPAddress: hostname,
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/webapi.py?action=add_host&request=%s&%s", client.URL, string(reqBytes), client.requestCredentials)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.httpClient.Do(req)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

// GetDowntimesForHost returns downtime entries for a given host
func (client Client) GetDowntimesForHost(host string) (Downtimes, error) {
	url := fmt.Sprintf("%s/view.py?host=%s&view_name=downtimes_of_host&%s", client.URL, host, client.requestCredentials)
	fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Downtimes{}, err
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return Downtimes{}, err
	}

	defer resp.Body.Close()

	downtimes := Downtimes{}
	err = downtimes.ParseFromReader(resp.Body)
	return downtimes, err
}

// ScheduleHostDowntime places a given host in downtime for a given amount of minutes
func (client Client) ScheduleHostDowntime(host string, minutes int, comment string) error {
	date := jodaTime.Format("YYYY-MM-dd", time.Now())
	timeParameters := fmt.Sprintf("_down_from_date=%s&_down_from_time=00:00&_down_to_date=%s&_down_to_time=00:00", date, date)
	minutesParameter := fmt.Sprintf("_down_minutes=%d", minutes)
	viewName := "view_name=allhosts"
	url := fmt.Sprintf("%s/view.py?_do_confirm=yes&_transid=-1&_do_actions=yes&host_regex=%s&%s&%s&%s&%s&_down_comment=%s&%s",
		client.URL, host, viewName, timeParameters, "&_down_from_now=From+now+for&_down_duration=02%3A00", minutesParameter,
		comment, client.requestCredentials)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	_, err = client.httpClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
