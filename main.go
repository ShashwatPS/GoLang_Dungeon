package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoLang_Dungeon/db"
	"github.com/gorilla/mux"
	"net/http"
)

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Hobby     string `json:"hobby"`
}

type Hobby struct {
	Hname string `json:"hname"`
}

// Route Handlers
func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fname := vars["fname"]
	lname := vars["lname"]
	hname := vars["hname"]

	if err := saveToDataBase(fname, lname, hname); err != nil {
		panic(err)
	}

	responseData := map[string]string{
		"Fiest Name": fname,
		"Last Name":  lname,
		"Hobby":      hname,
	}

	jsonResponse, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func FetchingUsers(w http.ResponseWriter, r *http.Request) {
	var user Hobby
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	HobbyName := user.Hname
	users, err := getAllUsers(HobbyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func GetUserFromBody(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	firstName := user.FirstName
	lastName := user.LastName
	hobby := user.Hobby

	if err := saveToDataBase(firstName, lastName, hobby); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User Posted Successfully"))
}

// Database Operations
func saveToDataBase(fname string, lname string, hobby string) error {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	createdUser, err := client.User.CreateOne(
		db.User.FirstName.Set(fname),
		db.User.LastName.Set(lname),
		db.User.Hobby.Set(hobby),
	).Exec(ctx)

	if err != nil {
		return err
	}

	result, _ := json.MarshalIndent(createdUser, "", "  ")
	fmt.Printf("created post: %s\n", result)

	return nil
}

func getAllUsers(hname string) ([]User, error) {
	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return nil, err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()

	users, err := client.User.FindMany(
		db.User.Hobby.Equals(hname),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		fmt.Println("No users found with hobby:", hname)
		return nil, nil
	}

	var mappedUsers []User
	for _, u := range users {
		mappedUsers = append(mappedUsers, User{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Hobby:     u.Hobby,
		})
	}

	result, _ := json.MarshalIndent(users, "", "  ")
	fmt.Printf("All users: %s\n", result)

	return mappedUsers, nil
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/first/{fname}/last/{lname}/hobby/{hname}", UserHandler).Methods("POST")
	r.HandleFunc("/users", FetchingUsers).Methods("GET")
	r.HandleFunc("/addUser", GetUserFromBody).Methods("POST")

	http.ListenAndServe(":4000", r)
}
