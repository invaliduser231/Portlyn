package acme

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"portlyn/internal/domain"
	"portlyn/internal/store"
)

type CertificateRenewer interface {
	SyncCertificate(ctx context.Context, item *domain.Certificate) (*domain.Certificate, error)
}

type AcmeWorker struct {
	store       CertificateStore
	metaStore   *store.CertificateStore
	renewer     CertificateRenewer
	interval    time.Duration
	renewWithin time.Duration
	logger      *slog.Logger
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

func NewAcmeWorker(store CertificateStore, metaStore *store.CertificateStore, renewer CertificateRenewer, interval, renewWithin time.Duration, logger *slog.Logger) *AcmeWorker {
	if interval <= 0 {
		interval = time.Minute
	}
	if renewWithin <= 0 {
		renewWithin = 30 * 24 * time.Hour
	}
	return &AcmeWorker{
		store:       store,
		metaStore:   metaStore,
		renewer:     renewer,
		interval:    interval,
		renewWithin: renewWithin,
		logger:      logger,
	}
}

func (w *AcmeWorker) Start(ctx context.Context) {
	if w.metaStore == nil || w.renewer == nil {
		return
	}
	if w.cancel != nil {
		return
	}
	runCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel
	w.wg.Add(1)
	go w.loop(runCtx)
}

func (w *AcmeWorker) Stop() {
	if w.cancel != nil {
		w.cancel()
		w.wg.Wait()
	}
}

func (w *AcmeWorker) loop(ctx context.Context) {
	defer w.wg.Done()
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		w.runOnce(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (w *AcmeWorker) runOnce(ctx context.Context) {
	pending, err := w.metaStore.ListPending(ctx, 32)
	if err != nil {
		if w.logger != nil {
			w.logger.Error("list pending certificates failed", "error", err)
		}
		return
	}
	for i := range pending {
		if _, err := w.renewer.SyncCertificate(ctx, &pending[i]); err != nil && w.logger != nil {
			w.logger.Error("pending certificate sync failed", "certificate_id", pending[i].ID, "error", err)
		}
	}

	expiring, err := w.metaStore.ListExpiringBefore(ctx, time.Now().UTC().Add(w.renewWithin))
	if err != nil {
		if w.logger != nil {
			w.logger.Error("list expiring certificates failed", "error", err)
		}
		return
	}
	for i := range expiring {
		if !expiring[i].IsAutoRenew {
			continue
		}
		expiring[i].Status = domain.CertificateStatusRenewing
		_ = w.metaStore.Update(ctx, &expiring[i])
		if _, err := w.renewer.SyncCertificate(ctx, &expiring[i]); err != nil && w.logger != nil {
			w.logger.Error("expiring certificate sync failed", "certificate_id", expiring[i].ID, "error", err)
		}
	}
}
