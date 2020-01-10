package utils

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

func MarshalMessageToAny(pb proto.Message) (any *any.Any, err error) {
	any, err = ptypes.MarshalAny(pb)
	if err != nil {
		logger.Errorf("[ptypes] marshalAny error: %v", err)
	}
	return
}
