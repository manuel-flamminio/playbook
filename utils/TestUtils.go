package utils

import (
	"math/rand"
	"playbook/entities"
	"playbook/filters"

	"github.com/google/uuid"
)

func GetRandomUsername() string {
	return RandSeq(10) + "@example.com"
}

func GetRandomString(size int) string {
	return RandSeq(size)
}

func GetRandomBool() bool {
	return rand.Int()%2 == 0
}

func GetRandomSortingType() filters.SortingType {
	n := rand.Intn(4)
	switch n {
	case 0:
		return filters.BestOfAllTime
	case 1:
		return filters.New
	case 2:
		return filters.Trending
	default:
		return filters.Random
	}
}

func GetRandomVote() entities.Vote {
	n := rand.Intn(3)
	switch n {
	case 0:
		return entities.Upvote
	case 1:
		return entities.Downvote
	default:
		return entities.None
	}
}

func GetRandomUUID() (uuid.UUID, error) {
	generatedUUID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	return generatedUUID, nil
}

func GetRandomTagName() string {
	return RandSeq(10)
}

func GetRandomTagDescription() string {
	return RandSeq(30)
}

func GetRandomValidUser() *entities.User {
	return &entities.User{
		Username:    GetRandomUsername(),
		DisplayName: GetRandomString(10),
		Password:    GetRandomString(18),
		ID:          uuid.Max,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
