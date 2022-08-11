package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
)

func main() {
	var plates arrayFlags
	flag.Var(&plates, "p", "a plate name to check (shorthand for plate)")
	flag.Var(&plates, "plate", "a plate name to check ")
	flag.Parse()

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()
	for _, plate := range plates {
		checkPlate(w, plate)
	}
}

func checkPlate(w *tabwriter.Writer, plate string) {
	encoded := url.QueryEscape(plate)
	resp, err := http.PostForm("https://personalizedplates.revenue.tn.gov/static/api/api.php",
		url.Values{
			"send[endpoint]": {fmt.Sprintf("/personalizedplates/verifyplate/%s/3210", encoded)},
			"send[type]":     {http.MethodGet}},
	)
	if err != nil {
		panic(err)
	}
	switch resp.StatusCode {
	case http.StatusOK:
		_, err := fmt.Fprintf(w, "%s:\tavailable\n", plate)
		if err != nil {
			panic(err)
		}
	case http.StatusUnprocessableEntity:
		printErr(w, plate, resp)
	case http.StatusBadRequest:
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		_, err = fmt.Fprintf(w, "%s:\tbad request -> %s\n", plate, string(bodyBytes))
		if err != nil {
			panic(err)
		}
	}
}

func printErr(w *tabwriter.Writer, plate string, resp *http.Response) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var errStruct []ErrStruct
	if err := json.Unmarshal(bodyBytes, &errStruct); err != nil {
		print(string(bodyBytes))
		panic(err)
	}
	_, err = fmt.Fprintf(w, "%s:\t%s\n", plate, errStruct[0].Err)
	if err != nil {
		panic(err)
	}
}

type ErrStruct struct {
	UserID int      `json:"userID"`
	Code   int      `json:"code"`
	Op     string   `json:"op"`
	Site   string   `json:"site"`
	Kind   int      `json:"kind"`
	Err    string   `json:"err"`
	Stack  []string `json:"stack"`
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
