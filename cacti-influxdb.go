package main

import (
	"log"
	"strings"

	"github.com/bradleyfalzon/cacti-influxdb/lib/db"
	"github.com/influxdb/influxdb/client"
)

func main() {

	// Load config defa

	dsn := "user:pass@tcp(127.0.0.1:3306)/dbName?charset=utf8"
	rraPath := "/var/www/html/cacti/rra/"
	rrdtoolPath := "/usr/bin/rrdtool"

	dbCacti, err := db.Connect(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbCacti.Close()

	c, err := client.NewClient(&client.ClientConfig{
		Username: "root",
		Password: "root",
		Database: "cacti",
	})
	if err != nil {
		log.Fatal(err)
	}

	// connect to cacti, find the rrd files for each tree/leaf we're monitoring

	// Find the order_key for graph_tree
	graphTreeItemsIDs := []string{"17", "25", "9", "8"}

	// Find cacti rrd files
	dataSources, err := getCactiRRDs(dbCacti, graphTreeItemsIDs)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Found %d data sources\n", len(dataSources))

	// load cacti rrd file
	for _, dataSource := range dataSources {

		rrdLastUpdate, err := rrdLastUpdate(rrdtoolPath, replacePathRRA(rraPath, dataSource.Data_source_path))

		if err != nil {
			log.Printf("Failed parsing DS: %s, err: %s\n", dataSource.Data_source_path, err)
			continue
		}

		dsName := strings.Split(dataSource.Data_source_path, "/")

		// Update influxdb
		insertInflux(c, dsName[len(dsName)-1], rrdLastUpdate)

	}

	// fetch recent data

}
