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

	docs, err := am.DB.Collection(authorizedMicrocontrollersCollectionName).
		Where("device_id", "==", deviceId).
		Limit(1).
		Documents(context.Background()).
		GetAll()

	if err != nil {
		return microcontroller, fmt.Errorf("error querying documents: %w", err)
	}

	if len(docs) == 0 {
		return microcontroller, fmt.Errorf("device not found")
	}

	if err := docs[0].DataTo(&microcontroller); err != nil {
		return microcontroller, fmt.Errorf("error unmarshalling document: %w", err)
	}

	return microcontroller, nil
}
