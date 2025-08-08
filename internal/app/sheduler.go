package app

import (
	"testYTask/internal/common"
	"testYTask/internal/domain/interfaces"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

// InitScheduler создаёт и инициализирует планировщик задач с заданными джобами.
//
// Параметры:
//   - uploadJob: задача, которая будет запущена по расписанию
//
// Возвращает:
//   - объект планировщика задач
//   - ошибку, если произошла ошибка при создании или инициализации планировщика
func InitScheduler(uploadJob interfaces.JobI) (gocron.Scheduler, error) {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		zap.L().Error(
			"Error loading timezone 'Asia/Almaty'",
			zap.Error(err),
		)
		return nil, err
	}

	scheduler, nErr := gocron.NewScheduler(
		gocron.WithLocation(location),
		gocron.WithStopTimeout(common.SchedulerStopTimeout),
	)
	if nErr != nil {
		zap.L().Error(
			"Error creating job scheduler",
			zap.Error(nErr),
		)
		return nil, nErr
	}

	if err = initJobs(scheduler, uploadJob); err != nil {
		zap.L().Error(
			"Error initializing scheduler jobs",
			zap.Error(err),
		)
		return nil, err
	}

	zap.L().Info("Successfully initialized scheduler")
	return scheduler, nil
}

// initJobs добавляет джобы в планировщик.
//
// Параметры:
//   - scheduler: объект планировщика задач
//   - uploadJob: задача, которая будет запущена по расписанию
//
// Возвращает ошибку, если добавление задачи завершилось неудачно.
func initJobs(scheduler gocron.Scheduler, uploadJob interfaces.JobI) error {
	if _, err := scheduler.NewJob(
		//gocron.DailyJob(
		//	1, // сколько раз в день запускать.
		//	gocron.NewAtTimes(
		//		gocron.NewAtTime(3, 0, 0),
		//	),
		//),
		gocron.DurationJob(30*time.Second),
		// for test use only
		//gocron.OneTimeJob(gocron.OneTimeJobStartImmediately()),
		gocron.NewTask(
			uploadJob.Run,
		),
	); err != nil {
		zap.L().Error(
			"Error initializing uploadJob job",
			zap.Error(err),
		)
		return err
	} else {
		zap.L().Info("Successfully initialized uploadJob job")
	}

	zap.L().Info("Successfully initialized all jobs")
	return nil
}
