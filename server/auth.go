package main

import (
    "crypto/subtle"
    "fmt"
    "html/template"
    "log"
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
            uri := "http://" + r.Host + "/login"
            http.Redirect(w, r, uri, 301)
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
    if r.Method == "GET" {
        t, _ := template.ParseFiles(
            "templates/base.gtpl",
            "templates/navbar.gtpl",
            "templates/login.gtpl",
        )
        if r.URL.Query().Get("success") != "" {
            t.Execute(w, map[string]string{"Message": "Login failed. Please try again."})
        } else {
            t.Execute(w, nil)
        }
    } else {
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

        // compare passwords here. definitely need to connect this to a datastore
        if subtle.ConstantTimeCompare([]byte(r.FormValue("password")), []byte("testpass")) == 1 {
            client.loggedIn = true
            storageMutex.Lock()
            sessionStore[cookie.Value] = client
            storageMutex.Unlock()

            // redirect to original page here
            log.Println("Redirecting after login")
            uri := "http://" + r.Host + "/edit"
            http.Redirect(w, r, uri, 301)
        } else {
            // If their login failed, we need to present a msg to let the user know
            uri := "http://" + r.Host + "/login?success=false"
            http.Redirect(w, r, uri, 301)
        }
    }
}

