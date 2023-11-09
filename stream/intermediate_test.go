package stream

import (
	"reflect"
	"testing"
)

func TestStream_Filter(t *testing.T) {
	type args struct {
		predicate func(int) bool
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want []int
	}{
		{
			name: "filter even numbers",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{predicate: func(i int) bool { return i%2 == 0 }},
			want: []int{2, 4, 6},
		},
		{
			name: "filter odd numbers",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{predicate: func(i int) bool { return i%2 != 0 }},
			want: []int{1, 3, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Filter(tt.args.predicate).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Filter().ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Map(t *testing.T) {
	type args struct {
		fn func(int) int
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want []int
	}{
		{
			name: "double each element",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{fn: func(i int) int { return i * 2 }},
			want: []int{2, 4, 6},
		},
		{
			name: "square each element",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{fn: func(i int) int { return i * i }},
			want: []int{1, 4, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Map(tt.args.fn).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Map().ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Sequential(t *testing.T) {
	tests := []struct {
		name string
		s    stream[int]
		want []int
	}{
		{
			name: "sequential",
			s:    OfSlice([]int{1, 2, 3}),
			want: []int{1, 2, 3},
		},
		{
			name: "sequential",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2),
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Take(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want []int
	}{
		{
			name: "take 2",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 2},
			want: []int{1, 2},
		},
		{
			name: "take 4",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 4},
			want: []int{1, 2, 3, 4},
		},
		{
			name: "take 6",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 6},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "take 7",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 7},
			want: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name: "take 0",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 0},
			want: []int{},
		},
		{
			name: "take -1",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: -1},
			want: []int{},
		},
		{
			name: "take 2 from empty stream",
			s:    OfSlice([]int{}),
			args: args{n: 2},
			want: []int{},
		},
		{
			name: "take 0 from empty stream",
			s:    OfSlice([]int{}),
			args: args{n: 0},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Take(tt.args.n).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Take().ToSlice() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestStream_Sorted(t *testing.T) {
	cmp := func(a, b int) int { return a - b }
	tests := []struct {
		name string
		s    stream[int]
		want []int
	}{
		{
			name: "sorted",
			s:    OfSlice([]int{3, 2, 1}),
			want: []int{1, 2, 3},
		},
		{
			name: "sorted",
			s:    OfSlice([]int{3, 2, 1}).Parallel(2),
			want: []int{1, 2, 3},
		},
		{
			name: "sorted",
			s:    OfSlice([]int{1, 2, 3}),
			want: []int{1, 2, 3},
		},
		{
			name: "sorted",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2),
			want: []int{1, 2, 3},
		},
		{
			name: "sorted",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2),
			want: []int{1, 2, 3},
		},
	}
	for i := 0; i < 30; i++ {
		arr := randomSlice(20000)
		want := make([]int, len(arr))
		copy(want, arr)
		SortSlice(want, cmp)

		s := OfSlice(arr)
		tests = append(tests, struct {
			name string
			s    stream[int]
			want []int
		}{
			name: "sorted",
			s:    s,
			want: want,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Parallel(4).Sorted(cmp).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Sorted().ToSlice() = %v, want %v", got[0:10], tt.want[0:10])
			}
		})
	}
}

func TestStream_FlatMap(t *testing.T) {
	type args struct {
		fn func(int) stream[int]
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want []int
	}{
		{
			name: "flat map",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{fn: func(i int) stream[int] { return OfSlice[int]([]int{i, i}) }},
			want: []int{1, 1, 2, 2, 3, 3},
		},
		{
			name: "flat map",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{fn: func(i int) stream[int] { return OfSlice([]int{i * 2, i * 2}) }},
			want: []int{2, 2, 4, 4, 6, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.FlatMap(tt.args.fn).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.FlatMap().ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Peek(t *testing.T) {

	tests := []struct {
		name string
		s    stream[int]
		want []int
	}{
		{
			name: "peek",
			s:    OfSlice([]int{1, 2, 3}),
			want: []int{1, 2, 3},
		},
		{
			name: "peek",
			s:    OfSlice([]int{1, 2, 2, 3}),
			want: []int{1, 2, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := []int{}
			_ = tt.s.Peek(func(i int) { visited = append(visited, i) }).ToSlice()
			if !reflect.DeepEqual(visited, tt.want) {
				t.Errorf("Stream.Peek(). visited = %v, want %v", visited, tt.want)
			}
		})
	}
}

func TestStream_FilterN(t *testing.T) {
	type args struct {
		n         int
		predicate func(int) bool
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want []int
	}{
		{
			name: "filter even numbers",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 3, predicate: func(i int) bool { return i%2 == 0 }},
			want: []int{2, 4, 6},
		},
		{
			name: "filter odd numbers",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 2, predicate: func(i int) bool { return i%2 != 0 }},
			want: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.FilterN(tt.args.n, tt.args.predicate).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.FilterN().ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Skip(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want []int
	}{
		{
			name: "skip 2",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 2},
			want: []int{3, 4, 5, 6},
		},
		{
			name: "skip 4",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 4},
			want: []int{5, 6},
		},
		{
			name: "skip 6",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 6},
			want: []int{},
		},
		{
			name: "skip 7",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 7},
			want: []int{},
		},
		{
			name: "skip 0",
			s:    OfSlice([]int{1, 2, 3, 4, 5, 6}),
			args: args{n: 0},
			want: []int{1, 2, 3, 4, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Skip(tt.args.n).ToSlice()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stream.Skip().ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}

}
