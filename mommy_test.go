package mommy

import (
	"io"
	"math/rand"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestGenerator(t *testing.T) {
	tests := []struct {
		name      string
		config    Responses
		response  ResponseType
		overrides Overrides
		expected  string
		seed      int64
	}{
		{
			name:     "positive mommy",
			config:   DefaultResponses,
			response: PositiveResponse,
			expected: "good girl~\nmommy's so proud of you~",
		},
		{
			name:     "positive daddy",
			config:   DefaultResponses,
			response: PositiveResponse,
			overrides: Overrides{
				VariableRole:    "daddy",
				VariablePronoun: "his",
			},
			expected: "daddy loves you~",
			seed:     10,
		},
		{
			name:      "positive mommy referring to boy",
			config:    DefaultResponses,
			response:  PositiveResponse,
			overrides: Overrides{VariableAffectionateTerm: "boy"},
			expected:  "mommy loves her cute little boy~",
		},
		{
			name:     "positive daddy referring to boy",
			config:   DefaultResponses,
			response: PositiveResponse,
			overrides: Overrides{
				VariableRole:             "daddy",
				VariablePronoun:          "his",
				VariableAffectionateTerm: "boy",
			},
			expected: "well done~!\ndaddy is so happy for you~",
			seed:     10,
		},
		{
			name:     "negative mommy",
			config:   DefaultResponses,
			response: NegativeResponse,
			expected: "mommy still loves you~",
		},
		{
			name:     "thirsty positive mommy",
			config:   DefaultResponses,
			response: PositiveResponse,
			overrides: Overrides{
				VariableMood: Thirsty,
			},
			expected: "*pats your butt*\nthat's a good girl~",
			seed:     5,
		},
		{
			name:     "thirsty negative mommy",
			config:   DefaultResponses,
			response: NegativeResponse,
			overrides: Overrides{
				VariableMood: Thirsty,
			},
			expected: "is mommy's little girl having trouble reaching the keyboard~?",
		},
		{
			name:     "yikes positive mommy",
			config:   DefaultResponses,
			response: PositiveResponse,
			overrides: Overrides{
				VariableMood: Yikes,
			},
			expected: "you're so good with your fingers~\nmommy knows where her toy should put them next~",
			seed:     5,
		},
		{
			name:     "yikes negative mommy",
			config:   DefaultResponses,
			response: NegativeResponse,
			overrides: Overrides{
				VariableMood: Yikes,
			},
			expected: "never forget you belong to mommy~",
		},
	}

	for _, test := range tests {
		test := test // make parallel tests work
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			rng := rand.New(rand.NewSource(test.seed))

			gen, err := NewGeneratorWithRandom(test.config, rng)
			assert.NoError(t, err)

			actual, err := gen.Generate(test.response, test.overrides)
			assert.NoError(t, err)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	rng := rand.New(rand.NewSource(0))
	gen, _ := NewGeneratorWithRandom(DefaultResponses, rng)

	b.Run("string", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.Generate(PositiveResponse, nil)
		}
	})
	b.Run("discard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gen.GenerateTo(io.Discard, PositiveResponse, nil)
		}
	})
}
