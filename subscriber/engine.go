package subscriber

import (
	"context"
	"log"

	"todo-service/common"
	"todo-service/common/async_job"
	goservice "todo-service/go-sdk"
	"todo-service/pubsub"
)

type subJob struct {
	Title string
	Hdl   func(ctx context.Context, message *pubsub.Message) error
}

type pbEngine struct {
	serviceCtx goservice.ServiceContext
}

func NewEngine(serviceCtx goservice.ServiceContext) *pbEngine {
	return &pbEngine{
		serviceCtx: serviceCtx,
	}
}

func (engine *pbEngine) Start() error {
	engine.startSubTopic(common.TopicUserLikedItem, false,
		IncreaseLikedCountWhenUserLikesItem(engine.serviceCtx),
		PushNotificationWhenUserLikesItem(engine.serviceCtx),
	)

	engine.startSubTopic(common.TopicUserUnlikedItem, true,
		DecreaseLikedCountWhenUserUnlikesItem(engine.serviceCtx),
	)

	return nil
}

type GroupJob interface {
	Run(ctx context.Context) error
}

func (engine *pbEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, jobs ...subJob) error {
	ps := engine.serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
	c, _ := ps.Subscribe(context.Background(), topic)
	for _, job := range jobs {
		log.Println("Setup subscriber for:", job.Title)
	}

	getJobHandler := func(job *subJob, message *pubsub.Message) async_job.JobHandler {
		return func(ctx context.Context) error {
			log.Printf("running job for %v. Value: %v", job.Title, message.Data())
			return job.Hdl(ctx, message)
		}
	}

	go func() {
		for {
			msg := <-c
			jobHdlArr := make([]async_job.Job, len(jobs))

			for i := range jobs {
				jobHdl := getJobHandler(&jobs[i], msg)
				jobHdlArr[i] = async_job.NewJob(jobHdl, async_job.WithName(jobs[i].Title))
			}

			group := async_job.NewGroup(isConcurrent, jobHdlArr...)
			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}
