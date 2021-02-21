package main

import (
    "fmt"

    "geoguessr/mailbox"
)

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

func main() {
    mail := &mailbox.MailBox{};
    err := mail.Init();

    if err != nil {
        panic(err);
    }

    fmt.Println(mail);
    err = signUp(mail);

    if err != nil {
        panic(err);
    }

    mail.Delete();
}
