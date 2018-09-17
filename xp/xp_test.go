package xp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	sample, err := os.Open("./sample.xp")
	if assert.NoError(t, err) {
		xp, err := Read(sample)
		if assert.NoError(t, err) {
			assert.Equal(t, 2, len(xp.Layers))

			assert.Equal(t, 5, xp.Width)
			assert.Equal(t, 5, xp.Height)

			assert.Equal(t, int('A'), xp.Layers[0].Cells[0][0].Char)
			assert.Equal(t, int('A'), xp.Layers[0].Cells[1][1].Char)

			assert.Equal(t, int('R'), xp.Layers[1].Cells[1][2].Char)
			assert.Equal(t, int('R'), xp.Layers[1].Cells[1][3].Char)
		}
	}
}
