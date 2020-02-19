package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	//go get github.com/gomodule/redigo/redis
	redigo "github.com/gomodule/redigo/redis"
)

//--------------------------------------------------------------------------------
//util
func GetIp() string {
	var ipString string

	ifaces, err := net.Interfaces()
	if err != nil {
		for _, i := range ifaces {
			addrs, err := i.Addrs()

			if err != nil {
				for _, addr := range addrs {
					var ip net.IP
					switch v := addr.(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}

					ipString = ip.String()
					// process IP address
				}
			}

		}
	}

	return ipString
}

func GetPublicIp() string {
	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	fmt.Printf("Getting IP address from  ipify ...\n")

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// /fmt.Printf("%s\n", ip)

	return string(ip)
}

//--------------------------------------------------------------------------------
// RankingManager ...
type RankingManager struct {
	sync.Mutex
	redisConnection redigo.Conn
}

var rankingManagerInstance *RankingManager
var rankingManagerOnce sync.Once

// GetRankingManager is singleton
func GetRankingManager() *RankingManager {
	rankingManagerOnce.Do(func() {
		rankingManagerInstance = &RankingManager{}
		rankingManagerInstance.Init()
	})
	return rankingManagerInstance
}

// Init ...
func (manager *RankingManager) Init() {
	ip := GetPublicIp()
	fmt.Print(ip)
	ip = ip + ":6379"

	passwordOption := redigo.DialPassword("")
	redisConnection, err := redigo.Dial("tcp", ip, passwordOption)
	//redisConnection, err := redigo.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	manager.redisConnection = redisConnection
}
func (manager *RankingManager) Destory() {
	manager.redisConnection.Close()
}
func (manager *RankingManager) SetScroe(uid string, score int32) {
	manager.redisConnection.Do("zadd", "userRanking", score, uid)
}
func (manager *RankingManager) GetHighscore(start int32, end int32, isDsc bool, useWithScores bool) ([]interface{}, error) {
	orderMethod := "ZRANGE"
	if isDsc {
		orderMethod = "ZREVRANGE"
	}
	return redigo.Values(manager.redisConnection.Do(orderMethod, "userRanking", start, end))
}
