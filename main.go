package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jakkrabig/captcha"
	"github.com/robbert229/jwt"
)

type apiHandler struct{}
type tokenHandler struct{}

func (apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path
	subStr := strings.Split(str, "/")
	fmt.Fprintln(w, getCaptcha(subStr[2]))

}

func (tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	phrase := r.Header.Get("token")
	subStr := strings.Split(phrase, "/")
	rs := captcha.Validate(subStr[0], subStr[1])

	if rs == false {
		w.WriteHeader(401)
		fmt.Fprint(w)
	}

	algorithm := jwt.HmacSha256("ri74krrIyXaQ1KAcZRhTUHak656iqTzQ")
	claims := jwt.NewClaim()
	claims.SetTime("exp", time.Now().Add(time.Minute*5))
	token, err := algorithm.Encode(claims)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, token)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/captcha/", apiHandler{})
	mux.Handle("/token/", tokenHandler{})
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){})
	http.ListenAndServe(":8000", mux)
}

func getCaptcha(str string) string {
	rs := strings.Split(str, "")
	v := make([]int, 4)
	for i, r := range rs {
		toInt, err := strconv.Atoi(r)
		if err != nil {
			log.Fatal(err)
		}
		v[i] = toInt
	}

	c := captcha.Captcha(v[0], v[1], v[2], v[3])
	return c.String()
}
