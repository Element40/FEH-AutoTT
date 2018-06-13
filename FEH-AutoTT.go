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
	run_cost    = 15
	new_points  = 441

	dif5 = "input tap 345 889"
	dif4 = "input tap 370 1070"

	threshold_max = 50000
	threshold_3 = 10000
	threshold_2 = 6000
	threshold_1 = 3000

	next_reward_max = 1250
	next_reward_3 = 1000
	next_reward_2 = 500
	next_reward_1 = 250

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
	var next_reward int
	fmt.Scanln(&next_reward)
	fmt.Print("Enter the difficulty, 4 or 5- ")
	var diff int
	fmt.Scanln(&diff)
	if diff == 4 {
		new_points = 352
		run_cost = 12
	}
	fmt.Print("Enter the number of seconds until the next stamina point (before pressing enter be on the difficulty select screen)- ")
	var next_stamina_point float64
	fmt.Scanln(&next_stamina_point)
	runAutoBattle(device, stamina, score, next_reward, diff, next_stamina_point)

}

func runAutoBattle(device *adb.Device, stamina int, current_score int, next_reward int, diff int, next_stamina_point float64) {
	continue_battle := true
	last_check_time := time.Now()
	seconds := 300.0 - next_stamina_point
	first_iteration := true
	for continue_battle {
		if current_score >= threshold_max {
			continue_battle = false
		}
		if !first_iteration{
			seconds += time.Since(last_check_time).Seconds()
			last_check_time = time.Now()
			if seconds >= 300 {
				stamina += 1
				seconds -= 300
			}
		}
		first_iteration = false
		sc := fmt.Sprintf("Current Score : %d", current_score)
		fmt.Println(sc)
		st := fmt.Sprintf("Current Stamina : %d", stamina)
		fmt.Println(st)
		if stamina < run_cost {
			refill(device, diff)
			stamina += 99
			time.Sleep(5 * time.Second)
			seconds = 0.0
		}
		fight(device, stamina, diff)
		fmt.Println()
		current_score += new_points

		if current_score >= next_reward {
			if current_score >= threshold_max {
				fmt.Println("Tempest Trials Completed")
				continue_battle = false
			} else if next_reward >= threshold_3 {
				next_reward += next_reward_max
			} else if next_reward >= threshold_2{
				next_reward += next_reward_3
			}else if next_reward >= threshold_1{
				next_reward += next_reward_2
			} else {
				next_reward += next_reward_1
			}
			rew := fmt.Sprintf("Updating Reward, Next Is Available @ %d Points", next_reward)
			fmt.Println(rew)
			redeem(device)
			time.Sleep(4 * time.Second)
		}

		stamina -= run_cost
	}
}

func fight(device *adb.Device, stamina int, diff int) {
	//Difficulty selection
	if diff == 4 {
		device.RunCommand(dif4) // This was for Hard 4 on my device
	} else {
		device.RunCommand(dif5) // This was for Hard 5 on my device
	}
	time.Sleep(4 * time.Second)                 //This tells the phone to wait an amount of milliseconds before proceeding
	device.RunCommand("input tap 355 870") // Confirming difficult choice
	time.Sleep(10 * time.Second)                //Longer pause, I used multiplier for each steps to fine tunes things
	st := fmt.Sprintf("Fight Started. Current Stamina : %d", stamina-run_cost)
	fmt.Println(st)
	device.RunCommand("input tap 85 1200") // Menu Open
	time.Sleep(4 * time.Second)
	device.RunCommand("input tap 399 665") //autobattle start
	time.Sleep(4 * time.Second)
	device.RunCommand("input tap 377 630") //Confirm autobattle start

	time.Sleep(1 * time.Minute) //This is the time it waited for all the battle to ends
	time.Sleep(30 * time.Second)
	if diff == 5 {
		time.Sleep(45* time.Second) //Wait longer for
	}

	device.RunCommand("input tap 363 1117")  // Click to quit from the last battle
	time.Sleep(12 * time.Second)
	device.RunCommand("input tap 363 1117") // Quitting from score screen
	time.Sleep(6 * time.Second)
}

func refill(device *adb.Device, diff int) {
	if diff == 4 {
		device.RunCommand(dif4) // This was for Hard 5 on my device
	} else {
		device.RunCommand(dif5) // This was for Hard 5 on my device
	}
	time.Sleep(4 * time.Second)                 //This tells the phone to wait an amount of milliseconds before proceeding
	device.RunCommand("input tap 355 870") // Confirming difficult choice
	time.Sleep(4 * time.Second)                 //This tells the phone to wait an amount of milliseconds before proceeding
	device.RunCommand("input tap 376 694") // Click Restore
	time.Sleep(8 * time.Second)                 //This tells the phone to wait an amount of milliseconds before proceeding
	device.RunCommand("input tap 369 740") //Close Window
}

func redeem(device *adb.Device) {
	device.RunCommand("input tap 377 818") // Redeem
	time.Sleep(4 * time.Second)                 //This tells the phone to wait an amount of milliseconds before proceeding
}
