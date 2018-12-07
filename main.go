package main

import (
	"fmt"

	"deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
	"gobot.io/x/gobot/platforms/keyboard"
	"image"
	"image/color"
	"gocv.io/x/gocv"
	"runtime"
)

func main() {
	runtime.LockOSThread()
	track := ddr.NewTrack()
	defer track.Close()

	inflight := false

	window := gocv.NewWindow("Drone")
	//drone := ddr.NewDrone(ddr.DroneReal, "../drone-camera-calibration-400.yaml")
	drone := ddr.NewDrone(ddr.DroneFake, "./camera-calibration.yaml")
	err := drone.Init()

	if err != nil {
		fmt.Printf("error while initializing drone: %v\n", err)
		return
	}

	for {
		frame := <-drone.VideoStream()

		// detect markers in this frame
		markers := track.GetMarkers(&frame)

		rings := track.ExtractRings(markers)
		for id, ring := range rings {
			pose := ring.EstimatePose(drone)
			gocv.PutText(&frame, fmt.Sprintf("%++v\n", pose),
				image.Pt(50,50),
				gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

			ring.Draw(&frame, pose, drone)
			_ = id
		}

		window.IMShow(frame)
		key := window.WaitKey(1)

		switch key {
		case keyboard.Spacebar: // space
			if inflight {
				drone.Land()
			} else {
				drone.TakeOff()
			}
			inflight = !inflight
		}

	}

}
