package id

import (
	"fmt"
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
