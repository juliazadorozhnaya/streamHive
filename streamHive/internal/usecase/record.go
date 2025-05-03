package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"streamHive/streamHive/internal/domain"
)

type useCase struct {
	recorder domain.Recorder
	storage  domain.Storage
	pub      domain.EventPublisher
	state    domain.StateStore
}

func New(rec domain.Recorder, store domain.Storage, pub domain.EventPublisher, state domain.StateStore) domain.UseCase {
	return &useCase{
		recorder: rec,
		storage:  store,
		pub:      pub,
		state:    state,
	}
}

// StartRecording запускает задачу записи потока по RTSP
// - Генерирует jobID
// - Сохраняет статус "started" ъ
// - Асинхронно запускает ffmpeg
// - После завершения обновляет статус и публикует событие
func (u *useCase) StartRecording(ctx context.Context, streamURL string) (string, error) {
	jobID := uuid.New().String()
	output := fmt.Sprintf("/tmp/%s.mp4", jobID)

	// сохраняем в Mongo, что задача началась
	err := u.state.SaveJob(jobID, streamURL, "started")
	if err != nil {
		return "", err
	}

	// запускаем запись в фоне (не блокирует основной поток)
	go func() {
		err := u.recorder.Record(ctx, streamURL, output)
		if err == nil {
			u.state.UpdateJobStatus(jobID, "finished")
		} else {
			u.state.UpdateJobStatus(jobID, "failed")
		}

		// публикуем событие завершения задачи
		_ = u.pub.Publish("recording.finished", []byte(fmt.Sprintf(`"%s"`, jobID)))
	}()

	return jobID, nil
}

// StopRecording завершает активную запись по jobID
// - вызывает .Stop() у Recorder
// - обновляет статус задачи в Mongo
func (u *useCase) StopRecording(jobID string) error {
	err := u.recorder.Stop(jobID)
	if err != nil {
		return err
	}
	return u.state.UpdateJobStatus(jobID, "stopped")
}

// GetJobStatus возвращает текущий статус задачи из хранилища
func (u *useCase) GetJobStatus(jobID string) (string, error) {
	return u.state.GetJobStatus(jobID)
}

// StartBackgroundConsumer запускает обработку входящих команд из NATS
// - command.start → вызывает StartRecording
// - command.stop  → вызывает StopRecording
func (u *useCase) StartBackgroundConsumer(ctx context.Context) {
	u.pub.Subscribe("command.start", func(data []byte) {
		var req struct {
			StreamURL string `json:"stream_url"`
		}
		_ = json.Unmarshal(data, &req)
		_, _ = u.StartRecording(ctx, req.StreamURL)
	})

	u.pub.Subscribe("command.stop", func(data []byte) {
		var req struct {
			JobID string `json:"job_id"`
		}
		_ = json.Unmarshal(data, &req)
		_ = u.StopRecording(req.JobID)
	})
}
