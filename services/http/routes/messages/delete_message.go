package messages

import (
	"github.com/gin-gonic/gin"
	"niromash-api/api"
)

func DeleteMessageRoute() *api.Route {
	return &api.Route{
		Path:            "/:scope",
		Method:          api.MethodDelete,
		IsAuthenticated: true,
		AuthenticateMiddleware: func(c *gin.Context, user api.User, service api.MainService) {
			if !user.HasPermission("messages.delete") {
				c.AbortWithStatusJSON(403, gin.H{
					"message": "You do not have the permission!",
				})
				return
			}

			return
		},
		Handler: func(c *gin.Context, user api.User, service api.MainService) {
			//scope := c.Params("scope")
			//key := c.Params("key")
			// Todo with postgres
			//coll := database.Mongo.Database("datasource").Collection("messages")
			//
			//filters := map[string]bson.D{
			//	"scopeOnly": {
			//		{
			//			"scope",
			//			bson.M{"$regex": primitive.Regex{Pattern: "^" + scope + "$", Options: "i"}},
			//		},
			//	},
			//	"scopeAndKey": {
			//		{
			//			"scope",
			//			bson.M{"$regex": primitive.Regex{Pattern: "^" + scope + "$", Options: "i"}},
			//		},
			//		{
			//			"translations",
			//			bson.M{
			//				"$elemMatch": bson.M{"key": primitive.Regex{Pattern: "^" + key + "$", Options: "i"}},
			//			},
			//		},
			//	},
			//}
			//
			//if len(key) > 0 {
			//	result := coll.FindOneAndUpdate(context.Background(), filters["scopeAndKey"], bson.M{
			//		"$pull": bson.M{
			//			"translations": bson.M{"key": key},
			//		},
			//	})
			//
			//	if err = result.Err(); err != nil {
			//		if err == mongo.ErrNoDocuments {
			//			return c.Status(400).SendString(fmt.Sprintf("No result found with scope: %s and key: %s does not exist!", scope, key))
			//		}
			//		return
			//	}
			//
			//	if err = database.Redis.Del(context.Background(), fmt.Sprintf("message:%s:%s", scope, key)).Err(); err != nil {
			//		return
			//	}
			//
			//	return c.SendString(fmt.Sprintf("The translation with the scope: %s with the key: %s has been deleted!", scope, key))
			//}
			//
			//cursor := coll.FindOneAndDelete(context.Background(), filters["scopeOnly"])
			//if err = cursor.Err(); err != nil {
			//	if err == mongo.ErrNoDocuments {
			//		return c.Status(400).SendString(fmt.Sprintf("The scope: %s does not exist!", scope))
			//	}
			//	return
			//}
			//
			//var keys []string
			//keys, err = database.Redis.Keys(context.Background(), "message:"+scope+"*").Result()
			//if err != nil {
			//	return
			//}
			//
			//for _, keyIterated := range keys {
			//	database.Redis.Del(context.Background(), keyIterated)
			//
			//}
			//return c.SendString(fmt.Sprintf("The whole scope: %s has been deleted", scope))
		},
	}
}
