package analytics

import "embed"

//go:embed queries/*.sql
var queriesFS embed.FS

func mustReadQuery(name string) string {
	data, err := queriesFS.ReadFile("queries/" + name)
	if err != nil {
		panic("failed to read embedded query " + name + ": " + err.Error())
	}
	return string(data)
}
