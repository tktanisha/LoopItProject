package society_repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"loopit/internal/db"
	"loopit/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SocietyDBRepo struct {
    db  *db.DynamoClient
}

func NewSocietyDBRepo(db *db.DynamoClient) *SocietyDBRepo {
    return &SocietyDBRepo{db: db}
}

// ✅ Create Society
func (r *SocietyDBRepo) Create(society models.Society) error {
    society.ID = time.Now().UnixNano()
    idSK := fmt.Sprintf("ID#%d", society.ID)

    item := map[string]types.AttributeValue{
        "pk":        &types.AttributeValueMemberS{Value: "SOCIETY"},
        "sk":        &types.AttributeValueMemberS{Value: idSK},
        "ID":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", society.ID)},
        "Name":      &types.AttributeValueMemberS{Value: society.Name},
        "Location":  &types.AttributeValueMemberS{Value: society.Location},
        "Pincode":   &types.AttributeValueMemberS{Value: society.Pincode},
        "CreatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
    }

    _, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
        TableName: aws.String(r.db.Table),
        Item:      item,
    })
    return err
}

// ✅ FindAll Societies
func (r *SocietyDBRepo) FindAll() ([]models.Society, error) {
    out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
        TableName:              aws.String(r.db.Table),
        KeyConditionExpression: aws.String("pk = :pk"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "SOCIETY"},
        },
    })
    if err != nil {
        return nil, err
    }

    var societies []models.Society
    if err := attributevalue.UnmarshalListOfMaps(out.Items, &societies); err != nil {
        return nil, err
    }
    return societies, nil
}

// ✅ FindByID
func (r *SocietyDBRepo) FindByID(id int64) (models.Society, error) {
    log.Print("start by id")
    key := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: "SOCIETY"},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
    }

    out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
        TableName: aws.String(r.db.Table),
        Key:       key,
    })
    if err != nil {
        return models.Society{}, err
    }
    if out.Item == nil {
        return models.Society{}, errors.New("society not found")
    }

    var society models.Society
    if err := attributevalue.UnmarshalMap(out.Item, &society); err != nil {
        return models.Society{}, err
    }
    log.Print("end in id=",society)
    return society, nil
}


func (r *SocietyDBRepo) Update(society models.Society) error {
    log.Print("start in update=",society)
    _, err := r.db.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "SOCIETY"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", society.ID)},
        },
        UpdateExpression: aws.String("SET #n = :name, #loc = :location, Pincode = :pincode"),
        ExpressionAttributeNames: map[string]string{
            "#n": "Name",
            "#loc": "Location",
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":name":    &types.AttributeValueMemberS{Value: society.Name},
            ":location": &types.AttributeValueMemberS{Value: society.Location},
            ":pincode": &types.AttributeValueMemberS{Value: society.Pincode},
        },
    })
    return err
}

// ✅ Delete Society
func (r *SocietyDBRepo) Delete(id int64) error {
    _, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "SOCIETY"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
        },
    })
    return err
}