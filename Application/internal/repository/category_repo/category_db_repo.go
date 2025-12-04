// FindAll fetches all categories from the database
// func (r *CategoryDBRepo) FindAll() ([]models.Category, error) {
// 	rows, err := r.db.Query("SELECT id, name, price, security FROM categories")
// 	if err != nil {
// 		if r.log != nil {
// 		}
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var categories []models.Category
// 	for rows.Next() {
// 		var c models.Category
// 		if err := rows.Scan(&c.ID, &c.Name, &c.Price, &c.Security); err != nil {
// 			if r.log != nil {
// 				r.log.Warning(fmt.Sprintf("Repo: Could not scan category row: %v", err))
// 			}
// 			continue
// 		}
// 		categories = append(categories, c)
// 	}
// 	return categories, nil
// }

// // FindByID fetches a category by its ID
// func (r *CategoryDBRepo) FindByID(id int) (models.Category, error) {
// 	row := r.db.QueryRow("SELECT id, name, price, security FROM categories WHERE id=$1", id)
// 	var c models.Category
// 	if err := row.Scan(&c.ID, &c.Name, &c.Price, &c.Security); err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			if r.log != nil {
// 				r.log.Warning(fmt.Sprintf("Repo: No category found in DB with id %d", id))
// 			}
// 			return models.Category{}, errors.New("category not found")
// 		}
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("Repo: DB error scanning category by id %d: %v", id, err))
// 		}
// 		return models.Category{}, err
// 	}
// 	return c, nil
// }

// // Create inserts a new category into the database
// func (r *CategoryDBRepo) Create(category models.Category) error {
// 	query := `
//     INSERT INTO categories (name, price, security)
//     VALUES ($1, $2, $3)
//     RETURNING id
//     `
// 	err := r.db.QueryRow(query, category.Name, category.Price, category.Security).Scan(&category.ID)
// 	if err != nil && r.log != nil {
// 		r.log.Error(fmt.Sprintf("Repo: DB error creating category '%s': %v", category.Name, err))
// 	}
// 	return err
// }

// // Update modifies an existing category in the database
// func (r *CategoryDBRepo) Update(category models.Category) error {
// 	query := `UPDATE categories SET name=$1, price=$2, security=$3 WHERE id=$4`
// 	result, err := r.db.Exec(query, category.Name, category.Price, category.Security, category.ID)
// 	if err != nil {
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("Repo: DB error updating category id %d: %v", category.ID, err))
// 		}
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("Repo: DB error fetching rows affected for category id %d: %v", category.ID, err))
// 		}
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		if r.log != nil {
// 			r.log.Warning(fmt.Sprintf("Repo: No category found to update with id %d", category.ID))
// 		}
// 		return errors.New("category not found")
// 	}
// 	return nil
// }

// func (r *CategoryDBRepo) Delete(id int) error {
// 	query := `DELETE FROM categories WHERE id=$1`
// 	result, err := r.db.Exec(query, id)
// 	if err != nil {
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("Repo: DB error deleting category id %d: %v", id, err))
// 		}
// 		return err
// 	}
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		if r.log != nil {
// 			r.log.Error(fmt.Sprintf("Repo: DB error fetching rows affected for delete category id %d: %v", id, err))
// 		}
// 		return err
// 	}
// 	if rowsAffected == 0 {
// 		if r.log != nil {
// 			r.log.Warning(fmt.Sprintf("Repo: No category found to delete with id %d", id))
// 		}
// 		return errors.New("category not found")
// 	}
// 	return nil
// }

package category_repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"loopit/internal/db"
	"loopit/internal/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CategoryDBRepo struct {
	db *db.DynamoClient
}

func NewCategoryDBRepo(db *db.DynamoClient) *CategoryDBRepo {
	return &CategoryDBRepo{db: db}
}

func (r *CategoryDBRepo) Create(category models.Category) error {
	category.ID = time.Now().UnixNano()
	idSK := fmt.Sprintf("ID#%d", category.ID)

	item := map[string]types.AttributeValue{
		"pk":        &types.AttributeValueMemberS{Value: "CATEGORY"},
		"sk":        &types.AttributeValueMemberS{Value: idSK},
		"ID":        &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", category.ID)},
		"Name":      &types.AttributeValueMemberS{Value: category.Name},
		"Price":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", category.Price)},
		"Security":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", category.Security)},
		"CreatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	_, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(r.db.Table),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

func (r *CategoryDBRepo) FindAll() ([]models.Category, error) {
	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "CATEGORY"},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	log.Print("get all=", out)

	var categories []models.Category
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &categories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
	}
	log.Print("after getting unmarshal", categories)
	return categories, nil
}

func (r *CategoryDBRepo) FindByID(id int64) (models.Category, error) {
	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "CATEGORY"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
	}

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key:       key,
	})
	if err != nil {
		return models.Category{}, fmt.Errorf("failed to get category: %w", err)
	}
	if out.Item == nil {
		return models.Category{}, errors.New("category not found")
	}

	var category models.Category
	if err := attributevalue.UnmarshalMap(out.Item, &category); err != nil {
		return models.Category{}, fmt.Errorf("failed to unmarshal category: %w", err)
	}
	return category, nil
}

func (r *CategoryDBRepo) Update(category models.Category) error {
	_, err := r.db.Client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "CATEGORY"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", category.ID)},
		},
		UpdateExpression: aws.String("SET #n = :name, Price = :price, Security = :security"),
		ExpressionAttributeNames: map[string]string{
			"#n": "Name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name":     &types.AttributeValueMemberS{Value: category.Name},
			":price":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", category.Price)},
			":security": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", category.Security)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (r *CategoryDBRepo) Delete(id int64) error {
	_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(r.db.Table),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "CATEGORY"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", id)},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (r *CategoryDBRepo) Save() error {
	return nil
}
