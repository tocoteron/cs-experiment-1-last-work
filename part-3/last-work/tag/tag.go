package tag

import (
	"encoding/csv"
	"os"
	"strconv"
)

// Tag is correspondign the db table 'tag'
type Tag struct {
	ID  int
	Tag string
}

func UnmarshalTag(data []string) (Tag, error) {
	id, err := strconv.Atoi(data[0])
	if err != nil {
		return Tag{}, err
	}

	tag := Tag{
		ID:  id,
		Tag: data[1],
	}

	return tag, nil
}

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

func LoadTags(path string) ([]Tag, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	tags, err := UnmarshalTags(lines)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
