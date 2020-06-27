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
	"sort"
	"text/template"

	"github.com/labstack/echo"
)

// GeoTagsPointerTable is mapping ID to GeoTag
type GeoTagsPointerTable map[int]*geotag.GeoTag

// TagSearchTable is mapping Tag to Geotag
type TagSearchTable map[string][]*geotag.GeoTag

// TemplateRenderer is a view templates renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func startWebServer(port string, tagSearchTable TagSearchTable) {
	e := echo.New()
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

func searchGeoTagsByTag(tagSearchTable TagSearchTable, tag string) []geotag.GeoTag {
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
		geotags, err = geotag.ReadGeoTagsFromCSV("samples/geotag.csv", 10500000, 1000)
		return err
	})

	fmt.Printf("%d, %f\n", len(geotags), t)

	t, err = measure.MeasureFuncTime(func() error {
		tags, err = tag.ReadTagsFromCSV("samples/tag.csv")
		return err
	})

	fmt.Printf("%d, %f\n", len(tags), t)

	geotagsPointerTable := GeoTagsPointerTable{}

	for i := 0; i < len(geotags); i++ {
		geotagsPointerTable[geotags[i].ID] = &geotags[i]
	}

	tagSearchTable := TagSearchTable{}

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

	startWebServer(*port, tagSearchTable)
}
