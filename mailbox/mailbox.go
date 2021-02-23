package mailbox

import (
    // "encoding/json"
    // "errors"
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

func (m *MailBox) Init() error {
    response, err := http.Get(API_BASE_URL + "?action=genRandomMailbox");
    if err != nil {
        return err;
    }

    defer response.Body.Close();
    raw, err := ioutil.ReadAll(response.Body);
    if err != nil {
        return err;
    }

    parts := strings.Split(string(raw[2:len(raw)-2]), "@");
    m.login, m.domain = parts[0], parts[1];

    return nil;
}

func (m *MailBox) Delete() error {
    return nil
}
