package main

// npm start

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"

	"fileo"
	"foldero"
	"imageo"
	"stringo"
)

var (
	fILEtYPES  = []string{".jpg", ".jpeg", ".png", ".gif", ".webm", ".mp4", ".mov", ".ico"}
	defaultCol = "3"
	pageFiles  = 40
)

// Row ..
type Row []Element

// Element ...
type Element struct {
	Path     string
	File     string
	FileName string
	Height   int
	Width    int
	IsImage  bool
	Columns  int
	Type     string
}
type data struct {
	Folders    []foldero.FolderInfo
	Row        []Element
	Columns    int
	Path       string
	BackPath   string
	Page       int
	TotalPages int
	NextPage   int
	PrevPage   int
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func removeLastPath(path string) string {
	for i := len(path); i > 1; i-- {
		// We skip the first /.
		if string(path[i-2]) == "/" {
			return path[:i-1]
		}
	}
	return "/"
}

func getQueryValue(values url.Values, value, defaultValue string) string {
	if val, ok := values[value]; ok {
		if len(val) > 0 {
			return val[0]
		}
	}

	return defaultValue
}

func getFilesInFolder(filePath string, columns int, page int) []Element {
	files, err := os.ReadDir(filePath)
	check(err)

	fFiles := fileo.FilterFiles(files, fILEtYPES)
	numFiles := len(fFiles)

	if numFiles == 0 {
		return []Element{}
	}

	numFilesOnPage := pageFiles

	// if there is not enough files to fill page, set number of files to how many files can fit.
	if (numFiles / pageFiles) <= page {
		numFilesOnPage = numFiles - (pageFiles * (page - 1))
		if numFilesOnPage <= 0 {
			return []Element{}
		}
	}

	row := make([]Element, numFilesOnPage)

	var wg sync.WaitGroup
	var mux sync.Mutex
	start := pageFiles * (page - 1)
	end := (pageFiles * (page - 1)) + numFilesOnPage
	for i, f := range fFiles[start:end] {
		wg.Add(1)
		go handle(&row, i, columns, f.Name(), filePath, &mux, &wg)
	}

	wg.Wait()
	return row
}

func handle(row *([]Element), index int, columns int, fileName, filePath string, mux *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	if stringo.StringInSlice(fileName, []string{".jpg", ".jpeg", ".png"}) {

		fileo.ReplaceCharacterInFileName(fmt.Sprintf("%s%s", filePath, fileName), "+", "_")
		fileo.ReplaceCharacterInFileName(fmt.Sprintf("%s%s", filePath, fileName), " ", "_")

		// removed temporarily to find a better solution
		// if _, err := os.Stat(fmt.Sprintf("%sthumbnail/%s", filePath, fileName)); os.IsNotExist(err) {
		// 	imageo.CreateThumbnail(filePath, fileName, filePath+"thumbnail/", 360, 360)
		// }

		height, width := imageo.GetImageDimensions(fileName, filePath)

		if height == 0 {
			height = 1
		}
		if width == 0 {
			width = 1
		}

		mux.Lock()
		(*row)[index] = Element{
			Path:     filePath,
			FileName: fileName,
			File:     filePath,
			Height:   height,
			Width:    width,
			Columns:  columns,
			IsImage:  true,
			Type:     "",
		}
		mux.Unlock()
		return
	}

	fileType := ""
	if stringo.StringContains(fileName, ".mp4") || stringo.StringContains(fileName, ".mov") {
		fileType = "video/mp4"
	} else if stringo.StringContains(fileName, ".webm") {
		fileType = "video/webm"
	}

	if stringo.StringInSlice(fileName, []string{".mp4", ".mov", ".webm"}) {

		mux.Lock()
		(*row)[index] = Element{
			Path:     filePath,
			FileName: fileName,
			Height:   0,
			Width:    0,
			IsImage:  false,
			Type:     fileType,
		}
		mux.Unlock()
	}
	return
}

// Routing

// Server implements the web server specification found at
type Server struct {
	sync.Mutex
	*mux.Router
	path   string
	folder string
}

type spaHandler struct {
	tpl    *template.Template
	server *Server
}

func (s *Server) stylesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		style, _ := mux.Vars(r)["style"]
		css, err := os.Open(fmt.Sprintf("./src/styles/%s", style))
		if err != nil {
			log.Fatal(err)
		}
		defer css.Close()
		w.Header().Set("Content-Type", "text/css")
		io.Copy(w, css)
	}
}

func (s *Server) scriptsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		script, _ := mux.Vars(r)["script"]
		src, err := os.Open(fmt.Sprintf("./src/scripts/%s", script))
		if err != nil {
			log.Fatal(err)
		}
		defer src.Close()
		w.Header().Set("Content-Type", "text/javascript")
		io.Copy(w, src)
	}
}

func (s *Server) imagesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		img, _ := mux.Vars(r)["img"]
		image, err := os.Open(fmt.Sprintf("./src/images/%s", img))
		if err != nil {
			log.Fatal(err)
		}
		defer image.Close()
		w.Header().Set("Content-Type", "image/png")
		io.Copy(w, image)
	}
}

func (s *Server) exitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		os.Exit(0)
	}
}

func (s spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	fmt.Println("URL Path: ", path)
	if path[len(path)-1] != '/' {
		path += "/"
	}
	s.server.path = path

	values := r.URL.Query()

	// path := getQueryValue(values, "path", pATH)
	columns, err := strconv.Atoi(getQueryValue(values, "columns", defaultCol))
	check(err)
	page, err := strconv.Atoi(getQueryValue(values, "page", "1"))
	check(err)

	files, err := os.ReadDir(path)
	check(err)

	totalPages := int(math.Ceil(float64(fileo.NumberOfFiles(files, fILEtYPES)) / float64(pageFiles)))
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = 1
	}
	prevPage := (page - 1)
	if prevPage <= 0 {
		prevPage = totalPages
	}

	data := data{
		Folders:    foldero.GetFolderInfo(path, fILEtYPES),
		Row:        getFilesInFolder(path, columns, page),
		Columns:    columns,
		Path:       path,
		BackPath:   removeLastPath(path),
		Page:       page,
		TotalPages: totalPages,
		NextPage:   nextPage,
		PrevPage:   prevPage,
	}

	s.tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func (s *Server) deleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file := r.FormValue("fileName")
		fmt.Println("File: ", file)
		fmt.Println("Delete file: ", s.path+file)
		err := fileo.RemoveFile(s.path + file)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) routes() {
	s.HandleFunc("/styles/{style}", s.stylesHandler()).Methods("GET")
	s.HandleFunc("/scripts/{script}", s.scriptsHandler()).Methods("GET")
	s.HandleFunc("/images/{img}", s.imagesHandler()).Methods("GET")
	s.HandleFunc("/exit", s.exitHandler()).Methods("GET")

	s.HandleFunc("/delete", s.deleteHandler()).Methods("POST")

	s.HandleFunc("/{file}", s.deleteHandler()).Methods("DELETE")

	spa := spaHandler{tpl: template.Must(template.ParseGlob("./src/templates/*.gohtml")), server: s}
	s.PathPrefix("/").Handler(spa)
}

func newServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		folder: "/",
		path:   "/",
	}
	s.routes()
	return s
}

func main() {
	srv := newServer()
	log.Fatal(http.ListenAndServe("localhost:8080", srv))
}
