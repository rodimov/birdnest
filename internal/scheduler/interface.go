package scheduler

type UseCase interface {
	StartScheduler() error
}
