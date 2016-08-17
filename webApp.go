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
    "strings"
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
    Posts []datastructure.Post
    Static string
}

func Home(w http.ResponseWriter, req *http.Request) {
    // Load all posts
    posts := dbutil.FindAll(c)

    var content Context
    for _, x := range posts {
        context := datastructure.Post{Author: x.Author,
                                      Title: x.Title,
                                      Content: x.Content[0:100] + "...",
                                      DateTime: x.DateTime}
        content.Posts = append(content.Posts, context)
    }

    render(w, "index", content)
}

func View(w http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[len("/view/"):]
    if req.Method == "GET" && title != "" {
        post := dbutil.Find(c, bson.M{"title": title})
        var content Context
        context := datastructure.Post{Author: post.Author,
                                      Title: post.Title,
                                      Content: post.Content,
                                      DateTime: post.DateTime}
        content.Posts = append(content.Posts, context)
        render(w, "view", content)
        return
    }
    http.Redirect(w, req, "/", http.StatusFound)
}

func Edit(w http.ResponseWriter, req *http.Request) {
    if req.Method == "GET" {
        var content Context
        context := datastructure.Post{Title: ""}
        content.Posts = append(content.Posts, context)
        render(w, "edit", content)
    } else {
        p := datastructure.Post{req.FormValue("author"),
                                req.FormValue("title"),
                                eq.FormValue("content"),
                                time.Now().UTC()}
        if p.Author != "" && p.Title != "" && p.Content != "" {
            dbutil.Insert(c, p)
            http.Redirect(w, req, "/view/"+req.FormValue("title"), http.StatusFound)
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
    http.HandleFunc("/view/", View)
    http.HandleFunc(STATIC_URL, StaticHandler)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
