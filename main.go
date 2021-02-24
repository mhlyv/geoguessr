package main

import (
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "regexp"
    "time"

    "geoguessr/mailbox"
)

const (
    PORT = 8080;
    LOG_FILE = "/tmp/geoguessr.log";
)

var (
    // compile a regex so it doesn't have to be done again
    hrefPattern = func () *regexp.Regexp{
        // this matches the first href="http://..." in the html mail
        pattern, err := regexp.Compile("href=\\\"([^\"]+)\\\"");
        if err != nil {
            panic(err);
        }
        return pattern;
    }();
)

func signUp(m *mailbox.MailBox) error {
    // sign up to geoguessr with the email address
    _, err := mailbox.Request(
        "POST",
        "https://www.geoguessr.com/api/v3/accounts/signup",
        fmt.Sprintf("{\"email\":\"%s\"}", m.GetAddr()),
        // set some necessary headers
        [][2]string{
            {"content-type", "application/json"},
            {"authority", "www.geoguessr.com"},
            {"user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4421.5 Safari/537.36"},
            {"accept", "*/*"},
        },
    );

    if err != nil {
        return fmt.Errorf("signUp: %v", err);
    }

    return nil;
}

func getVerificationId(m *mailbox.MailBox) (int, error) {
    slept := time.Duration(0);
    delay := 1 * time.Millisecond;
    limit := 3 * time.Second;

    // wait for the verification email to arrive
    for true {
        time.Sleep(delay);
        slept += delay;

        ids, err := m.GetMessageIds();
        if err != nil {
            return 0, fmt.Errorf("getVerificationId: %v", err);
        }

        if len(ids) > 0 {
            // return the id of the message
            return ids[0], nil;
        }

        // error on timeout
        if slept > limit {
            return 0, errors.New("email timed out");
        }
    }

    panic("reached unreachable");
    return 0, nil;
}

func getGeoguessrUrl() (string, error) {
    // init temporary mail
    mail := &mailbox.MailBox{};
    err := mail.Init();

    if err != nil {
        return "", fmt.Errorf("getGeoguessrUrl: %v", err);
    }

    // sign up to geoguessr
    err = signUp(mail);

    if err != nil {
        return "", fmt.Errorf("getGeoguessrUrl: %v", err);
    }

    // get verification email id
    id, err := getVerificationId(mail);
    if err != nil {
        return "", fmt.Errorf("getGeoguessrUrl: %v", err);
    }

    // get the raw email data
    msg, err := mail.ReadMessage(id);
    if err != nil {
        return "", fmt.Errorf("getGeoguessrUrl: %v", err);
    }

    // get the verification url from the raw message
    matched := hrefPattern.FindStringSubmatch(msg);
    if len(matched) < 2 {
        return "", errors.New("url extraction failed");
    }

    return matched[1], nil;
}

func handler(w http.ResponseWriter, r *http.Request) {
    logFile, err := os.OpenFile(LOG_FILE, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666);
    if err != nil {
        panic(err);
    }
    defer logFile.Close();
    log.SetOutput(io.MultiWriter(os.Stdout, logFile));

    // disallow caching of the page
    w.Header().Set("Cache-Control", "no-cache,no-store,max-age=0");

    // get verification url
    url, err := getGeoguessrUrl();
    if err != nil {
        log.Println(r, err);
        fmt.Fprintf(w, "Internal error: '%v'\nTry again!\n", err);
    } else {
        log.Println(r, url);
        // redirect user
        http.Redirect(w, r, url, 301);
    }
}

func main() {
    logFile, err := os.OpenFile(LOG_FILE, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666);
    if err != nil {
        panic(err);
    }
    defer logFile.Close();
    log.SetOutput(io.MultiWriter(os.Stdout, logFile));

    http.HandleFunc("/", handler);
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil));
}
