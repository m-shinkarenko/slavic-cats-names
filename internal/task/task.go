package task

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const dateLayout = "20060102150405"

type Task struct {
	difficulty int
	date       time.Time
	rand       []byte

	target *big.Int
}

func NewTask(difficulty int) (*Task, error) {
	randVal := make([]byte, 16)
	_, err := rand.Read(randVal)
	if err != nil {
		return nil, errors.Join(err, errors.New("rand value generation error"))
	}

	return &Task{
		difficulty: difficulty,
		date:       time.Now().Truncate(time.Second).UTC(),
		rand:       randVal,
		target:     new(big.Int).Lsh(big.NewInt(1), uint(sha1.Size*8-difficulty)),
	}, nil
}

func ParseTask(task []byte) (*Task, error) {
	parts := strings.Split(string(task), ":")
	if len(parts) != 3 {
		return nil, errors.New("incorrect task format")
	}

	dif, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.Join(err, errors.New("incorrect difficulty format"))
	}

	date, err := time.Parse(dateLayout, parts[1])
	if err != nil {
		return nil, errors.Join(err, errors.New("incorrect date format"))
	}

	rand, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.Join(err, errors.New("incorrect rand format"))
	}

	return &Task{
		difficulty: dif,
		date:       date,
		rand:       rand,
		target:     new(big.Int).Lsh(big.NewInt(1), uint(sha1.Size*8-dif)),
	}, nil
}

func (t *Task) ToBytes() []byte {
	return []byte(
		fmt.Sprintf("%d:%s:%s",
			t.difficulty,
			t.date.Format(dateLayout),
			base64.StdEncoding.EncodeToString(t.rand)),
	)
}

func (t *Task) Solve() int64 {
	for nonce := int64(0); nonce <= math.MaxInt64; nonce++ {
		if t.Validate(nonce) {
			return nonce
		}
	}

	return 0
}

func (t *Task) Validate(nonce int64) bool {
	var intHash big.Int

	hash := sha1.Sum(t.initNonce(nonce))
	intHash.SetBytes(hash[:])

	return intHash.Cmp(t.target) == -1
}

func (t *Task) initNonce(nonce int64) []byte {
	return []byte(
		fmt.Sprintf("%d:%s:%s:%d",
			t.difficulty,
			t.date.Format(dateLayout),
			base64.StdEncoding.EncodeToString(t.rand),
			nonce),
	)
}
