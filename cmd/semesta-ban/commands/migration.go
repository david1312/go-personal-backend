package commands

import (
	"libra-internal/bootstrap"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql" // MySQL driver
	_ "github.com/golang-migrate/migrate/v4/source/file"    // File source driver
	"github.com/spf13/cobra"
)

func init() {
	registerCommand(migrateUpCommand)
	registerCommand(migrateDownCommand)
}

func migrateUpCommand(dep *bootstrap.Dependency) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-up",
		Short: "Apply all up migrations",
		Long:  `This command applies all the database migrations in the "up" direction.`,
		Run: func(cmd *cobra.Command, args []string) {
			// cfg := dep.GetConfig()
			m, err := migrate.New(
				"file://files/db_migration/", // Replace with the correct path to your migration files
				"mysql://mysql_user:mysql_password@tcp(localhost:3306)/sunmoris_customer",            // Replace with your database URL
			)
			if err != nil {
				log.Fatalf("failed to create migrate instance: %v", err)
			}

			if err := m.Up(); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("no change in migrations")
				} else {
					log.Fatalf("failed to apply up migrations: %v", err)
				}
			} else {
				log.Println("successfully applied all up migrations")
			}
		},
	}
}

func migrateDownCommand(dep *bootstrap.Dependency) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate-down",
		Short: "Revert the last migration",
		Long:  `This command reverts the last applied database migration.`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := dep.GetConfig()
			m, err := migrate.New(
				"file://path/to/migrations", // Replace with the correct path to your migration files
				cfg.Database.Write,            // Replace with your database URL
			)
			if err != nil {
				log.Fatalf("failed to create migrate instance: %v", err)
			}

			if err := m.Down(); err != nil {
				if err == migrate.ErrNoChange {
					log.Println("no migrations to revert")
				} else {
					log.Fatalf("failed to apply down migration: %v", err)
				}
			} else {
				log.Println("successfully reverted the last migration")
			}
		},
	}
}