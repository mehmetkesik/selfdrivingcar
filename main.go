package main

import (
	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"image"
	"image/color"
	"os"
)

const (
	WIDTH  = 640
	HEIGHT = 480
)

type Scene struct {
	window     *sdlcanvas.Window
	canvas     *canvas.Canvas
	frameCount int64
	bg         image.Image
	car        *Car
	epoch      int
}

func (self *Scene) getBGPixel(x, y int) color.Color {
	return self.bg.At(x, y)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var scene = new(Scene)
	var err error
	scene.window, scene.canvas, err = sdlcanvas.CreateWindow(WIDTH, HEIGHT, "Yapay sinir ağları ile kendi kendini süren araba")
	if err != nil {
		panic(err)
	}
	defer scene.window.Destroy()

	start(scene)
	restart(scene)

	//fps sınırı - 1 olursa monitöre göre 0 olursa sınırsız
	err = sdl.GLSetSwapInterval(1);
	checkErr(err)

	scene.window.MainLoop(func() {
		scene.frameCount++
		update(scene)
	})

	//oyun kapatılınca
	scene.car.brain.Save("brain.json")

}

func start(scene *Scene) {
	road, err := os.Open("asset/road.png")
	checkErr(err)
	defer road.Close()

	scene.bg, _, err = image.Decode(road)
	checkErr(err)

	scene.car = newCar(scene)

	if _, err := os.Stat("brain.json"); !os.IsNotExist(err) {
		scene.car.brain.Load("brain.json")
	}
	iconSurface, err := img.Load("asset/icon.png")
	checkErr(err)
	defer iconSurface.Free()
	scene.window.Window.SetIcon(iconSurface)
}

func restart(scene *Scene) {
	scene.epoch = 0
	scene.frameCount = 0
	scene.car.isOne = true
}

func update(scene *Scene) {
	scene.frameCount++
	w, h := float64(scene.canvas.Width()), float64(scene.canvas.Height())
	scene.canvas.SetFillStyle("#fff")
	scene.canvas.FillRect(0, 0, w, h)

	drawBackground(scene)

	scene.car.paint()

	ls := scene.car.leftSensor()
	fs := scene.car.frontSensor()
	rs := scene.car.rightSensor()

	if ls == 1 && fs == 1 && rs == 1 {
		scene.car.direction = 1;
		if scene.car.isOne {
			scene.car.reward() //eğer sensörler 1 e ulaşmışsa tek seferlik ödüllendiriyoruz.
		}
		scene.car.isOne = false
		scene.car.clearInputs();
		return;
	} else {
		scene.car.isOne = true
	}

	if ls == 0 || fs == 0 || rs == 0 {
		scene.epoch++
		scene.car.direction = 1;
		if scene.frameCount > scene.car.prevLifeTime {
			scene.car.reward() //ödüllendirme
		} else {
			scene.car.panish() //cezalandırma
		}
		scene.car.prevLifeTime = scene.frameCount
		scene.car.reset()
		restart(scene)
		return
	}

	result := scene.car.brain.Predict([]float64{ls, fs, rs})

	//result[0] sola dön, result[1] hiçbirşey yapma-devam, result[2] sağa dön
	if result[0] > result[1] {
		if (result[0] > result[2]) {
			//sola dön
			scene.car.steerLeft();
			scene.car.direction = 0;
		} else {
			//sağa dön
			scene.car.steerRight();
			scene.car.direction = 2;
		}
	} else {
		if (result[1] > result[2]) {
			//düz git
			scene.car.direction = 1;
		} else {
			//sağa dön
			scene.car.steerRight();
			scene.car.direction = 2;
		}
	}

	scene.car.addInput([]float64{ls, fs, rs, float64(scene.car.direction)});
}

func drawBackground(scene *Scene) {
	scene.canvas.DrawImage(scene.bg, 0, 0)
}
