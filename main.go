package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/volons/hive/controllers"
	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/db"
	"github.com/volons/hive/libs/websocket"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
	"github.com/volons/hive/models/config"
	"github.com/volons/hive/platform"

	"github.com/gorilla/mux"
	//_ "net/http/pprof"
)

func main() {
	log.SetFlags(log.Lshortfile)
	//log.SetOutput(ioutil.Discard) // Disable logs

	//
	// Init configuration
	//
	configFilePath := flag.String("config", "", "path to the config file")
	flag.Parse()
	config.Read(*configFilePath)
	conf := config.Get()

	//
	// Init database
	//
	err := db.Init(conf.Database)
	if err != nil {
		log.Println(err)
	}

	//
	// Init fence
	//
	fence := models.FenceData{}
	err = db.Get("fence", &fence)
	if err != nil {
		log.Printf("Cound not get fence: %v", err)
	} else {
		err = models.SetFence(fence)
		if err != nil {
			log.Printf("Cound not set fence: %v", err)
		}
	}

	//
	// Init connection to the Volons platform if configured
	//
	if conf.VolonsPlatform != "" {
		go platform.Platform.Run(conf.VolonsPlatform)
	}

	//
	// Init routes
	//
	router := mux.NewRouter()

	// Default page
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Init vehicle websocket listener
	vehicleWS := websocket.NewServer(messages.NewParser(func(typ string) interface{} {
		switch typ {
		case "position", "goto":
			return &models.Position{}
		case "battery":
			return &models.Battery{}
		case "rc":
			return &models.Rc{}
		case "fence":
			return &models.Fence{}
		case "webrtc:sdp":
			return &models.SessionDescription{}
		case "webrtc:icecandidate":
			return &models.IceCandidate{}
		case "webrtc:start", "takeoff", "land", "rtl":
			return &struct{}{}
		case "caps":
			return &models.Caps{}
		default:
			return &libs.JSONObject{}
		}
	}))
	vehicleWS.SetConnectionListener(controllers.Vehicle)
	router.Handle("/vehicle", vehicleWS)

	// Init admin websocket listener
	adminWS := websocket.NewServer(messages.NewParser(func(typ string) interface{} {
		switch typ {
		case "position", "goto":
			return &models.Position{}
		case "battery":
			return &models.Battery{}
		case "rc":
			return &models.Rc{}
		case "fence:set":
			return &models.FenceData{}
		case "webrtc:sdp":
			return &models.SessionDescription{}
		case "webrtc:icecandidate":
			return &models.IceCandidate{}
		case "webrtc:start", "takeoff", "land", "rtl":
			return &struct{}{}
		case "caps":
			return &models.Caps{}
		default:
			return &libs.JSONObject{}
		}
	}))
	adminWS.SetConnectionListener(controllers.Admin)
	router.Handle("/admin", adminWS)

	// Init webrtc websocket
	//ws = new(websocket.Server)
	//ws.SetConnectionListener(controllers.WebRTC.ConnectionListener)
	//router.Handle("/webrtc", ws)

	//log.Printf("Listening on %s\n", conf.HTTPAddr)
	//corsHandler := handlers.CORS(
	//	handlers.AllowedMethods([]string{"GET", "POST", "DELETE"}),
	//	handlers.AllowedOrigins([]string{"*"}),
	//)(router)

	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	//vehicles.NewDummyVehicle("127.0.0.1:12345")

	//http.ListenAndServe(conf.HTTPAddr, corsHandler)

	log.Printf("Listening on %s\n", conf.HTTPAddr)
	err = http.ListenAndServe(conf.HTTPAddr, router)
	if err != nil {
		panic(err.Error())
	}
}

//func handler(fn func(http.ResponseWriter, *http.Request) (string, error)) func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		result, err := fn(w, r)
//		if err != nil {
//			http.Error(w, err.Error(), 500)
//			return
//		}
//
//		fmt.Fprintf(w, result)
//	}
//}
