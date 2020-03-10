package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20200304000001_role_init, Down20200304000001_role_init)
}

func Up20200304000001_role_init(tx *sql.Tx) error {
	_, err := tx.Exec(`
    CREATE TABLE roles (
				id BIGINT NOT NULL AUTO_INCREMENT,
				policy VARCHAR(100) NOT NULL,
				tenant VARCHAR(100) NOT NULL,
        subject VARCHAR(100) NOT NULL,
        created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT primary_role PRIMARY KEY (id)
    );
  `)
	if err != nil {
		return err
	}
	return nil
}

func Down20200304000001_role_init(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS roles;")
	if err != nil {
		return err
	}
	return nil
}
