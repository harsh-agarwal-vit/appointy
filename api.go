package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userdetails struct {
	uid      string
	name     string
	email    string
	password string
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello astaxie!") // write data to response
}

func main() {
	var detail userdetails

	http.HandleFunc("/", sayhelloName)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("login.gtpl")
			t.Execute(w, nil)
		} else {
			r.ParseForm()
			detail.uid = "debugger"
			detail.name = r.FormValue("name")
			detail.email = r.FormValue("email")
			detail.password = r.FormValue("password")
		}

	})
	err1 := http.ListenAndServe(":9090", nil)
	if err1 != nil {
		log.Fatal("ListenAndServe: ", err1)
	}

	fmt.Println(detail)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://harsh:1234567890@cluster0.crhis.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	instagramDatabase := client.Database("instagram")
	usersCollection := instagramDatabase.Collection("users")
	reults := usersCollection.InsertOne(ctx,detail)
}
