package aggregator

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type MyServer struct {
	http.Server
	shutdownReq      chan bool
	reqCount         uint32
	newOrderChan     chan Order
	requestOrderChan chan OrderRequest
	reviewOrderChan  chan OrderReview
	finOrderChan     chan Order
}

func StartServer(listenPort string, newOrderChan chan Order, requestOrderChan chan OrderRequest, reviewOrderChan chan OrderReview, finOrderChan chan Order) {
	server := NewServer(listenPort, newOrderChan, requestOrderChan, reviewOrderChan, finOrderChan)

	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	<-done
	log.Println("DONE!")
}

func NewServer(listenPort string, _newOrderChan chan Order, _requestOrderChan chan OrderRequest, _reviewOrderChan chan OrderReview, _finOrderChan chan Order) *MyServer {
	s := &MyServer{
		Server: http.Server{
			Addr:         listenPort,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq:      make(chan bool),
		newOrderChan:     _newOrderChan,
		requestOrderChan: _requestOrderChan,
		reviewOrderChan:  _reviewOrderChan,
		finOrderChan:     _finOrderChan,
	}

	router := mux.NewRouter()

	//register handlers
	router.HandleFunc("/shutdown", s.ShutdownHandler)
	router.HandleFunc("/order/new", s.receiveNewOrder).Methods("POST")
	router.HandleFunc("/order/review", s.receiveReview).Methods("POST")
	router.HandleFunc("/order/forward/request", s.forwardRequest).Methods("POST")
	router.HandleFunc("/order/finished", s.receiveFinishedOrder).Methods("POST")

	//set http server handler
	s.Handler = router

	return s
}

func (s *MyServer) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))

	//Do nothing if shutdown request already issued
	//if s.reqCount == 0 then set to 1, return true otherwise false
	if !atomic.CompareAndSwapUint32(&s.reqCount, 0, 1) {
		log.Printf("Shutdown through API call in progress...")
		return
	}

	go func() {
		s.shutdownReq <- true
	}()
}

func (s *MyServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *MyServer) receiveNewOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Invalid request")
	}

	var parsedRequest Order
	json.Unmarshal(reqBody, &parsedRequest)

	s.newOrderChan <- parsedRequest

}

func (s *MyServer) receiveReview(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Invalid request")
	}

	var parsedRequest OrderReview
	json.Unmarshal(reqBody, &parsedRequest)

	s.reviewOrderChan <- parsedRequest
}

func (s *MyServer) forwardRequest(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Invalid request")
	}

	var parsedRequest OrderRequest
	json.Unmarshal(reqBody, &parsedRequest)
	s.requestOrderChan <- parsedRequest

}

func (s *MyServer) receiveFinishedOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Invalid request")
	}

	var parsedRequest Order
	json.Unmarshal(reqBody, &parsedRequest)
	s.finOrderChan <- parsedRequest
}
