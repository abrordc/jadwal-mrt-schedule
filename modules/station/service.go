package station

import (
	"net/http"
	"time"

	"github.com/abrordc/jadwal-mrt-schedule/common/client"

	"encoding/json"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
}

type service struct {
	client *http.Client
}

func NewService() Service {
	return &service {
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func(s *service) GetAllStation() (response []StationResponse, err error) {
	// service
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	//hit url
	byteResponse,err := client.DoRequest(s.client,url) 

	if err != nil {
		return
	}

	var stations []Station
	err = json.Unmarshal(byteResponse,&stations)

	for _,item := range stations {
		response = append(response, StationResponse{
			Id: item.Id,
			Name: item.Name,
		})
	}


	return
}