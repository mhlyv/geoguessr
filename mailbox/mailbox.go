package mailbox

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
    "strings"
)

const (
    API_BASE_URL = "https://www.developermail.com/api/v1/mailbox"
)

type MailBox struct {
    name string
    token string
}

func (m *MailBox) GetAddr() string {
    // API specific
    return m.name + "@developermail.com";
}

func Request(method string, url string, data string, header [][2]string) ([]byte, error) {
    client := &http.Client{};
    request, err := http.NewRequest(method, url, strings.NewReader(data));

    if err != nil {
        return []byte{}, err;
    }

    request.ContentLength = int64(len(data));
    for i := range header {
        request.Header.Add(header[i][0], header[i][1]);
    }
    response, err := client.Do(request);

    if err != nil {
        return []byte{}, err;
    }

    defer response.Body.Close();
    contents, err := ioutil.ReadAll(response.Body);

    if err != nil {
        return []byte{}, err;
    }

    return contents, nil;
}

func (m *MailBox) Init() error {
    var response map[string]interface{};
    raw, err := Request("PUT", API_BASE_URL, "",
        [][2]string{{"accept", "application/json"}});

    if err != nil {
        return err;
    }

    err = json.Unmarshal(raw, &response);

    if err != nil {
        return err;
    }

    if response["success"].(bool) != true {
        return errors.New(response["errors"].(string));
    }

    m.name = response["result"].(map[string]interface{})["name"].(string);
    m.token = response["result"].(map[string]interface{})["token"].(string);

    return nil;
}

func (m *MailBox) Delete() error {
    var response map[string]interface{};
    raw, err := Request("DELETE", API_BASE_URL + "/" + m.name, "",
        [][2]string{{"accept", "application/json"}, {"X-MailboxToken", m.token}});

    if err != nil {
        return err;
    }

    err = json.Unmarshal(raw, &response);

    if err != nil {
        return err;
    }

    if response["success"].(bool) != true {
        return errors.New(response["errors"].(string));
    }

    return nil;
}
