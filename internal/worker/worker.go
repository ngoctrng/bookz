package worker

import (
	"github.com/hibiken/asynq"
	"github.com/ngoctrng/bookz/internal/book/repository"
	"github.com/ngoctrng/bookz/internal/book/tasks"
	"github.com/ngoctrng/bookz/internal/book/usecases"
	"github.com/ngoctrng/bookz/internal/events"
	"github.com/ngoctrng/bookz/pkg/config"
	"gorm.io/gorm"
)

type Worker struct {
	s   *asynq.Server
	mux *asynq.ServeMux
}

func New(c *config.Config, db *gorm.DB) *Worker {
	server := makeServer(c)

	bTask := initBookTaskHandlers(db)
	mux := asynq.NewServeMux()
	mux.HandleFunc(events.ProposalAcceptedEvent, bTask.ProposalAcceptedEventHandler)

	return &Worker{
		s:   server,
		mux: mux,
	}
}

func initBookTaskHandlers(db *gorm.DB) *tasks.Handler {
	r := repository.New(db)
	u := usecases.NewService(r)
	h := tasks.NewHandler(u)
	return h
}

func (w *Worker) Run() error {
	return w.s.Run(w.mux)
}

func makeServer(c *config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
		},
		asynq.Config{
			Concurrency: 3, // TODO: make this configurable from env
			Queues: map[string]int{
				"default": 3,
			},
		},
	)
}
