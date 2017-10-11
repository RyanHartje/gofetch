package main

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRequest(t *testing.T){
    called := false
    test_func := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        called = true
    })
    w := &httptest.ResponseRecorder{}
    r, err := http.NewRequest("GET", "/", nil)
    if err != nil{
        t.Fatal(err)
    }
    test_func.ServeHTTP(w, r)

    // assert.Equal(t, false)
    // assert.NotNil(t, err)
    assert.True(t, called)
}

func TestRestrictedPage(t *testing.T){
    // w := &httptest.ResponseRecorder{}
    _, err := http.NewRequest("GET", "/edit", nil)
    if err != nil{
        t.Fatal(err)
    }
}

func TestLogin(t *testing.T){
    w := &httptest.ResponseRecorder{}
    _, err := http.NewRequest("GET", "/edit?username=test&password=testpass", nil)
    if err != nil{
        t.Fatal(err)
    }
    assert.True(t, strings.Contains(w.Body.String(), "logging in"))
    fmt.Println(w.Body.String())
}
