package main

import (
	"log"
	"strconv"
	"time"
)

func main() {
	log.Println("Starting")

	motor, err := FindMiniTacho()
	if err != nil {
		log.Fatal(err)
	}

	FatalOnErr(SetValue(motor.Path, "position", "0"))
	FatalOnErr(SetValue(motor.Path, "run_mode", "position"))
	FatalOnErr(SetValue(motor.Path, "stop_mode", "hold"))
	FatalOnErr(SetValue(motor.Path, "regulation_mode", "on"))
	FatalOnErr(SetValue(motor.Path, "ramp_up_sp", "30"))
	FatalOnErr(SetValue(motor.Path, "ramp_down_sp", "30"))
	FatalOnErr(SetValue(motor.Path, "pulses_per_second_sp", "1000"))

	irPath, err := FindDevice("/sys/class/msensor", "port_name", "in3")
	if err != nil {
		log.Fatal("IR not found:", err)
	}

	m := []int{}

	for i := 0; i <= 360; i += 10 {
		res, err := stepMeasure(motor.Path, i, irPath)
		if err != nil {
			log.Fatal(err)
		}

		m = append(m, res)
	}

	log.Println("Measurements:", m)

	FatalOnErr(SetValue(motor.Path, "position_sp", "0"))

	FatalOnErr(SetValue(motor.Path, "run", "1"))

}

func stepMeasure(mpath string, position int, irPath string) (int, error) {
	if err := SetValue(mpath, "position_sp", strconv.Itoa(position)); err != nil {
		return 0, err
	}

	if err := SetValue(mpath, "run", "1"); err != nil {
		return 0, err
	}

	for {
		val, err := GetValue(mpath, "run")
		if err != nil {
			return 0, err
		}
		if val == "0" {
			break
		}

		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)

	val, err := GetValue(irPath, "value0")
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(val)
}

func FatalOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
