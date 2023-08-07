// nolint
package pipe

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetKeyValue(t *testing.T) {
	cases := []struct {
		name      string
		line      string
		delimeter rune
		k         string
		v         string
		respError string
	}{
		{
			name:      "Success",
			line:      " id: test",
			delimeter: ':',
			k:         "id",
			v:         "test",
		},
		{
			name:      "Success w space",
			line:      " id  :   test  ",
			delimeter: ':',
			k:         "id",
			v:         "test",
		},
		{
			name:      "Success w tab & space",
			line:      `			 id  :   test  `,
			delimeter: ':',
			k:         "id",
			v:         "test",
		},
		{
			name:      "Success w tab & space 1",
			line:      `			 id  :   test value  `,
			delimeter: ':',
			k:         "id",
			v:         "test value",
		},
		{
			name:      "Not found",
			line:      " id - test",
			delimeter: ':',
			k:         "id",
			v:         "test",
			respError: ErrNotFound.Error(),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			h := Helper{}

			gotKey, gotValue, err := h.GetKeyValue(tc.line, tc.delimeter)

			if err == nil {
				require.NoError(t, err)

				require.Equal(t, gotKey, tc.k)
				require.Equal(t, gotValue, tc.v)
			} else {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.respError)

				require.Empty(t, gotKey)
				require.Empty(t, gotValue)
			}
		})
	}
}
