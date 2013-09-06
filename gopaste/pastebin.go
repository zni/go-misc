package main

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/xml"
    "log"
    "strings"
)

const pastebinURL = "http://pastebin.com/api/api_post.php"
const pastebinLoginURL = "http://pastebin.com/api/api_login.php"

const (
    PUBLIC   = iota
    UNLISTED = iota
    PRIVATE  = iota
)

type PasteBinInfo struct {
    APIKey  string
    UserKey string
}

func PasteBin(apiKey string) (*PasteBinInfo) {
    return &PasteBinInfo{apiKey, ""}
}

type Paste struct {
    Key         string   `xml:"paste_key"`
    Date        string   `xml:"paste_date"`
    Title       string   `xml:"paste_title"`
    Size        int64    `xml:"paste_size"`
    ExpireDate  string   `xml:"paste_expire_date"`
    Private     bool     `xml:"paste_private"`
    FormatLong  string   `xml:"paste_format_long"`
    FormatShort string   `xml:"paste_format_short"`
    URL         string   `xml:"paste_url"`
    Hits        int64    `xml:"paste_hits"`
}

type Pastes struct {
    XMLName     xml.Name `xml:"root"`
    Pastes []Paste `xml:"paste"`
}

// Get the currently trending pastes.
func (pasteBin *PasteBinInfo) GetTrendingPastes() (*Pastes) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_option", "trends")

    resp, err := http.PostForm(pastebinURL, query)
    if err != nil {
        log.Fatal("unable to fetch trending pastes:", err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("could not read response body:", err)
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        log.Fatal("could not login:", result)
    }

    pastes := &Pastes{}
    result = "<root>" + result + "</root>"

    err = xml.Unmarshal([]byte(result), pastes)
    if err != nil {
        log.Fatal("could not parse xml:", err)
    }

    return pastes
}

// Login to pastebin and populate the UserKey in PasteBinInfo.
func (pasteBin *PasteBinInfo) UserLogin(user, password string) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_user_name", user)
    query.Add("api_user_password", password)

    resp, err := http.PostForm(pastebinLoginURL, query)
    if err != nil {
        log.Fatal("could not login:", err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("could not read response body:", err)
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        log.Fatal("could not login:", result)
    }

    pasteBin.UserKey = result
}

type User struct {
    XMLName     xml.Name `xml:"user"`

    Name        string `xml:"user_name"`
    FormatShort string `xml:"user_format_short"`
    Expiration  string `xml:"user_expiration"`
    AvatarURL   string `xml:"user_avatar_url"`
    Private     string `xml:"user_private"`
    Website     string `xml:"user_website"`
    Email       string `xml:"user_email"`
    Location    string `xml:"user_location"`
    AccountType string `xml:"user_account_type"`
}

// Get information about the logged in user.
func (pasteBin *PasteBinInfo) UserInfo() (*User) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_user_key", pasteBin.UserKey)
    query.Add("api_option", "userdetails")

    resp, err := http.PostForm(pastebinURL, query)
    if err != nil {
        log.Fatal("unable to fetch user info:", err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal("could not read response body:", err)
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        log.Fatal("could not get user information:", result)
    }

    user := &User{}
    err = xml.Unmarshal([]byte(result), user)
    if err != nil {
        log.Fatal("could not parse xml:", err)
    }

    return user
}

/*
func (pasteBin *PasteBinInfo) AnonymousPaste(content string) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
}
*/
