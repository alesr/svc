package svc

import (
	"reflect"
	"strings"
	"testing"

	"github.com/caarlos0/env/v6"
	"github.com/stretchr/testify/require"
)

func TestLoadFromEnv(t *testing.T) {
	err := LoadFromEnv(map[string]string{})
	require.Error(t, err)
	require.Equal(t, env.ErrNotAStructPtr, err)

	test := struct {
		StrVal           string `env:"strVal"`
		IntVal           int    `env:"intVal"`
		EmptyRequiredVal string `env:"emptyVal" validate:"required"`
	}{}

	t.Setenv("strVal", "testStrVal")
	t.Setenv("intVal", "123")

	err = LoadFromEnv(&test)
	require.Error(t, err)
	require.Equal(t, "Key: 'EmptyRequiredVal' Error:Field validation for 'EmptyRequiredVal' failed on the 'required' tag", err.Error())
	require.Equal(t, "testStrVal", test.StrVal)
	require.Equal(t, 123, test.IntVal)
}

func TestLoadFromEnvWithParsers(t *testing.T) {
	test := struct {
		MapVal map[string]string `env:"mapVal" envDefault:""`
	}{}

	t.Setenv("mapVal", "testKey:testVal")

	err := LoadFromEnvWithParsers(&test, map[reflect.Type]env.ParserFunc{
		reflect.TypeOf(map[string]string{}): func(v string) (interface{}, error) {
			items := strings.Split(v, ":")
			return map[string]string{items[0]: items[1]}, nil
		},
	})
	require.NoError(t, err)
	require.Equal(t, map[string]string{"testKey": "testVal"}, test.MapVal)
}
