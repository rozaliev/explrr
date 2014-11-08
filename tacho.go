package main

import "errors"

type TachoMotor struct {
	Path string
}

func FindMiniTacho() (*TachoMotor, error) {
	path, err := FindDevice("/sys/class/tacho-motor", "type", "minitacho")
	if err != nil {
		return nil, errors.New("Mini Tacho error:" + err.Error())
	}

	return &TachoMotor{Path: path}, nil

}
