package main

import (
  "os"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var (
	APP_PORT       = "3000"
	API_RATE_LIMIT = 5.00
)

func getEnvOrDefault(key, defaultValue string) string {
  val := os.Getenv(key)
  if val == "" {
    val = defaultValue
  }
  return val
}

func PrayerNames() []string {
	return []string{"Fajr", "Dhuhr", "Asr", "Maghrib", "Isha"}
}

func main() {
	e := echo.New()

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(API_RATE_LIMIT))))

	// CORS restricted with a custom function to allow origins
	// and with the GET, PUT, POST or DELETE methods allowed.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(_ string) (bool, error) {
			// for now, always return true
			return true, nil
		},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", func(c echo.Context) error {
		ipAddr := getEnvOrDefault("DEBUG_IP", c.RealIP())
		responseType := c.QueryParam("response-type")
    isResponseJson := responseType == "json"

    // debug log
    //log.Printf("acceptType: %v, contentType: %v", acceptType, contentType)
    //for name, values := range c.Request().Header {
      //for _, value := range values {
        //log.Printf("%v: %v", name, value)
      //}
    //}

    plaintextUserAgentKeywords := []string{"curl", "httpie"}
		userAgent := c.Request().Header.Get("User-Agent")
    isResponsePlaintext := false
    for _, keyword := range plaintextUserAgentKeywords {
      if strings.Contains(userAgent, keyword) {
        isResponsePlaintext = true
        break
      }
    }

    if !isResponseJson && !isResponsePlaintext {
      return c.File("static/pages/index.html")
    }

		loc, err := getIpInfo(ipAddr)
		if err != nil {
			return err
		}

		locTz, _ := time.LoadLocation(loc.Timezone)
		curTime := time.Now().In(locTz)

		prayerTimes, err := getPrayerTimes(curTime.Format("02-01-2006"), loc)
		if err != nil {
			return err
		}

		nextPrayer, nextPrayerUntil, err := getPrayerTimeCountdown(curTime, loc, prayerTimes)
		if err != nil {
			return err
		}

    // avoid client-side caching: https://stackoverflow.com/a/9886945/2496217
    c.Response().Header().Set(echo.HeaderCacheControl, "max-age=0, no-cache, must-revalidate, proxy-revalidate")
		if isResponseJson {
			return responseJson(c, curTime, loc, prayerTimes, nextPrayer, nextPrayerUntil)
		}

		return responsePlaintext(c, curTime, loc, prayerTimes, nextPrayer, nextPrayerUntil)
	})

	// https://echo.labstack.com/docs/error-handling
	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	//   return func(c echo.Context) error {
	//     // https://sorcererxw.com/en/articles/go-echo-error-handing
	//     err := next(c)
	//     log.Printf("encounter error %v", err)
	//     return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	//   }
	// })

	e.Logger.Fatal(e.Start(":" + APP_PORT))
}

func responsePlaintext(c echo.Context, curTime time.Time, loc IpInfo, prayerTimes PrayerTimes, nextPrayer string, nextPrayerUntil time.Duration) error {
	respText := fmt.Sprintf(`KapanSolat
==========
detected location: %v, %v, %v
current local time: %v
next prayer: %v (%v remaining)
==========
prayer times for %v
`, loc.City, loc.Region, loc.Country, curTime.Format("15:04"), nextPrayer, nextPrayerUntil.Round(time.Minute).String(), curTime.Format("02-01-2006"))
	for _, i := range PrayerNames() {
		respText += i + ": " + prayerTimes[i] + "\n"
	}

	return c.String(http.StatusOK, respText)
}

func responseJson(c echo.Context, curTime time.Time, loc IpInfo, prayerTimes PrayerTimes, nextPrayer string, nextPrayerUntil time.Duration) error {
	return c.JSON(http.StatusOK, echo.Map{
		"current_location":    loc,
		"current_time":        curTime,
		"prayer_times":        prayerTimes,
		"next_prayer":         nextPrayer,
		"time_to_next_prayer": nextPrayerUntil.Round(time.Minute).String(),
	})

}

func jsonRequest(req *http.Request, res interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    log.Printf("encounter error: %v", err)
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("http response unsuccessful. got status code %v", resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//log.Printf(string(respBody))

	err = json.Unmarshal(respBody, &res)
	if err != nil {
		return err
	}

	return nil
}

type IpInfo struct {
	Addr     string  `json:"query"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Country  string  `json:"country"`
	Region   string  `json:"regionName"`
	City     string  `json:"city"`
	Timezone string  `json:"timezone"`
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

// type PrayerTimes struct {
// Location IpInfo
// Data struct {
// Timings struct {
// Fajr     string `json:"Fajr"`
// Dhuhr    string `json:"Dhuhr"`
// Asr      string `json:"Asr"`
// Maghrib  string `json:"Maghrib"`
// Isha     string `json:"Isha"`
// } `json:"timings"`
// } `json:"data"`
// }
//
//	type PrayerTimes struct {
//	  Location IpInfo
//	  Fajr     string `json:"Fajr"`
//	  Dhuhr    string `json:"Dhuhr"`
//	  Asr      string `json:"Asr"`
//	  Maghrib  string `json:"Maghrib"`
//	  Isha     string `json:"Isha"`
//	}
type PrayerTimes map[string]string

func getPrayerTimes(date string, location IpInfo) (PrayerTimes, error) {
	res := make(PrayerTimes)
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

	for _, i := range PrayerNames() {
		res[i] = timings[i].(string)
	}
	// res = PrayerTimes{
	//   Location: location,
	//   Fajr: timings["Fajr"].(string),
	//   Dhuhr: timings["Dhuhr"].(string),
	//   Asr: timings["Asr"].(string),
	//   Maghrib: timings["Maghrib"].(string),
	//   Isha: timings["Isha"].(string),
	// }

	return res, nil
}

func parseTime(datetime, loc string) (time.Time, error) {
	var emptyTime time.Time
	const format = "02-01-2006 15:04"
	tz, err := time.LoadLocation(loc)
	if err != nil {
		return emptyTime, err
	}
	return time.ParseInLocation(format, datetime, tz)
}

func getPrayerTimeCountdown(curTime time.Time, loc IpInfo, prayerTimes PrayerTimes) (string, time.Duration, error) {
	date := curTime.Format("02-01-2006 ")
	nextPrayer := ""
	var duration time.Duration

	for _, i := range PrayerNames() {
		parsedTime, err := parseTime(date+prayerTimes[i], loc.Timezone)
		if err != nil {
			return nextPrayer, duration, err
		}

		if curTime.Before(parsedTime) {
			nextPrayer = i
			duration = parsedTime.Sub(curTime)
			break
		}
	}
	return nextPrayer, duration, nil
}
