package store

import (
	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/db"
	"github.com/volons/hive/libs/pubsub"
	"github.com/volons/hive/models"
)

type queue struct {
	//sync.RWMutex
	*pubsub.Topic

	//users []models.QueueItem
}

func newQueue() *queue {
	return &queue{
		Topic: pubsub.NewTopic(),
	}
}

func (q *queue) Set(users []interface{}) {
	var queue []models.QueueItem
	for _, user := range users {
		if u, ok := user.(map[string]interface{}); ok {
			usr := libs.JSONObject(u)
			id, _ := usr.GetString("id")
			name, _ := usr.GetString("name")

			queue = append(queue, models.QueueItem{
				ID:   id,
				Name: name,
			})
		}
	}

	db.Set("queue", queue)
	//q.Lock()
	//q.users = queue
	//q.Unlock()

	q.Publish(nil)
}

// JSON returns the list of users in a json serializable format
func (q *queue) JSON() []models.QueueItem {
	queue := []models.QueueItem{}
	db.Get("queue", &queue)
	return queue

	//q.RLock()
	//defer q.RUnlock()
	//return q.users
}
