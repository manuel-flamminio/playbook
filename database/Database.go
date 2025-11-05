package database

import (
	"log"
	"math"
	"os"
	"playbook/constants"
	"playbook/entities"
	"playbook/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

type StatisticsCountResult struct {
	Vote  entities.Vote
	Total int
}

func MigrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&entities.User{})
	db.AutoMigrate(&entities.Tag{})
	db.AutoMigrate(&entities.PickupLine{})
	db.AutoMigrate(&entities.Reaction{})
}

func GetDatabaseConnection() *Database {
	db := ConnectToDatabase()
	return &Database{
		db: db,
	}
}

func GetTestDatabaseConnection(db *gorm.DB) *Database {
	return &Database{
		db: db,
	}
}

func InitializeAndGetTestDatabase() (*gorm.DB, func()) {
	db, cleanup := ConnectToTestDatabase(utils.RandSeq(8))
	MigrateDatabase(db)
	return db, cleanup
}

func InitializeAndGetTestDatabaseTransaction() (*gorm.DB, func()) {
	db, dbCleanup := InitializeAndGetTestDatabase()
	tx := db.Begin()
	transactionCleanUp := func() {
		tx.Rollback()
		dbCleanup()
	}
	return tx, transactionCleanUp
}

func ConnectToTestDatabase(dbName string) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(dbName+".db"), &gorm.Config{})
	if err != nil {
		panic("Could not connect to database")
	}

	cleanup := func() {
		os.Remove(dbName + ".db")
	}

	return db, cleanup
}

func getGormLogger() logger.Interface {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,         // Don't include params in the SQL log
			Colorful:                  false,        // Disable color
		},
	)
	return newLogger
}

func ConnectToDatabase() *gorm.DB {
	dsn := os.Getenv(constants.DATABASE_DSN)
	if dsn == "" {
		dsn = "root:testeroni@tcp(127.0.0.1:3306)/playbook?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: getGormLogger(),
	})
	if err != nil {
		panic("Could not connect to database")
	}

	return db
}

func (d *Database) CreateUser(user *entities.User) error {
	result := d.db.Create(user)
	return result.Error
}

func (d *Database) UpdateUser(user *entities.User) error {
	result := d.db.Save(user)
	return result.Error
}

func (d *Database) GetUserByUsername(username string) (*entities.User, error) {
	user := &entities.User{}
	result := d.db.Where("username = ?", username).First(user)
	return user, result.Error
}

func (d *Database) GetUserById(userId uuid.UUID) (*entities.User, error) {
	user := &entities.User{}
	result := d.db.First(&user, "id = ?", userId)
	return user, result.Error
}

func (d *Database) GetUserByUserItem(user *entities.User) (*entities.User, error) {
	result := d.db.Where(user).First(user)
	return user, result.Error
}

func (d *Database) GetTagById(tagId string) (*entities.Tag, error) {
	tag := &entities.Tag{}
	result := d.db.First(&tag, "id = ?", tagId)
	return tag, result.Error
}

func (d *Database) GetTagsByUserId(userId uuid.UUID) ([]*entities.Tag, error) {
	var tags []*entities.Tag
	result := d.db.Find(&tags, "user_id = ?", userId)
	return tags, result.Error
}

func (d *Database) GetListByTagIdAndUserIdList(tagsToFind [][]uuid.UUID) ([]*entities.Tag, error) {
	var tags []*entities.Tag
	result := d.db.Find(&tags, "(id, user_id) in ?", tagsToFind)
	return tags, result.Error
}

func (d *Database) CreateTag(tag *entities.Tag) error {
	result := d.db.Create(tag)
	return result.Error
}

func (d *Database) UpdateTag(tag *entities.Tag) error {
	result := d.db.Save(tag)
	return result.Error
}

func (d *Database) DeleteTag(tag *entities.Tag) error {
	result := d.db.Where("user_id = ?", tag.UserId).Delete(tag)
	return result.Error
}

func (d *Database) CreatePickupLine(pickupLine *entities.PickupLine) error {
	result := d.db.Create(pickupLine)
	return result.Error
}

func (d *Database) CreateReaction(reaction *entities.Reaction) error {
	result := d.db.Create(reaction)
	return result.Error
}

func (d *Database) DeleteReaction(reaction *entities.Reaction) error {
	result := d.db.Delete(reaction)
	return result.Error
}

func (d *Database) UpdateReaction(reaction *entities.Reaction) error {
	return d.db.Save(reaction).Error
}

func (d *Database) UpdatePickupLine(pickupLine *entities.PickupLine) error {
	tx := d.db.Begin()
	err := updatePickupLineTransaction(pickupLine, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func updatePickupLineTransaction(pickupLine *entities.PickupLine, tx *gorm.DB) error {
	err := tx.Model(pickupLine).Association("Tags").Replace(pickupLine.Tags)
	if err != nil {
		return err
	}

	err = tx.Save(pickupLine).Error
	return err
}

func (d *Database) DeletePickupLine(pickupLine *entities.PickupLine) error {
	result := d.db.Where("user_id = ?", pickupLine.User.ID).Delete(pickupLine)
	return result.Error
}

func (d *Database) GetPickupLineByIdAndUserId(pickupLineId uuid.UUID, userId uuid.UUID) (*entities.PickupLine, error) {
	pickupLine := &entities.PickupLine{}
	result := d.db.Preload("User").Preload("Tags").Preload("Reactions", "user_id = ?", userId).First(&pickupLine, "id = ? and user_id = ?", pickupLineId, userId)
	return pickupLine, result.Error
}

func (d *Database) GetPickupLineForVisibilityCheck(pickupLineId uuid.UUID) (*entities.PickupLine, error) {
	pickupLine := &entities.PickupLine{}
	result := d.db.Select("Visible", "UserID").First(&pickupLine, "id = ?", pickupLineId)
	return pickupLine, result.Error
}

// TODO: statistic cache directly in the database or method extraction
func (d *Database) GetStatisticsById(pickupLineId uuid.UUID) (*entities.Statistic, error) {
	statistics := &entities.Statistic{}
	var countResult []StatisticsCountResult
	err := d.db.Raw(`
		select r.vote, count(*) as total
		from playbook.reactions r
		where r.pickup_line_id = ?
		group by r.vote
	`, pickupLineId).Scan(&countResult).Error
	if err != nil {
		return nil, err
	}
	total := 0
	for _, result := range countResult {
		switch result.Vote {
		case entities.None:
			continue
		case entities.Upvote:
			statistics.NumberOfSuccesses = result.Total
			total += result.Total
		case entities.Downvote:
			statistics.NumberOfFailures = result.Total
			total += result.Total
		}
	}
	statistics.NumberOfTries = total
	if total > 0 {
		statistics.SuccessPercentage = math.Round(float64(statistics.NumberOfSuccesses)/float64(statistics.NumberOfTries)*100) / 100
	}

	return statistics, nil
}

func (d *Database) GetPickupLineById(pickupLineId uuid.UUID) (*entities.PickupLine, error) {
	pickupLine := &entities.PickupLine{}
	result := d.db.First(&pickupLine, "id = ?", pickupLineId)
	return pickupLine, result.Error
}

func (d *Database) GetReactionByPickupLineIdAndUserId(pickupLineId uuid.UUID, userId uuid.UUID) (*entities.Reaction, error) {
	reaction := &entities.Reaction{}
	result := d.db.First(&reaction, "pickup_line_id = ? and user_id = ?", pickupLineId, userId)
	return reaction, result.Error
}

func (d *Database) DeleteUser(user *entities.User) error {
	result := d.db.Delete(user)
	return result.Error
}
