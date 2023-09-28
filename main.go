package main

import (
	"auth/models"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var usersCollection *mongo.Collection

func main() {
	godotenv.Load()
	URI := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(URI))
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected successfully")
	usersCollection = client.Database("authentication").Collection("users")
	fs := http.FileServer(http.Dir("./view"))
	http.Handle("/", fs)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/failed-signin", failedSignin)
	http.HandleFunc("/profile", successSignin)
	log.Fatal(http.ListenAndServe(":5500", nil))
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	fullName := r.FormValue("fullName")
	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")
	user := models.User{
		FullName: fullName,
		Email:    email,
		Password: password,
		Username: username,
	}
	var search models.User
	if usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&search); search.Username == username {
		fmt.Fprintln(w, "Username already taken")
		return
	}
	if usersCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&search); search.Email == email {
		fmt.Fprintln(w, "Email already taken")
		return
	}
	_, err := usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	http.ServeFile(w, r, "./view/index.html")
}

func failedSignin(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./view/index.html")
}

var data models.PageTemplate

func successSignin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./view/profile.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	var user models.User
	err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil || user.Username != username || user.Password != password {
		http.Redirect(w, r, "/failed-signin", http.StatusTemporaryRedirect)
		return
	}
	data = models.PageTemplate{
		Username: user.Username,
		Fullname: user.FullName,
		Email:    user.Email,
	}
	http.Redirect(w, r, "/profile", http.StatusTemporaryRedirect)
}
