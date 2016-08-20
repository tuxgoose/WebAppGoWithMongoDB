/*
 *  Example of a blog using GoLang and MongoDB with MGO driver.
 *
 *  This example was developed based on the "Golang Web Apps" tutorial amongst others.
 *  Link: http://reinbach.com/golang-webapps-1.html
 */

package main

// Import all necessary packages
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

// "Path" to the static files (e.g. CSS files, JS files, Images)
const STATIC_URL string = "/static/"
const STATIC_ROOT string = "static/"

var c * mgo.Collection = nil

// Context struct based on the "Post" struct from the datastructure package
type Context struct {
    Posts []datastructure.Post
    Static string
}

/*
 *  Function for the Home page
 *  The home page should show a list with all posts in the database
 */
func Home(w http.ResponseWriter, req *http.Request) {
    // Load all posts
    posts := dbutil.FindAll(c)

    // Creates an "object" of the Context structure
    var context Context

    // For each post append its values in the "Posts" vector from the Context structure
    for _, x := range posts {
        content := datastructure.Post{Author:   x.Author,
                                      Title:    x.Title,
                                      Content:  x.Content[0:100] + "...", // Shows a summary of the post (only the first 100 characters)
                                      DateTime: x.DateTime}
        // Append the current post
        context.Posts = append(context.Posts, content)
    }

    // Calls the function to render the home page (index)
    render(w, "index", context)
}

/*
 *  Function for the View page
 *  The View page should show a single post (title, content and author)
 */
func View(w http.ResponseWriter, req *http.Request) {
    // Gets the post title passed by the GET method
    title := req.URL.Path[len("/view/"):]

    // If the request method is GET and the title is not empty
    if req.Method == "GET" && title != "" {
        // Search in the database based on the title of the post
        post := dbutil.Find(c, bson.M{"title": title})

        // Creates an "object" of the Context structure
        var context Context

        // Get the values from the post variable
        content := datastructure.Post{Author:   post.Author,
                                      Title:    post.Title,
                                      Content:  post.Content,
                                      DateTime: post.DateTime}
        // Append the post content
        context.Posts = append(context.Posts, content)

        // Calls the function to render the View page
        render(w, "view", context)

        return
    }
    // Else, redirect to the home page (index)
    http.Redirect(w, req, "/", http.StatusFound)
}

/*
 *  Function for the Edit page
 *  The edit page should provide a way to the user enter with the post data (title, content and author)
 */
func Edit(w http.ResponseWriter, req *http.Request) {
    // If the method is GET
    if req.Method == "GET" {
        // Creates an "object" of the Context structure
        var context Context
        // Creates an empty "object"
        content := datastructure.Post{Title: ""}
        context.Posts = append(context.Posts, content)
        // Calls the function to render the Edit page
        render(w, "edit", content)
    } else {
        // If the method is not GET
        // Creates an object with the values passed by the user over the http.Request
        p := datastructure.Post{req.FormValue("author"),
                                req.FormValue("title"),
                                req.FormValue("content"),
                                time.Now().UTC()}
        // If the title, content and author are filled, insert the post in the database
        if p.Author != "" && p.Title != "" && p.Content != "" {
            // Insert the post "object" in the collection
            dbutil.Insert(c, p)
            // Redirect to the View page, passing by the GET method the title of the inserted post
            http.Redirect(w, req, "/view/"+req.FormValue("title"), http.StatusFound)
        }
    }
}

/*
 *  Function used to render the pages
 */
func render(w http.ResponseWriter, tmpl string, context Context) {
    // Fill the static variable in the context struct (passed by parameter)
    context.Static = STATIC_URL
    // Creates a template list, based on the base template and the template passed by parameter
    tmpl_list := []string{"templates/base.html",
        fmt.Sprintf("templates/%s.html", tmpl)}
    // Creates a new template and parses the template definitions from the named files
    t, err := template.ParseFiles(tmpl_list...)
    // If any error occurs, show it
    if err != nil {
        log.Println("Template parsing error: ", err)
    }
    // Applies a parsed template to the specified data object
    err = t.Execute(w, context)
    // If any error occurs, show it
    if err != nil {
        log.Println("Template executing error: ", err)
    }
}

/*
 *  Function used to deal with the static files
 */
func StaticHandler(w http.ResponseWriter, req *http.Request) {
    static_file := req.URL.Path[len(STATIC_URL):]
    // If the path is not empty
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

/*
 *  Main function: connects to the database and opens a session to work in the "posts" collection.
 */
func main() {
    // Establish a connection, obtain a session
    session := dbutil.Connect("localhost")
    // Ensure that the session will be closed
    defer session.Close()

    c = session.DB("blog").C("posts")

    // Assigns each page to each function
    http.HandleFunc("/", Home)
    http.HandleFunc("/edit/", Edit)
    http.HandleFunc("/view/", View)
    http.HandleFunc(STATIC_URL, StaticHandler)

    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
