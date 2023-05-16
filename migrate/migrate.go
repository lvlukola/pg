package migrate

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"os"
	"time"
)

func Create(path, name string) (string, error) {
	log.Println("Migration create run...")

	now := time.Now()
	sec := now.Unix()

	fileNameUp := fmt.Sprintf("%s/%d_%s.up.sql", path, sec, name)
	if err := createFile(fileNameUp); err != nil {
		log.Println(err)
		return "", err
	}

	fileNameDown := fmt.Sprintf("%s/%d_%s.down.sql", path, sec, name)
	if err := createFile(fileNameDown); err != nil {
		log.Println(err)
		return "", err
	}

	message := fmt.Sprintf("Files created: %s, %s", fileNameUp, fileNameDown)
	log.Println(message)

	return message, nil
}

func createFile(fileName string) error {
	f, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}

func RunMigrateUp(source, dsn string) (string, error) {
	log.Println("Migration up run...")

	m, err := migrate.New(source, dsn+"?sslmode=disable")

	if err != nil {
		log.Println(err)
		return "", nil
	}

	if err := m.Up(); err == migrate.ErrNoChange {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("No new migration. Current version: %d", ver)
		log.Println(message)
		return message, nil
	} else if err != nil {
		log.Println(err)
		return "", err
	} else {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("Migrate up complete. Current version %d ", ver)
		log.Println(message)
		return message, nil
	}
}

func RunMigrateDown(source, dsn string) (string, error) {
	log.Println("Migration down run...")

	m, err := migrate.New(source, dsn+"?sslmode=disable")

	if err != nil {
		return "", err
	}

	if err := m.Steps(-1); err == migrate.ErrNoChange {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("No down migration. Current version: %d", ver)
		log.Println(message)
		return message, nil
	} else if err != nil {
		return "", err
	} else {
		ver, _, _ := m.Version()
		message := fmt.Sprintf("Migrate down complete. Current version %d ", ver)
		log.Println(message)
		return message, nil
	}

}
