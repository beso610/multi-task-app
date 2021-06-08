package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository interface {
	CreateTask(parentID uuid.UUID, args CreateTaskArgs) (*Task, error)
	UpdateTask(uid uuid.UUID, args UpdateTaskArgs) (*Task, error)
	DeleteTask(uid uuid.UUID) error
	GetTask(uid uuid.UUID) (*Task, error)
	GetTasksByUserID() ([]*Task, error)
}

type CreateTaskArgs struct {
	UserID uuid.UUID
	Name   string
}

type UpdateTaskArgs struct {
	Name     string
	Finished bool
}

func (repo *DBRepository) CreateTask(parentID uuid.UUID, args CreateTaskArgs) (*Task, error) {
	uid := uuid.Must(uuid.NewUUID())
	task := &Task{
		ID:       uid,
		UserID:   args.UserID,
		Name:     args.Name,
		ParentID: parentID,
		Finished: false,
	}

	err := repo.DB.Create(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (repo *DBRepository) UpdateTask(taskID uuid.UUID, args UpdateTaskArgs) (*Task, error) {
	updateData := Task{}
	if args.Name != "" {
		updateData.Name = args.Name
	}
	if args.Finished {
		updateData.Finished = args.Finished
	}

	err := repo.DB.Model(&Task{}).Where(&Task{ID: taskID}).Updates(updateData).Error
	if err != nil {
		return nil, err
	}

	return repo.GetTask(taskID)
}

func (repo *DBRepository) DeleteTask(id uuid.UUID) error {
	err := repo.DB.Delete(&Task{ID: id}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *DBRepository) GetTask(id uuid.UUID) (*Task, error) {
	task := &Task{}
	err := repo.DB.Where("id = ?", id).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (repo *DBRepository) GetTasksByUserID(userID uuid.UUID) ([]*Task, error) {
	tasks := []*Task{}
	err := repo.DB.Where(&Task{UserID: userID}).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
