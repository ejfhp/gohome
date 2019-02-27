package gohome

import (
	"fmt"
)

type What string
type Who string
type Where string
type Dimension string
type Value string

const (
	Scenario    Who = "0"
	Lightning   Who = "1"
	Automation  Who = "2"
	Load        Who = "3"
	Temperature Who = "4"
	Alarm       Who = "5"
	CEN         Who = "15"
	Custom      Who = "9"
)

const (
	TURN_OFF         What = "0"
	TURN_ON          What = "1"
	SET_20           What = "2"
	SET_30           What = "3"
	SET_40           What = "4"
	SET_50           What = "5"
	SET_60           What = "6"
	SET_70           What = "7"
	SET_80           What = "8"
	SET_90           What = "9"
	SET_100          What = "10"
	ON_1_MIN         What = "11"
	ON_2_MIN         What = "12"
	ON_3_MIN         What = "13"
	ON_4_MIN         What = "14"
	ON_5_MIN         What = "15"
	ON_15_MIN        What = "16"
	ON_30_SEC        What = "17"
	ON_0_5_SEC       What = "18"
	BLINK_ON_0_5_SEC What = "20"
	BLINK_ON_1_SEC   What = "21"
	BLINK_ON_1_5_SEC What = "22"
	BLINK_ON_2_SEC   What = "23"
	BLINK_ON_2_5_SEC What = "24"
	BLINK_ON_3_SEC   What = "25"
	BLINK_ON_3_5_SEC What = "26"
	BLINK_ON_4_SEC   What = "27"
	BLINK_ON_4_5_SEC What = "28"
	BLINK_ON_5_SEC   What = "29"
	UP_ONE_LEVEL     What = "30"
	DOWN_ONE_LEVEL   What = "31"
	JOLLY            What = "1000"
)

const (
	GENERAL Where = "0"
)

func (w What) Token() string {
	return fmt.Sprintf("*%s", w)
}
