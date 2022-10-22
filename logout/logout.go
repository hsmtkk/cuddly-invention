package logout

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/hsmtkk/cuddly-invention/common"
)

func init() {
	functions.HTTP("logout", logout)
}

func logout(w http.ResponseWriter, r *http.Request) {
	code, err := logout2(w, r)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
		return
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func logout2(w http.ResponseWriter, r *http.Request) (int, error) {
	projectID, err := common.RequiredEnv("PROJECT_ID")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	sessionCollection, err := common.RequiredEnv("SESSION_COLLECTION")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	client, err := firestore.NewClient(r.Context(), projectID)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return http.StatusBadRequest, fmt.Errorf("sessionID cookie does not exist")
	}
	sessionID := cookie.Value

	if _, err := client.Collection(sessionCollection).Doc(sessionID).Delete(r.Context()); err != nil {
		return 0, fmt.Errorf("firestore.DocumentRef.Delete failed; %w", err)
	}

	return 0, nil
}
