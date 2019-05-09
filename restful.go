package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/ethacctdb/api"
)

const restfulPort = 7210

//RestfulServer ...
type RestfulServer struct {
	serv *http.Server
}

//Start starts restful service
func (r *RestfulServer) Start() {
	go func() {
		r.serv.ListenAndServe()
	}()
}

//Stop stops restful service
func (r *RestfulServer) Stop() {
	r.serv.Shutdown(context.Background())
}

//NewRestfulServer creates a restful service
func NewRestfulServer() *RestfulServer {
	router := gin.New()
	v1Group := router.Group("/v1")
	accGroup := v1Group.Group("/account")
	miscGroup := v1Group.Group("/misc")

	accGroup.GET("/txnscnt/:addr", api.GetAccountTotalTxnsCount)
	accGroup.GET("/txns/:addr", api.GetAccountTxns)

	miscGroup.GET("/hashinfo/:hash", api.GetHashInfo)

	serv := &http.Server{Addr: ":" + strconv.FormatInt(restfulPort, 10), Handler: router}

	r := &RestfulServer{
		serv: serv,
	}

	return r
}
