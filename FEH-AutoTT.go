package main

import (
	"flag"
	"fmt"
	adb "github.com/yosemite-open/go-adb"
	"os"
	"time"
)

var (
	port = flag.Int("p", adb.AdbPort, "")

	cli         *adb.Adb
	runCost    = 15
	newPoints  = 441
	battleTimer = 90

	thresholdMax = 50000
	threshold3 = 10000
	threshold2 = 6000
	threshold1 = 3000

	nextRewardMax = 1250
	nextReward3 = 1000
	nextReward2 = 500
	nextReward1 = 250

	sealRewards = []int{6000, 10000, 20000}
	characterRewards = []int{1000, 15000}

	dif5 = "input tap 345 889"
	dif4 = "input tap 370 1070"
	confirmDifficulty = "input tap 355 870"
	openMenu = "input tap 85 1200"
	clickAutoBattle = "input tap 399 665"
	confirmAutoBattle = "input tap 377 630"
	quitScreen = "input tap 363 1117"
	redeemNormalRewardsButton = "input tap 377 818"
	redeemSealRewardsButton = "input tap 383 767"
	redeemCharacterRewardsButton = "input tap 377 799"
	restoreStamina = "input tap 376 694"
	closeRestoredWindow = "input tap 369 740"


)

func main() {
	flag.Parse()
	var err error

	cli, err = adb.NewWithConfig(adb.ServerConfig{
		Port: *port,
	})
	cli.StartServer()

	serials, err := cli.ListDeviceSerials()
	fmt.Println(serials)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	device := cli.Device(adb.DeviceWithSerial(serials[0]))

	fmt.Print("Enter current stamina- ")
	var stamina int
	fmt.Scanln(&stamina)
	fmt.Print("Enter current score- ")
	var score int
	fmt.Scanln(&score)
	fmt.Print("Enter the point value of the next reward to unlock- ")
	var nextReward int
	fmt.Scanln(&nextReward)
	fmt.Print("Enter the difficulty, 4 or 5- ")
	var diff int
	fmt.Scanln(&diff)
	if diff == 4 {
		newPoints = 352
		runCost = 12
	}
	fmt.Print("Enter the number of seconds until the next stamina point (before pressing enter be on the difficulty select screen)- ")
	var nextStaminaPoint float64
	fmt.Scanln(&nextStaminaPoint)
	runAutoBattle(device, stamina, score, nextReward, diff, nextStaminaPoint)

}

func runAutoBattle(device *adb.Device, stamina int, currentScore int, nextReward int, diff int, nextStaminaPoint float64) {
	continueBattle := true
	lastCheckTime := time.Now()
	seconds := 300.0 - nextStaminaPoint
	firstIteration := true
	for continueBattle {
		if currentScore >= thresholdMax {
			continueBattle = false
		}
		if !firstIteration{
			seconds += time.Since(lastCheckTime).Seconds()
			lastCheckTime = time.Now()
			if seconds >= 300 {
				stamina += 1
				seconds -= 300
			}
		}
		firstIteration = false
		sc := fmt.Sprintf("Current Score : %d", currentScore)
		fmt.Println(sc)
		st := fmt.Sprintf("Current Stamina : %d", stamina)
		fmt.Println(st)
		if stamina < runCost {
			refill(device, diff)
			stamina += 99
			wait(5)
			seconds = 0.0
		}
		fight(device, stamina, diff)
		fmt.Println()
		currentScore += newPoints

		if currentScore >= nextReward {
			oldReward := nextReward
			if currentScore >= thresholdMax {
				fmt.Println("Tempest Trials Completed")
				continueBattle = false
			} else if nextReward >= threshold3 {
				nextReward += nextRewardMax
			} else if nextReward >= threshold2{
				nextReward += nextReward3
			}else if nextReward >= threshold1{
				nextReward += nextReward2
			} else {
				nextReward += nextReward1
			}
			rew := fmt.Sprintf("Updating Reward, Next Is Available @ %d Points", nextReward)
			fmt.Println(rew)
			redeem(device, oldReward)
			wait(4)
		}

		stamina -= runCost
	}
}

func fight(device *adb.Device, stamina int, difficulty int) {
	//Difficulty selection
	if difficulty == 4 {
		device.RunCommand(dif4)
	} else {
		device.RunCommand(dif5)
	}
	wait(4)
	device.RunCommand(confirmDifficulty)
	wait(10)
	st := fmt.Sprintf("Fight Started. Current Stamina : %d", stamina-runCost)
	fmt.Println(st)
	device.RunCommand(openMenu)
	wait(4)
	device.RunCommand(clickAutoBattle)
	wait(4)
	device.RunCommand(confirmAutoBattle)

	// Wait for battling to end, adjust as needed for your characters
	wait(battleTimer)
	if difficulty == 5 {
		additionalTime := battleTimer/2
		wait(additionalTime)
	}

	device.RunCommand(quitScreen)
	wait(12)
	device.RunCommand(quitScreen)
	wait(6)
}

func refill(device *adb.Device, diff int) {
	if diff == 4 {
		device.RunCommand(dif4)
	} else {
		device.RunCommand(dif5)
	}
	wait(4)
	device.RunCommand(confirmDifficulty)
	wait(4)
	device.RunCommand(restoreStamina)
	wait(8)
	device.RunCommand(closeRestoredWindow)
}

func redeem(device *adb.Device, oldReward int) {
	for _, element := range sealRewards {
		if element == oldReward {
			device.RunCommand(redeemSealRewardsButton)
			wait(4)
			return
		}
	}
	for _, element := range characterRewards {
		if element == oldReward {
			device.RunCommand(redeemCharacterRewardsButton)
			wait(4)
			return
		}
	}
	device.RunCommand(redeemNormalRewardsButton)
	wait(4)
}

func wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}