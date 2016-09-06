package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"server/models"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/julienschmidt/httprouter"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

// NewUserController returns a instance of UserController structure
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

// GetUser retrives an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub a user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB("rest_example").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	fmt.Fprintf(w, "%s", uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Stub a user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	// Add an Id
	u.ID = bson.NewObjectId()

	// Add default avatar
	u.Avatar = "default.png"

	// Write the user to mongo
	uc.session.DB("rest_example").C("users").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	fmt.Fprintf(w, "%s", uj)
}

// UpdateAvatar update a user avatar
func (uc UserController) UpdateAvatar(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Stub a user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB("rest_example").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		fmt.Println(err)
		return
	}

	// the FormFile function takes in the POST input id file
	file, header, err := r.FormFile("avatar")
	if err != nil {
		fmt.Println(err)
		fmt.Println(header)
		fmt.Println(file)
		return
	}
	defer file.Close()

	if err := os.Remove(fmt.Sprintf("./avatar/%s", u.Avatar)); err != nil {
		fmt.Println(err)
	}

	avatarArr := strings.Split(header.Filename, ".")

	avatarFile := fmt.Sprintf("%s.%s", id, avatarArr[len(avatarArr)-1])

	out, err := os.Create("./avatar/" + avatarFile)
	if err != nil {
		fmt.Println("Unable to create the file for writing. Check your write access privilege")
		return
	}
	defer out.Close()

	// write the content from POST to the file
	_, err = io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
	}

	colQuerier := bson.M{"_id": oid}
	update := bson.M{"$set": bson.M{"avatar": avatarFile}}

	// Update avatar
	if err := uc.session.DB("rest_example").C("users").Update(colQuerier, update); err != nil {
		w.WriteHeader(404)
		fmt.Println(err)
		return
	}

	// Fetch user
	if err := uc.session.DB("rest_example").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		fmt.Println(err)
		return
	}

	// header.Filename
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	fmt.Fprintf(w, "%s", uj)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("id")

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(id)

	// Remove user
	if err := uc.session.DB("rest_example").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}