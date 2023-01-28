package models

// CreationResp creation response struct
type CreationResp struct {
	InsertedID string `json:"inserted_id"`
}

// MultiCreationResp multi creation response struct
type MultiCreationResp struct {
	InsertedIDs []string `json:"inserted_ids"`
}
