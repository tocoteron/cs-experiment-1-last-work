package main

import (
	"sort"

	//	"github.com/tokoroten-lab/cs-experiment-1-last-work/geotag"
	"github.com/tokoroten-lab/cs-experiment-1-last-work/geotag"
	"github.com/tokoroten-lab/cs-experiment-1-last-work/tag"
)

func main() {
	geotags, err := geotag.ReadGeoTagsFromCSV("data/geotag.csv", 10400000, 10000)
	if err != nil {
		panic(err)
	}

	tags, err := tag.ReadTagsFromCSV("data/tag.csv", 22820000, 10000)
	if err != nil {
		panic(err)
	}

	idSearchTable := geotag.IDSearchTable{}

	for i := 0; i < len(geotags); i++ {
		idSearchTable[geotags[i].ID] = &geotags[i]
	}

	tagSearchTable := make(geotag.TagSearchTable, len(tags))

	for _, tag := range tags {
		tagSearchTable[tag.Tag] = append(tagSearchTable[tag.Tag], idSearchTable[tag.ID])
	}

	for tag, geotags := range tagSearchTable {
		sort.Slice(geotags, func(i, j int) bool {
			return geotags[i].Time > geotags[j].Time
		})

		if len(geotags) > 100 {
			tagSearchTable[tag] = geotags[:100]
		}
	}

	geotag.WriteTagSearchTableToCSV("data/tag-search-table.csv", tagSearchTable)
}
