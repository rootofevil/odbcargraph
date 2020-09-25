package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart"
)

//exposes "chart"

var path string
var saveto string
var force bool

var files []template.URL

type CarData struct {
	Second float64
	Pid    string
	Value  float64
	Units  string
}

type CarDataSeries struct {
	Second []float64
	Pid    string
	Value  []float64
	Units  string
}

type FilesLinks struct {
	file template.URL
}

func main() {
	flag.StringVar(&path, "p", "2020-09-21 17-09-05(1).csv", "path to file")
	flag.StringVar(&saveto, "s", "result", "path to result")
	flag.BoolVar(&force, "f", false, "replace files")
	flag.Parse()
	savePath := saveto + "/" + strings.Split(path, ".")[0]
	fmt.Println("Save to: ", savePath)
	err := os.MkdirAll(savePath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	data := readCsv(path)

	for name, values := range data {
		// fmt.Println(name)
		filename := savePath + "/" + strings.ReplaceAll(name, "/", "") + ".png"
		files = append(files, template.URL(filename))
		if _, err := os.Stat(filename); os.IsNotExist(err) || force {
			fmt.Println(filename)
			buf := drawChart(name, values)
			saveChart(filename, buf)
		}
	}

	http.HandleFunc("/", httpHand)
	http.HandleFunc("/result/", httpImage)
	http.ListenAndServe(":3000", nil)
}

func readCsv(path string) map[string]CarDataSeries {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// FullData := make(map[string][]CarData)
	FullData := make(map[string]CarDataSeries)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		data := lineToData(scanner.Text())
		if _, ok := FullData[data.Pid]; !ok {
			FullData[data.Pid] = CarDataSeries{Pid: data.Pid, Units: data.Units}
		}
		temp := FullData[data.Pid]
		temp.Second = append(temp.Second, data.Second)
		temp.Value = append(temp.Value, data.Value)
		FullData[data.Pid] = temp
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	if _, ok := FullData["PID"]; ok {
		delete(FullData, "PID")
	}
	// fmt.Printf("%+v\n", FullData)
	return FullData
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

func drawChart(name string, values CarDataSeries) *bytes.Buffer {

	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: values.Second,
				YValues: values.Value,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})

	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Fatal(err)
	}

	return buffer
}

func saveChart(filename string, buffer *bytes.Buffer) {
	err := ioutil.WriteFile(filename, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func httpHand(w http.ResponseWriter, r *http.Request) {
	// t := template.Must(template.New("html-tmpl").Parse(`{{range $element := .}} {{$element}} {{end}}`))
	t, err := template.ParseFiles("page.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, files)
	if err != nil {
		log.Fatal(err)
	}
}

func httpImage(w http.ResponseWriter, r *http.Request) {
	// log.Println(r.RequestURI)
	fp, _ := url.QueryUnescape(strings.TrimLeft(r.RequestURI, "/"))
	// fp := "result/2020-09-21 17-09-05(1)/Использовано топлива.png"
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		log.Println("Can not find file")
	}
	http.ServeFile(w, r, fp)
}
