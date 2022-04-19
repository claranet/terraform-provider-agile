package tools

import agile "github.com/claranet/agilec-go-client/client"

func ExtractSliceOfStrings(itemsRaw []interface{}) []*string {
	items := make([]*string, len(itemsRaw))
	for i, raw := range itemsRaw {
		items[i] = agile.String(raw.(string))
	}
	return items
}

func CreateSliceOfStrings(itemsRaw []*string) []string {
	items := make([]string, len(itemsRaw))
	for i, raw := range itemsRaw {
		items[i] = *raw
	}
	return items
}
