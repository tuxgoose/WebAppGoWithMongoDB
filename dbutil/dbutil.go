
package dbutil

import (
    "log"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/kelvins/webAppGoWithMongoDB/datastructure"
)

// Stablish a connection
func Connect(url string) *mgo.Session {
    session, err := mgo.Dial(url)
    if err != nil {
        log.Fatal(err)
    }
    return session
}

// Insert the object into the collection
func Insert(c* mgo.Collection, p datastructure.Post) bool {
    err := c.Insert(&p)
    if err != nil {
        log.Println("Error in the Insert function: ", err)
        return false
    }
    return true
}

// Get all the data from the collection
func FindAll(c* mgo.Collection) []datastructure.Post {
    var results []datastructure.Post
    err := c.Find(nil).All(&results)
    if err != nil {
        log.Println("Error ni the FindAll function: ", err)
    }
    return results
}

// Search for an element
func Find(c* mgo.Collection, query bson.M) datastructure.Post {
    var result datastructure.Post
    err := c.Find(query).One(&result)
    if err != nil {
        log.Println("Error in the Find function: ", err)
    }
    return result
}
