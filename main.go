package main

import (
  //"log"
	"net/http"
  "io/ioutil"
  "encoding/json"

	"github.com/labstack/echo/v4"
)

var APP_PORT = "3000"

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"ip_addr": c.RealIP(),
		})
	})
	e.Logger.Fatal(e.Start(":" + APP_PORT))
}

type IpInfo struct {
  Addr     string `json:"query"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Country  string `json:"country"`
	Region   string `json:"regionName"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

func jsonRequest(req *http.Request, res interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    return err
	}

  respBody, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
	}

  err = json.Unmarshal(respBody, &res)
  if err != nil {
    return err
	}

  return nil
}

func getIpInfo(ipAddr string) (IpInfo, error) {
	var res IpInfo
  url := "http://ip-api.com/json/" + ipAddr

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
    return res, err
	}

	// req.Header = http.Header{
	// 	"Host":          {"www.host.com"},
	// 	"Content-Type":  {"application/json"},
	// 	"Authorization": {"Bearer Token"},
	// }

	err = jsonRequest(req, &res)
	if err != nil {
    return res, err
	}

	return res, nil
}
