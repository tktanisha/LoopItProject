// package society_repo

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"loopit/internal/db"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// )

// type SocietyDBRepo struct {
// 	db  *db.DynamoClient
// 	log logger.LoggerInterface
// }

// func NewSocietyDBRepo(db *db.DynamoClient, log logger.LoggerInterface) *SocietyDBRepo {
// 	return &SocietyDBRepo{db: db, log: log}
// }

// // Create inserts a new society into the database
// func (r *SocietyDBRepo) Create(society models.Society) error {
// 	query := `
// 	INSERT INTO societies (name, location, pincode)
// 	VALUES ($1, $2, $3)
// 	RETURNING id
// 	`
// 	err := r.db.QueryRow(query, society.Name, society.Location, society.Pincode).Scan(&society.ID)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Repo: Failed to insert society (name=%s): %v", society.Name, err))
// 		return err
// 	}
// 	r.log.Info(fmt.Sprintf("Repo: Society inserted successfully (id=%d, name=%s)", society.ID, society.Name))
// 	return nil
// }

// // FindAll returns all societies
// func (r *SocietyDBRepo) FindAll() ([]models.Society, error) {
// 	rows, err := r.db.Query("SELECT id, name, location, pincode FROM societies")
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Repo: Failed to query societies: %v", err))
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var societies []models.Society
// 	for rows.Next() {
// 		var s models.Society
// 		if err := rows.Scan(&s.ID, &s.Name, &s.Location, &s.Pincode); err != nil {
// 			r.log.Warning(fmt.Sprintf("Repo: Failed to scan society row: %v", err))
// 			continue
// 		}
// 		societies = append(societies, s)
// 	}
// 	r.log.Info(fmt.Sprintf("Repo: Retrieved %d societies", len(societies)))
// 	return societies, nil
// }

// // FindByID returns a society by its ID
// func (r *SocietyDBRepo) FindByID(id int) (models.Society, error) {
// 	row := r.db.QueryRow("SELECT id, name, location, pincode FROM societies WHERE id=$1", id)
// 	var s models.Society
// 	if err := row.Scan(&s.ID, &s.Name, &s.Location, &s.Pincode); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			r.log.Warning(fmt.Sprintf("Repo: Society not found with id=%d", id))
// 			return models.Society{}, errors.New("society not found")
// 		}
// 		r.log.Error(fmt.Sprintf("Repo: Failed to scan society by id=%d: %v", id, err))
// 		return models.Society{}, err
// 	}
// 	r.log.Info(fmt.Sprintf("Repo: Society found (id=%d, name=%s)", s.ID, s.Name))
// 	return s, nil
// }

// // update modifies an existing society in the database
// func (r *SocietyDBRepo) Update(society models.Society) error {
// 	query := `UPDATE societies SET name=$1, location=$2, pincode=$3 WHERE id=$4`
// 	result, err := r.db.Exec(query, society.Name, society.Location, society.Pincode, society.ID)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Repo: Failed to update society (id=%d): %v", society.ID, err))
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Repo: Failed to get rows affected for society update (id=%d): %v", society.ID, err))
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		r.log.Warning(fmt.Sprintf("Repo: No society found to update with id=%d", society.ID))
// 		return errors.New("society not found")
// 	}
// 	r.log.Info(fmt.Sprintf("Repo: Society updated successfully (id=%d)", society.ID))
// 	return nil
// }

// // Delete removes a society from the database by its ID
// func (r *SocietyDBRepo) Delete(id int) error {
// 	query := `DELETE FROM societies WHERE id=$1`
// 	result, err := r.db.Exec(query, id)
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Repo: Failed to delete society (id=%d): %v", id, err))
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		r.log.Error(fmt.Sprintf("Repo: Failed to get rows affected for society delete (id=%d): %v", id, err))
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		r.log.Warning(fmt.Sprintf("Repo: No society found to delete with id=%d", id))
// 		return errors.New("society not found")
// 	}
// 	r.log.Info(fmt.Sprintf("Repo: Society deleted successfully (id=%d)", id))
// 	return nil
// }

// // Save is a no-op for Postgres because changes are applied immediately
//
//	func (r *SocietyDBRepo) Save() error {
//		return nil
//	}
package society_repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"loopit/internal/db"
	"loopit/internal/models"
	"loopit/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SocietyDBRepo struct {
    db  *db.DynamoClient
    log logger.LoggerInterface
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
    return society, nil
}

// ✅ Update Society
func (r *SocietyDBRepo) Update(society models.Society) error {
    _, err := r.db.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "SOCIETY"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", society.ID)},
        },
        UpdateExpression: aws.String("SET #n = :name, Location = :location, Pincode = :pincode"),
        ExpressionAttributeNames: map[string]string{
            "#n": "Name",
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