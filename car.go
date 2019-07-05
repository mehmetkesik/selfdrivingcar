package main

import (
	"./gonn"
	"math"
)

type Car struct {
	x            float64
	y            float64
	w            float64
	h            float64
	scene        *Scene
	angle        int
	speed        int
	sensorLength int
	inputs       [1000][]float64
	inputIndex   int
	brain        *gonn.MyNetwork
	direction    int
	prevLifeTime int64
	isOne        bool
}

func newCar(scene *Scene) *Car {
	car := new(Car)
	car.scene = scene
	car.w = 25
	car.h = 15
	car.speed = 1
	car.sensorLength = 65
	car.brain = new(gonn.MyNetwork)
	car.brain.Init([]int{3, 7, 3})
	car.reset()
	return car
}

func (self *Car) paint() {
	delayrate := 1.0
	tx := self.x + self.w/2
	ty := self.y + self.h/2
	radianAngle := float64(self.angle) * math.Pi / 180
	x := delayrate * float64(self.speed) * math.Cos(radianAngle)
	y := delayrate * float64(self.speed) * math.Sin(radianAngle)
	self.x = self.x + x
	self.y = self.y + y
	self.scene.canvas.Translate(tx, ty)
	self.scene.canvas.Rotate(radianAngle)
	self.scene.canvas.Translate(-tx, -ty)
	self.scene.canvas.SetFillStyle("#000000")
	//araba çizimi
	self.scene.canvas.FillRect(math.Round(self.x), math.Round(self.y), self.w, self.h)
	self.scene.canvas.SetFillStyle("#ff0000")
	//lamba çizimleri
	self.scene.canvas.FillRect(self.x+self.w-7, self.y, 5, 5)
	self.scene.canvas.FillRect(self.x+self.w-7, self.y+self.h-5, 5, 5)
	self.scene.canvas.SetTransform(1, 0, 0, 1, 0, 0)
}

func (self *Car) steerLeft() {
	self.angle -= 1
	if (self.angle < 0) {
		self.angle = 360 + self.angle
	}
	self.angle = self.angle % 360
}

func (self *Car) steerRight() {
	self.angle += 1;
	self.angle = self.angle % 360;
}

func (self *Car) leftSensor() float64 {
	return self.sideSensor(false)
}

func (self *Car) rightSensor() float64 {
	return self.sideSensor(true)
}

func (self *Car) sideSensor(side bool) float64 {
	m := 0.0
	if side {
		m = float64((self.h / 2) / (self.w / 2))
	} else {
		m = float64(-(self.h / 2) / (self.w / 2))
	}

	aci := math.Atan(m)
	radianAngle := (float64(self.angle) * (math.Pi / 180)) + aci
	if (radianAngle < 0) {
		radianAngle = (math.Pi * 2) + radianAngle
	}

	tut := self.w / 2
	x := 0.0
	y := 0.0
	touch := false;
	i := 1
	sensorfark := 25
	for ; i <= self.sensorLength-sensorfark; i++ {
		farkx := (tut + float64(i)) * math.Cos(radianAngle)
		farky := (tut + float64(i)) * math.Sin(radianAngle)
		x = float64(self.x+(self.w/2)) + farkx
		y = float64(self.y+(self.h/2)) + farky
		id := self.scene.getBGPixel(int(math.Round(x)), int(math.Round(y)))
		r, g, b, _ := id.RGBA()
		if (r == 0 && g == 0 && b == 0) {
			touch = true
			self.scene.canvas.BeginPath()
			self.scene.canvas.SetStrokeStyle("#0000ffff")
			self.scene.canvas.Arc(x, y, 3, 0, math.Pi*2, false)
			self.scene.canvas.Stroke()
			break;
		}
	}

	if (!touch) {
		self.scene.canvas.BeginPath()
		self.scene.canvas.SetStrokeStyle("#0000ffff")
		self.scene.canvas.Arc(x, y, 3, 0, math.Pi*2, false)
		self.scene.canvas.Stroke()
	}

	self.scene.canvas.SetTransform(1, 0, 0, 1, 0, 0)
	return float64(i-1) / float64(self.sensorLength-sensorfark)
}

func (self *Car) frontSensor() float64 {
	radianAngle := float64(self.angle) * float64(math.Pi) / 180
	tut := self.w / 2
	x := 0.0
	y := 0.0
	touch := false
	i := 1
	for ; i <= self.sensorLength; i++ {
		farkx := (tut + float64(i)) * math.Cos(radianAngle)
		farky := (tut + float64(i)) * math.Sin(radianAngle)
		x = float64(self.x+(self.w/2)) + farkx
		y = float64(self.y+(self.h/2)) + farky
		id := self.scene.getBGPixel(int(math.Round(x)), int(math.Round(y)))
		r, g, b, _ := id.RGBA()
		if r == 0 && g == 0 && b == 0 {
			touch = true;
			self.scene.canvas.BeginPath()
			self.scene.canvas.SetStrokeStyle("#0000ff")
			self.scene.canvas.Arc(x, y, 3, 0, math.Pi*2, false)
			self.scene.canvas.Stroke()
			break;
		}
	}

	if (!touch) {
		self.scene.canvas.BeginPath()
		self.scene.canvas.SetStrokeStyle("#0000ff")
		self.scene.canvas.Arc(x, y, 3, 0, math.Pi*2, false)
		self.scene.canvas.Stroke()
	}

	self.scene.canvas.SetTransform(1, 0, 0, 1, 0, 0)
	return float64(i-1) / float64(self.sensorLength)
}

func (self *Car) reset() {
	self.x = 50
	self.y = 70
	self.angle = 0
	self.clearInputs()
	self.paint()
}

func (self *Car) addInput(input []float64) {
	self.inputs[self.inputIndex] = input
	self.inputIndex++
	if self.inputIndex == 1000 {
		self.inputIndex = 0
	}
}

func (self *Car) clearInputs() {
	var i [1000][]float64
	self.inputs = i
	self.inputIndex = 0
}

func (self *Car) reward() {
	for _, input := range self.inputs {
		if len(input) == 0 {
			continue
		}
		propArray := make([]float64, 0)
		if input[3] == 0 { //sola dönmüşse
			propArray = []float64{1, 0, 0}
		} else if input[3] == 1 { //düz gitmişse
			propArray = []float64{0, 1, 0}
		} else if input[3] == 2 { //sağa dönmüşse
			propArray = []float64{0, 0, 1}
		}
		self.brain.SingleTrain(input[:len(input)-1], propArray)
	}
}

func (self *Car) panish() {
	for _, input := range self.inputs {
		if len(input) == 0 {
			continue
		}
		rateArray := self.inputs[self.inputIndex-1]
		rate := (rateArray[0] + rateArray[1] + rateArray[2]) / 6
		propArray := make([]float64, 0)
		if input[3] == 0 { //sola dönmüşse
			propArray = []float64{rate, 1 - rate, 1}
		} else if input[3] == 1 { //düz gitmişse
			propArray = []float64{1 - rate, rate, 1 - rate}
		} else if input[3] == 2 { //sağa dönmüşse
			propArray = []float64{1 - rate, 1 - rate, rate}
		}
		self.brain.SingleTrain(input[:len(input)-1], propArray)
	}
}
