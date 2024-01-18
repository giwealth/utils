package main

import (
	"database/sql"
	"fmt"

	"github.com/giwealth/utils/queue"
)

var (
	initSQL = `
-- 创建函数
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE 
        data json;
        notification json;
    
    BEGIN
    
        -- 根据操作类型将新行或新行转换为JSON.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;
        
        -- 将通知构造为JSON字符串.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
        
                        
        -- 执行 pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);
        
        -- 结果被忽略, 因为这是一个AFTER触发器.
        RETURN NULL;
    END;
    
$$ LANGUAGE plpgsql;

-- 创建队列表
CREATE TABLE IF NOT EXISTS queues (
	id SERIAL,
	event_id TEXT
);

-- 创建触发器
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'queues_notify_event') THEN
        CREATE TRIGGER queues_notify_event  
        AFTER INSERT ON queues
        FOR EACH ROW EXECUTE PROCEDURE notify_event();
    END IF;
END
$$;
`
)

func main() {
	var dsn string = "dbname=test_db user=postgres password=postgres host=127.0.0.1 sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(initSQL)
	if err != nil {
		panic(err)
	}

	queue, err := queue.NewQueue(dsn, "events")
	if err != nil {
		panic(err)
	}

	notice := queue.WaitForNotification()
	for {
		msg := <-notice
		fmt.Println(msg)
		// 返回结果: {"table" : "queues", "action" : "INSERT", "data" : {"id":5,"event_id":"55"}}
	}
}
