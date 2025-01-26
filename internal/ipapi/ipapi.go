package ipapi

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

var (
	cachedIP  *Response
	cachedAt  time.Time
	cachedMut sync.Mutex
)

type Response struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
	Query       string  `json:"query"`
}

func GetIP() (*Response, error) {
	cachedMut.Lock()
	defer cachedMut.Unlock()

	if cachedIP != nil && time.Since(cachedAt) < time.Minute {
		return cachedIP, nil
	}

	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := &Response{}
	if err = json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	cachedIP = data
	cachedAt = time.Now()
	return data, nil
}
