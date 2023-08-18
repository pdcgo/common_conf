package pdc_application

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type EventType string

const (
	StartEvent  = "start"
	FinishEvent = "finish"
)

type EventLogModel struct {
	TrackId     string              `json:"track_id" bigquery:"track_id"`
	License     string              `json:"license" bigquery:"license"`
	Bot         string              `json:"bot" bigquery:"bot"`
	Version     string              `json:"version" bigquery:"version"`
	Program     string              `json:"program" bigquery:"program"`
	Hostname    string              `json:"hostname" bigquery:"hostname"`
	Status      EventType           `json:"status" bigquery:"status"`
	Description bigquery.NullString `json:"description" bigquery:"description"`
	Timestamp   civil.DateTime      `json:"created_at" bigquery:"created_at"`
	ErrMsg      string              `json:"err_msg" bigquery:"err_msg"`
}

func (app *PdcApplication) CreateEventClient(cfg *AppFileConfig, logname string) (startEventLog func() error, endEventLog func(err error) error, err error) {

	startEventLog = func() error { return nil }
	endEventLog = func(err error) error { return nil }

	if len(app.Credential) == 0 {
		log.Println("event start not sending to cloud")
		return startEventLog, endEventLog, nil

	}

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, cfg.ProjectID, option.WithCredentialsJSON(app.Credential))

	if err != nil {
		return startEventLog, endEventLog, err
	}

	dataset := client.Dataset("event")
	table := dataset.Table("log")

	insert := table.Inserter()
	hostname, _ := os.Hostname()

	// creating inserter
	startEventLog = func() error {
		err := insert.Put(ctx, &EventLogModel{
			TrackId:   uuid.New().String(),
			License:   cfg.Lisensi.Email,
			Version:   app.Version,
			Hostname:  hostname,
			Bot:       logname,
			Program:   "golang",
			Status:    StartEvent,
			Timestamp: civil.DateTimeOf(time.Now()),
		})

		if err != nil {
			log.Println(err)
		}
		return err
	}

	endEventLog = func(err error) error {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		err = insert.Put(ctx, EventLogModel{
			TrackId:   uuid.New().String(),
			License:   cfg.Lisensi.Email,
			Version:   app.Version,
			Hostname:  hostname,
			Bot:       logname,
			Program:   "golang",
			Status:    FinishEvent,
			Timestamp: civil.DateTimeOf(time.Now()),
			ErrMsg:    errMsg,
		})

		if err != nil {
			log.Println(err)
		}
		return err
	}

	return startEventLog, endEventLog, nil
}
