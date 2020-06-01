package trn

import (
	"database/sql/driver"
	"encoding/base32"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Following in the style of URN / ARN:
//  * https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
//  * https://tools.ietf.org/html/rfc2141
//  * trn:PARTITION:SERVICE:REGION:ACCOUNT:PREFIX/UUID
// Example:
//  * `trn:topple:content:sfo2:12341234:content/3e84977e-5e9a-4494-97a3-3ca15b427569`
const Format = `trn:%v:%v:%v:%v:%v/%v`

type TRN string

func NewTRN(partition, service, region, account, prefix string) TRN {
	// TODO: validate that none of the input contain colons
	id, err := uuid.NewRandom()
	if err != nil {
		// random has exhausted entropy?
		panic(err)
	}
	return TRN(fmt.Sprintf(Format, partition, service, region, account, prefix, id.String()))
}

const charset = `0123456789`

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func slowRand(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func IsValid(t TRN) bool {
	parts := strings.SplitN(string(t), `:`, 6)
	return len(parts) == 6 && parts[0] == `trn`
}

func NewSlowTRN(partition, service, region, account, prefix string) TRN {
	// TODO: validate that none of the input contain colons
	id := slowRand(10)
	return TRN(fmt.Sprintf(Format, partition, service, region, account, prefix, id))
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

func (t TRN) Components() (id, partition, service, region, account, resource string) {
	parts := strings.SplitN(string(t), `:`, 6)
	if len(parts) != 6 {
		panic(`invalid TRN format`)
	}
	return parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]

}

func (t TRN) ID() string {
	i, _, _, _, _, _ := t.Components()
	return i
}

func (t TRN) Partition() string {
	_, p, _, _, _, _ := t.Components()
	return p
}

func (t TRN) Service() string {
	_, _, s, _, _, _ := t.Components()
	return s
}

func (t TRN) Region() string {
	_, _, _, r, _, _ := t.Components()
	return r
}

func (t TRN) Account() string {
	_, _, _, _, a, _ := t.Components()
	return a
}

func (t TRN) Resource() string {
	_, _, _, _, _, r := t.Components()
	return r
}

var (
	ErrNotStringType = fmt.Errorf(`trn can only be decoded from a string type`)
)

func (t *TRN) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return ErrNotStringType
	}
	var err error
	*t, err = Decode(s)
	return err
}

func (t TRN) Value() (driver.Value, error) {
	return t.Encode(), nil
}

type ServiceIdentifier int

const (
	Metadata ServiceIdentifier = iota
	Ingress
	Content
	Broadcast
	Account
	Workspace
)

var serviceNames = []string{
	`metadata`,
	`ingress`,
	`content`,
	`broadcast`,
	`account`,
	`workspace`,
}

func (s ServiceIdentifier) String() string {
	return serviceNames[s]
}
func ParseServiceIdentifier(i string) (ServiceIdentifier, error) {
	for k, v := range serviceNames {
		if i == v {
			return ServiceIdentifier(k), nil
		}
	}
	return ServiceIdentifier(-1), fmt.Errorf(`invalid service name`)
}
