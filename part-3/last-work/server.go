package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"text/template"

	"github.com/tokoroten-lab/cs-experiment-1/part-3/last-work/geotag"

	"github.com/labstack/echo"
)

// ResponseCache is cache of http response data
type ResponseCache map[string][]byte

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

func renderWithCache(e *echo.Echo, c echo.Context, cache ResponseCache, tag string, tagSearchTable geotag.TagSearchTable) error {
	response, isCaching := cache[tag]
	if isCaching {
		return c.HTMLBlob(http.StatusOK, response)
	}

	geotags := tagSearchTable[tag]
	buf := new(bytes.Buffer)

	err := e.Renderer.Render(buf, "search.html", map[string]interface{}{
		"geotags": geotags,
	}, c)

	if err != nil {
		return err
	}

	cache[tag] = buf.Bytes()

	return c.HTMLBlob(http.StatusOK, buf.Bytes())
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

func startWebServer(isDebug bool, port uint, tagSearchTable geotag.TagSearchTable) {
	cache := ResponseCache{}

	e := echo.New()
	e.HideBanner = true

	if isDebug {
		e.HTTPErrorHandler = customHTTPErrorHandler
	}

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
	e.Renderer = renderer

	e.GET("/search", func(c echo.Context) error {
		tag := c.QueryParam("tag")
		return renderWithCache(e, c, cache, tag, tagSearchTable)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func main() {
	port := flag.Uint("port", 8080, "Port number of web server")
	isDebug := flag.Bool("debug", false, "Debug mode flag")
	flag.Parse()

	geotags, err := geotag.ReadCompressedGeoTagsFromCSV("data/minimum-geotag.csv", 10500000, 1000)
	if err != nil {
		panic(err)
	}

	runtime.GC()
	debug.FreeOSMemory()

	if *isDebug {
		printMemory()
	}

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

	if *isDebug {
		printMemory()
	}

	fmt.Println(len(tagSearchTable), "tags loaded")

	if *isDebug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	startWebServer(*isDebug, *port, tagSearchTable)
}
