
package dbutil

import (
    "log"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/kelvins/webAppGoWithMongoDB/datastructure"
)

func Connect(url string) *mgo.Session {
    session, err := mgo.Dial(url)
    if err != nil {
        log.Fatal(err)
    }
    return session
}

func Insert(c* mgo.Collection, p datastructure.Content) bool {
    err := c.Insert(&p)
    if err != nil {
        log.Println("Could not insert the object")
        log.Println(err)
        return false
    }
    return true
}

func FindAll(c* mgo.Collection) []datastructure.Content {
    var results []datastructure.Content
    err := c.Find(nil).All(&results)
    if err != nil {
        log.Println("Empty")
    }
    return results
}

func Find(c* mgo.Collection, query bson.M) datastructure.Content {
    var result datastructure.Content
    err := c.Find(query).One(&result)
    if err != nil {
        log.Println("Could not find the element")
    }
    return result
}
