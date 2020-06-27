package main

import (
	"cs-experiment-1/part-3/last-work/geotag"
	"cs-experiment-1/part-3/last-work/measure"
	"cs-experiment-1/part-3/last-work/tag"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"runtime"
	"sort"
	"text/template"
	"unsafe"

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

// GeoTagsPointerTable is mapping ID to GeoTag
type GeoTagsPointerTable map[uint64]*geotag.GeoTag

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

	e.HTTPErrorHandler = customHTTPErrorHandler

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("view/*.html")),
	}
	e.Renderer = renderer

	e.GET("/search", func(c echo.Context) error {
		tag := c.QueryParam("tag")
		geotags := searchGeoTagsByTag(tagSearchTable, tag)

		return c.Render(http.StatusOK, "search.html", map[string]interface{}{
			"geotags": geotags,
		})
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

func searchGeoTagsByTag(tagSearchTable geotag.TagSearchTable, tag string) []geotag.GeoTag {
	geotagPointers := tagSearchTable[tag]
	geotags := make([]geotag.GeoTag, len(geotagPointers))

	for i := 0; i < len(geotags); i++ {
		geotags[i] = *geotagPointers[i]
	}

	return geotags
}

func main() {
	port := flag.String("port", "1323", "Port number of web server")
	flag.Parse()

	var err error
	var t float64
	var geotags []geotag.GeoTag
	var tags []tag.Tag

	t, err = measure.MeasureFuncTime(func() error {
		geotags, err = geotag.ReadCompressedGeoTagsFromCSV("samples/compressed-geotag.csv", 10500000, 1000)
		return err
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d geotags loaded, %f(sec)\n", len(geotags), t)
	fmt.Println("bytes of geotags[0] is", unsafe.Sizeof(geotags[0]))

	printMemory()

	t, err = measure.MeasureFuncTime(func() error {
		tags, err = tag.ReadTagsFromCSV("samples/tag.csv", 23000000, 1000)
		return err
	})

	fmt.Printf("%d tags loaded, %f(sec)\n", len(tags), t)
	fmt.Println("bytes of tags[0] is", unsafe.Sizeof(tags[0]))

	printMemory()

	geotagsPointerTable := GeoTagsPointerTable{}

	for i := 0; i < len(geotags); i++ {
		geotagsPointerTable[geotags[i].ID] = &geotags[i]
	}

	tagSearchTable := geotag.TagSearchTable{}

	for i := 0; i < len(tags); i++ {
		tagSearchTable[tags[i].Tag] = append(tagSearchTable[tags[i].Tag], geotagsPointerTable[tags[i].ID])
	}

	for i := 0; i < len(tags); i++ {
		targetGeotags := tagSearchTable[tags[i].Tag]
		sort.Slice(targetGeotags, func(a, b int) bool {
			return targetGeotags[a].Time > targetGeotags[b].Time
		})

		last := int(math.Min(float64(len(targetGeotags)), 100))
		tagSearchTable[tags[i].Tag] = targetGeotags[:last]
	}

	err = geotag.WriteTagSearchTableToCSV("samples/tag-search-table.csv", tagSearchTable)
	if err != nil {
		panic(err)
	}

	tags = nil

	runtime.GC()

	printMemory()

	startWebServer(*port, tagSearchTable)
}
