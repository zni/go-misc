package gopastebin

import (
    "net/url"
    "strings"
    "testing"
)

const apiKey = "Your API Key"
const username = "Your Username"
const password = "Your Password"

var paste *PasteBinInfo = PasteBin(apiKey)
var pasteURL *url.URL

func TestUserLogin(t *testing.T) {
    err := paste.UserLogin(username, password)
    var out error = nil
    if err != nil {
        t.Errorf("UserLogin(%v, %v) = %v, want %v", username, password, err, out)
    }
}

func TestUserInfo(t *testing.T) {
    userInfo, err := paste.UserInfo()
    var out error = nil
    if err != nil {
        t.Errorf("UserInfo() = %v, want %v", err, out)
    }

    if userInfo.Name != username {
        t.Errorf("UserInfo.Name = %v, want %v", userInfo.Name, username)
    }
}

func TestTrendingPastes(t *testing.T) {
    trendingPastes, err := paste.TrendingPastes()
    var out error = nil
    if err != nil {
        t.Errorf("TrendingPastes() = %v, want %v", err, out)
    }

    numberTrendingPastes := 18
    lengthTrendingPastes := len(trendingPastes.Pastes)
    if lengthTrendingPastes != numberTrendingPastes {
        t.Errorf("len(Pastes.Pastes) = %v, want %v", lengthTrendingPastes,
            numberTrendingPastes)
    }
}

func TestAnonymousPaste(t *testing.T) {
    options := &PasteOptions{Unlisted, TenMinutes, "text", "Untitled"}
    content := "Ein Zwei Drei November Oscar"
    _, err := paste.AnonymousPaste(&content, options)
    var out error = nil
    if err != nil {
        t.Errorf("AnonymousPaste(%v, %v) = %v, want %v", content,
            options, err, out)
    }
}

func TestUserPaste(t *testing.T) {
    options := &PasteOptions{Private, TenMinutes, "text", "Numbers"}
    content := "Ein Zwei Drei November Oscar"
    var out error = nil
    var err error
    pasteURL, err = paste.UserPaste(&content, options)
    if err != nil {
        t.Errorf("UserPaste(%v, %v) = %v, want %v", content,
            options, err, out)
    }
}

func TestDeletePaste(t *testing.T) {
    key := strings.TrimPrefix(pasteURL.Path, "/")
    err := paste.DeletePaste(key)
    var out error = nil
    if err != nil {
        t.Error("DeletePaste(%v) = %v, want %v", key, err, out)
    }
}
