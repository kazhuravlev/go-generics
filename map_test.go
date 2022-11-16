package just_test

import (
	"sort"
	"strconv"
	"testing"

	"github.com/kazhuravlev/just"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapMerge(t *testing.T) {
	t.Parallel()

	alwaysTen := func(...int) int { return 10 }

	table := []struct {
		name string
		m1   map[int]int
		m2   map[int]int
		fn   func(...int) int
		exp  map[int]int
	}{
		{
			name: "empty_nil",
			m1:   nil,
			m2:   nil,
			fn:   alwaysTen,
			exp:  map[int]int{},
		},
		{
			name: "empty_len0",
			m1:   map[int]int{},
			m2:   map[int]int{},
			fn:   alwaysTen,
			exp:  map[int]int{},
		},
		{
			name: "merge_all_keys",
			m1:   map[int]int{1: 1},
			m2:   map[int]int{2: 2},
			fn:   alwaysTen,
			exp:  map[int]int{1: 10, 2: 10},
		},
		{
			name: "merge_all_keys_duplicated",
			m1:   map[int]int{1: 1, 2: 2},
			m2:   map[int]int{2: 2, 1: 1},
			fn:   alwaysTen,
			exp:  map[int]int{1: 10, 2: 10},
		},
		{
			name: "merge_all_keys_m1_empty",
			m1:   map[int]int{},
			m2:   map[int]int{2: 2, 1: 1},
			fn:   alwaysTen,
			exp:  map[int]int{1: 10, 2: 10},
		},
		{
			name: "merge_all_keys_m2_empty",
			m1:   map[int]int{2: 2, 1: 1},
			m2:   map[int]int{},
			fn:   alwaysTen,
			exp:  map[int]int{1: 10, 2: 10},
		},
		{
			name: "merge_all_keys_get_biggest",
			m1:   map[int]int{1: 10, 2: 11},
			m2:   map[int]int{1: 11, 2: 10},
			fn:   just.Max[int],
			exp:  map[int]int{1: 11, 2: 11},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapMerge(row.m1, row.m2, func(k, a, b int) int { return row.fn(a, b) })
			require.EqualValues(t, row.exp, res)
		})
	}
}

func TestMapFilterKeys(t *testing.T) {
	t.Parallel()

	alwaysTrue := func(_ int) bool { return true }
	alwaysFalse := func(_ int) bool { return false }

	table := []struct {
		name string
		m    map[int]int
		fn   func(int) bool
		exp  map[int]int
	}{
		{
			name: "empty_nil",
			m:    nil,
			fn:   alwaysTrue,
			exp:  map[int]int{},
		},
		{
			name: "empty_len0",
			m:    map[int]int{},
			fn:   alwaysTrue,
			exp:  map[int]int{},
		},
		{
			name: "should_copy_all_kv",
			m:    map[int]int{1: 1, 2: 2},
			fn:   alwaysTrue,
			exp:  map[int]int{1: 1, 2: 2},
		},
		{
			name: "should_ignore_all",
			m:    map[int]int{1: 1, 2: 2},
			fn:   alwaysFalse,
			exp:  map[int]int{},
		},
		{
			name: "keep_only_key_1_or_2",
			m:    map[int]int{1: 10, 2: 20, 3: 100, 4: -1},
			fn: func(k int) bool {
				return k == 1 || k == 2
			},
			exp: map[int]int{1: 10, 2: 20},
		},
		{
			name: "keep_only_even_keys",
			m:    map[int]int{1: 10, 2: 2, 3: 100, 4: -1},
			fn: func(k int) bool {
				return k%2 == 0
			},
			exp: map[int]int{2: 2, 4: -1},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapFilterKeys(row.m, row.fn)
			require.EqualValues(t, row.exp, res)
		})
	}
}

func TestMapFilterValues(t *testing.T) {
	t.Parallel()

	alwaysTrue := func(_ int) bool { return true }
	alwaysFalse := func(_ int) bool { return false }

	table := []struct {
		name string
		m    map[int]int
		fn   func(int) bool
		exp  map[int]int
	}{
		{
			name: "empty_nil",
			m:    nil,
			fn:   alwaysTrue,
			exp:  map[int]int{},
		},
		{
			name: "empty_len0",
			m:    map[int]int{},
			fn:   alwaysTrue,
			exp:  map[int]int{},
		},
		{
			name: "should_copy_all_kv",
			m:    map[int]int{1: 1, 2: 2},
			fn:   alwaysTrue,
			exp:  map[int]int{1: 1, 2: 2},
		},
		{
			name: "should_ignore_all",
			m:    map[int]int{1: 1, 2: 2},
			fn:   alwaysFalse,
			exp:  map[int]int{},
		},
		{
			name: "keep_only_values_gte_20",
			m:    map[int]int{1: 10, 2: 20, 3: 100, 4: -1},
			fn: func(v int) bool {
				return v >= 20
			},
			exp: map[int]int{2: 20, 3: 100},
		},
		{
			name: "keep_only_even_values",
			m:    map[int]int{1: 10, 2: 2, 3: 100, 4: -1},
			fn: func(v int) bool {
				return v%2 == 0
			},
			exp: map[int]int{2: 2, 3: 100, 1: 10},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapFilterValues(row.m, row.fn)
			require.EqualValues(t, row.exp, res)
		})
	}
}

func TestMapGetKeys(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		m    map[int]int
		exp  []int
	}{
		{
			name: "empty_nil",
			m:    nil,
			exp:  nil,
		},
		{
			name: "empty_len0",
			m:    map[int]int{},
			exp:  nil,
		},
		{
			name: "case1",
			m:    map[int]int{1: 11, 2: 22, 3: 33},
			exp:  []int{1, 2, 3},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapGetKeys(row.m)
			sort.Ints(res)
			require.EqualValues(t, row.exp, res)
		})
	}
}

func TestMapGetValues(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		m    map[int]int
		exp  []int
	}{
		{
			name: "empty_nil",
			m:    nil,
			exp:  nil,
		},
		{
			name: "empty_len0",
			m:    map[int]int{},
			exp:  nil,
		},
		{
			name: "case1",
			m:    map[int]int{1: 11, 2: 22, 3: 33},
			exp:  []int{11, 22, 33},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapGetValues(row.m)
			sort.Ints(res)
			require.EqualValues(t, row.exp, res)
		})
	}
}

func TestMapPairs(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		m    map[int]int
		exp  []just.KV[int, int]
	}{
		{
			name: "empty_nil",
			m:    nil,
			exp:  []just.KV[int, int]{},
		},
		{
			name: "empty_len0",
			m:    map[int]int{},
			exp:  []just.KV[int, int]{},
		},
		{
			name: "case1",
			m:    map[int]int{1: 11, 2: 22, 3: 33},
			exp: []just.KV[int, int]{
				{Key: 1, Val: 11},
				{Key: 2, Val: 22},
				{Key: 3, Val: 33},
			},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapPairs(row.m)
			require.EqualValues(t, row.exp, just.SliceSortCopy(res, func(a, b just.KV[int, int]) bool { return a.Key < b.Key }))
		})
	}
}

func TestMapDefaults(t *testing.T) {
	t.Parallel()

	table := []struct {
		name              string
		in, defaults, exp map[int]int
	}{
		{
			name:     "empty",
			in:       nil,
			defaults: nil,
			exp:      map[int]int{},
		},
		{
			name:     "empty_defaults",
			in:       map[int]int{1: 1, 2: 2},
			defaults: nil,
			exp:      map[int]int{1: 1, 2: 2},
		},
		{
			name:     "defaults_will_not_rewrite_src",
			in:       map[int]int{1: 1, 2: 2},
			defaults: map[int]int{1: 11, 2: 22},
			exp:      map[int]int{1: 1, 2: 2},
		},
		{
			name:     "defaults_will_extend_non_exists_keys",
			in:       map[int]int{1: 1, 2: 2},
			defaults: map[int]int{2: 22, 3: 33},
			exp:      map[int]int{1: 1, 2: 2, 3: 33},
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapDefaults(row.in, row.defaults)
			assert.Equal(t, row.exp, res)
		})
	}
}

func TestMapContainsKeysAny(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		in   map[int]int
		keys []int
		exp  bool
	}{
		{
			name: "empty_both",
			in:   nil,
			keys: nil,
			exp:  false,
		},
		{
			name: "empty_keys",
			in:   map[int]int{1: 1, 2: 2},
			keys: nil,
			exp:  false,
		},
		{
			name: "empty_in",
			in:   nil,
			keys: []int{1, 2, 3},
			exp:  false,
		},
		{
			name: "one_key_is_exists",
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{1, 100},
			exp:  true,
		},
		{
			name: "all_keys_not_exists",
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{100, 200, 300},
			exp:  false,
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapContainsKeysAny(row.in, row.keys)
			assert.Equal(t, row.exp, res)
		})
	}
}

func TestMapContainsKeysAll(t *testing.T) {
	t.Parallel()

	table := []struct {
		name string
		in   map[int]int
		keys []int
		exp  bool
	}{
		{
			name: "empty_both",
			in:   nil,
			keys: nil,
			exp:  false,
		},
		{
			name: "empty_keys",
			in:   map[int]int{1: 1, 2: 2},
			keys: nil,
			exp:  false,
		},
		{
			name: "empty_in",
			in:   nil,
			keys: []int{1, 2, 3},
			exp:  false,
		},
		{
			name: "one_key_is_exists",
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{1, 100},
			exp:  false,
		},
		{
			name: "all_keys_not_exists",
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{100, 200, 300},
			exp:  false,
		},
		{
			name: "all_keys_is_exists",
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{1, 2},
			exp:  true,
		},
	}

	for _, row := range table {
		row := row
		t.Run(row.name, func(t *testing.T) {
			t.Parallel()

			res := just.MapContainsKeysAll(row.in, row.keys)
			assert.Equal(t, row.exp, res)
		})
	}
}

func TestMapMap(t *testing.T) {
	in := map[int]int{
		1: 11,
		2: 22,
	}
	res := just.MapMap(in, func(k, v int) (string, string) {
		return strconv.Itoa(k), strconv.Itoa(v)
	})
	exp := map[string]string{
		"1": "11",
		"2": "22",
	}

	require.Equal(t, exp, res)
}

func TestMapMapErr(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		in := map[int]string{
			1: "11",
			2: "22",
		}
		res, err := just.MapMapErr(in, func(k int, v string) (int, int, error) {
			vInt, err := strconv.Atoi(v)
			return k, vInt, err
		})
		exp := map[int]int{
			1: 11,
			2: 22,
		}

		require.NoError(t, err)
		require.Equal(t, exp, res)
	})

	t.Run("fail", func(t *testing.T) {
		in := map[int]string{
			1: "11",
			2: "not-a-number",
		}
		res, err := just.MapMapErr(in, func(k int, v string) (int, int, error) {
			vInt, err := strconv.Atoi(v)
			return k, vInt, err
		})

		require.Error(t, err)
		require.Empty(t, res)
	})
}

func TestMapApply(t *testing.T) {
	var callCounter int
	just.MapApply(map[int]int{
		1: 1,
		2: 2,
	}, func(k, v int) {
		callCounter += 1
	})

	require.Equal(t, 2, callCounter)
}

func TestMapJoin(t *testing.T) {
	table := []struct {
		maps []map[int]int
		exp  map[int]int
	}{
		{
			maps: nil,
			exp:  map[int]int{},
		},
		{
			maps: []map[int]int{},
			exp:  map[int]int{},
		},
		{
			maps: []map[int]int{
				{1: 1},
				{2: 2},
				{3: 3},
			},
			exp: map[int]int{1: 1, 2: 2, 3: 3},
		},
		{
			maps: []map[int]int{
				{1: 1},
				{1: 2},
				{1: 3},
			},
			exp: map[int]int{1: 3},
		},
		{
			maps: []map[int]int{
				{1: 1},
				{},
				nil,
			},
			exp: map[int]int{1: 1},
		},
	}

	for _, row := range table {
		t.Run("", func(t *testing.T) {
			res := just.MapJoin(row.maps...)
			require.Equal(t, row.exp, res)
		})
	}
}

func TestMapGetDefault(t *testing.T) {
	table := []struct {
		in     map[int]int
		key    int
		defVal int
		exp    int
	}{
		{
			in:     nil,
			key:    10,
			defVal: 7,
			exp:    7,
		},
		{
			in:     map[int]int{1: 1},
			key:    1,
			defVal: 2,
			exp:    1,
		},
		{
			in:     map[int]int{1: 1},
			key:    10,
			defVal: 7,
			exp:    7,
		},
	}

	for _, row := range table {
		t.Run("", func(t *testing.T) {
			res := just.MapGetDefault(row.in, row.key, row.defVal)
			assert.Equal(t, row.exp, res)
		})
	}
}

func TestMapNotNil(t *testing.T) {
	t.Parallel()

	table := []struct {
		in  map[int]int
		exp map[int]int
	}{
		{
			in:  nil,
			exp: map[int]int{},
		},
		{
			in:  map[int]int{},
			exp: map[int]int{},
		},
		{
			in:  map[int]int{1: 2},
			exp: map[int]int{1: 2},
		},
	}

	for _, row := range table {
		t.Run("", func(t *testing.T) {
			res := just.MapNotNil(row.in)
			assert.Equal(t, row.exp, res)
		})
	}
}

func TestMapDropKeys(t *testing.T) {
	t.Parallel()

	table := []struct {
		in   map[int]int
		keys []int
		exp  map[int]int
	}{
		{
			in:   nil,
			keys: nil,
			exp:  nil,
		},
		{
			in:   map[int]int{},
			keys: []int{1, 2, 3},
			exp:  map[int]int{},
		},
		{
			in:   map[int]int{1: 1},
			keys: []int{2, 3, 4},
			exp:  map[int]int{1: 1},
		},
		{
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{1, 2, 3, 4},
			exp:  map[int]int{},
		},
		{
			in:   map[int]int{1: 1, 2: 2},
			keys: []int{1, 1, 1, 1, 1, 1},
			exp:  map[int]int{2: 2},
		},
	}

	for _, row := range table {
		t.Run("", func(t *testing.T) {
			just.MapDropKeys(row.in, row.keys...)
			assert.Equal(t, row.exp, row.in)
		})
	}
}
