package main

import (
	"log"
	"time"

	"github.com/influxdb/influxdb/client"
)

func insertInflux(c *client.Client, tableName string, results rrdResults) (err error) {

	start := time.Now()

	results.dsK = append(results.dsK, "time")
	results.dsV = append(results.dsV, results.ts*1000)

	log.Printf("%#v", results.dsV)

	series := &client.Series{
		Name:    tableName,
		Columns: results.dsK,
		Points:  [][]interface{}{results.dsV},
	}
	if err := c.WriteSeries([]*client.Series{series}); err != nil {
		return err
	}

	log.Println("Wrote to data in:", time.Since(start))

	return

}
