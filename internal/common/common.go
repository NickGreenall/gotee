//common provides the encoder and decoder interfaces.
package common

//Encoder is an abstractation for JSON like encoder.
type Encoder interface {
	Encode(interface{}) error
}

//Decoder is an abstractation for JSON like decoder.
type Decoder interface {
	Decode(interface{}) error
}
