package main

import (
	"time"

	"github.com/airenas/go-app/pkg/goapp"
	"github.com/airenas/music-transcription-service/internal/pkg/file"
	"github.com/airenas/music-transcription-service/internal/pkg/limiter"
	"github.com/airenas/music-transcription-service/internal/pkg/service"
	"github.com/airenas/music-transcription-service/internal/pkg/transcriber"
	"github.com/labstack/gommon/color"
	"github.com/pkg/errors"
)

func main() {
	goapp.StartWithDefault()

	data := service.Data{}
	data.Port = goapp.Config.GetInt("port")

	var err error
	goapp.Log.Infof("Temp dir: %s", goapp.Config.GetString("tempDir"))
	data.Saver, err = file.NewSaver(goapp.Config.GetString("tempDir"))
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init file saver"))
	}
	data.Worker, err = transcriber.NewWorker(goapp.Config.GetString("app.cmd"))
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init transcriber wrapper"))
	}
	data.Limiter, err = limiter.NewCount(10, time.Second*2)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't init rate limiter"))
	}
	printBanner()

	err = service.StartWebServer(&data)
	if err != nil {
		goapp.Log.Fatal(errors.Wrap(err, "can't start the service"))
	}
}

var (
	version string
)

func printBanner() {
	banner := `
     __  ___           _     
    /  |/  /_  _______(_)____
   / /|_/ / / / / ___/ / ___/
  / /  / / /_/ (__  ) / /__  
 /_/  /_/\__,_/____/_/\___/  
   __                                  _ __             
  / /__________ _____  _______________(_) /_  ___  _____
 / __/ ___/ __ ` + "`" + `/ __ \/ ___/ ___/ ___/ / __ \/ _ \/ ___/
/ /_/ /  / /_/ / / / (__  ) /__/ /  / / /_/ /  __/ /    
\__/_/   \__,_/_/ /_/____/\___/_/  /_/_.___/\___/_/ v: %s 

%s
________________________________________________________                                                 

`
	cl := color.New()
	cl.Printf(banner, cl.Red(version), cl.Green("https://github.com/airenas/music-transcription-service"))
}
