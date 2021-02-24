package main

// npm start

import (
	"flag"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"

	"github.com/wigsnes/imageViewer/packages/fileo"
	"github.com/wigsnes/imageViewer/packages/foldero"
	"github.com/wigsnes/imageViewer/packages/imageo"
	"github.com/wigsnes/imageViewer/packages/stringo"
)

var (
	fILEtYPES  = []string{".jpg", ".jpeg", ".png", ".gif", ".webm", ".mp4", ".mov"}
	defaultCol = "3"
	pageFiles  = 90
)

// Row ..
type Row []Element

// Element ...
type Element struct {
	Path     string
	File     string
	FileName string
	Height   string
	IsImage  bool
	Columns  int
	NewCol   bool
	EndCol   bool
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
	files, err := ioutil.ReadDir(filePath)
	check(err)

	numFiles := fileo.NumberOfFiles(files)

	if numFiles == 0 {
		return []Element{}
	}

	if (numFiles / pageFiles) <= page {
		numFiles = numFiles - (pageFiles * (page - 1))
	}

	if numFiles > pageFiles {
		numFiles = pageFiles
	}

	if numFiles < 0 {
		return []Element{}
	}

	numInColums := int(math.Floor(float64(numFiles / columns)))

	row := make([]Element, numFiles)

	var wg sync.WaitGroup
	var mux sync.Mutex
	newCol := false
	endCol := false
	index := 0
	colNum := 0
	ii := 0
	for i, f := range files {
		if f.IsDir() {
			ii++
			continue
		}

		if !stringo.StringInSlice(f.Name(), fILEtYPES) {
			ii++
			continue
		}

		if (i - ii) <= numFiles*(page-1)-1 {
			continue
		}
		if (i - ii) >= (numFiles*(page-1))+numFiles {
			break
		}

		if numInColums == 0 {
			endCol = true
			newCol = true
		} else {
			if index%numInColums == numInColums-1 && colNum != columns {
				endCol = true
			}
			if index%numInColums == 0 && colNum != columns {
				newCol = true
				colNum++
			}
		}
		wg.Add(1)
		go handle(&row, index, newCol, endCol, columns, f.Name(), filePath, &mux, &wg)

		newCol = false
		endCol = false
		index++
	}

	wg.Wait()
	return row
}

func handle(row *([]Element), index int, newCol bool, endCol bool, columns int, fileName, filePath string, mux *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	if !stringo.StringInSlice(fileName, fILEtYPES) {
		return
	}

	if stringo.StringInSlice(fileName, []string{".jpg", ".jpeg", ".png"}) {

		fileo.ReplaceCharacterInFileName(fmt.Sprintf("%s%s", filePath, fileName), "+", "_")
		fileo.ReplaceCharacterInFileName(fmt.Sprintf("%s%s", filePath, fileName), " ", "_")

		// removed temporarily to find a better solution
		// if _, err := os.Stat(fmt.Sprintf("%sthumbnail/%s", filePath, fileName)); os.IsNotExist(err) {
		// 	imageo.CreateThumbnail(filePath, fileName, filePath+"thumbnail/", 360, 360)
		// }

		height := imageo.GetImageHeight(fileName, filePath)
		mux.Lock()
		(*row)[index] = Element{
			Path:     filePath,
			FileName: fileName,
			File:     filePath,
			Height:   height,
			Columns:  columns,
			IsImage:  true,
			NewCol:   newCol,
			EndCol:   endCol,
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
			Height:   "",
			IsImage:  false,
			NewCol:   newCol,
			EndCol:   endCol,
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
		fmt.Println(r)
		style, _ := mux.Vars(r)["style"]
		fmt.Println(style)
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
		fmt.Println(r)
		script, _ := mux.Vars(r)["script"]
		fmt.Println(script)
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
	if path[len(path)-1] != '/' {
		path += "/"
	}
	s.server.path = path

	filePath := fmt.Sprintf("%s%s", s.server.folder, path)
	fmt.Println("FilePath: ", filePath)

	values := r.URL.Query()

	// path := getQueryValue(values, "path", pATH)
	columns, err := strconv.Atoi(getQueryValue(values, "columns", defaultCol))
	check(err)
	page, err := strconv.Atoi(getQueryValue(values, "page", "1"))
	check(err)

	files, err := ioutil.ReadDir(filePath)
	check(err)

	totalPages := int(math.Ceil(float64(fileo.NumberOfFiles(files)) / float64(pageFiles)))
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = 1
	}
	prevPage := (page - 1)
	if prevPage <= 0 {
		prevPage = totalPages
	}

	data := data{
		Folders:    foldero.GetFolderInfo(filePath, path),
		Row:        getFilesInFolder(filePath, columns, page),
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
		file, _ := mux.Vars(r)["file"]
		fmt.Println("Delete file: ", s.folder+s.path+file)
		err := fileo.RemoveFile(s.folder + s.path + file)
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

	s.HandleFunc("/{file}", s.deleteHandler()).Methods("DELETE")

	spa := spaHandler{tpl: template.Must(template.ParseGlob("src/templates/*.gohtml")), server: s}
	s.PathPrefix("/").Handler(spa)
}

func newServer(path string) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		folder: path,
	}
	s.routes()
	return s
}

func main() {
	fmt.Println("Start!")
	path := flag.String("path", "C:\\", "path of the folder you want to open")
	flag.Parse()

	fmt.Println(*path)
	srv := newServer(*path)
	log.Fatal(http.ListenAndServe("localhost:8080", srv))
}
