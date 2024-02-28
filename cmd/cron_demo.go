package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"todo-service/common"
	goservice "todo-service/go-sdk"
	"todo-service/plugin/sdk_gorm"
)

/*
update todo_items ti inner join (select item_id, count(item_id) as `count`
                                 from user_like_items
                                 group by item_id) c on c.item_id = ti.id
set ti.liked_count = c.count
*/

var cronDemoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run demo cron job",
	Run: func(cmd *cobra.Command, args []string) {
		service := goservice.New(
			goservice.WithName("social-todo-list"),
			goservice.WithVersion("1.0.0"),
			goservice.WithInitRunnable(sdk_gorm.NewGormDB("main.mysql", common.PluginDBMain)),
		)

		if err := service.Init(); err != nil {
			log.Fatalln(err)
		}

		db := service.MustGet(common.PluginDBMain).(*gorm.DB)
		log.Println("I am demo cron with DB connection:", db)
	},
}
