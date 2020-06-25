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

	tagSearchTable := map[string][]*geotag.GeoTag{}

	for i := 0; i < len(tags); i++ {
		tagSearchTable[tags[i].Tag] = append(tagSearchTable[tags[i].Tag], geotagsPointerTable[tags[i].ID])
	}

	var searchTag string

	for {
		fmt.Print("Search: ")
		fmt.Scan(&searchTag)

		if searchTag == "exit" || searchTag == "quit" {
			break
		}

		fmt.Println(tagSearchTable[searchTag])
	}
}
