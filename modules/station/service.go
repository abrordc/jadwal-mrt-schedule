package station

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/abrordc/jadwal-mrt-schedule/common/client"

	"encoding/json"
)

type Service interface {
	GetAllStation() (response []StationResponse, err error)
	CheckSchedulesByStation(id string) (response []ScheduleResponse, err error)
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

func (s *service) CheckSchedulesByStation(id string) (response []ScheduleResponse, err error){
	url := "https://www.jakartamrt.co.id/id/val/stasiuns"

	//hit url
	byteResponse,err := client.DoRequest(s.client,url) 

	if err != nil {
		return
	}
	var schedule []Schedule
	err = json.Unmarshal(byteResponse,&schedule)
	if err != nil {
		return
	}

	//schedule selected by id station

	var scheduleSelected Schedule
	for _,item := range schedule {
		if item.StationId == id {
			scheduleSelected= item
			break
		}
	}

	if scheduleSelected.StationId == ""{
		err = errors.New("station not found")
		return
	}

	response, err = ConvertDataToResponse(scheduleSelected)
	if err != nil {
		return
	}

	return


}

func ConvertDataToResponse(schedule Schedule) (response []ScheduleResponse, err error){
	var (
		LebakBulusTripName = "Station Lebak Bulus Grab"
		BundaranHITripName = "Station Bundaran HI Bank DKI"
	)

	scheduleLebakBulus := schedule.ScheduleLebakBulus
	scheduleBundaranHI := schedule.ScheduleBundaranHi

	scheduleLebakBulusParsed,err := ConvertScheduleToTimeFormat(scheduleLebakBulus)
	if err != nil {
		return
	}

	scheduleBundaranHIParsed,err := ConvertScheduleToTimeFormat(scheduleBundaranHI)
	if err != nil {
		return
	}

	//convert to response
	for _,item := range scheduleLebakBulusParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: LebakBulusTripName,
				Time: item.Format("15:04"),
			})
		}
	}

	for _,item := range scheduleBundaranHIParsed {
		if item.Format("15:04") > time.Now().Format("15:04") {
			response = append(response, ScheduleResponse{
				StationName: BundaranHITripName,
				Time: item.Format("15:04"),
			})
		}
	}

	return
	
}

func ConvertScheduleToTimeFormat(schedule string)(response []time.Time, err error){
	var (
		parsedTime time.Time
		schedules = strings.Split(schedule,",")
	)

	for _,item := range schedules {
		trimedTime := strings.TrimSpace(item)
		if trimedTime == "" {
			continue
		}
		parsedTime, err = time.Parse("15:04", trimedTime)
		if err != nil {
			err = errors.New("invalid time format"+ trimedTime)
			return
		}
		response = append(response, parsedTime)
	}
	return
}
