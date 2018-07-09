package BLC

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Rwq_TXOutputs struct {
	Rwq_Outputs []Rwq_TXOutput
}

//  序列化 TXOutputs
func (outs Rwq_TXOutputs) Rwq_Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 反序列化 TXOutputs
func Rwq_DeserializeOutputs(data []byte) Rwq_TXOutputs {
	var outputs Rwq_TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}
