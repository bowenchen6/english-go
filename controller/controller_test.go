// controller_test.go
package controller_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"

	"golang.org/x/net/publicsuffix"
)

func TestGetCookie(t *testing.T) {
	// Start a server to give us cookies.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie("Flavor"); err != nil {
			http.SetCookie(w, &http.Cookie{Name: "Flavor", Value: "Chocolate Chip"})
		} else {
			cookie.Value = "Oatmeal Raisin"
			http.SetCookie(w, cookie)
		}
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	if _, err = client.Get(u.String()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("After 1st request:")
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	if _, err = client.Get(u.String()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("After 2nd request:")
	for _, cookie := range jar.Cookies(u) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
}
