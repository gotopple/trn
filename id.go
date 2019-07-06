package id

import (
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Following in the style of URN / ARN:
//  * https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
//  * https://tools.ietf.org/html/rfc2141
// trn:PARTITION:SERVICE:REGION:ACCOUNT:PREFIX/UUID
const format = `trn:%v:%v:%v:%v:%v/%v`

type TRN string

func NewTRN(partition, service, region, account, prefix string) TRN {
	id, err := uuid.NewRandom()
	if err != nil {
		// random has exhausted entropy?
		panic(err)
	}
	return TRN(fmt.Sprintf(format, partition, service, region, account, prefix, id.String()))
}

func Decode(trn string) (TRN, error) {
	o, err := base32.StdEncoding.DecodeString(trn)
	if err != nil {
		return TRN(``), err
	}
	// TODO: validate actual TRN
	return TRN(o), err
}

func (t TRN) Encode() string {
	return base32.StdEncoding.EncodeToString([]byte(t))
}

func (t TRN) ID() string {
	parts := strings.SplitN(string(t), `:`, 6)
	return parts[0]
}

func (t TRN) Partition() string {
	parts := strings.SplitN(string(t), `:`, 6)
	return parts[1]
}

func (t TRN) Service() string {
	parts := strings.SplitN(string(t), `:`, 6)
	return parts[2]
}

func (t TRN) Region() string {
	parts := strings.SplitN(string(t), `:`, 6)
	return parts[3]
}

func (t TRN) Account() string {
	parts := strings.SplitN(string(t), `:`, 6)
	return parts[4]
}

func (t TRN) Resource() string {
	parts := strings.SplitN(string(t), `:`, 6)
	return parts[5]
}

type ServiceIdentifier int

const (
	Metadata ServiceIdentifier = iota
	Ingress
	Content
	Broadcast
)

var ServiceNames = []string{
	`metadata`,
	`ingress`,
	`content`,
	`broadcast`,
}

func (s ServiceIdentifier) String() string {
	return ServiceNames[s]
}
func ParseServiceIdentifier(i string) (ServiceIdentifier, error) {
	for k, v := range ServiceNames {
		if i == v {
			return ServiceIdentifier(k), nil
		}
	}
	return ServiceIdentifier(-1), fmt.Errorf(`invalid service name`)
}
