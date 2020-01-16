package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func MarshalMessageToAny(pb proto.Message) (any *any.Any, err error) {
	any, err = ptypes.MarshalAny(pb)
	if err != nil {
		logger.Errorf("[ptypes] marshalAny error: %s", err.Error())
	}
	return
}

func UnMarshalAnyToMessage(any *any.Any, pb proto.Message) (err error) {
	if err = ptypes.UnmarshalAny(any, pb); err != nil {
		logger.Errorf("[ptypes] unmarshalAny error: %s", err.Error())
	}
	return
}
