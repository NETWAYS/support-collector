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

func TestNewExtensionMatch(t *testing.T) {
	re := NewExtensionMatch("ini")

	assert.True(t, re.MatchString(`test.ini`))
	assert.True(t, re.MatchString(`/etc/bla/test.ini`))
	assert.False(t, re.MatchString(`test.txt`))
}

func TestNewCommandMatch(t *testing.T) {
	assert.True(t, NewCommandMatch("icinga2", "daemon").MatchString("icinga2 daemon"))
}
