package cookietest

import (
	"net/http"
	"strconv"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("cookietest", cookieTest)
}

func cookieTest(w http.ResponseWriter, r *http.Request) {
	countCookie, err := r.Cookie("count")
	if err == nil {
		count, err := strconv.Atoi(countCookie.Value)
		if err != nil {
			count = 1
		}
		countCookie.Value = strconv.Itoa(count + 1)
		http.SetCookie(w, countCookie)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("count: " + strconv.Itoa(count)))
		return
	} else {
		countCookie = &http.Cookie{
			Name:  "count",
			Value: "1",
		}
		http.SetCookie(w, countCookie)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("count: 1"))
	}
}
