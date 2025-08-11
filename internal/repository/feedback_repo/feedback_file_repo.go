package feedback_repo

import (
	"loopit/internal/models"
	"loopit/internal/storage"
)

type FeedBackFileRepo struct {
	feedbackFile string
	feedbacks    []models.Feedback
}

func NewFeedBackFileRepo(filePath string) (*FeedBackFileRepo, error) {
	feedbacks, err := storage.ReadJSONFile[models.Feedback](filePath)
	if err != nil {
		return nil, err
	}
	return &FeedBackFileRepo{
		feedbackFile: filePath,
		feedbacks:    feedbacks,
	}, nil
}

func (r *FeedBackFileRepo) CreateFeedback(feedback models.Feedback) error {
	feedback.ID = len(r.feedbacks) + 1
	r.feedbacks = append(r.feedbacks, feedback)
	return nil
}

func (r *FeedBackFileRepo) GetAllFeedbacks() ([]models.Feedback, error) {
	if len(r.feedbacks) == 0 {
		return nil, nil // or return an empty slice if preferred
	}
	return r.feedbacks, nil
}

func (r *FeedBackFileRepo) Save() error {
	return storage.WriteJSONFile(r.feedbackFile, r.feedbacks)
}
