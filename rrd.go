package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type rrdResults struct {
	ts  uint64
	dsK []string
	//dsV []float64
	dsV []interface{}
}

func rrdLastUpdate(rrdtoolPath, filename string) (results rrdResults, err error) {

	//output, err := exec.Command(cmd).Output()
	output, err := exec.Command(rrdtoolPath, "lastupdate", filename).Output()
	if err != nil {
		log.Printf("Command failed: %s lastupdate %s\n", rrdtoolPath, filename)
		return
	}

	// Remove leading space and trailing newline
	outStr := strings.Trim(string(output), " \n")

	// remove double newline
	outStr = strings.Replace(outStr, "\n\n", "\n", 1)

	outSplit := strings.Split(outStr, "\n")

	header := strings.Split(outSplit[0], " ")

	// split body to separate timestamp from results
	body := strings.Split(outSplit[1], ": ")

	// Get the timestamp
	results.ts, err = strconv.ParseUint(body[0], 10, 32)
	if err != nil {
		return
	}

	rraValues := strings.Split(body[1], " ")

	for k, v := range header {

		rraValue, err := strconv.ParseFloat(rraValues[k], 64)
		if err != nil {
			return results, err
		}

		results.dsK = append(results.dsK, v)
		results.dsV = append(results.dsV, rraValue)

	}

	return
}

func replacePathRRA(rraPath, dataSource string) string {
	rraPath = strings.TrimRight(rraPath, "/")
	return strings.Replace(dataSource, "<path_rra>", rraPath, 1)
}
