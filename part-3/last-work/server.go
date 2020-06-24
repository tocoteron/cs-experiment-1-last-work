package main

import (
	"cs-experiment-1/part-3/last-work/geotag"
	"cs-experiment-1/part-3/last-work/tag"
	"fmt"
)

func main() {
	geotags, err := geotag.ReadGeoTagsFromCSV("samples/huge_geotag.csv")
	if err != nil {
		panic(err)
	}

	tags, err := tag.ReadTagsFromCSV("samples/huge_tag.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println("Length of geotags", len(geotags))
	fmt.Println("Length of tags", len(tags))
}
