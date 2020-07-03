package tag

import (
	"os"
	"strconv"

	"github.com/tokoroten-lab/cs-experiment-1/part-3/last-work/csvutil"
)

// Tag is correspondign the db table 'tag'
type Tag struct {
	ID  uint64
	Tag string
}

// UnmarshalTag converts []string to Tag
func UnmarshalTag(data []string) (Tag, error) {
	id, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return Tag{}, err
	}

	tag := Tag{
		ID:  id,
		Tag: data[1],
	}

	return tag, nil
}

// UnmarshalTags converts [][]string to []Tag by using UnmarshalTag function
func UnmarshalTags(data [][]string) ([]Tag, error) {
	var tags []Tag

	for i := 0; i < len(data); i++ {
		tag, err := UnmarshalTag(data[i])
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// ReadTagsFromCSV marshal []Tag from CSV file corresponding the path
func ReadTagsFromCSV(path string, capacity int, buffsize int) ([]Tag, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	tags := make([]Tag, 0, capacity)

	for record := range csvutil.AsyncReadCSV(reader, buffsize) {
		tag, err := UnmarshalTag(record)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
