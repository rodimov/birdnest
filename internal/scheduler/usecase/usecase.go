package usecase

import (
	dronesPackage "birdnest/internal/drones"
	dronesDelivery "birdnest/internal/drones/delivery"
	dronesModel "birdnest/internal/drones/model"
	"birdnest/internal/events"
	pilotsPackage "birdnest/internal/pilots"
	pilotsDelivery "birdnest/internal/pilots/delivery"
	pilotsModel "birdnest/internal/pilots/model"
	"birdnest/internal/scheduler"
	"birdnest/logger"
	"birdnest/pkg"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"
)

type Report struct {
	XMLName    xml.Name `xml:"report"`
	DeviceInfo Device   `xml:"deviceInformation"`
	Snapshots  Capture  `xml:"capture"`
}

type Device struct {
	XMLName  xml.Name `xml:"deviceInformation"`
	Interval int      `xml:"updateIntervalMs"`
}

type Capture struct {
	XMLName           xml.Name `xml:"capture"`
	SnapshotTimestamp string   `xml:"snapshotTimestamp,attr"`
	Drones            []Drone  `xml:"drone"`
}

type Drone struct {
	XMLName      xml.Name `xml:"drone"`
	SerialNumber string   `xml:"serialNumber"`
	PositionX    float64  `xml:"positionX"`
	PositionY    float64  `xml:"positionY"`
}

type PilotData struct {
	ID          string `json:"pilotId"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	CreatedDt   string `json:"createdDt"`
	Email       string `json:"email"`
}

type usecase struct {
	dronesUseCase           dronesPackage.UseCase
	pilotsUseCase           pilotsPackage.UseCase
	mutex                   sync.Mutex
	ndzClearTimeMin         int64
	dronesListURL           string
	pilotsInfoURL           string
	lastDronesListTimestamp time.Time
	eventsChan              chan *events.EventMessage
}

func createNewViolationEvent(pilot *pilotsModel.Pilot) *events.EventMessage {
	return &events.EventMessage{
		EventType: "new_violation",
		Data:      pilotsDelivery.ToPilotDto(pilot),
	}
}

func createNewDroneEvent(drone *dronesModel.Drone) *events.EventMessage {
	return &events.EventMessage{
		EventType: "new_drone",
		Data:      dronesDelivery.ToDroneDto(drone),
	}
}

func createDeletePilotEvent(pilot *pilotsModel.Pilot) *events.EventMessage {
	return &events.EventMessage{
		EventType: "delete_pilot",
		Data:      pilotsDelivery.ToPilotDto(pilot),
	}
}

func NewUseCase(dronesUseCase dronesPackage.UseCase,
	pilotsUseCase pilotsPackage.UseCase,
	ndzClearTimeMin int64,
	dronesListURL string,
	pilotsInfoURL string,
	eventsChan chan *events.EventMessage) scheduler.UseCase {

	return &usecase{
		dronesUseCase:   dronesUseCase,
		pilotsUseCase:   pilotsUseCase,
		ndzClearTimeMin: ndzClearTimeMin,
		dronesListURL:   dronesListURL,
		pilotsInfoURL:   pilotsInfoURL,
		eventsChan:      eventsChan,
	}
}

func (u *usecase) StartScheduler() error {
	go u.updateDronesList()

	time.Sleep(time.Duration(2) * time.Second)

	go u.hardRemoveOldDrones()
	go u.permanentlyDeleteGarbage()

	return nil
}

func (u *usecase) updateDronesList() {
	for {
		xmlData, err := u.makeDronesListRequest()

		if err != nil {
			logger.AppLogger.Error("Can't get list of drones! " + err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		var report Report
		err = xml.Unmarshal(xmlData, &report)

		if err != nil {
			logger.AppLogger.Error("Can't read list of drones!" + err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		if err != nil {
			logger.AppLogger.Error("Can't parse snapshot time!" + err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		err = u.createDrones(&report.Snapshots)

		if err != nil {
			logger.AppLogger.Error("Can't create drones!" + err.Error())
			time.Sleep(5 * time.Second)
			continue
		}

		snapshotTime, err := time.Parse(time.RFC3339, report.Snapshots.SnapshotTimestamp)

		if err != nil {
			logger.AppLogger.Error("Can't parse snapshot time!" + err.Error())
		} else {
			u.lastDronesListTimestamp = snapshotTime
		}

		time.Sleep(time.Millisecond * time.Duration(report.DeviceInfo.Interval))
	}
}

func (u *usecase) makeDronesListRequest() ([]byte, error) {
	response, err := http.Get(u.dronesListURL)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *usecase) createDrones(capture *Capture) error {
	ctx := context.Background()

	snapshotTime, err := time.Parse(time.RFC3339, capture.SnapshotTimestamp)

	if err != nil {
		return err
	}

	for _, drone := range capture.Drones {
		droneModelItem := dronesModel.Drone{
			ID:        drone.SerialNumber,
			PositionX: drone.PositionX,
			PositionY: drone.PositionY,
			LastSeen:  snapshotTime,
		}

		if !isInNDZ(drone.PositionX, drone.PositionY) {
			//u.eventsChan <- createNewDroneEvent(&droneModelItem)
			continue
		}

		u.mutex.Lock()
		_, err = u.dronesUseCase.Create(ctx, &droneModelItem)

		if err != nil {
			u.mutex.Unlock()
			return err
		}

		//pilotInfo, err := u.createPilot(&droneModelItem)
		_, err := u.createPilot(&droneModelItem)

		if err == nil {
			//u.eventsChan <- createNewViolationEvent(pilotInfo)
		} else {
			logger.AppLogger.Error("Can't send new violation event: no pilot info")
		}

		u.mutex.Unlock()

		go u.removeDrone(&droneModelItem)
	}

	return nil
}

func isInNDZ(x float64, y float64) bool {
	const XC float64 = 250000
	const YC float64 = 250000
	const R float64 = 100000

	return math.Pow(x-XC, 2)+math.Pow(y-YC, 2) < math.Pow(R, 2)
}

func (u *usecase) removeDrone(drone *dronesModel.Drone) {
	savingDuration := time.Minute * time.Duration(u.ndzClearTimeMin)

	time.Sleep(savingDuration)
	ctx := context.Background()

	u.mutex.Lock()
	defer u.mutex.Unlock()

	isDroneRemoved, err := u.dronesUseCase.IsRemoved(ctx, drone.ID)

	if err != nil {
		logger.AppLogger.Error("Can't check drone with id = " + drone.ID)
		logger.AppLogger.Error(err.Error())
		return
	}

	originalLastSeen := drone.LastSeen
	drone, err = u.dronesUseCase.GetById(ctx, drone.ID)

	if err != nil {
		logger.AppLogger.Error("Can't get drone with id = " + drone.ID)
		logger.AppLogger.Error(err.Error())
		return
	}

	if !isDroneRemoved && (originalLastSeen == drone.LastSeen) {
		_, err := u.dronesUseCase.DeleteById(ctx, drone.ID)

		if err != nil {
			logger.AppLogger.Error("Can't delete drone with id = " + drone.ID)
			logger.AppLogger.Error(err.Error())
		} else {
			//pilotInfo, err := u.pilotsUseCase.GetByDroneId(ctx, drone.ID)
			//
			//if err == nil {
			//	u.eventsChan <- createDeletePilotEvent(pilotInfo)
			//} else {
			//	logger.AppLogger.Error("Can't send delete pilot event: no pilot info")
			//}
		}
	}
}

func (u *usecase) createPilot(drone *dronesModel.Drone) (*pilotsModel.Pilot, error) {
	ctx := context.Background()
	isExists := u.pilotsUseCase.IsExistsByDroneId(ctx, drone.ID)

	if isExists {
		return u.pilotsUseCase.GetByDroneId(ctx, drone.ID)
	}

	jsonData, err := u.makeGetPilotInfoRequest(drone.ID)

	if err != nil {
		logger.AppLogger.Info("Get pilot request is failed with a message: " + err.Error())
		logger.AppLogger.Info("Saving unknown pilot information to database")

		pilot := &pilotsModel.Pilot{
			ID:      drone.ID,
			DroneID: drone.ID,
			Drone:   drone,
		}

		return u.pilotsUseCase.Create(ctx, pilot)
	}

	var pilotData PilotData
	err = json.Unmarshal(jsonData, &pilotData)

	if err != nil {
		return nil, err
	}

	registrationTime, _ := time.Parse(time.RFC3339, pilotData.CreatedDt)

	pilot := &pilotsModel.Pilot{
		ID:               pilotData.ID,
		FirstName:        pkg.StringToNullString(pilotData.FirstName),
		LastName:         pkg.StringToNullString(pilotData.LastName),
		Phone:            pkg.StringToNullString(pilotData.PhoneNumber),
		Email:            pkg.StringToNullString(pilotData.Email),
		RegistrationTime: pkg.TimeToNullTime(registrationTime),
		DroneID:          drone.ID,
		Drone:            drone,
	}

	return u.pilotsUseCase.Create(ctx, pilot)
}

func (u *usecase) makeGetPilotInfoRequest(droneId string) ([]byte, error) {
	response, err := http.Get(u.pilotsInfoURL + droneId)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("request is failed with error code: %v", response.StatusCode)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (u *usecase) hardRemoveOldDrones() {
	for {
		u.mutex.Lock()
		ctx := context.Background()
		drones, err := u.dronesUseCase.GetAll(ctx)

		if err != nil {
			u.mutex.Unlock()
			logger.AppLogger.Error("Can't get drones list to clear table")
			time.Sleep(time.Minute * time.Duration(2))
			continue
		}

		for _, drone := range drones {
			timeToHardRemove := drone.LastSeen.Add(time.Minute * time.Duration(2*u.ndzClearTimeMin))

			if timeToHardRemove.Before(u.lastDronesListTimestamp) {
				_, err := u.dronesUseCase.DeleteById(ctx, drone.ID)

				if err != nil {
					logger.AppLogger.Error("Can't delete drone with id = " + drone.ID)
					logger.AppLogger.Error(err.Error())
				} else {
					logger.AppLogger.Info("Drone with id = " + drone.ID + " was hard removed")
					//pilotInfo, err := u.pilotsUseCase.GetByDroneId(ctx, drone.ID)
					//
					//if err == nil {
					//	u.eventsChan <- createDeletePilotEvent(pilotInfo)
					//} else {
					//	logger.AppLogger.Error("Can't send delete pilot event: no pilot info")
					//}
				}
			}
		}
		u.mutex.Unlock()

		time.Sleep(time.Minute * time.Duration(u.ndzClearTimeMin))
	}
}

func (u *usecase) permanentlyDeleteGarbage() {
	for {
		u.mutex.Lock()
		ctx := context.Background()
		drones, err := u.dronesUseCase.GetAllDeleted(ctx)

		if err != nil {
			u.mutex.Unlock()
			logger.AppLogger.Error("Can't get deleted drones list to clear garbage")
			time.Sleep(time.Minute * time.Duration(2))
			continue
		}

		for _, drone := range drones {
			timeToPermanentRemove := drone.DeletedAt.Time.Add(time.Minute * time.Duration(2*u.ndzClearTimeMin))
			now := time.Now()
			now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(),
				now.Minute(), now.Day(), now.Nanosecond(), time.UTC)

			if timeToPermanentRemove.Before(now) {
				_, err := u.dronesUseCase.DeletePermanentlyById(ctx, drone.ID)

				if err != nil {
					logger.AppLogger.Error("Can't delete permanently drone with id = " + drone.ID)
					logger.AppLogger.Error(err.Error())
				}
			}
		}
		u.mutex.Unlock()

		time.Sleep(time.Minute * time.Duration(u.ndzClearTimeMin))
	}
}
