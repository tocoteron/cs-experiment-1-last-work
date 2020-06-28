package main

import (
	"cs-experiment-1/part-3/last-work/geotag"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"text/template"

	"github.com/labstack/echo"
)

var mem runtime.MemStats

func toKB(bytes uint64) float64 {
	return float64(bytes) / 1024
}

func toMB(bytes uint64) float64 {
	return toKB(bytes) / 1024
}

func toGB(bytes uint64) float64 {
	return toMB(bytes) / 1024
}

func printMemory() {
	runtime.ReadMemStats(&mem)
	fmt.Println("-")
	fmt.Printf("Alloc      %f(GB)\n", toGB(mem.Alloc))
	fmt.Printf("HeapAlloc  %f(GB)\n", toGB(mem.HeapAlloc))
	fmt.Printf("TotalAlloc %f(GB)\n", toGB(mem.TotalAlloc))
	fmt.Printf("HeapSys    %f(GB)\n", toGB(mem.HeapSys))
	fmt.Printf("Sys        %f(GB)\n", toGB(mem.Sys))
	fmt.Println("-")
}

// TemplateRenderer is a view templates renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	c.Logger().Error(err)
}

func startWebServer(port string, tagSearchTable geotag.TagSearchTable) {
	e := echo.New()

	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
	e.Renderer = renderer

	e.GET("/search", func(c echo.Context) error {
		tag := c.QueryParam("tag")
		geotags := tagSearchTable[tag]

		return c.Render(http.StatusOK, "search.html", map[string]interface{}{
			"geotags": geotags,
		})
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func main() {
	port := flag.String("port", "1323", "Port number of web server")
	flag.Parse()

	geotags, err := geotag.ReadCompressedGeoTagsFromCSV("data/minimum-geotag.csv", 10500000, 1000)
	if err != nil {
		panic(err)
	}

	runtime.GC()
	debug.FreeOSMemory()

	printMemory()
	fmt.Println(len(geotags), "geotags loaded")

	idSearchTable := geotag.IDSearchTable{}

	for i := 0; i < len(geotags); i++ {
		idSearchTable[geotags[i].ID] = &geotags[i]
	}

	tagSearchTable, err := geotag.ReadTagSearchTableFromCSV("data/tag-search-table.csv", idSearchTable)
	if err != nil {
		panic(err)
	}

	runtime.GC()
	debug.FreeOSMemory()

	printMemory()
	fmt.Println(len(tagSearchTable), "tags loaded")

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	startWebServer(*port, tagSearchTable)
}
