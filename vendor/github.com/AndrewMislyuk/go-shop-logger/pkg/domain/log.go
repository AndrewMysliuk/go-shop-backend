package domain

import (
	"errors"
	"time"
)

const (
	ENTITY_USER    = "USER"
	ENTITY_PRODUCT = "PRODUCT"
	ENTITY_FILE    = "FILE"

	ACTION_CREATE      = "CREATE"
	ACTION_UPDATE      = "UPDATE"
	ACTION_GET         = "GET"
	ACTION_DELETE      = "DELETE"
	ACTION_REGISTER    = "REGISTER"
	ACTION_LOGIN       = "LOGIN"
	ACTION_UPLOAD_FILE = "UPLOAD_FILE"
)

var (
	entities = map[string]LogRequest_Entities{
		ENTITY_USER:    LogRequest_USER,
		ENTITY_PRODUCT: LogRequest_PRODUCT,
		ENTITY_FILE:    LogRequest_FILE,
	}

	actions = map[string]LogRequest_Actions{
		ACTION_CREATE:      LogRequest_CREATE,
		ACTION_UPDATE:      LogRequest_UPDATE,
		ACTION_GET:         LogRequest_GET,
		ACTION_DELETE:      LogRequest_DELETE,
		ACTION_REGISTER:    LogRequest_REGISTER,
		ACTION_LOGIN:       LogRequest_LOGIN,
		ACTION_UPLOAD_FILE: LogRequest_UPLOAD_FILE,
	}
)

type LogItem struct {
	Entity    string    `bson:"entity"`
	Action    string    `bson:"action"`
	EntityID  string    `bson:"entity_id"`
	Timestamp time.Time `bson:"timestamp"`
}

func ToPbEntity(entity string) (LogRequest_Entities, error) {
	val, ex := entities[entity]
	if !ex {
		return 0, errors.New("wrong entity")
	}

	return val, nil
}

func ToPbAction(action string) (LogRequest_Actions, error) {
	val, ex := actions[action]
	if !ex {
		return 0, errors.New("wrong action")
	}

	return val, nil
}
