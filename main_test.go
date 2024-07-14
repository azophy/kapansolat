package main

import (
  "log"
	"testing"
)

const (
	ipKemenag = "103.7.13.247" // kemenag.go.id
)

// func TestGetIpInfoLocalhost(t *testing.T) {
// 	_, err := getIpInfo("127.0.0.1")
// 	if err == nil {
// 		t.Error("expecting error, got success")
// 	}
// }

func TestGetIpInfoKemenag(t *testing.T) {
	res, _ := getIpInfo(ipKemenag)
  log.Printf("result: %v\n", res)
	if res.Country != "Indonesia"  {
		t.Errorf("expecting country Indonesia, got %v", res.Country)
	}
}

func TestGetPrayerTime(t *testing.T) {
  // result from above getIp query
  location := IpInfo{
    Addr: "103.7.13.24",
    Country: "Indonesia",
    Region: "Jakarta",
    City: "Jakarta",
    Lat: -6.17189,
    Lon: 106.834,
    Timezone: "Asia/Jakarta",
  }

	res, _ := getPrayerTimes("14-07-2024", location)
  log.Printf("result: %v\n", res)
  if res.Isha != "19:06"  {
    t.Errorf("isha time result doesn't match")
  }
}

