package mailbox

import (
    "encoding/json"
    // "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

const (
    API_BASE_URL = "https://www.1secmail.com/api/v1/"
)

type MailBox struct {
    login string
    domain string
}

func (m *MailBox) GetAddr() string {
    return m.login + "@" + m.domain;
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

func getBody(url string) ([]byte, error) {
    response, err := http.Get(url);
    if err != nil {
        return []byte{}, err;
    }

    defer response.Body.Close();
    raw, err := ioutil.ReadAll(response.Body);
    if err != nil {
        return raw, err;
    }

    return raw, nil;
}

func (m *MailBox) Init() error {
    raw, err := getBody(API_BASE_URL + "?action=genRandomMailbox");
    if err != nil {
        return err;
    }

    /// this works but it's not right
    // parts := strings.Split(string(raw[2:len(raw)-2]), "@");
    // m.login, m.domain = parts[0], parts[1];

    var addresses []string
    err = json.Unmarshal(raw, &addresses);
    if err != nil {
        return err;
    }

    split := strings.Split(addresses[0], "@");
    m.login, m.domain = split[0], split[1];

    return nil;
}

func (m *MailBox) GetMessageIds() ([]int, error) {
    ids := []int{};
    raw, err := getBody(API_BASE_URL + "?action=getMessages" +
        "&login=" + m.login + "&domain=" + m.domain);
    if err != nil {
        return ids, err;
    }

    var mails []map[string]interface{};
    err = json.Unmarshal(raw, &mails);
    if err != nil {
        return ids, err;
    }

    for i := range mails {
        ids = append(ids, int(mails[i]["id"].(float64)));
    }

    return ids, nil;
}

func (m *MailBox) ReadMessage(id int) (string, error) {
    raw, err := getBody(API_BASE_URL + "?action=readMessage" +
        "&login=" + m.login + "&domain=" + m.domain +
        fmt.Sprintf("&id=%d", id));

    if err != nil {
        return string(raw), err;
    }

    var msg map[string]interface{};
    err = json.Unmarshal(raw, &msg);
    if err != nil {
        return "", err;
    }

    return msg["body"].(string), nil;
}
