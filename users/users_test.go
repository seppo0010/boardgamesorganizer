package users

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testCreateUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreateUser(&ExternalUser{
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
	ID, err := f.GetOrCreateUser(&ext)
	assert.NoError(err)
	ID2, err := f.GetOrCreateUser(&ext)
	assert.NoError(err)
	assert.Equal(ID, ID2)
}

func testCreateSecondUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreateUser(&ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	ID2, err := f.GetOrCreateUser(&ExternalUser{
		ID:     "DEF",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	assert.NotEqual(ID, ID2)
}

func testCreateSecondUserSource(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreateUser(&ExternalUser{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	ID2, err := f.GetOrCreateUser(&ExternalUser{
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
	ID, err := f.GetOrCreateUser(ext)
	assert.NoError(err)
	ext2, err := f.GetExternalUser(ID)
	assert.NoError(err)
	assert.EqualValues(ext, ext2)
}

func testGetNoExistingUser(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext, err := f.GetExternalUser("3")
	assert.Nil(ext)
	assert.Equal(err, UserNotFound)
	ext, err = f.GetExternalUser("A")
	assert.Nil(ext)
	assert.Equal(err, UserNotFound)
}

func testCreateGroup(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreateGroup(&ExternalGroup{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	assert.NotEmpty(ID)
}

func testCreateGetGroup(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext := ExternalGroup{
		ID:     "ABC",
		Source: SourceTelegram,
	}
	ID, err := f.GetOrCreateGroup(&ext)
	assert.NoError(err)
	ID2, err := f.GetOrCreateGroup(&ext)
	assert.NoError(err)
	assert.Equal(ID, ID2)
}

func testCreateSecondGroup(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreateGroup(&ExternalGroup{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	ID2, err := f.GetOrCreateGroup(&ExternalGroup{
		ID:     "DEF",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	assert.NotEqual(ID, ID2)
}

func testCreateSecondGroupSource(t *testing.T, f Factory) {
	assert := assert.New(t)
	ID, err := f.GetOrCreateGroup(&ExternalGroup{
		ID:     "ABC",
		Source: SourceTelegram,
	})
	assert.NoError(err)
	ID2, err := f.GetOrCreateGroup(&ExternalGroup{
		ID:     "ABC",
		Source: SourceCustom,
	})
	assert.NoError(err)
	assert.NotEqual(ID, ID2)
}

func testGetExistingGroup(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext := &ExternalGroup{
		ID:     "ABC",
		Source: SourceTelegram,
	}
	ID, err := f.GetOrCreateGroup(ext)
	assert.NoError(err)
	ext2, err := f.GetExternalGroup(ID)
	assert.NoError(err)
	assert.EqualValues(ext, ext2)
}

func testGetNoExistingGroup(t *testing.T, f Factory) {
	assert := assert.New(t)
	ext, err := f.GetExternalGroup("3")
	assert.Nil(ext)
	assert.Equal(err, GroupNotFound)
	ext, err = f.GetExternalGroup("A")
	assert.Nil(ext)
	assert.Equal(err, GroupNotFound)
}
