package types

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PostType struct {
	ID       bson.ObjectId `bson:"_id" json:"_id"`
	Title    string        `bson:"title" json:"title"`
	Time     time.Time     `bson:"time" json:"time"`
	Views    uint32        `bson:"views" json:"views"`
	Tags     []string      `bson:"tags" json:"tags"`
	Contents string        `bson:"contents" json:"contents"`
}

var postType = graphql.NewObject(graphql.ObjectConfig{
	Name: "post",
	Fields: graphql.Fields{
		"_id":      &graphql.Field{Type: graphql.String},
		"title":    &graphql.Field{Type: graphql.String},
		"time":     &graphql.Field{Type: graphql.DateTime},
		"views":    &graphql.Field{Type: graphql.Int},
		"tags":     &graphql.Field{Type: graphql.NewList(graphql.String)},
		"contents": &graphql.Field{Type: graphql.String},
	},
})

var Post = &graphql.Field{
	Type: graphql.NewList(postType),
	Args: graphql.FieldConfigArgument{
		"_id": &graphql.ArgumentConfig{Type: graphql.String},
		"tag": &graphql.ArgumentConfig{Type: graphql.String},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		db, err := dbConnect()
		if err != nil {
			return nil, err
		}
		defer db.Close()

		if _id, ok := p.Args["_id"].(string); ok {
			var result PostType
			if err := db.DB("blog").C("post").FindId(bson.ObjectIdHex(_id)).One(&result); err != nil {
				return nil, err
			}
			result.Views++
			if err := db.DB("blog").C("post").UpdateId(result.ID, result); err != nil {
				return nil, err
			}
			fmt.Println(time.Now().String()[:19], "{blog(_id:"+_id+")}")
			return []PostType{result}, nil
		}

		if tag, ok := p.Args["tag"].(string); ok {
			var result []PostType
			if err := db.DB("blog").C("post").Find(bson.M{"tags": tag}).Sort("-time").All(&result); err != nil {
				return nil, err
			}
			fmt.Println(time.Now().String()[:19], "{blog(tag:"+tag+")}")
			return result, nil
		}

		var result []PostType
		if err := db.DB("blog").C("post").Find(nil).Sort("-time").Limit(20).All(&result); err != nil {
			return nil, err
		}
		fmt.Println(time.Now().String()[:19], "{blog}")
		return result, nil
	},
}

var Tags = &graphql.Field{
	Type:        graphql.NewList(graphql.String),
	Description: "유니크한 태그 리스트 전체를 불러옴",
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		db, err := dbConnect()
		if err != nil {
			return nil, err
		}
		defer db.Close()

		var result []string
		if err := db.DB("blog").C("post").Find(nil).Distinct("tags", &result); err != nil {
			return nil, err
		}
		fmt.Println(time.Now().String()[:19], "{tags}")
		return result, nil
	},
}

var PostCreate = &graphql.Field{
	Type: postType,
	Args: graphql.FieldConfigArgument{
		"title":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"tags":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewList(graphql.String))},
		"contents": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		title := p.Args["title"].(string)
		var tags []string
		for _, v := range p.Args["tags"].([]interface{}) {
			tags = append(tags, v.(string))
		}
		contents := p.Args["contents"].(string)

		db, err := dbConnect()
		if err != nil {
			return nil, err
		}
		defer db.Close()

		var result = PostType{
			bson.NewObjectId(),
			title,
			time.Now(),
			0,
			tags,
			contents,
		}

		if err := db.DB("blog").C("post").Insert(result); err != nil {
			return nil, err
		}
		fmt.Println(time.Now().String()[:19], "{create(title:"+title+")}")
		return result, nil
	},
}

var PostUpdate = &graphql.Field{
	Type: postType,
	Args: graphql.FieldConfigArgument{
		"_id":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"title":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		"tags":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewList(graphql.String))},
		"contents": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		_id := p.Args["_id"].(string)
		title := p.Args["title"].(string)
		var tags []string
		for _, v := range p.Args["tags"].([]interface{}) {
			tags = append(tags, v.(string))
		}
		contents := p.Args["contents"].(string)

		db, err := dbConnect()
		if err != nil {
			return nil, err
		}
		defer db.Close()

		var result PostType

		if err := db.DB("blog").C("post").FindId(bson.ObjectIdHex(_id)).One(&result); err != nil {
			return nil, err
		}

		result.Title = title
		result.Time = time.Now()
		result.Tags = tags
		result.Contents = contents

		if err := db.DB("blog").C("post").UpdateId(result.ID, result); err != nil {
			return nil, err
		}
		fmt.Println(time.Now().String()[:19], "{update(title:"+title+")}")
		return result, err
	},
}

var PostDelete = &graphql.Field{
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		"_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		_id := p.Args["_id"].(string)

		db, err := dbConnect()
		if err != nil {
			return nil, err
		}
		defer db.Close()

		if err := db.DB("blog").C("post").RemoveId(bson.ObjectIdHex(_id)); err != nil {
			return nil, err
		}
		fmt.Println(time.Now().String()[:19], "{delete(_id:"+_id+")}")
		return "done", nil
	},
}
