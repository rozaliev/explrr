package main

import (
	"log"

	"runtime"
	"strconv"
	"time"
)

func main() {
	log.Println("Starting")
	runtime.GOMAXPROCS(2)

	leftEngine, err := FindDevice("/sys/class/tacho-motor", "port_name", "outC")
	if err != nil {
		log.Fatal("Left engine not found:", err)
	}

	rightEngine, err := FindDevice("/sys/class/tacho-motor", "port_name", "outD")
	if err != nil {
		log.Fatal("Right engine not found:", err)
	}

	irSensor, err := FindDevice("/sys/class/msensor", "port_name", "in3")
	if err != nil {
		log.Fatal("IR not found:", err)
	}

	touchSensor, err := FindDevice("/sys/class/msensor", "port_name", "in2")
	if err != nil {
		log.Fatal("TOUCH not found:", err)
	}

	log.Println("Inited")
MAIN:
	for {

		if wasTouched(touchSensor, leftEngine, rightEngine) {
			log.Println("Was touched")
			continue
		}

		log.Println("Was NOT touched")

		touched, err := GetValue(touchSensor, "value0")
		if err != nil {
			log.Fatal("Ir sensor error:", err)
		}

		if touched == "1" {
			Stop(leftEngine, rightEngine)
			Run(leftEngine, rightEngine, "-100", time.Millisecond*1000)
			TurnLeft(leftEngine, rightEngine, "100", time.Millisecond*500)
			continue MAIN
		}

		dist, err := GetValue(irSensor, "value0")
		if err != nil {
			log.Fatal("Ir sensor error:", err)
		}

		distInt, err := strconv.Atoi(dist)
		if err != nil {
			log.Fatal("Ir sensor value error:", err)
		}

		if distInt < 60 {
			Stop(leftEngine, rightEngine)
			TurnLeft(leftEngine, rightEngine, "100", time.Millisecond*500)
		} else {
			RunInf(leftEngine, rightEngine, "100")
		}

		time.Sleep(time.Millisecond * 20)
	}
}

func wasTouched(touchSensor, leftEngine, rightEngine string) bool {
	touched, err := GetValue(touchSensor, "value0")
	if err != nil {
		log.Fatal("Ir sensor error:", err)
	}

	if touched == "1" {
		Run(leftEngine, rightEngine, "-100", time.Millisecond*1000)
		TurnLeft(leftEngine, rightEngine, "100", time.Millisecond*500)
		return true
	}

	return false
}

func Run(leftEngine, rightEngine, speed string, timeout time.Duration) {
	log.Println("Run", speed, timeout)
	FatalOnErr(SetValue(leftEngine, "duty_cycle_sp", speed))
	FatalOnErr(SetValue(rightEngine, "duty_cycle_sp", speed))
	FatalOnErr(SetValue(leftEngine, "run", "1"))
	FatalOnErr(SetValue(rightEngine, "run", "1"))
	<-time.After(timeout)
	FatalOnErr(SetValue(leftEngine, "run", "0"))
	FatalOnErr(SetValue(rightEngine, "run", "0"))
}

func RunInf(leftEngine, rightEngine, speed string) {
	log.Println("RunInf", speed)
	FatalOnErr(SetValue(leftEngine, "duty_cycle_sp", speed))
	FatalOnErr(SetValue(rightEngine, "duty_cycle_sp", speed))
	FatalOnErr(SetValue(leftEngine, "run", "1"))
	FatalOnErr(SetValue(rightEngine, "run", "1"))
}

func Stop(leftEngine, rightEngine string) {
	log.Println("STOP")
	FatalOnErr(SetValue(leftEngine, "run", "0"))
	FatalOnErr(SetValue(rightEngine, "run", "0"))
}

func TurnLeft(leftEngine, rightEngine, speed string, timeout time.Duration) {
	log.Println("TurnLeft", speed, timeout)
	FatalOnErr(SetValue(leftEngine, "duty_cycle_sp", "-"+speed))
	FatalOnErr(SetValue(rightEngine, "duty_cycle_sp", speed))
	FatalOnErr(SetValue(leftEngine, "run", "1"))
	FatalOnErr(SetValue(rightEngine, "run", "1"))
	<-time.After(timeout)
	FatalOnErr(SetValue(leftEngine, "run", "0"))
	FatalOnErr(SetValue(rightEngine, "run", "0"))
}

func TurnRight(leftEngine, rightEngine, speed string, timeout time.Duration) {
	log.Println("TurnRight", speed, timeout)
	FatalOnErr(SetValue(leftEngine, "duty_cycle_sp", speed))
	FatalOnErr(SetValue(rightEngine, "duty_cycle_sp", "-"+speed))
	FatalOnErr(SetValue(leftEngine, "run", "1"))
	FatalOnErr(SetValue(rightEngine, "run", "1"))
	<-time.After(timeout)
	FatalOnErr(SetValue(leftEngine, "run", "0"))
	FatalOnErr(SetValue(rightEngine, "run", "0"))
}

func FatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
