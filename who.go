package gohome

import (
	"strings"
)

type Who struct {
	Code    string
	Desc    string
	Actions map[string]string
}

var actions_1 = map[string]string{
	"0":    "TURN_OFF",
	"1":    "TURN_ON",
	"2":    "SET_20",
	"3":    "SET_30",
	"4":    "SET_40",
	"5":    "SET_50",
	"6":    "SET_60",
	"7":    "SET_70",
	"8":    "SET_80",
	"9":    "SET_90",
	"10":   "SET_100",
	"11":   "ON_1_MIN",
	"12":   "ON_2_MIN",
	"13":   "ON_3_MIN",
	"14":   "ON_4_MIN",
	"15":   "ON_5_MIN",
	"16":   "ON_15_MIN",
	"17":   "ON_30_SEC",
	"18":   "ON_0_5_SEC",
	"20":   "BLINK_ON_0_5_SEC",
	"21":   "BINK_ON_1_SEC",
	"22":   "BLINK_ON_1_5_SEC",
	"23":   "BLINK_ON_2_SEC",
	"24":   "BLINK_ON_2_5_SEC",
	"25":   "BLINK_ON_3_SEC",
	"26":   "BLINK_OM_3_5_SEC",
	"27":   "BLINK_ON_4_SEC",
	"28":   "BLINK_ON_4_5_SEC",
	"29":   "BLINK_ON_5_SEC",
	"30":   "UP_ONE_LEVEL",
	"31":   "DOWN_ONE_LEVEL",
	"1000": "JOLLY",
}

var none = &Who{Code: "", Desc: "", Actions: map[string]string{}}

var allWho = map[string]*Who{
	"1": &Who{Code: "1", Desc: "LIGHT", Actions: actions_1},
}

func NewWho(who string) *Who {
	w, ok := allWho[who]
	if !ok {
		return none
	}
	return w
}

func (w Who) WhatFromDesc(text string) (What, error) {
	for k, v := range w.Actions {
		if v == strings.ToUpper(text) {
			return What{Code: k, Desc: v}, nil
		}
	}
	return What{}, ErrWhatNotFound
}

func (w Who) WhatFromCode(code string) (What, error) {
	desc, ok := w.Actions[code]
	if !ok {
		return What{}, ErrWhatNotFound
	}
	return What{Code: code, Desc: desc}, nil
}
