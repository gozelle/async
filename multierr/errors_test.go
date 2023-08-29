package multierr

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/gozelle/testify/require"
)

func TestErrors(t *testing.T) {
	errs := Errors{}
	require.NoError(t, errs.Error())
	wg := sync.WaitGroup{}
	n := 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			errs.AddError(fmt.Errorf("error_%d", i))
		}(i)
	}
	wg.Wait()
	require.Error(t, errs.Error())
	for i := 0; i < n; i++ {
		require.True(t, strings.Contains(errs.String(), fmt.Sprintf("error_%d", i)))
	}
}
