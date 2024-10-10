package main

// short url checker
// todo: set run params for "deep" redirection level

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s [url]\n", os.Args[0])
		os.Exit(1)
	}
	fqdn := os.Args[1]
	getHead(fqdn)
}

func getHead(fqdn string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	count := 0
	fmt.Printf("[%d] original url:\t%s\n", count, fqdn)
	rdr := fqdn
	for {
		count++
		res, err := client.Head(rdr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer res.Body.Close()
		rdr = res.Header.Get("Location")
		if len(rdr) == 0 {
			break
		}
		fmt.Printf("[%d] redirect to:\t%+v\n", count, rdr)

		h, err := url.Parse(rdr)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		hostname := strings.TrimPrefix(h.Hostname(), "www.")
		fmt.Printf("WHOIS:\t\t%+v\n", hostname)
		fmt.Println()
		whoIs(hostname)
	}
}

func whoIs(domain string) {
	//fmt.Printf("[%s]", domain)
	con, err := net.Dial("tcp", "whois.iana.org:43")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer con.Close()
	con.SetDeadline(time.Now().Add(time.Second * 5))
	_, err = con.Write([]byte(domain + "\r\n"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(con)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
