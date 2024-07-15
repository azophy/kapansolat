package main

import (
  //"log"
  "fmt"
  "time"
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

type IpInfo struct {
  Addr     string `json:"query"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Country  string `json:"country"`
	Region   string `json:"regionName"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
  TimeLoc  *time.Location // *Location
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

  // load timezone for time package (https://pkg.go.dev/time#LoadLocation)
  res.TimeLoc, err = time.LoadLocation(res.Timezone)

	return res, nil
}

//type PrayerTimes struct {
  //Location IpInfo
  //Data struct {
    //Timings struct {
      //Fajr     string `json:"Fajr"`
      //Dhuhr    string `json:"Dhuhr"`
      //Asr      string `json:"Asr"`
      //Maghrib  string `json:"Maghrib"`
      //Isha     string `json:"Isha"`
    //} `json:"timings"`
  //} `json:"data"`
//}
type PrayerTimes struct {
  Location IpInfo
  Fajr     string `json:"Fajr"`
  Dhuhr    string `json:"Dhuhr"`
  Asr      string `json:"Asr"`
  Maghrib  string `json:"Maghrib"`
  Isha     string `json:"Isha"`
}

func getPrayerTimes(date string, location IpInfo) (PrayerTimes, error) {
  res := PrayerTimes{Location: location}
  url := fmt.Sprintf("http://api.aladhan.com/v1/timings/%s?method=20&latitude=%v&longitude=%v", date, location.Lat, location.Lon)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
    return res, err
	}

  var resp interface{}
	err = jsonRequest(req, &resp)
	if err != nil {
    return res, err
	}

  timings := resp.(map[string]interface{})["data"].(map[string]interface{})["timings"].(map[string]interface{})

  res = PrayerTimes{
    Location: location,
    Fajr: timings["Fajr"].(string),
    Dhuhr: timings["Dhuhr"].(string),
    Asr: timings["Asr"].(string),
    Maghrib: timings["Maghrib"].(string),
    Isha: timings["Isha"].(string),
  }

	return res, nil
}

//func countDownNextPrayer(curTime time.Time, timing PrayerTimes) (string, string, bool) {
//}
