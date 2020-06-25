package main

import (
	"cs-experiment-1/part-3/last-work/geotag"
	"cs-experiment-1/part-3/last-work/tag"
	"flag"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
)

// TemplateRenderer is a view templates renderer
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func startWebServer(port string, tagSearchTable map[string][]*geotag.GeoTag) {
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

func searchGeoTagsByTag(tagSearchTable map[string][]*geotag.GeoTag, tag string) []geotag.GeoTag {
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

	geotags, err := geotag.ReadGeoTagsFromCSV("samples/geotag.csv")
	if err != nil {
		panic(err)
	}

	tags, err := tag.ReadTagsFromCSV("samples/tag.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println("Length of geotags", len(geotags))
	fmt.Println("Length of tags", len(tags))

	geotagsPointerTable := map[int]*geotag.GeoTag{}

	for i := 0; i < len(geotags); i++ {
		geotagsPointerTable[geotags[i].ID] = &geotags[i]
	}

	tagSearchTable := map[string][]*geotag.GeoTag{}

	for i := 0; i < len(tags); i++ {
		tagSearchTable[tags[i].Tag] = append(tagSearchTable[tags[i].Tag], geotagsPointerTable[tags[i].ID])
	}

	startWebServer(*port, tagSearchTable)
}
