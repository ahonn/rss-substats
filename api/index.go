package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type StatData struct {
	TotalSubs        int               `json:"totalSubs"`
	SubsInEachSource map[string]int    `json:"subsInEachSource"`
	FailedSources    map[string]string `json:"failedSources"`
}

type StatResponse struct {
	Status int       `json:"status"`
	Data   *StatData `json:"data"`
}

func getSubStatusURL(src string, key string) (*url.URL, error) {
	statsUrl, err := url.Parse("https://api.spencerwoo.com/substats")
	if err != nil {
		return nil, err
	}
	query := statsUrl.Query()
	query.Set("source", src)
	query.Set("queryKey", key)
	statsUrl.RawQuery = query.Encode()
	return statsUrl, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	source := strings.Split(r.URL.Query().Get("source"), ",")
	queryKey := strings.Split(r.URL.Query().Get("queryKey"), ",")

	data := StatData{0, map[string]int{}, map[string]string{}}

	size := len(source) * len(queryKey)
	ch := make(chan StatResponse, size)
	defer close(ch)

	for i, s := range source {
		for j, k := range queryKey {
			go func(src string, key string, i int, j int) {
				record := src + "#" + key

				statsUrl, err := getSubStatusURL(src, key)
				if err != nil {
					data.FailedSources[record] = err.Error()
				}

				res, err := http.Get(statsUrl.String())
				if err != nil {
					data.FailedSources[record] = err.Error()
				}
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					data.FailedSources[record] = err.Error()
				}

				var result StatResponse
				if err := json.Unmarshal(body, &result); err != nil {
					data.FailedSources[record] = err.Error()
				}

				ch <- result
			}(s, k, i, j)
		}
	}

	for i := 0; i < size; i += 1 {
		result := <-ch
		data.TotalSubs += int(result.Data.TotalSubs)
    for k, v := range result.Data.SubsInEachSource {
      data.SubsInEachSource[k] += v;
    }
	}

	response := StatResponse{200, &data}
	fmt.Fprint(w, response.Status)
	res, _ := json.Marshal(response)
	fmt.Fprint(w, string(res))
}
