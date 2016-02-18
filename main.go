package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tarm/serial"
)

type handler func(w http.ResponseWriter, r *http.Request)

var sc *SerialComs

var boud int
var name, port string

func init() {
	flag.IntVar(&boud, "boud", 115200, "Boud rate")
	flag.StringVar(&name, "name", "/dev/ttys000", "Serial port name")
	flag.StringVar(&port, "port", "8000", "Server Port number")
	flag.Parse()
}

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func GetOnly(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h(w, r)
			return
		}
		http.Error(w, "get only", http.StatusMethodNotAllowed)
	}
}

func PostOnly(h handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			h(w, r)
			return
		}
		http.Error(w, "post only", http.StatusMethodNotAllowed)
	}
}

func HandlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.PostForm.Get("data")
	log.Println(data)
	sc.Write(data)
	io.WriteString(w, fmt.Sprintf("sent %s\n", data))
}

func main() {
	sc = &SerialComs{}
	sc.OpenSerial(name, boud)

	http.HandleFunc("/send", PostOnly(HandlePost))
	http.ListenAndServe(":"+port, nil)
}

//SerialComs for serial comminucation
type SerialComs struct {
	// Connection
	Conn *serial.Port
}

func (sc *SerialComs) OpenSerial(name string, baud int) {
	c := &serial.Config{Name: name, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Serial Connected\n")
	sc.Conn = s
}

func (sc *SerialComs) Write(s string) {
	_, err := sc.Conn.Write([]byte(s))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Writing\n")

}

func (sc *SerialComs) Read() {
	buf := make([]byte, 128)
	n, err := sc.Conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%q", buf[:n])
}
