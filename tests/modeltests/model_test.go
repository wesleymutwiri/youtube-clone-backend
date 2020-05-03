package modeltests

import (
	"api/youtube-clone-backend/api/controllers"
	"api/youtube-clone-backend/api/models"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	TestDbDriver := os.Getenv("TestDbDriver")
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbUser"), os.Getenv("TestDbName"), os.Getenv("TestDbPassword"))
	server.DB, err = gorm.Open(TestDbDriver, DBURL)

	if err != nil {
		fmt.Printf("Cannot connect to %s database \n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {
	refreshUserTable()
	user := models.User{
		Username: "Pet",
		Email:    "pet@email.com",
		Password: "password",
	}
	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("Cannot seed users table: %v", err)
	}
	return user, nil
}

func seedUsers() error {
	users := []models.User{
		models.User{
			Username: "Steven something",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Username: "Chuck Norris",
			Email:    "chuck@gmail.com",
			Password: "password",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}
