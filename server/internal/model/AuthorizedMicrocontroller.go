package model

type AuthorizedMicrocontroller struct {
	DeviceId   string `json:"device_id" firestore:"device_id"`
	HMACSecret string `json:"hmac_secret" firestore:"hmac_secret"`
	IsActive   bool   `json:"is_active" firestore:"is_active"`
}
