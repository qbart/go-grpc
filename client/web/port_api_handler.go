package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/qbart/go-grpc/pb"
	"go.uber.org/zap"
)

type PortAPIHandler struct {
	Logger            *zap.Logger
	PortDomainService pb.PortDomainServiceClient
}

func (h PortAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	port, err := h.PortDomainService.Get(ctx, &pb.PortId{Id: params["id"]})

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if err != nil {
			h.Logger.Error("Encoding error", zap.Error(err))
		}
		return
	}

	if err = json.NewEncoder(w).Encode(port); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		if err != nil {
			h.Logger.Error("Encoding error", zap.Error(err))
		}
		return
	}
}
