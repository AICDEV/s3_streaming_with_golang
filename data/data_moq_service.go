package data

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	pb "github.com/aicdev/s3_streaming_with_golang/proto"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

type TestDataServiceInterface interface {
	CreateTestData(int)
	GetTestData() []*pb.User
}

type testDataService struct {
	testData []*pb.User
}

func NewTestDataService() TestDataServiceInterface {
	rand.Seed(time.Now().Unix())
	return &testDataService{}
}

func (tsd *testDataService) CreateTestData(limit int) {
	for i := 0; i < limit; i++ {
		email := fmt.Sprintf("%s@%s.%s", getRandomStringBytes(8), getRandomStringBytes(4), getRandomStringBytes(2))
		testUser := &pb.User{
			Email: email,
			Id:    GetHash(email),
		}

		tsd.createTestTransactions(testUser)
		tsd.testData = append(tsd.testData, testUser)
	}
}

func (tsd *testDataService) GetTestData() []*pb.User {
	return tsd.testData
}

func (tsd *testDataService) createTestTransactions(testUser *pb.User) {
	for i := 0; i < 5000; i++ {
		testUser.Transactions = append(testUser.Transactions, &pb.Transaction{
			Anmount:         rand.Float64(),
			Recipient:       fmt.Sprintf("%s@%s.%s", getRandomStringBytes(8), getRandomStringBytes(4), getRandomStringBytes(2)),
			Reason:          getRandomStringBytes(20),
			VerificationKey: fmt.Sprintf("%s-%s", getRandomStringBytes(4), getRandomStringBytes(8)),
			UserId:          testUser.GetId(),
			Currency:        pb.Transaction_USD,
		})
	}
}

// inspired from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func getRandomStringBytes(limit int) string {
	b := make([]byte, limit)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetHash(t string) string {
	sum := sha512.Sum512([]byte(t))
	return hex.EncodeToString(sum[:])
}
