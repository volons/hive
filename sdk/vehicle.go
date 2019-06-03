package sdk

// sdkVehicle is a vehicle that connects through the native SDK
//type sdkVehicle struct {
//	models.VehicleCore
//
//	autopilot messages.Line
//
//	i      VehicleI // Interface to native vehicle
//	parser messages.Parser
//
//	done chan bool
//}

// newSDKVehicle creates a new sdkVehicle
//func newSDKVehicle(id string, i VehicleI) *sdkVehicle {
//	vehicle := &sdkVehicle{
//		VehicleCore: models.NewVehicleCore(
//			id, nil, models.Caps{},
//		),
//		i:         i,
//		autopilot: messages.NewLine(true),
//		parser: messages.NewParser(func(typ string) interface{} {
//			switch typ {
//			case "position", "goto":
//				return &models.Position{}
//			case "battery":
//				return &models.Battery{}
//			case "rc":
//				return &models.Rc{}
//			case "fence":
//				return &models.Fence{}
//			case "webrtc:sdp":
//				return &models.SessionDescription{}
//			case "webrtc:icecandidate":
//				return &models.IceCandidate{}
//			case "webrtc:start", "takeoff", "land", "rtl":
//				return &struct{}{}
//			case "caps":
//				return &models.Caps{}
//			default:
//				return &libs.JSONObject{}
//			}
//		}),
//	}
//
//	i.SetListener(vehicle)
//
//	ap := autopilot.Get(id)
//	ap.ConnectVehicle(vehicle.autopilot)
//
//	store.Vehicles.Add(vehicle)
//
//	return vehicle
//}

//func (v *sdkVehicle) run() {
//	for {
//		select {
//		case msg := <-v.autopilot.Recv():
//			v.Send(msg)
//		}
//	}
//}

// Send sends a message to the vehicle
//func (v *sdkVehicle) Send(msg messages.Message) error {
//	done := make(chan error, 1)
//	json, err := msg.ToJSON()
//
//	if err != nil {
//		done <- err
//	} else {
//		go v.i.HandleMessage(json, cb(func(err error) {
//			done <- err
//		}))
//	}
//
//	return <-done
//}
//
//func (v *sdkVehicle) OnMessage(json string) {
//	msg, err := v.parser.Parse([]byte(json))
//	if err != nil {
//		log.Println("error parsing json message")
//		return
//	}
//
//	v.autopilot.Send(msg)
//}
//
//// Done returns the done channel
//func (v *sdkVehicle) Done() <-chan bool {
//	return v.done
//}
//
//// OnDisconnected removes the vehicle from the list of connected devices
//func (v *sdkVehicle) Remove() {
//	store.Vehicles.Remove(v)
//	close(v.done)
//}

//// GetCaps returns the vehicle's capabilities
//func (v *sdkVehicle) GetCaps() models.Caps {
//	return v.caps
//}

// SetCaps sets the vehicle's capabilities
//func (v *sdkVehicle) SetCaps(json string) error {
//	var err error
//	v.caps, err = models.ParseCaps(json)
//	return err
//}

// TakeOff tells the vehcile to takeoff
//func (v *sdkVehicle) TakeOff() error {
//	done := make(chan error)
//	go v.i.TakeOff(cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}

// Land tells the vehcile to land
//func (v *sdkVehicle) Land() error {
//	done := make(chan error)
//	go v.i.Land(cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}

// RTL tells the vehcile to return home and land
//func (v *sdkVehicle) RTL() error {
//	done := make(chan error)
//	go v.i.RTL(cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}

// GoTo sets the vehicle's target position
//func (v *sdkVehicle) GoTo(pos models.Position) error {
//	done := make(chan error)
//	go v.i.GoTo(pos.Lat(), pos.Lon(), pos.RelAlt(), cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}

// Rc sets the vehicle's rc values
//func (v *sdkVehicle) Rc(rc *models.Rc) error {
//	done := make(chan error)
//
//	if rc == nil {
//		go v.i.Rc(0, 0, 0, 0, 0,
//			cb(func(err error) {
//				done <- err
//			}),
//		)
//	} else {
//		go v.i.Rc(
//			rc.Roll(),
//			rc.Pitch(),
//			rc.Yaw(),
//			rc.Throttle(),
//			rc.Gimbal(),
//			cb(func(err error) {
//				done <- err
//			}),
//		)
//	}
//
//	return <-done
//}

// StartWebRTCSession starts the WebRTC live video feed connection
//func (v *sdkVehicle) StartWebRTCSession() error {
//	done := make(chan error)
//	go v.i.StartWebRTCSession(cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}

// SendSDP sends the user's WebRTC session description to the vehicle
//func (v *sdkVehicle) SendSDP(sdp interfaces.SessionDescription) error {
//	done := make(chan error)
//	go v.i.SendSDP(sdp.Type, sdp.Sdp, cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}

// SendIceCandidate sends a user's WebRTC ice candidate to the vehicle
//func (v *sdkVehicle) SendIceCandidate(candidate interfaces.IceCandidate) error {
//	done := make(chan error)
//	go v.i.SendIceCandidate(
//		candidate.SdpMLineIndex,
//		candidate.SdpMid,
//		candidate.Candidate,
//		cb(func(err error) {
//			done <- err
//		}),
//	)
//
//	return <-done
//}

//OnWebRTCSessionDescription should be called with a WebRTC offer or answer to be sent to the client
//func (v *sdkVehicle) OnWebRTCSessionDescription(typ, desc string) {
//	sdp := interfaces.SessionDescription{
//		Type: typ,
//		Sdp:  desc,
//	}
//
//	v.Trigger("sdp", sdp)
//}

//OnWebRTCIceCandidate should be called with a WebRTC ice candidate to be sent to the client
//func (v *sdkVehicle) OnWebRTCIceCandidate(sdpMLineIndex float64, sdpMid, desc string) {
//	candidate := interfaces.IceCandidate{
//		SdpMLineIndex: sdpMLineIndex,
//		SdpMid:        sdpMid,
//		Candidate:     desc,
//	}
//
//	v.Trigger("icecandidate", candidate)
//}

//OnBattery should be called to update the vehicle's battery status
//func (v *sdkVehicle) OnBattery(voltage, current, percent float64) {
//	v.Trigger("battery", models.NewBattery(voltage, current, percent))
//}

//OnReport should be called to send a vehicle text status message to the connected clients
//func (v *sdkVehicle) OnReport(report string) {
//	v.Trigger("report", report)
//}

//OnPosition should be called to update the vehicle's position
//func (v *sdkVehicle) OnPosition(lat, lon, alt, hdg, relAlt, vx, vy, vz float64) {
//	v.Trigger("position", models.NewPosition(
//		lat, lon, alt, relAlt, vx, vy, vz, hdg,
//	))
//}

//OnUserMessage transfers a user's webrtc message to the connected fwdClient
//func (v *sdkVehicle) OnUserMessage(msg string) {
//	v.Trigger("usermsg", []byte(msg))
//}

//func (v *sdkVehicle) SendMessageToUser(msg string) error {
//	done := make(chan error)
//	go v.i.SendMessageToUser(msg, cb(func(err error) {
//		done <- err
//	}))
//
//	return <-done
//}
