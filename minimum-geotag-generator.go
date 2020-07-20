package main

import (
	"fmt"

	"github.com/tokoroten-lab/cs-experiment-1-last-work/geotag"
)

func main() {
	geotags, err := geotag.ReadGeoTagsFromCSV("data/geotag.csv", 10400000, 10000)
	if err != nil {
		panic(err)
	}

	idSearchTable := geotag.IDSearchTable{}

	for i := 0; i < len(geotags); i++ {
		idSearchTable[geotags[i].ID] = &geotags[i]
	}

	tagSearchTable, err := geotag.ReadTagSearchTableFromCSV("data/tag-search-table.csv", idSearchTable)
	if err != nil {
		panic(err)
	}

	minimumIDSearchTable := geotag.IDSearchTable{}

	for _, geotagPointers := range tagSearchTable {
		for _, geotagPointer := range geotagPointers {
			minimumIDSearchTable[geotagPointer.ID] = geotagPointer
		}
	}

	minimumGeotags := []geotag.GeoTag{}

	for _, geotagPointer := range minimumIDSearchTable {
		minimumGeotags = append(minimumGeotags, *geotagPointer)
	}

	fmt.Println(len(minimumGeotags))

	geotag.WriteGeoTagsToCSV("data/minimum-geotag.csv", minimumGeotags)
}
