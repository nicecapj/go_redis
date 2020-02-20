package main

import (
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
)

func main() {

	rankingManager := GetRankingManager()
	defer rankingManager.Destory()

	rankingManager.SetScroe("nic", 300)
	rankingManager.SetScroe("jason", 200)
	rankingManager.SetScroe("lee", 100)

	highScore, err := rankingManager.GetHighscore(0, 10, true, false)
	if err == nil {
		fmt.Println(highScore)
	}

	count := len(highScore)
	for x := 0; x < count; x++ {
		fmt.Println(redigo.String(highScore[x], err))
	}
}
