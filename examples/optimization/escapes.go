package optimization

type Data struct {
	str string
}

var (
	sliceItems []Data
	mapItems   map[string]Data
)

var allStrings = []string{
	"PTURHXIgySk",
	"sHd5OAxx9QTZ!?-",
	"eFSOlI",
	"KB?AS-mkbNddHV_OrI",
	"9dzc29nTkkJ_HrF37l9LJ5eW4k84LklWjra",
	"oOAfa P-l -tA!GAdiCEojyDcM_bguX",
	"jttibOlm57LAv.87T13XN8?s.ASilDP!PQlCJq oH",
	"paQ!Hdmeh4UCfURP1IL5UXNoI3IVmtV11",
	"Wxd34 dBPWT3FFT9ete1I2T_LGTCHN,Wal7D9gu7IgX1UE",
	"iUdyZwS.dQ3f",
	"84WEH5GhoIFopQ",
}

func sliceFillDyn(n int) []Data {
	items := make([]Data, 0)
	for i := 0; i < n; i++ {
		items = append(items, Data{allStrings[i]})
	}

	return items
}

func sliceFillFix(n int) []Data {
	items := make([]Data, 0, n)
	for i := 0; i < n; i++ {
		items = append(items, Data{allStrings[i]})
	}

	return items
}

func sliceSeek(items []Data, seek string) int {
	for i, item := range items {
		if seek == item.str {
			return i
		}
	}

	return -1
}

func mapFillDyn(n int) map[string]Data {
	items := make(map[string]Data)
	for i := 0; i < n; i++ {
		items[allStrings[i]] = Data{allStrings[i]}
	}

	return items
}

func mapFillFix(n int) map[string]Data {
	items := make(map[string]Data, n)
	for i := 0; i < n; i++ {
		items[allStrings[i]] = Data{allStrings[i]}
	}

	return items
}

func mapSeek(items map[string]Data, seek string) (Data, bool) {
	data, has := items[seek]
	return data, has
}
