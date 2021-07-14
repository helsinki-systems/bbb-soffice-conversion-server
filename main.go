package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/acaloiaro/go-libreofficekit"
)

var office *libreofficekit.Office
var officeMux = sync.Mutex{}

func httpError(logMsg string, w http.ResponseWriter) {
	http.Error(w, "Internal Server Error", 500)
	log.Print(logMsg)
}

func tempFile(prefix string, w http.ResponseWriter) (name string, file *os.File) {
	f, _ := ioutil.TempFile("", prefix)
	fn := f.Name()
	log.Print(fn)
	return fn, f
}

func convert(w http.ResponseWriter, r *http.Request) {
	inFileName, inFile := tempFile("bbb-conv-in-", w)
	outFileName, outFile := tempFile("bbb-conv-out-", w)
	defer inFile.Close()
	defer outFile.Close()
	defer os.Remove(inFileName)
	defer os.Remove(outFileName)

	r.ParseMultipartForm(128000000) // use up to 128MB of RAM, before creating temporary files
	formFile, _, err := r.FormFile("file")
	if err != nil {
		httpError("Error getting document from request", w)
		return
	} else {
		log.Print("Got file input from form request")
	}
	formFileBytes, err := ioutil.ReadAll(formFile)
	inFile.Write(formFileBytes)

	officeMux.Lock()
	defer officeMux.Unlock()
	document, err := office.LoadDocument(inFile.Name())
	if err != nil {
		httpError("Failed to open document", w)
		return
	} else {
		log.Print("Document loaded successfully")
	}

	t := r.FormValue("type")
	if t == "" {
		t = "pdf"
	}

	err = document.SaveAs(outFileName, t, "")
	if err != nil {
		httpError("Failed to save document", w)
		document.Close()
		return
	} else {
		log.Print("Document saved successfully")
	}

	out, err := ioutil.ReadFile(outFileName)
	if err != nil {
		httpError("Failed to load converted document", w)
		document.Close()
		return
	} else {
		log.Printf("Successfully read converted document")
	}

	_, err = w.Write(out)
	if err != nil {
		log.Print("Failed to write http response")
	} else {
		log.Print("Converted document served")
	}
	document.Close()
}

func main() {
	var socket, lodir string

	flag.StringVar(&socket, "s", "/run/bbb-soffice-conversion-server/sock", "IP to listen on")
	flag.StringVar(&lodir, "lo", "libreoffice-directory", "Directory from which to load LibreOffice installation")

	flag.Parse()

	var err error
	office, err = libreofficekit.NewOffice(lodir)
	if err != nil {
		log.Fatal("Failed to load LibreOffice")
	} else {
		log.Print("Loaded LibreOffice")
	}
	defer office.Close()

	http.HandleFunc("/", convert)
	log.Print(socket)
	server := http.Server{}
	os.Remove(socket)
	unixListener, err := net.Listen("unix", socket)
	if err != nil {
		log.Fatal("Error creating socket: ", err)
		return
	}
	server.Serve(unixListener)

}
