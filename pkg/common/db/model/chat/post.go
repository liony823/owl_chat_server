package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/common/db/dbutil"
	"github.com/openimsdk/chat/pkg/common/db/table/chat"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPost(db *mongo.Database) (chat.PostInterface, error) {
	coll := db.Collection("post")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "post_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &Post{coll: coll}, nil
}

type Post struct {
	coll *mongo.Collection
}

func (g *Post) sortByCreateTime() bson.D {
	return bson.D{{Key: "create_time", Value: -1}}
}

func (o *Post) sortByPinedAndCreateTime() bson.D {
	return bson.D{
		{Key: "is_pined", Value: -1},
		{Key: "create_time", Value: -1},
	}
}

func (o *Post) Create(ctx context.Context, posts []*chat.PostDB) error {
	for i, post := range posts {
		if post.CreateTime.IsZero() {
			posts[i].CreateTime = time.Now()
		}
		if post.UpdateTime.IsZero() {
			posts[i].UpdateTime = time.Now()
		}
	}
	return mongoutil.InsertMany(ctx, o.coll, posts)
}

func (o *Post) Delete(ctx context.Context, postIDs []string) error {
	if len(postIDs) == 0 {
		return nil
	}
	return mongoutil.DeleteMany(ctx, o.coll, bson.M{"post_id": bson.M{"$in": postIDs}})
}
func (o *Post) Take(ctx context.Context, postID string) (*chat.Post, error) {
	filter := bson.M{"post_id": postID}
	results, err := mongoutil.Aggregate[*chat.Post](ctx, o.coll, GetAggregationPipeline(ctx, filter))
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return results[0], nil
}

func (o *Post) UpdateByMap(ctx context.Context, postID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	filter := bson.M{"post_id": postID}
	data["update_time"] = time.Now()
	return mongoutil.UpdateOne(ctx, o.coll, filter, bson.M{"$set": data}, false)
}

func (o *Post) GetPostsByCursorAndUserIDs(ctx context.Context, cursor int64, userIDs []string, count int64) ([]*chat.Post, string, error) {
	filter := bson.M{
		"user_id": bson.M{"$in": userIDs},
		"$or": []bson.M{
			{"comment_post_id": nil},
			{"comment_post_id": ""},
			{"comment_post_id": bson.M{"$exists": false}},
		},
	}
	sort := o.sortByCreateTime()
	return dbutil.FindPageWithCursor[*chat.Post](ctx, o.coll, cursor, "CreateTime", -1, count, filter, sort, GetAggregationPipeline(ctx))
}

func (o *Post) GetPostsByCursorAndUser(ctx context.Context, cursor int64, userID string, count int64) ([]*chat.Post, string, error) {
	filter := bson.M{"user_id": userID}
	sort := o.sortByPinedAndCreateTime()
	return dbutil.FindPageWithCursor[*chat.Post](ctx, o.coll, cursor, "CreateTime", -1, count, filter, sort, GetAggregationPipeline(ctx))
}

func (o *Post) GetPostsByCursorAndPostIDs(ctx context.Context, cursor int64, postIDs []string, count int64) ([]*chat.Post, string, error) {
	filter := bson.M{"post_id": bson.M{"$in": postIDs}}
	sort := o.sortByCreateTime()
	return dbutil.FindPageWithCursor[*chat.Post](ctx, o.coll, cursor, "CreateTime", -1, count, filter, sort, GetAggregationPipeline(ctx))
}

func (o *Post) GetCommentPostsByPostID(ctx context.Context, cursor int64, postID string, count int64) ([]*chat.Post, string, error) {
	filter := bson.M{"comment_post_id": postID}
	sort := o.sortByCreateTime()
	return dbutil.FindPageWithCursor[*chat.Post](ctx, o.coll, cursor, "CreateTime", -1, count, filter, sort, GetAggregationPipeline(ctx))
}

func (o *Post) GetCommentPostIDsByUser(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{"user_id": userID, "comment_post_id": bson.M{
		"$exists": true,
		"$nin":    []interface{}{nil, ""},
	}}
	return mongoutil.Find[string](ctx, o.coll, filter, options.Find().SetProjection(bson.M{"post_id": 1, "_id": 0}))
}

func (o *Post) GetPinnedPostByUserID(ctx context.Context, userID string) (*chat.Post, error) {
	filter := bson.M{"user_id": userID, "is_pined": constant.Pinned}
	return mongoutil.FindOne[*chat.Post](ctx, o.coll, filter)
}

// 通过forward_post_id获取post
func (o *Post) GetPostByForwardPostID(ctx context.Context, userID, forwardPostID string) (*chat.Post, error) {
	filter := bson.M{"user_id": userID, "forward_post_id": forwardPostID}
	return mongoutil.FindOne[*chat.Post](ctx, o.coll, filter)
}

func (o *Post) GetFollowedUserIDs(ctx context.Context, userID string) ([]string, error) {
	userRelationColl := o.coll.Database().Collection("friend_relation")

	filter := bson.M{"owner_user_id": userID, "is_following": 1, "is_blocked": 0}

	cursor, err := userRelationColl.Find(ctx, filter, options.Find().SetProjection(bson.M{"related_user_id": 1, "_id": 0}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		RelatedUserID string `bson:"related_user_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	followedUserIDs := make([]string, len(results))
	for i, result := range results {
		followedUserIDs[i] = result.RelatedUserID
	}

	return followedUserIDs, nil
}

func (o *Post) GetSubscriberUserIDs(ctx context.Context, userID string) ([]string, error) {
	// 假设user_relation集合在同一个数据库中
	userRelationColl := o.coll.Database().Collection("friend_relation")

	filter := bson.M{"owner_user_id": userID, "is_subscribed": 1, "is_blocked": 0}

	cursor, err := userRelationColl.Find(ctx, filter, options.Find().SetProjection(bson.M{"related_user_id": 1, "_id": 0}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		RelatedUserID string `bson:"related_user_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	subscriberUserIDs := make([]string, len(results))
	for i, result := range results {
		subscriberUserIDs[i] = result.RelatedUserID
	}

	return subscriberUserIDs, nil

}

func GetAggregationPipeline(ctx context.Context, filter ...bson.M) mongo.Pipeline {
	opUserID, _ := mctx.CheckUser(ctx)
	var _pipeline []bson.D
	if len(filter) > 0 {
		_pipeline = append(_pipeline, bson.D{{Key: "$match", Value: filter[0]}})
	}
	_pipeline = append(_pipeline,
		lookupUserInfo(),
		unwindUserInfo(),
		lookupRelations(),
		lookupPost(ForwardPost, opUserID, 2),
		unwindPost(ForwardPost),
		lookupPost(CommentPost, opUserID, 2),
		unwindPost(CommentPost),
		lookupPost(RefPost, opUserID, 2),
		unwindPost(RefPost),
		lookupAtUserInfo(),
		lookupCommentPostsByUserID(opUserID),
		addFields(opUserID),
	)

	return _pipeline
}

type PostLookupType string

const (
	CommentPost PostLookupType = "comment"
	ForwardPost PostLookupType = "forward"
	RefPost     PostLookupType = "ref"
)

func lookupPost(lookupType PostLookupType, opUserID string, maxDepth int) bson.D {

	var matchField string
	var asField string

	switch lookupType {
	case CommentPost:
		matchField = "comment_post_id"
		asField = "comment_post"
	case ForwardPost:
		matchField = "forward_post_id"
		asField = "forward_post"
	case RefPost:
		matchField = "ref_post_id"
		asField = "ref_post"
	}

	if maxDepth <= 0 {
		return bson.D{{Key: "$addFields", Value: bson.D{{Key: asField, Value: nil}}}}
	}

	pipeline := bson.A{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$expr", Value: bson.D{{Key: "$eq", Value: bson.A{"$post_id", "$$postId"}}}},
			}},
		},
		lookupUserInfo(),
		unwindUserInfo(),
		lookupRelations(),
		lookupPost(ForwardPost, opUserID, maxDepth-1),
		unwindPost(ForwardPost),
		lookupPost(CommentPost, opUserID, maxDepth-1),
		unwindPost(CommentPost),
		lookupPost(RefPost, opUserID, maxDepth-1),
		unwindPost(RefPost),
		lookupAtUserInfo(),
		lookupCommentPostsByUserID(opUserID),
		addFields(opUserID),
	}

	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "post"},
			{Key: "let", Value: bson.D{{Key: "postId", Value: "$" + matchField}}},
			{Key: "pipeline", Value: pipeline},
			{Key: "as", Value: asField},
		}},
	}
}

func unwindPost(lookupType PostLookupType) bson.D {

	var path string

	switch lookupType {
	case CommentPost:
		path = "$comment_post"
	case ForwardPost:
		path = "$forward_post"
	case RefPost:
		path = "$ref_post"
	}

	return bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: path},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}
}

func lookupCommentPostsByUserID(opUserID string) bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "post"},
			{Key: "let", Value: bson.D{{Key: "postId", Value: "$post_id"}}},
			{Key: "pipeline", Value: bson.A{
				bson.D{
					{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$and", Value: bson.A{
								bson.D{{Key: "$eq", Value: bson.A{"$comment_post_id", "$$postId"}}},
								bson.D{{Key: "$eq", Value: bson.A{"$user_id", opUserID}}},
							}},
						}},
					}},
				},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "as", Value: "comment_counts"},
		}},
	}
}

func lookupUserInfo() bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "attribute"},
			{Key: "localField", Value: "user_id"},
			{Key: "foreignField", Value: "user_id"},
			{Key: "as", Value: "user_info"},
		}},
	}
}

func unwindUserInfo() bson.D {
	return bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$user_info"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}},
	}
}

func lookupRelations() bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "user_post_relation"},
			{Key: "localField", Value: "post_id"},
			{Key: "foreignField", Value: "post_id"},
			{Key: "as", Value: "relations"},
		}},
	}
}

func lookupAtUserInfo() bson.D {
	return bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "attribute"},
			{Key: "let", Value: bson.D{{Key: "atUserIds", Value: "$at_user_ids"}}},
			{Key: "pipeline", Value: bson.A{
				bson.D{
					{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$cond", Value: bson.A{
								bson.D{{Key: "$ne", Value: bson.A{"$$atUserIds", nil}}},
								bson.D{{Key: "$in", Value: bson.A{"$user_id", "$$atUserIds"}}},
								false,
							}},
						}},
					}},
				},
				bson.D{
					{Key: "$project", Value: bson.D{
						{Key: "user_id", Value: 1},
						{Key: "nickname", Value: 1},
						{Key: "account", Value: 1},
						{Key: "address", Value: 1},
						{Key: "face_url", Value: 1},
					}},
				},
			}},
			{Key: "as", Value: "at_user_info_list"},
		}},
	}
}

func addFields(opUserID string) bson.D {
	return bson.D{
		{Key: "$addFields", Value: bson.D{
			{Key: "comment_count", Value: getCommentCountField()},
			{Key: "like_count", Value: getCountField("is_liked")},
			{Key: "forward_count", Value: getCountField("is_forwarded")},
			{Key: "is_liked", Value: getIsField("is_liked", opUserID)},
			{Key: "is_collected", Value: getIsField("is_collected", opUserID)},
			{Key: "is_commented", Value: getIsField("is_commented", opUserID)},
			{Key: "is_forwarded", Value: getIsField("is_forwarded", opUserID)},
		}},
	}
}

func getCommentCountField() bson.D {
	return bson.D{
		{Key: "$ifNull", Value: bson.A{
			bson.D{{Key: "$first", Value: "$comment_counts.count"}},
			0,
		}},
	}
}

func getCountField(fieldName string) bson.D {
	return bson.D{
		{Key: "$size", Value: bson.D{
			{Key: "$filter", Value: bson.D{
				{Key: "input", Value: "$relations"},
				{Key: "as", Value: "relation"},
				{Key: "cond", Value: bson.D{
					{Key: "$eq", Value: bson.A{fmt.Sprintf("$$relation.%s", fieldName), 1}},
				}},
			}},
		}},
	}
}

func getIsField(fieldName string, opUserID string) bson.D {
	return bson.D{
		{Key: "$cond", Value: bson.D{
			{Key: "if", Value: bson.D{
				{Key: "$gt", Value: bson.A{
					bson.D{{Key: "$size", Value: bson.D{
						{Key: "$filter", Value: bson.D{
							{Key: "input", Value: "$relations"},
							{Key: "as", Value: "relation"},
							{Key: "cond", Value: bson.D{
								{Key: "$and", Value: bson.A{
									bson.D{{Key: "$eq", Value: bson.A{"$$relation.user_id", opUserID}}},
									bson.D{{Key: "$eq", Value: bson.A{"$$relation.post_id", "$post_id"}}},
									bson.D{{Key: "$eq", Value: bson.A{fmt.Sprintf("$$relation.%s", fieldName), 1}}},
								}},
							}},
						}},
					}}},
					0,
				}},
			}},
			{Key: "then", Value: 1},
			{Key: "else", Value: 0},
		}},
	}
}
