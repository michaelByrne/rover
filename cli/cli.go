package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"rover/getrover"
)

type cliApp struct {
	numDays     int
	maxPerDay   int
	roverName   string
	camera      string
	roverGetter RoverGetter
}

type RoverGetter interface {
	GetRover(name string, numDays int) (*getrover.Rover, error)
}

func CLI(args []string, roverGetter RoverGetter) error {
	app := cliApp{
		roverGetter: roverGetter,
	}

	err := app.fromArgs(args)
	if err != nil {
		return err
	}

	err = app.run()
	if err != nil {
		return err
	}

	return nil
}

func (app *cliApp) fromArgs(args []string) error {
	fl := flag.NewFlagSet("rover", flag.ContinueOnError)
	fl.IntVar(
		&app.numDays, "days", 10, "number of days to fetch images for",
	)
	fl.IntVar(
		&app.maxPerDay, "dailyMax", 3, "max number per day",
	)
	fl.StringVar(
		&app.roverName, "name", "curiosity", "rover to fetch images for",
	)
	fl.StringVar(
		&app.camera, "camera", "", "camera to fetch images for",
	)

	if err := fl.Parse(args); err != nil {
		return err
	}

	return nil
}

func (app *cliApp) run() error {
	images, err := app.roverGetter.GetRover(app.roverName, app.numDays)
	if err != nil {
		return err
	}

	err = printJSON(images)
	if err != nil {
		return err
	}

	// To test cache
	//cachedImages, err := app.roverGetter.GetRover(app.roverName, app.numDays)
	//if err != nil {
	//	return err
	//}
	//
	//err = printJSON(cachedImages)
	//if err != nil {
	//	return err
	//}

	return nil
}

func printJSON(rover *getrover.Rover) error {
	j, err := json.MarshalIndent(rover, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(j))

	return nil
}
