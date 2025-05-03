package domain

import "context"

// Recorder отвечает за захват видео потока и остановку записи по jobID
type Recorder interface {
	Record(ctx context.Context, streamURL, outputPath string) error
	Stop(jobID string) error
}

// Storage сохраняет видеофайл в локальное хранилище и запускает фоновую очистку старых файлов
type Storage interface {
	Save(jobID string, file []byte) error
	StartCleanupWorker()
}

// EventPublisher для взаимодействие с NATS
type EventPublisher interface {
	Publish(subject string, payload []byte) error
	Subscribe(subject string, handler func(data []byte)) error
}

// StateStore для MongoDB
type StateStore interface {
	SaveJob(jobID, streamURL, status string) error
	UpdateJobStatus(jobID, status string) error
	GetJobStatus(jobID string) (string, error)
}

// UseCase для бизнес-логики по управлению задачами
type UseCase interface {
	StartRecording(ctx context.Context, streamURL string) (string, error)
	StopRecording(jobID string) error
	GetJobStatus(jobID string) (string, error)
	StartBackgroundConsumer(ctx context.Context)
}
