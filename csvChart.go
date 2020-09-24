package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

//exposes "chart"

var path string

type CarData struct {
	Second float64
	Pid    string
	Value  float64
	Units  string
}

func main() {
	flag.StringVar(&path, "p", "2020-09-21 17-09-05(1).csv", "path to file")
	flag.Parse()
	readCsv(path)

	// x := []time.Time{

	// 	time.Now().AddDate(0, 0, -3),
	// 	time.Now().AddDate(0, 0, -2),
	// 	time.Now().AddDate(0, 0, -1),
	// 	time.Now()}

	// y := []float64{1.0, 2.0, 3.0, 4.0}
	// graph := chart.Chart{
	// 	Series: []chart.Series{
	// 		chart.TimeSeries{
	// 			XValues: x,
	// 			YValues: y,
	// 		},
	// 	},
	// }

	// // fmt.Println(graph.Series[0]

	// buffer := bytes.NewBuffer([]byte{})

	// err := graph.Render(chart.PNG, buffer)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = ioutil.WriteFile("test.png", buffer.Bytes(), 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func readCsv(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	FullData := make(map[string][]CarData)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := lineToData(scanner.Text())
		FullData[data.Pid] = append(FullData[data.Pid], data)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", FullData)
}

func lineToData(line string) CarData {
	data := CarData{}
	array := strings.Split(line, ";")
	if s, err := strconv.ParseFloat(strings.Trim(array[0], "\""), 64); err == nil {
		data.Second = s
	}
	data.Pid = strings.Trim(array[1], "\"")
	if s, err := strconv.ParseFloat(strings.Trim(array[2], "\""), 64); err == nil {
		data.Value = s
	}
	data.Units = strings.Trim(array[3], "\"")

	return data
}
