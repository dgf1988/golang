package hoetom

import (
	"io/ioutil"
	"net/http"
	"time"
)

var client http.Client

var (
	NetTimeout    int
	NetTrytime    int
	NetTrydelayed int
)

func init() {
	NetTimeout = Atoi(ConfigGetDefault("nettimeout", "30"))
	NetTrytime = Atoi(ConfigGetDefault("nettry", "3"))
	NetTrydelayed = Atoi(ConfigGetDefault("nettrydelayed", "3"))
	client.Timeout = time.Duration(NetTimeout) * time.Second
}

func Get(urlget string) (string, int) {
	var err error
	var resp *http.Response
	for i := 0; i < NetTrytime; i++ {
		resp, err = client.Get(urlget)
		if err != nil {
			time.Sleep(time.Duration(NetTrydelayed) * time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		return err.Error(), ErrCode
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err.Error(), ErrCode
	}
	return coding.ConvertString(string(body)), resp.StatusCode
}
