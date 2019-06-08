package users

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testCreateUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreate(&ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	assert.NotEmpty(ID)
}

func testCreateGetUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext := ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	}
	ID, err := f.GetOrCreate(&ext)
	assert.NoError(err)
	ID2, err := f.GetOrCreate(&ext)
	assert.NoError(err)
	assert.Equal(ID, ID2)
}

func testCreateSecondUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreate(&ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	ID2, err := f.GetOrCreate(&ExternalUser{
		ID:     "DEF",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	assert.NotEqual(ID, ID2)
}

func testCreateSecondUserSource(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreate(&ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	ID2, err := f.GetOrCreate(&ExternalUser{
		ID:     "ABC",
		Source: SourceCustom,
	})
	assert.NoError(err)
	assert.NotEqual(ID, ID2)
}

func testGetExistingUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext := &ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	}
	ID, err := f.GetOrCreate(ext)
	assert.NoError(err)
	ext2, err := f.GetExternal(ID)
	assert.NoError(err)
	assert.EqualValues(ext, ext2)
}

func testGetNoExistingUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext, err := f.GetExternal("A")
	assert.Nil(ext)
	assert.Equal(err, UserNotFound)
}
