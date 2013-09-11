package gopaste

import (
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/xml"
    "strings"
)

const pastebinURL = "http://pastebin.com/api/api_post.php"
const pastebinLoginURL = "http://pastebin.com/api/api_login.php"

const (
    Public   = iota
    Unlisted = iota
    Private  = iota
)

type PasteError struct {
    ErrorString string
}

func (p *PasteError) Error() string {
    return p.ErrorString
}

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
    XMLName xml.Name `xml:"root"`
    Pastes  []Paste  `xml:"paste"`
}

// Get the currently trending pastes.
func (pasteBin *PasteBinInfo) GetTrendingPastes() (*Pastes, error) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_option", "trends")

    resp, err := http.PostForm(pastebinURL, query)
    if err != nil {
        return nil, &PasteError{"unable to fetch trending pastes: " + err.Error()}
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, &PasteError{"could not read response: " + err.Error()}
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        return nil, &PasteError{"could not get trending pastes: " + result}
    }

    pastes := &Pastes{}
    result = "<root>" + result + "</root>"

    err = xml.Unmarshal([]byte(result), pastes)
    if err != nil {
        return nil, &PasteError{"could not parse trending pastes: " + err.Error()}
    }

    return pastes, nil
}

// Login to pastebin and populate the UserKey in PasteBinInfo.
func (pasteBin *PasteBinInfo) UserLogin(user, password string) (error) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_user_name", user)
    query.Add("api_user_password", password)

    resp, err := http.PostForm(pastebinLoginURL, query)
    if err != nil {
        return &PasteError{"could not login: " + err.Error()}
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return &PasteError{"could not read response body: " + err.Error()}
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        return &PasteError{"could not login: " + result}
    }

    pasteBin.UserKey = result
    return nil
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

func parseUserInfo(userInfo []byte) (*User, error) {
    user := &User{}
    err := xml.Unmarshal(userInfo, user)
    if err != nil {
        return nil, &PasteError{"could not parse user information:" + err.Error()}
    }

    return user, nil
}

// Get information about the logged in user.
func (pasteBin *PasteBinInfo) UserInfo() (*User, error) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_user_key", pasteBin.UserKey)
    query.Add("api_option", "userdetails")

    resp, err := http.PostForm(pastebinURL, query)
    if err != nil {
        return nil, &PasteError{"request failed: " + err.Error()}
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, &PasteError{"could not read response: " + err.Error()}
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        return nil, &PasteError{"could not get user information: " + result}
    }

    return parseUserInfo([]byte(result))
}

/*
func (pasteBin *PasteBinInfo) AnonymousPaste(content string) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
}
*/
