package model

type SessionModel struct {
	UserID string `firestore:"userid"`
}

type CountModel struct {
	Count int `firestore:"count"`
}