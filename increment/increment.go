package increment

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("increment", increment)
}

const incrementHTML = `
<html>
 <body>
  <p>count: %d</p>
 </body>
</html>
`

func increment(w http.ResponseWriter, r *http.Request) {
	count, code, err := increment2(w, r)
	if err != nil {
		w.WriteHeader(code)
		w.Write([]byte(err.Error()))
		fmt.Println(err.Error())
		return
	}
	html := fmt.Sprintf(incrementHTML, count)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func increment2(w http.ResponseWriter, r *http.Request) (int, int, error) {
	projectID, err := requiredEnv("PROJECT_ID")
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	sessionCollection, err := requiredEnv("SESSION_COLLECTION")
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	countCollection, err := requiredEnv("COUNT_COLLECTION")
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		return 0, http.StatusBadRequest, fmt.Errorf("sessionID cookie does not exist")
	}
	sessionID := cookie.Value

	userID, err := getUserID(r.Context(), projectID, sessionCollection, sessionID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	count, err := updateCount(r.Context(), projectID, countCollection, userID)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}

	return count, http.StatusOK, nil
}

type SessionData struct {
	UserID string `firestore:"userid"`
}

func getUserID(ctx context.Context, projectID, sessionCollection, sessionID string) (string, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	snap, err := client.Collection(sessionCollection).Doc(sessionID).Get(ctx)
	if err != nil {
		return "", fmt.Errorf("sessionID %s does not exist", sessionID)
	}
	var data SessionData
	if err := snap.DataTo(&data); err != nil {
		return "", fmt.Errorf("firestore.DocumentSnapshot.DataTo failed; %w", err)
	}
	return data.UserID, nil
}

type CountData struct {
	Count int `firestore:"count"`
}

func updateCount(ctx context.Context, projectID, countCollection, userID string) (int, error) {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return 0, fmt.Errorf("firestore.NewClient failed; %w", err)
	}
	defer client.Close()
	data := CountData{Count: 0}
	snap, err := client.Collection(countCollection).Doc(userID).Get(ctx)
	if err == nil {
		if err := snap.DataTo(&data); err != nil {
			return 0, fmt.Errorf("firestore.DocumentSnapshot.DataTo failed; %w", err)
		}
	}
	data.Count += 1
	if _, err := client.Collection(countCollection).Doc(userID).Set(ctx, data); err != nil {
		return 0, fmt.Errorf("firestore.DocumentRef.Set failed; %w", err)
	}
	return data.Count, nil
}

func requiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("you must define %s environment variable", key)
	}
	return val, nil
}
