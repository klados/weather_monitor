package repository

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/klados/weather_monitor/internal/model"
)

const authorizedMicrocontrollersCollectionName = "authorized_microcontrollers"

type AuthorizedMicrocontrollers struct {
	DB *firestore.Client
}

func (am *AuthorizedMicrocontrollers) GetAuthorizedMicrocontrollerByDeviceId(deviceId string) (model.AuthorizedMicrocontroller, error) {
	var microcontroller model.AuthorizedMicrocontroller

	doc, err := am.DB.Collection(authorizedMicrocontrollersCollectionName).
		Doc(deviceId).
		Get(context.Background())

	if err != nil {
		return microcontroller, fmt.Errorf("error retrieving document: %w", err)
	}

	if !doc.Exists() {
		return microcontroller, fmt.Errorf("device not found")
	}

	if err := doc.DataTo(&microcontroller); err != nil {
		return microcontroller, fmt.Errorf("error unmarshalling document: %w", err)
	}

	return microcontroller, nil
}
