package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoLang_Dungeon/db"
	"github.com/gorilla/mux"
	"net/http"
)

func BookHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/first/{fname}/last/{lname}/hobby/{hname}", BookHandler).Methods("POST")

	http.ListenAndServe(":3000", r)
}
