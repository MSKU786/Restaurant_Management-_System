package controller

import (
	"context"
	"log"
	"net/http"
	"restaurant-managment-system/database"
	"restaurant-managment-system/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type OrderItemPack struct {
	Table_id *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem");



func GetOrderItems() gin.HandlerFunc{
		return func(c *gin.Context) {
				var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);
				
				result, err := orderItemCollection.Find(context.TODO(), bson.M{});
				defer cancel();	

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing order items"})
					return
				}	

				var allOrderItems []bson.M;

				if err = result.All(ctx, &allOrderItems) ; err != nil {
					log.Fatal(err);
					return
				}	

				c.JSON(http.StatusOK, allOrderItems);
		}
}

func GetOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

			orderItemId := c.Param("order_item_id");

			var orderItem models.OrderItem;

			err := orderItemCollection.FindOne(ctx, bson.M{"order_item": orderItemId}).Decode(&orderItem);	

			defer cancel();

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing order items"})
				return
			}

			c.JSON(http.StatusOK, orderItem);
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc{
	return func(c *gin.Context) {
			orderId := c.Param("order_id");

			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

			result, err := orderItemCollection.Find(context.TODO(), bson.M{"order_id": orderId});

			defer cancel();
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing order items"})
				return;
			}

			var allOrderItems []bson.M;

			if err = result.All(ctx, &allOrderItems) ; err != nil {
				log.Fatal(err);
				return;
			}

			c.JSON(http.StatusOK, allOrderItems);

	}
}

func ItemsByOrder(id string) (OrderItmes []primitive.M, err error) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "order_id", Value: id}}}};
		lookupstage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "food"}, {Key: "localfield", Value: "food_id"}, {Key: "foreginfield", Value: "food_id"}, {Key: "as", Value: "food"}}}}
		unwindstage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

		lookupOrderstage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "order"}, {Key: "localfield", Value: "order_id"}, {Key: "foreginfield", Value: "order_id"}, {Key: "as", Value: "order"}}}}
		unwindStageOrder := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$order"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

		lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "table"}, {Key: "localfield", Value: "order.table_id"},{Key: "foreginfield", Value: "table_id"}, {Key: "as", Value: "table"}}}};
		unbindTableStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$table"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "amount", Value: "$food.price" },
				{Key: "totalCount", Value: 1},
				{Key: "food_name", Value: "$food.name"},
				{Key: "food_image", Value: "$food.food_image"},
				{Key: "table_number", Value: "$table.table_number"},
				{Key: "table_id", Value: "$table.table_id"},
				{Key: "order_id", Value: "$order.order_id"},
				{Key: "price", Value: "$food.price"},
				{Key: "quantity", Value: 1},
			},
		}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: bson.D{{Key: "order_id", Value: "$order_id"}, {Key: "table_id", Value: "$table_id"}, {Key: "table_number", Value: "$table_number"}}}, {Key: "payment_due", Value: bson.D{"$sum", "$amount"}}, {Key: "total_count", Value: bson.D{{"$sum", "$totalCount"}}}, {Key: "order_items", Value: bson.D{{"$push", bson.D{{"food_name", "$food_name"}}}}}}}}}}
		
	}

func CreateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
			var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second);

			var orderItemPack OrderItemPack
			var order models.Order;

			if err := c.BindJSON(&orderItemPack); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return;
			}


			order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));

			orderItemsToBeInserted := []interface{}{};
			order.Table_id = orderItemPack.Table_id;
			order_id := OrderItemOrderCreator(order)

			for _, orderItem := range orderItemPack.Order_items {
				orderItem.Order_id = order_id;

				validationErr := validate.Struct(orderItem);

				if validationErr != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
					return;
				}
				
				orderItem.ID = primitive.NewObjectID();
				orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
				orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));	
				orderItem.Order_item_id = orderItem.ID.Hex();
				var num = toFixed(*orderItem.Unit_price, 2)
				orderItem.Unit_price = &num;
				orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem);

			}

			insertedOrderItem, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted);

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return;
			}	

			defer cancel();

			c.JSON(http.StatusOK, insertedOrderItem);
	}
}

func UpdateOrderItem() gin.HandlerFunc{
	return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second);

			var orderItem models.OrderItem;
			
			orderItemId := c.Param("order_item_id");

			filter := bson.M{"order_item_id": orderItemId};

			var updateObj primitive.D;

			if orderItem.Unit_price != nil {
				updateObj = append(updateObj, bson.E{Key: "unit_price", Value: *&orderItem.Unit_price});
			}

			if orderItem.Quantity != 0 {
				updateObj = append(updateObj, bson.E{Key: "quantity", Value: *&orderItem.Quantity});
			}

			if orderItem.Food_id != nil {
				updateObj = append(updateObj, bson.E{Key: "food_id", Value: *orderItem.Food_id});
			}
			
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339));
			updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at});

			upsert := true;
			opt := options.UpdateOptions{Upsert: &upsert};


			result , err := orderItemCollection.UpdateOne(
				ctx,
				filter,
				bson.D{
					{Key: "$set", Value: updateObj},
				},
				&opt,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return;
			}

			defer cancel()

			c.JSON(http.StatusOK, result);
	}
}
