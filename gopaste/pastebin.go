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
    // We wrap the result in a throwaway root tag, so that the pastes can be
    // unmarshaled.
    result = "<root>" + result + "</root>"

    err = xml.Unmarshal([]byte(result), pastes)
    if err != nil {
        return nil, &PasteError{"could not parse trending pastes: " + err.Error()}
    }

    return pastes, nil
}

// Login to pastebin and populate the UserKey in PasteBinInfo.
func (pasteBin *PasteBinInfo) UserLogin(user, password string) error {
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
        return nil, &PasteError{"could not parse user information: " + err.Error()}
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

// Paste access level.
type Access int
const (
    Public Access = iota
    Unlisted
    Private
)

// Paste expiration time length.
type Time string
const (
    Never Time = "N"
    TenMinutes = "10M"
    OneHour    = "1H"
    OneDay     = "1D"
    OneWeek    = "1W"
    TwoWeeks   = "2W"
    OneMonth   = "1M"
)


// Paste options, made up of an access level, expiration time, and format.
// Format is one of: http://pastebin.com/api#5
type PasteOptions struct {
    Privacy    Access
    Expiration Time
    Format     string
}

// Default paste options. Expiration length is set to 'Never', Privacy is
// set to 'Public', and Format is set to 'text'.
func DefaultOptions() *PasteOptions {
    return &PasteOptions{Public, Never, "text"}
}

// Make an anonymous paste. An anonymous paste can only be Public or Unlisted.
func (pasteBin *PasteBinInfo) AnonymousPaste(content *string, options *PasteOptions) (*string, error) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_option", "paste")
    query.Add("api_paste_private", string(options.Privacy))
    query.Add("api_paste_expire_date", string(options.Expiration))
    query.Add("api_paste_format", options.Format)
    query.Add("api_paste_code", *content)

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
        return nil, &PasteError{"could not create paste: " + result}
    }

    return &result, nil
}

// Create a paste under a pastebin username. You must login with 'UserLogin' first.
func (pasteBin *PasteBinInfo) UserPaste(content *string, options *PasteOptions) (*string, error) {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_option", "paste")
    query.Add("api_user_key", pasteBin.UserKey)
    query.Add("api_paste_private", string(options.Privacy))
    query.Add("api_paste_expire_date", string(options.Expiration))
    query.Add("api_paste_format", options.Format)
    query.Add("api_paste_code", *content)

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
        return nil, &PasteError{"could not create paste: " + result}
    }

    return &result, nil
}

// Delete a paste owned by the logged in user, where pasteKey is the bit after
// pastebin.com/ in the address of the paste.
func (pasteBin *PasteBinInfo) DeletePaste(pasteKey string) error {
    query := url.Values{}
    query.Add("api_dev_key", pasteBin.APIKey)
    query.Add("api_user_key", pasteBin.UserKey)
    query.Add("api_paste_key", pasteKey)
    query.Add("api_option", "delete")

    resp, err := http.PostForm(pastebinURL, query)
    if err != nil {
        return &PasteError{"request failed: " + err.Error()}
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return &PasteError{"could not read response: " + err.Error()}
    }

    result := string(body)
    if strings.Contains(result, "Bad API request") {
        return &PasteError{"could not delete paste: " + result}
    }

    return nil
}
