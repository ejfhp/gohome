package own

import (
	"fmt"
)

const Lightning Who = "1"

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

//NewCommand build a new Command to send to the home plant
func NewCommand(who Who, what What, where Where) Command {
	cmd := fmt.Sprintf("*%s*%s*%s##", who, what, where)
	return Command(cmd)
}
