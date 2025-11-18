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

//go:generate mockgen -source=interface.go -destination=../../mock/mock_user_repo.go -package=mock

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

    // Helper struct for unmarshalling
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

    // Convert RoleString to enums.Role
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

// func (r *UserDBRepo) FindAll(filters models.UserFilter) ([]*models.User, error) {
//     var out *dynamodb.QueryOutput
//     var err error

//     if filters.SocietyID != "" {
//         out, err = r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
//             TableName:              aws.String(r.db.Table),
//             KeyConditionExpression: aws.String("pk = :pk"),
//             ExpressionAttributeValues: map[string]types.AttributeValue{
//                 ":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SOCIETY#%s", filters.SocietyID)},
//             },
//         })
//     } else if filters.Role != "" {
//         out, err = r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
//             TableName:              aws.String(r.db.Table),
//             KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :rolePrefix)"),
//             ExpressionAttributeValues: map[string]types.AttributeValue{
//                 ":pk":         &types.AttributeValueMemberS{Value: "USER"},
//                 ":rolePrefix": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROLE#%s", filters.Role)},
//             },
//         })
//     } else {
//         out, err = r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
//              TableName:              aws.String(r.db.Table),
//              KeyConditionExpression: aws.String("pk = :pk"),
//              ExpressionAttributeValues: map[string]types.AttributeValue{
//                 ":pk": &types.AttributeValueMemberS{Value: "USER"},
//     },
// })
//     }

//     if err != nil {
//         return nil, fmt.Errorf("failed to query users: %w", err)
//     }

//     var users []*models.User
//     if err := attributevalue.UnmarshalListOfMaps(out.Items, &users); err != nil {
//         return nil, fmt.Errorf("failed to unmarshal users: %w", err)
//     }
//     log.Print("get all user=",users)

//     // Apply search filter in-memory
//     if filters.Search != "" {
//         var filtered []*models.User
//         for _, u := range users {
//             if strings.Contains(strings.ToLower(u.FullName), strings.ToLower(filters.Search)) ||
//                 strings.Contains(strings.ToLower(u.Email), strings.ToLower(filters.Search)) {
//                 filtered = append(filtered, u)
//             }
//         }
//         return filtered, nil
//     }

//     return users, nil
// }
func (r *UserDBRepo) FindAll(filters models.UserFilter) ([]*models.User, error) {
    var out *dynamodb.QueryOutput
    var err error

    if filters.SocietyID != "" {
        out, err = r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
            TableName:              aws.String(r.db.Table),
            KeyConditionExpression: aws.String("pk = :pk"),
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SOCIETY#%s", filters.SocietyID)},
            },
        })
    } else if filters.Role != "" {
        out, err = r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
            TableName:              aws.String(r.db.Table),
            KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :rolePrefix)"),
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":pk":         &types.AttributeValueMemberS{Value: "USER"},
                ":rolePrefix": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROLE#%s", filters.Role)},
            },
        })
    } else {
       out, err = r.db.Client.Query(context.TODO(), &dynamodb.QueryInput{
            TableName:              aws.String(r.db.Table),
            KeyConditionExpression: aws.String("pk = :pk AND begins_with(sk, :skPrefix)"),
            ExpressionAttributeValues: map[string]types.AttributeValue{
                ":pk":       &types.AttributeValueMemberS{Value: "USER"},
                ":skPrefix": &types.AttributeValueMemberS{Value: "ID#"},
            },
})

    }

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
            continue // or log error
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

    // Apply search filter in-memory
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
    _, err := r.db.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
        TableName: aws.String(r.db.Table),
        Key: map[string]types.AttributeValue{
            "pk": &types.AttributeValueMemberS{Value: "USER"},
            "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", userID)},
        },
    })
    if err != nil {
        log.Print(err)
        return fmt.Errorf("failed to delete user: %w", err)
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
        ID           int64      `dynamodbav:"ID"`
        FullName     string     `dynamodbav:"FullName"`
        Email        string     `dynamodbav:"Email"`
        PhoneNumber  string     `dynamodbav:"PhoneNumber"`
        Address      string     `dynamodbav:"Address"`
        PasswordHash string     `dynamodbav:"PasswordHash"`
        SocietyID    int64      `dynamodbav:"SocietyID"`
        RoleString   string     `dynamodbav:"Role"` // Use string for unmarshalling the raw value
        CreatedAt    time.Time  `dynamodbav:"CreatedAt"`
        PK           string     `dynamodbav:"pk"`
        SK           string     `dynamodbav:"sk"`
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
    user.ID = time.Now().UnixNano() // int64 ID
    role := user.Role.String()

    // Common attributes
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

    // Item 1: USER - ID
    itemByID := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: "USER"},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ID#%d", user.ID)},
    }
    for k, v := range commonAttrs {
        itemByID[k] = v
    }

    // Item 2: USER - EMAIL
    itemByEmail := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: "USER"},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("EMAIL#%s", user.Email)},
    }
    for k, v := range commonAttrs {
        itemByEmail[k] = v
    }

    // Item 3: USER - ROLE
    itemByRole := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: "USER"},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROLE#%s#USER#%d", role, user.ID)},
    }
    for k, v := range commonAttrs {
        itemByRole[k] = v
    }

    // Item 4: SOCIETY - USER
    itemBySociety := map[string]types.AttributeValue{
        "pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("SOCIETY#%d", user.SocietyID)},
        "sk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#ID#%d", user.ID)},
    }
    for k, v := range commonAttrs {
        itemBySociety[k] = v
    }

    // Write all items
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


//uuid walal create

// func (r *UserDBRepo) Create(user *models.User) error {
//     user.ID = uuid.New().String()
//     idSK := fmt.Sprintf("ID#%s", user.ID)
//     emailSK := fmt.Sprintf("EMAIL#%s", user.Email)
//     role := user.Role.String()

//     // Common attributes
//     commonAttrs := map[string]types.AttributeValue{
//         "ID":           &types.AttributeValueMemberS{Value: user.ID},
//         "FullName":     &types.AttributeValueMemberS{Value: user.FullName},
//         "Email":        &types.AttributeValueMemberS{Value: user.Email},
//         "PhoneNumber":  &types.AttributeValueMemberS{Value: user.PhoneNumber},
//         "Address":      &types.AttributeValueMemberS{Value: user.Address},
//         "PasswordHash": &types.AttributeValueMemberS{Value: user.PasswordHash},
//         "SocietyID":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", user.SocietyID)},
//         "Role":         &types.AttributeValueMemberS{Value: role},
//         "CreatedAt":    &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
//     }

//     // First item: pk=USER, sk=ID#<UUID>
//     itemByID := map[string]types.AttributeValue{
//         "pk": &types.AttributeValueMemberS{Value: "USER"},
//         "sk": &types.AttributeValueMemberS{Value: idSK},
//     }
//     for k, v := range commonAttrs {
//         itemByID[k] = v
//     }

//     // Second item: pk=USER, sk=EMAIL#<email>
//     itemByEmail := map[string]types.AttributeValue{
//         "pk": &types.AttributeValueMemberS{Value: "USER"},
//         "sk": &types.AttributeValueMemberS{Value: emailSK},
//     }
//     for k, v := range commonAttrs {
//         itemByEmail[k] = v
//     }

//     // Write both items
//     _, err := r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
//         TableName: &r.db.Table,
//         Item:      itemByID,
//     })
//     if err != nil {
//         return fmt.Errorf("failed to create user by ID: %w", err)
//     }

//     _, err = r.db.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
//         TableName: &r.db.Table,
//         Item:      itemByEmail,
//     })
//     if err != nil {
//         return fmt.Errorf("failed to create user by Email: %w", err)
//     }

//     return nil
// }