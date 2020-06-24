package main

import (
	"fmt"

	"cs-experiment-1/part-3/last-work/geotag"
	"cs-experiment-1/part-3/last-work/tag"
)

func main() {
	geotags, err := geotag.LoadGeoTags("samples/geotag.csv")
	if err != nil {
		panic(err)
	}

	tags, err := tag.LoadTags("samples/tag.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println(geotags)
	fmt.Println(tags)
}
