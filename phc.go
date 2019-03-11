package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "phc"
	app.Usage = "http client for p"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "file name",
		},
		cli.StringFlag{
			Name:  "page,  p",
			Usage: "page count",
		},
		cli.StringFlag{
			Name:  "type, t",
			Value: "text",
			Usage: "type; (text|line)",
		},
	}

	app.Action = func(c *cli.Context) error {
		fileArg := c.String("file")
		pageArg := c.String("page")
		typeArg := c.String("type")
		if fileArg == "" || pageArg == "" || typeArg == "" {
			cli.ShowAppHelpAndExit(c, 1)
		} else if typeArg != "line" && typeArg != "text" {
			cli.ShowAppHelpAndExit(c, 1)
		}

		pageCount, err := strconv.Atoi(pageArg)
		if err != nil || pageCount < 1 {
			cli.ShowAppHelpAndExit(c, 1)
		}

		return exec(fileArg, pageArg, typeArg)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func exec(file string, page string, typeValue string) error {
	targetURL := "http://localhost:8080/" + typeValue
	values := url.Values{}
	values.Add("f", file)
	values.Add("p", page)

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return err
	}

	req.URL.RawQuery = values.Encode()
	client := &http.Client{
		Timeout: time.Duration(3) * time.Minute,
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	readResponse(res)

	return nil
}

func readResponse(resp *http.Response) {
	b, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println(string(b))
	}
}
