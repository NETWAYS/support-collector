package obfuscate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRegexpKeyValue(t *testing.T) {
	re := NewRegexpKeyValue("password")

	assert.True(t, re.MatchString(`password = "test"`))
	assert.True(t, re.MatchString(`  password = "test"`))
	assert.True(t, re.MatchString(`;password = "test"`))
	assert.True(t, re.MatchString(`  //password = "test"`))
	assert.True(t, re.MatchString(`Password = "test"`))
	assert.True(t, re.MatchString(`db_password = "test"`))
	assert.False(t, re.MatchString(`user = "test"`))
	assert.False(t, re.MatchString(`password`))
}
