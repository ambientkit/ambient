package websession

import (
	"encoding/json"
	"time"

	"github.com/josephspurrier/ambient/lib/envdetect"
)

// SessionDatabase -
type SessionDatabase struct {
	Records map[string]SessionData `json:"db"`
}

// SessionData -
type SessionData struct {
	ID     string    `json:"id"`
	Data   []byte    `json:"data"`
	Expire time.Time `json:"expire"`
}

// Load -
func (sd *SessionDatabase) Load(ss Sessionstorer, en Encrypter) error {
	b, err := ss.Load()
	if err != nil {
		return err
	}

	b, err = en.Decrypt(b)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, sd)
	if err != nil {
		return err
	}

	if sd.Records == nil {
		sd.Records = make(map[string]SessionData)
	}

	return nil
}

// Save -
func (sd *SessionDatabase) Save(ss Sessionstorer, en Encrypter) error {
	var b []byte
	var err error

	if envdetect.RunningLocalDev() {
		// Indent so the data is easy to read.
		b, err = json.MarshalIndent(sd, "", "    ")
	} else {
		b, err = json.Marshal(sd)
	}

	if err != nil {
		return err
	}

	b, err = en.Encrypt(b)
	if err != nil {
		return err
	}

	err = ss.Save(b)
	if err != nil {
		return err
	}

	return nil
}
