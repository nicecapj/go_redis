package main

import (
	"fmt"
)

func main() {

	rankingManager := GetRankingManager()
	defer rankingManager.Destory()

	rankingManager.SetScroe("nic", 300)
	rankingManager.SetScroe("jason", 200)
	rankingManager.SetScroe("lee", 100)

	highScore, err := rankingManager.GetHighscore(0, 10, true, false)
	if err != nil {
		fmt.Print(highScore)
	}
}
