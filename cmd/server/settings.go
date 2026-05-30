package main

import (
	"context"
	"fmt"

	"portlyn/internal/config"
	"portlyn/internal/store"
)

func runSettingsSubcommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: portlyn settings sync")
	}
	switch args[0] {
	case "sync":
		return runSettingsSync()
	default:
		return fmt.Errorf("unknown subcommand %q (expected: sync)", args[0])
	}
}

func runSettingsSync() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	db, err := store.NewDatabase(cfg)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	appSettingsStore := store.NewAppSettingsStore(db)
	appSettingsStore.SetDataEncryptionSecrets(cfg.DataEncryptionSecrets())
	if err := appSettingsStore.SeedDefaults(context.Background(), cfg); err != nil {
		return fmt.Errorf("seed defaults: %w", err)
	}
	drifts, err := appSettingsStore.SyncFromEnv(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("sync from env: %w", err)
	}
	if len(drifts) == 0 {
		fmt.Println("settings sync: no drift, nothing to do")
		return nil
	}
	fmt.Printf("settings sync: applied %d change(s) from environment to database:\n", len(drifts))
	for _, d := range drifts {
		fmt.Printf("  %s: %s -> %s\n", d.Field, d.DBValue, d.EnvValue)
	}
	return nil
}
