package user_repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"loopit/internal/db"
	"loopit/internal/enums"
	"loopit/internal/models"
	"loopit/internal/repository/lender_repo"

	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type UserDBRepo struct {
	db         *db.DynamoClient
	lenderRepo lender_repo.LenderRepo
}

func NewUserDBRepo(db *db.DynamoClient, lenderRepo lender_repo.LenderRepo) *UserDBRepo {
	return &UserDBRepo{db: db, lenderRepo: lenderRepo}
}

func (r *UserDBRepo) FindByID(userID int64) (*models.User, error) {
	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", userID)},
	}

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &r.db.Table,
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if out.Item == nil {
		return nil, errors.New("user not found")
	}

	type userHelper struct {
		ID           int64     `dynamodbav:"ID"`
		FullName     string    `dynamodbav:"FullName"`
		Email        string    `dynamodbav:"Email"`
		PhoneNumber  string    `dynamodbav:"PhoneNumber"`
		Address      string    `dynamodbav:"Address"`
		PasswordHash string    `dynamodbav:"PasswordHash"`
		SocietyID    int64     `dynamodbav:"SocietyID"`
		RoleString   string    `dynamodbav:"Role"` // Raw string from DynamoDB
		CreatedAt    time.Time `dynamodbav:"CreatedAt"`
		PK           string    `dynamodbav:"pk"`
		SK           string    `dynamodbav:"sk"`
	}

	var helper userHelper
	if err := attributevalue.UnmarshalMap(out.Item, &helper); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	role, err := enums.ParseRole(helper.RoleString)
	if err != nil {
		return nil, fmt.Errorf("invalid role '%s': %w", helper.RoleString, err)
	}

	user := models.User{
		ID:           helper.ID,
		FullName:     helper.FullName,
		Email:        helper.Email,
		PhoneNumber:  helper.PhoneNumber,
		Address:      helper.Address,
		PasswordHash: helper.PasswordHash,
		SocietyID:    helper.SocietyID,
		Role:         role,
		CreatedAt:    helper.CreatedAt,
		PK:           helper.PK,
		SK:           helper.SK,
	}

	log.Printf("FindByID result: %+v", user)
	return &user, nil
}

func (r *UserDBRepo) FindAll(filters models.UserFilter) ([]*models.User, error) {
	var pk, skPrefix string

	if filters.SocietyID != "" {
		pk = fmt.Sprintf("SOCIETY#%s", filters.SocietyID)
		skPrefix = "USER#ID#"
	} else {
		pk = "USER"
		if filters.Role != "" {
			skPrefix = fmt.Sprintf("ROLE#%s", filters.Role)
		} else {
			skPrefix = "ID#"
		}
	}

	out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(r.db.Table),
		KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":       &types.AttributeValueMemberS{Value: pk},
			":skPrefix": &types.AttributeValueMemberS{Value: skPrefix},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	type userHelper struct {
		ID           int64     `dynamodbav:"ID"`
		FullName     string    `dynamodbav:"FullName"`
		Email        string    `dynamodbav:"Email"`
		PhoneNumber  string    `dynamodbav:"PhoneNumber"`
		Address      string    `dynamodbav:"Address"`
		PasswordHash string    `dynamodbav:"PasswordHash"`
		SocietyID    int64     `dynamodbav:"SocietyID"`
		RoleString   string    `dynamodbav:"Role"`
		CreatedAt    time.Time `dynamodbav:"CreatedAt"`
		PK           string    `dynamodbav:"pk"`
		SK           string    `dynamodbav:"sk"`
	}

	var helpers []userHelper
	if err := attributevalue.UnmarshalListOfMaps(out.Items, &helpers); err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	var users []*models.User
	for _, h := range helpers {
		role, err := enums.ParseRole(h.RoleString)
		if err != nil {
			continue
		}
		users = append(users, &models.User{
			ID:           h.ID,
			FullName:     h.FullName,
			Email:        h.Email,
			PhoneNumber:  h.PhoneNumber,
			Address:      h.Address,
			PasswordHash: h.PasswordHash,
			SocietyID:    h.SocietyID,
			CreatedAt:    h.CreatedAt,
			PK:           h.PK,
			SK:           h.SK,
			Role:         role,
		})
	}

	if filters.Search != "" {
		var filtered []*models.User
		for _, u := range users {
			if strings.Contains(strings.ToLower(u.FullName), strings.ToLower(filters.Search)) ||
				strings.Contains(strings.ToLower(u.Email), strings.ToLower(filters.Search)) {
				filtered = append(filtered, u)
			}
		}
		return filtered, nil
	}

	return users, nil
}

func (r *UserDBRepo) DeleteByID(userID int64) error {

	user, err := r.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	keys := []map[string]types.AttributeValue{
		{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", userID)},
		},
		{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", user.Email)},
		},
		{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROLE#%s#USER#%d", user.Role.String(), userID)},
		},
		{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SOCIETY#%d", user.SocietyID)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#ID#%d", userID)},
		},
	}

	for _, key := range keys {
		_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
			TableName: aws.String(r.db.Table),
			Key:       key,
		})
		if err != nil {
			log.Printf("failed to delete item: %v", err)
		}
	}

	if user.Role == 1 {
		out, err := r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
			TableName:              aws.String(r.db.Table),
			KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk":       &types.AttributeValueMemberS{Value: "PRODUCT"},
				":skPrefix": &types.AttributeValueMemberS{Value: fmt.Sprintf("LENDER#%d", userID)},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to query products for lender: %w", err)
		}

		for _, item := range out.Items {
			var product struct {
				ID int64 `dynamodbav:"ID"`
			}
			if err := attributevalue.UnmarshalMap(item, &product); err != nil {
				log.Printf("failed to unmarshal product: %v", err)
				continue
			}

			_, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
				TableName: aws.String(r.db.Table),
				Key: map[string]types.AttributeValue{
					"pk": &types.AttributeValueMemberS{Value: "PRODUCT"},
					"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("PRODUCT#%d", product.ID)},
				},
			})
			if err != nil {
				log.Printf("failed to delete product: %v", err)
			}
		}
	}

	return nil
}

func (r *UserDBRepo) FindByEmail(email string) (*models.User, error) {

	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", email)},
	}

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &r.db.Table,
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb GetItem failed: %w", err)
	}
	if out.Item == nil {
		return nil, errors.New("user not found")
	}

	type userHelper struct {
		ID           int64     `dynamodbav:"ID"`
		FullName     string    `dynamodbav:"FullName"`
		Email        string    `dynamodbav:"Email"`
		PhoneNumber  string    `dynamodbav:"PhoneNumber"`
		Address      string    `dynamodbav:"Address"`
		PasswordHash string    `dynamodbav:"PasswordHash"`
		SocietyID    int64     `dynamodbav:"SocietyID"`
		RoleString   string    `dynamodbav:"Role"`
		CreatedAt    time.Time `dynamodbav:"CreatedAt"`
		PK           string    `dynamodbav:"pk"`
		SK           string    `dynamodbav:"sk"`
	}

	var helper userHelper
	if err := attributevalue.UnmarshalMap(out.Item, &helper); err != nil {
		return nil, fmt.Errorf("unmarshal failed into helper struct: %w", err)
	}

	role, err := enums.ParseRole(helper.RoleString)
	if err != nil {
		return nil, fmt.Errorf("role parsing failed for value '%s': %w", helper.RoleString, err)
	}

	user := models.User{
		ID:           helper.ID,
		FullName:     helper.FullName,
		Email:        helper.Email,
		PhoneNumber:  helper.PhoneNumber,
		Address:      helper.Address,
		PasswordHash: helper.PasswordHash,
		SocietyID:    helper.SocietyID,
		CreatedAt:    helper.CreatedAt,
		PK:           helper.PK,
		SK:           helper.SK,
		Role:         role,
	}

	return &user, nil
}

func (r *UserDBRepo) Create(user *models.User) error {
	user.ID = time.Now().UnixNano()
	role := user.Role.String()

	commonAttrs := map[string]types.AttributeValue{
		"ID":           &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.ID)},
		"FullName":     &types.AttributeValueMemberS{Value: user.FullName},
		"Email":        &types.AttributeValueMemberS{Value: user.Email},
		"PhoneNumber":  &types.AttributeValueMemberS{Value: user.PhoneNumber},
		"Address":      &types.AttributeValueMemberS{Value: user.Address},
		"PasswordHash": &types.AttributeValueMemberS{Value: user.PasswordHash},
		"SocietyID":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.SocietyID)},
		"Role":         &types.AttributeValueMemberS{Value: role},
		"CreatedAt":    &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	itemByID := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", user.ID)},
	}
	for k, v := range commonAttrs {
		itemByID[k] = v
	}

	itemByEmail := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", user.Email)},
	}
	for k, v := range commonAttrs {
		itemByEmail[k] = v
	}

	itemByRole := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROLE#%s#USER#%d", role, user.ID)},
	}
	for k, v := range commonAttrs {
		itemByRole[k] = v
	}

	itemBySociety := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SOCIETY#%d", user.SocietyID)},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#ID#%d", user.ID)},
	}
	for k, v := range commonAttrs {
		itemBySociety[k] = v
	}

	items := []map[string]types.AttributeValue{itemByID, itemByEmail, itemByRole, itemBySociety}
	for _, item := range items {
		_, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: &r.db.Table,
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed to create user item: %w", err)
		}
	}

	return nil
}

func (r *UserDBRepo) BecomeLender(userID int64) error {
	key := map[string]types.AttributeValue{
		"pk": &types.AttributeValueMemberS{Value: "USER"},
		"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", userID)},
	}

	out, err := r.db.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(r.db.Table),
		Key:       key,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}
	if out.Item == nil {
		return fmt.Errorf("user not found")
	}

	type userHelper struct {
		ID        int64  `dynamodbav:"UserID"`
		FullName  string `dynamodbav:"FullName"`
		Email     string `dynamodbav:"Email"`
		SocietyID int64  `dynamodbav:"SocietyID"`
	}
	var u userHelper
	if err := attributevalue.UnmarshalMap(out.Item, &u); err != nil {
		return fmt.Errorf("failed to unmarshal user: %w", err)
	}

	role := enums.RoleLender.String()
	commonAttrs := map[string]types.AttributeValue{
		"UserID":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", u.ID)},
		"FullName":  &types.AttributeValueMemberS{Value: u.FullName},
		"Email":     &types.AttributeValueMemberS{Value: u.Email},
		"SocietyID": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", u.SocietyID)},
		"Role":      &types.AttributeValueMemberS{Value: role},
		"CreatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
	}

	items := []map[string]types.AttributeValue{
		mergeMap(commonAttrs, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", u.ID)},
		}),
		mergeMap(commonAttrs, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", u.Email)},
		}),
		mergeMap(commonAttrs, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROLE#%s#USER#%d", role, u.ID)},
		}),
		mergeMap(commonAttrs, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SOCIETY#%d", u.SocietyID)},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#ID#%d", u.ID)},
		}),
		mergeMap(commonAttrs, map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: "USER"},
			"sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("NAME#%s", u.FullName)},
		}),
	}

	for _, item := range items {
		_, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(r.db.Table),
			Item:      item,
		})
		if err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}
	}

	lender := models.Lender{
		ID:            u.ID,
		IsVerified:    true,
		TotalEarnings: 0.0,
	}
	if err := r.lenderRepo.Create(&lender); err != nil {
		return fmt.Errorf("failed to create lender: %w", err)
	}

	return nil
}
func mergeMap(base, extra map[string]types.AttributeValue) map[string]types.AttributeValue {
	merged := make(map[string]types.AttributeValue, len(base)+len(extra))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range extra {
		merged[k] = v
	}
	return merged
}
