package main

import (
	"cs-experiment-1/part-3/last-work/geotag"
	"cs-experiment-1/part-3/last-work/tag"
	"fmt"
)

func main() {
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

	fmt.Println(geotagsPointerTable[geotags[0].ID])
}
