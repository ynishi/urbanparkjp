package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"sync"

	"io/ioutil"

	"context"
	"time"

	"github.com/urfave/cli"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"

	"github.com/ynishi/urbanparkjp"
	"encoding/json"
)

func main() {

	app := cli.NewApp()
	app.Name = "urbanparkjp"
	app.Usage = "urbanparkjp [--file file.xml] [--dryrun] [--table talbename] [--region region] [--json]"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "dryrun, d",
			Usage: "dryrun(not put dynamo)",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Load xml from `FILE`",
			Value: "file.xml",
		},
		cli.StringFlag{
			Name:  "region, r",
			Usage: "Region for dynamo",
			Value: "ap-northeast-1",
		},
		cli.StringFlag{
			Name:  "table, t",
			Usage: "set dynamo table to `TABLE`",
			Value: "parks",
		},
		cli.BoolFlag{
			Name:  "json, j",
			Usage: "json output(not put dynamo)",
		},
	}

	app.Action = func(c *cli.Context) error {

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		var t dynamo.Table

		if !(c.GlobalBool("dryrun") || c.GlobalBool("json")) {
			sess, err := session.NewSession(&aws.Config{Region: aws.String(c.GlobalString("region"))})
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			log.Printf("region: %v\n", c.GlobalString("region"))

			dynamo.RetryTimeout = 30 * time.Minute
			db := dynamo.New(sess)
			t = db.Table(c.GlobalString("table"))
			log.Printf("table: %v\n", c.GlobalString("table"))
		}

		xmldoc, err := ioutil.ReadFile(c.GlobalString("file"))
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}

		ds := urbanparkjp.Dataset{}
		err = xml.Unmarshal(xmldoc, &ds)
		if err != nil {
			log.Println(err)
			os.Exit(3)
		}

		ps := make(map[string]*urbanparkjp.Posf64)
		for i := range ds.Points {
			posf64, err := urbanparkjp.PosToPosf64(ds.Points[i].Pos)
			ps[ds.Points[i].Id] = posf64
			if err != nil {
				log.Println(err)
				os.Exit(4)
			}
		}

		urbanparkjp.SetParksLoc(ds.Parks, ps)

		if c.GlobalBool("json") {
			jsonBytes, _ := json.Marshal(ds.Parks)
			fmt.Println(string(jsonBytes))
		} else {
			var wg sync.WaitGroup
			semaphore := make(chan int, 20)
			for _, park := range ds.Parks {
				wg.Add(1)
				go func(park2 urbanparkjp.Park) {
					defer wg.Done()
					semaphore <- 1
					if c.Bool("dryrun") {
						dry(park2, ps)
					} else {
						do(park2, ps, t)
					}
					<-semaphore
				}(park)
			}
			wg.Wait()
		}

		return nil
	}

	app.Run(os.Args)
}

func do(park urbanparkjp.Park, ps map[string]*urbanparkjp.Posf64, t dynamo.Table) {
	err := t.Put(park).Run()
	if err != nil {
		log.Println(err)
		os.Exit(5)
	}
	log.Printf("park: %v\n", park)
}

func dry(park urbanparkjp.Park, ps map[string]*urbanparkjp.Posf64) {
	fmt.Printf("park: %v\n", park)
}