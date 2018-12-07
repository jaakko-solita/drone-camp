package main

import (
	"fmt"

	"deus.solita.fi/Solita/projects/drone_code_camp/repositories/git/ddr.git"
	"gobot.io/x/gobot/platforms/keyboard"
	"image"
	"image/color"
	"gocv.io/x/gocv"
//	"github.com/go-gl/mathgl/mgl32"
	"runtime"
)
const (
	xetaisyys = 0.0
	yetaisyys = 0.0
	zetaisyys = 1.0
)
func main() {
	runtime.LockOSThread()
	track := ddr.NewTrack()
	defer track.Close()

	inflight := false

	window := gocv.NewWindow("Drone")
	drone := ddr.NewDrone(ddr.DroneReal, "./drone-camera-calibration-720.yaml")
	//drone := ddr.NewDrone(ddr.DroneFake, "./camera-calibration.yaml")
	err := drone.Init()

	if err != nil {
		fmt.Println("error while initializing drone: %v\n", err)
		return
	}
	
	var counter = 0
	for {
		frame := <-drone.VideoStream()

		// detect markers in this frame
		markers := track.GetMarkers(&frame)

		rings := track.ExtractRings(markers)
		if(counter > 100){
			for id, ring := range rings {
				pose := ring.EstimatePose(drone)
				// ringFace := pose.Rotation.Mul3x1(mgl32.Vec3{0.0, 0.0, 1.0})
				position := drone.CameraToDroneMatrix().Mul3x1(pose.Position)
			
				gocv.PutText(&frame, fmt.Sprintf("%++v\n", pose.Position),
					image.Pt(50,50),
					gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

				if (position[0]< -xetaisyys) {
					gocv.PutText(&frame, fmt.Sprintf("%++v\n", position),
									image.Pt(100,100),
									gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

					drone.Left(10)
					fmt.Println("vasuri");
				}
				if (position[0]>xetaisyys) {
					gocv.PutText(&frame, fmt.Sprintf("%++v\n", position),
									image.Pt(100,100),
									gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

					drone.Right(10)
					fmt.Println("oikee");
				}

				if (position[1]< -yetaisyys) {
					gocv.PutText(&frame, fmt.Sprintf("%++v\n", position),
									image.Pt(100,100),
									gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

					drone.Up(10)
					fmt.Println("ylös");
				}
				if (position[1]>yetaisyys) {
				gocv.PutText(&frame, fmt.Sprintf("%++v\n", position),
								image.Pt(100,100),
								gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

				drone.Down(10)
				fmt.Println("alas");
				}

				if (position[2]< zetaisyys) {
					gocv.PutText(&frame, fmt.Sprintf("%++v\n", position),
									image.Pt(100,100),
									gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

					drone.Backward(10)
					fmt.Println("taakse", position[2]);
				}
				if (position[2]>zetaisyys) {
				gocv.PutText(&frame, fmt.Sprintf("%++v\n", position[2]),
								image.Pt(100,100),
								gocv.FontHersheySimplex, 0.8, color.RGBA{0, 0, 0, 0}, 2)

				drone.Forward(10)
				fmt.Println("eteen", position[2]);
				}

				ring.Draw(&frame, pose, drone)
				_ = id
			}
		}
		if(len(rings) < 1){
			fmt.Println("Resetoidaan")
			drone.Hover()
		}
		fmt.Println(counter)
		window.IMShow(frame)
		key := window.WaitKey(1)

		switch key {
		case keyboard.Spacebar: // space
			if inflight {
				drone.Land()
				counter = 0
			} else {
				drone.TakeOff()
				counter = 1;
			}
			inflight = !inflight
		}
		if(counter > 0){
		counter++
		}
	}

}
