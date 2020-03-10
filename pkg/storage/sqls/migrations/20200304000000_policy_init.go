package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up20200304000000_policy_init, Down20200304000000_policy_init)
}

func Up20200304000000_policy_init(tx *sql.Tx) error {
	_, err := tx.Exec(`
    CREATE TABLE policies (
        id BIGINT NOT NULL AUTO_INCREMENT,
        uuid VARCHAR(255) NOT NULL UNIQUE,
        subject VARCHAR(100) NOT NULL,
        tenant VARCHAR(100) NOT NULL,
        resource VARCHAR(100) NOT NULL,
        action VARCHAR(100) NOT NULL,
        effect VARCHAR(5) NOT NULL, 
        active SMALLINT(1) DEFAULT 1, 
        expired TIMESTAMP NULL DEFAULT NULL,
        created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        CONSTRAINT primary_policy PRIMARY KEY (id)
    );
  `)
	if err != nil {
		return err
	}
	return nil
}

func Down20200304000000_policy_init(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS policies;")
	if err != nil {
		return err
	}
	return nil
}
