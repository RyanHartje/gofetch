package main

import (
   "net/http" 
   "sync"
)

var sessionStore map[string]Client
var storageMutex sync.RWMutex

type Client struct {
        loggedIn bool
}

type AuthenticationMiddleware struct {
	wrappedHandler http.Handler
}

func (h AuthenticationMiddleware) ServeHTTP (w http.ResponseWriter, r *http.Request){
    cookie, err := r.Cookie("session")
    if err != nil {
        if err != http.ErrNoCookie {
            fmt.Fprint(w, err)
            return
        } else {
            err = nil
        }

}

func Authenticate(h http.Handler) AuthenticationMiddleware {
	return AuthenticationMiddleware{h}
}

func ProcessLogin(w http.ResponseWriter, r *http.Request){
    cookie, err := r.Cookie("session")
    if err != nil {
        if err != http.ErrNoCookie {
            fmt.Fprint(w, err)
            return
        } else {
            err = nil
        }
    }
    var present bool
    var client Client
    if cookie != nil {
        storageMutex.RLock()
        client, present = sessionStore[cookie.Value]
        storageMutex.RUnlock()
    } else {
        present = false
    }

    if present == false {
        cookie = &http.Cookie{
            Name: "session",
            Value: uuid.NewV4().String(),
        }
        client = Client{false}
        storageMutex.Lock()
        sessionStore[cookie.Value] = client
        storageMutex.Unlock()
    }
    http.SetCookie(w, cookie)
}

