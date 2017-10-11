package main

import (
    "crypto/subtle"
    "fmt"
    "html/template"
    "net/http"
    "sync"

    "github.com/satori/go.uuid"
)

var sessionStore map[string]Client
var storageMutex sync.RWMutex

type Client struct {
        loggedIn bool
}

type AuthenticationMiddleware struct {
	wrappedHandler http.Handler
}

// ServeHTTP is implemented so that AuthetnicationMiddleware can be used as an http.Handler
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

    // Now we check to see if this cookie is saved in our Map. If not, we generate a new one
    var present bool
    var client Client
    if cookie != nil {
        storageMutex.RLock()
        client, present = sessionStore[cookie.Value]
        storageMutex.RUnlock()
    } else {
        present = false
    }

	// If the cookie didn't exist, let's create it
    if present == false {
        cookie = &http.Cookie{
            Name: "session",
            Value: uuid.NewV4().String(),
        }
        client = Client{false}
        storageMutex.Lock()
        sessionStore[cookie.Value] = client
        storageMutex.Unlock()

	// Lastly, we need to set the cookie to the response writer, and handle redirects to either login or the original request uri
	    http.SetCookie(w, cookie)
        if client.loggedIn == false {
            t, _ := template.ParseFiles(
                "templates/base.gtpl",
                "templates/navbar.gtpl",
                "templates/login.gtpl",
            )

            t.Execute(w, nil)
            return
        }
        if client.loggedIn == true {
            h.wrappedHandler.ServeHTTP(w, r)
            return
        }
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

    err = r.ParseForm()
    if err != nil {
        fmt.Fprint(w, err)
        return
    }

    if subtle.ConstantTimeCompare([]byte(r.FormValue("password")), []byte("testpass")) == 1 {
        client.loggedIn = true
        fmt.Fprintln(w, "Thank you for logging in.")
        storageMutex.Lock()
        sessionStore[cookie.Value] = client
        storageMutex.Unlock()
    } else {
        fmt.Fprintln(w, "Wrong password.")
    }
}

