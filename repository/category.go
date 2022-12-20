package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	cate := []entity.Category{}
	err := r.db.Where("user_id = ?", id).Find(&cate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cate, nil
		}
		return cate, err
	}

	return cate, nil // TODO: replace this
}

func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	result := r.db.Create(&category).Error
	if result != nil {
		return 0, result
	}
	return category.ID, nil // TODO: replace this
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	err := r.db.CreateInBatches(categories, 5).Error
	return err // TODO: replace this
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	cate := entity.Category{}
	err := r.db.Where("id = ?", id).First(cate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return cate, nil
		}
		return cate, err
	}
	return cate, nil // TODO: replace this
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	err := r.db.Where("id = ?", category.ID).Updates(&category).Error
	return err // TODO: replace this
}

func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	cate := entity.Category{}
	err := r.db.Where("id = ?", id).Delete(&cate).Error
	return err // TODO: replace this
}
