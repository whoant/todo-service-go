package cmd

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"todo-service/common"
	"todo-service/demogrpc/demo"
	goservice "todo-service/go-sdk"
	"todo-service/memcache"
	"todo-service/middleware"
	ginitem "todo-service/module/item/transport/gin"
	"todo-service/module/upload"
	userStorage "todo-service/module/user/storage"
	ginuser "todo-service/module/user/transport/gin"
	userLikeStorage "todo-service/module/user_like_item/storage"
	ginuserlike "todo-service/module/user_like_item/transport/gin"
	"todo-service/module/user_like_item/transport/rpc"
	"todo-service/plugin/app_redis"
	"todo-service/plugin/nats"
	"todo-service/plugin/rpc_caller"
	"todo-service/plugin/sdk_gorm"
	"todo-service/plugin/simple"
	"todo-service/plugin/token_provider/jwt"
	"todo-service/subscriber"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("todo-service"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdk_gorm.NewGormDB("main.mysql", common.PluginDBMain)),
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple", "simple")),
		goservice.WithInitRunnable(jwt.NewJWTProvider(common.PluginJWT)),
		//goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
		goservice.WithInitRunnable(nats.NewNATS(common.PluginPubSub)),
		goservice.WithInitRunnable(rpc_caller.NewApiItemCaller(common.PluginItemAPI)),
		goservice.WithInitRunnable(app_redis.NewRedisDB("redis", common.PluginRedis)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start social TODO service",
	Run: func(cmd *cobra.Command, args []string) {
		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		// Set up gRPC
		address := "0.0.0.0:50052"
		lis, err := net.Listen("tcp", address)
		if err != nil {
			serviceLogger.Fatalln(err)
		}

		serviceLogger.Infof("Server is listening on %v", address)
		s := grpc.NewServer()
		db := service.MustGet(common.PluginDBMain).(*gorm.DB)
		store := userLikeStorage.NewSQLStore(db)
		demo.RegisterItemLikeServiceServer(s, rpc.NewRPCService(store))

		go func() {
			if err := s.Serve(lis); err != nil {
				serviceLogger.Fatalln(err)
			}
		}()

		opts := grpc.WithTransportCredentials(insecure.NewCredentials())
		cc, err := grpc.Dial(address, opts)
		if err != nil {
			serviceLogger.Fatalln(err)
		}

		client := demo.NewItemLikeServiceClient(cc)

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {

			//Example
			log.Println(service.MustGet("simple").(interface {
				GetValue() string
			}).GetValue())
			//

			engine.Use(middleware.Recover())
			engine.Static("/static", "./static")
			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			authStorage := userStorage.NewMySQLStorage(db)
			authCache := memcache.NewUserCaching(memcache.NewRedisCaching(service), authStorage)
			middlewareAuth := middleware.RequiredAuth(authCache, service)

			v1 := engine.Group("/v1")
			{
				v1.POST("/register", ginuser.Register(service))
				v1.POST("/login", ginuser.Login(service))
				v1.PUT("/upload", upload.UploadHandler(service))
				v1.GET("/user", middlewareAuth, ginuser.GetUser())

				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("/", ginitem.CreateNewItemHandler(service))
					items.GET("/", ginitem.ListItemHandler(service, client))
					items.GET("/:id", ginitem.GetItemHandler(service))
					items.PATCH("/:id", ginitem.UpdateItemHandler(service))
					items.DELETE("/:id", ginitem.DeleteItemHandler(service))

					items.POST("/:id/like", ginuserlike.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlike.UnlikeItem(service))
					items.GET("/:id/liked-users", ginuserlike.ListUserLiked(service))
				}

				rpc := v1.Group("rpc")
				{
					rpc.POST("/get_item_likes", ginuserlike.GetItemLikes(service))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		je, err := jaeger.NewExporter(jaeger.Options{
			AgentEndpoint: "localhost:6831",
			Process:       jaeger.Process{ServiceName: "Todo-List-Service"},
		})
		if err != nil {
			log.Println(err)
		}

		trace.RegisterExporter(je)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

		_ = subscriber.NewEngine(service).Start()
		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	rootCmd.AddCommand(cronDemoCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
