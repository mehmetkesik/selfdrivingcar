# self-driving car
self-driving car with deep learning

![pic](https://github.com/mehmetkesik/selfdrivingcar/blob/master/asset/pic.png)

# Usage
**sdl2** graphic library was usedin this project.
Before compiling the project, **sd2** must be installed and the https://github.com/tfriedel6/canvas library should be installed.
later `go build` and run..
<br/><br/>
The file **brain.json** in the main folder is the file where the artificial neural network is saved.
if the **brain.json** file is deleted, the training starts from the beginning.

# Training
The values of the sensors of the car are given as input to the artificial neural network and decide whether to go right, left or straight as output.
<br/><br/>
The car's sensor data is kept in a sequence until **1,1,1** and when the sensor data is **1,1,1** the inputs are rewarded and the sensor data is cleared. if the car crashes without sensor **1,1,1** the entries up to that time will be penalized. this way the car learns to go without crashing.
