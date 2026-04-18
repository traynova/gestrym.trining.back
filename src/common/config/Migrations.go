package config

import (
	"fmt"
	"gestrym-training/src/common/models"
	"gestrym-training/src/common/utils"
)

var logger = utils.NewLogger()

func MigrateDB() (IDatabaseConnection, error) {
	connection := NewPostgresConnection()
	db := connection.GetDB()

	//Se agregan los modelos de base de datos
	err := db.AutoMigrate(
		&models.Exercise{},
		&models.Food{},
		&models.FoodCategory{},
		&models.Workout{},
		&models.WorkoutExercise{},
		&models.WorkoutSet{},
	)

	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] Error al migrar las entidades: %s", err.Error()))
		return nil, err
	}

	logger.Info("[OK] Todas las migraciones completadas exitosamente")
	return connection, nil
}
