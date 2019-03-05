package gohome

import (
	"os"
  "flag"
	"github.com/pkg/errors"
)

func main() {
	flag.
	//command line must be WHO WHAT WHERE

}

func getPlant(file string) (*Plant, error) {
	config, err := os.Open("gohome.conf")
	if err != nil {
		return nil, errors.Wrapf(err, "cannot open configuration file %s", file)
	}
	defer config.Close()
	plant := NewPlant(config)
	return plant, nil
}
