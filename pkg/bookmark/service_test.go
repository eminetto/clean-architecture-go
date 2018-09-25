package bookmark

import (
	"testing"

	"github.com/eminetto/clean-architecture-go/pkg/entity"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	repo := NewInmemRepository()
	service := NewService(repo)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	id, err := service.Store(b)
	assert.Nil(t, err)
	assert.True(t, entity.IsValidID(id.String()))
}

func TestSearchAndFindAll(t *testing.T) {
	repo := NewInmemRepository()
	service := NewService(repo)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	b2 := &entity.Bookmark{
		Name:        "Google",
		Description: "Google",
		Link:        "http://google.com",
		Tags:        []string{"search", "engine"},
		Favorite:    false,
	}
	bID, _ := service.Store(b)
	_, _ = service.Store(b2)

	t.Run("search", func(t *testing.T) {
		c, err := service.Search("minetto")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(c))
		assert.Equal(t, "Elton Minetto", c[0].Name)

		c, err = service.Search("bing")
		assert.Equal(t, entity.ErrNotFound, err)
		assert.Nil(t, c)
	})
	t.Run("find all", func(t *testing.T) {
		all, err := service.FindAll()
		assert.Nil(t, err)
		assert.Equal(t, 2, len(all))
	})

	t.Run("find", func(t *testing.T) {
		saved, err := service.Find(bID)
		assert.Nil(t, err)
		assert.Equal(t, b.Name, saved.Name)
	})
}

func TestDelete(t *testing.T) {
	repo := NewInmemRepository()
	service := NewService(repo)
	b := &entity.Bookmark{
		Name:        "Elton Minetto",
		Description: "Minetto's page",
		Link:        "http://www.eltonminetto.net",
		Tags:        []string{"golang", "php", "linux", "mac"},
		Favorite:    true,
	}
	b2 := &entity.Bookmark{
		Name:        "Google",
		Description: "Google",
		Link:        "http://google.com",
		Tags:        []string{"search", "engine"},
		Favorite:    false,
	}
	bID, _ := service.Store(b)
	b2ID, _ := service.Store(b2)

	err := service.Delete(bID)
	assert.Equal(t, entity.ErrCannotBeDeleted, err)

	err = service.Delete(b2ID)
	assert.Nil(t, err)
	_, err = service.Find(b2ID)
	assert.Equal(t, entity.ErrNotFound, err)

}
