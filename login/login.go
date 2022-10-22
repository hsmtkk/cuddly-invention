package login

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	"github.com/hsmtkk/cuddly-invention/common/env"
	"github.com/hsmtkk/cuddly-invention/common/model"
)

func init() {
	functions.HTTP("login", login)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handleGet(w, r)
	} else if r.Method == http.MethodPost {
		handlePost(w, r)
	}
}

const loginHTML = `
<html>
 <body>
  <form method="POST" action="/login">
   <input type="text" name="userID">
   <input type="submit" value="submit">
  </form>
 </body>
</html>
`

func handleGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(loginHTML)); err != nil {
		fmt.Printf("http.ResponseWriter.Write failed; %s\n", err.Error())
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	sessionID, code, err := handlePost2(w, r)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
		return
	}
	cookie := http.Cookie{
		Name:  "sessionID",
		Value: sessionID,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/increment", http.StatusMovedPermanently)
}

func handlePost2(w http.ResponseWriter, r *http.Request) (string, int, error) {
	userID := r.FormValue("userID")
	if userID == "" {
		return "", http.StatusBadRequest, fmt.Errorf("userID does not exist")
	}
	sessionID := uuid.NewString()

	projectID, err := env.RequiredEnv("PROJECT_ID")
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	collection, err := env.RequiredEnv("SESSION_COLLECTION")
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	client, err := firestore.NewClient(r.Context(), projectID)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	defer client.Close()

	if _, _, err := client.Collection(collection).Add(r.Context(), model.SessionModel{
		UserID: userID,
	}); err != nil {
		return "", http.StatusInternalServerError, err
	}
	return sessionID, http.StatusMovedPermanently, nil
}
