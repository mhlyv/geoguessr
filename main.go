package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "regexp"
    "time"

    "geoguessr/mailbox"
)

const PORT = 8080;

// this is ugly but it has to be done, otherwise
// the regex needs to be recompiled again and again
var hrefPattern = func () *regexp.Regexp{
    pattern, err := regexp.Compile("href=\\\"([^\"]+)\\\"");
    if err != nil {
        panic(err);
    }
    return pattern;
}();

func signUp(m *mailbox.MailBox) error {
    content, err := mailbox.Request(
        "POST",
        "https://www.geoguessr.com/api/v3/accounts/signup",
        fmt.Sprintf("{\"email\":\"%s\"}", m.GetAddr()),
        [][2]string{
            {"content-type", "application/json"},
            {"authority", "www.geoguessr.com"},
            {"user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4421.5 Safari/537.36"},
            {"accept", "*/*"},
        },
    );

    if err != nil {
        return err;
    }

    fmt.Println(string(content));

    return nil;
}

func getVerificationId(m *mailbox.MailBox) (int, error) {
    slept := time.Duration(0);
    delay := 1 * time.Millisecond;
    limit := 3 * time.Second;

    for true {
        time.Sleep(delay);
        slept += delay;

        ids, err := m.GetMessageIds();
        if err != nil {
            panic(err);
        }

        if len(ids) > 0 {
            return ids[0], nil;
        }

        if slept > limit {
            return 0, errors.New("email timed out");
        }
    }

    panic("reached unreachable");
    return 0, nil;
}

func getGeoguessrUrl() string {
    mail := &mailbox.MailBox{};
    err := mail.Init();

    if err != nil {
        panic(err);
    }

    fmt.Printf("'%s'\n", mail.GetAddr());
    err = signUp(mail);

    if err != nil {
        panic(err);
    }

    id, err := getVerificationId(mail);
    if err != nil {
        panic(err);
    }

    msg, err := mail.ReadMessage(id);
    if err != nil {
        panic(err);
    }

    matched := hrefPattern.FindStringSubmatch(msg);
    if len(matched) < 2 {
        panic("regex failed");
    }

    return matched[1];
}

func handler(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, getGeoguessrUrl(), 301);
}

func main() {
    http.HandleFunc("/", handler);
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil));
}
