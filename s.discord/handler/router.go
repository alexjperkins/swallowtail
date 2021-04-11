package handler

import "github.com/monzo/typhon"

var router = typhon.Router{}

func init() {
	router.PUT("/post-to-channel", POSTToChannel)
}

// Service returns a function which handles requests for this service
func Service() typhon.Service {
	return router.Serve()
}
