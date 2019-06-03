package store

//Users stores the users by token and vehicleID
var Users = newUsers()

// Vehicles contains the list of vehicles
var Vehicles = newVehicleList()

// Positions contains the list of vehicle positions
var Positions = newPositionList()

//Queue stores the users waiting for their turn
var Queue = newQueue()

type key interface{}
