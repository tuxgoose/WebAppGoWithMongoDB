/*
 *  Based on the tutorial:
 *  http://reinbach.com/golang-webapps-1.html
 */

package main

import (
    "fmt"
    "log"
    "io"
    "time"
  	"net/http"
    "html/template"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/kelvins/webAppGoWithMongoDB/datastructure"
    "github.com/kelvins/webAppGoWithMongoDB/dbutil"
)

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

var c * mgo.Collection = nil

type Context struct {
    Author string
    Title string
    Content string
    DateTime time.Time
    Static string
}

func Home(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/"):]
    if req.Method == "GET" && title != "" {
        post := dbutil.Find(c, bson.M{"title": title})

        context := Context{Author: post.Author,
                           Title: post.Title,
                           Content: post.Content,
                           DateTime: post.DateTime}
        render(w, "index", context)
        return
    }
    // Load all posts
    post := dbutil.Find(c, bson.M{"title": "MySecondPost"})
    //posts := dbutil.FindAll(c)

    /*for _, x := range posts {
        fmt.Println("############################################")
        fmt.Println("Name:",    x.Name)
        fmt.Println("Surname:", x.Surname)
        fmt.Println("Address:", x.Address)
        fmt.Println("Phone:",   x.Phone)
        fmt.Println("############################################")
        fmt.Println("")
    }*/

    context := Context{Author: post.Author,
                       Title: post.Title,
                       Content: post.Content,
                       DateTime: post.DateTime}
    render(w, "index", context)
}

func Edit(w http.ResponseWriter, req *http.Request) {
    if req.Method == "GET" {
        context := Context{}
        render(w, "edit", context)
    } else {
        p := datastructure.Post{req.FormValue("author"),
                                   req.FormValue("title"),
                                   req.FormValue("content"),
                                   time.Now().UTC()}
        if p.Author != "" && p.Title != "" && p.Content != "" {
            dbutil.Insert(c, p)
            http.Redirect(w, req, "/"+req.FormValue("title"), http.StatusFound)
        }
    }
}

func render(w http.ResponseWriter, tmpl string, context Context) {
    context.Static = STATIC_URL
    tmpl_list := []string{"templates/base.html",
        fmt.Sprintf("templates/%s.html", tmpl)}
    t, err := template.ParseFiles(tmpl_list...)
    if err != nil {
        log.Println("Template parsing error: ", err)
    }
    err = t.Execute(w, context)
    if err != nil {
        log.Println("Template executing error: ", err)
    }
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
    static_file := req.URL.Path[len(STATIC_URL):]
    if len(static_file) != 0 {
        f, err := http.Dir(STATIC_ROOT).Open(static_file)
        if err == nil {
            content := io.ReadSeeker(f)
            http.ServeContent(w, req, static_file, time.Now(), content)
            return
        }
    }
    http.NotFound(w, req)
}

func main() {
    // Establish a connection, obtain a session
    session := dbutil.Connect("localhost")
    // Ensure that the session will be closed
    defer session.Close()

    c = session.DB("blog").C("posts")

    http.HandleFunc("/", Home)
    http.HandleFunc("/edit/", Edit)
    http.HandleFunc(STATIC_URL, StaticHandler)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
