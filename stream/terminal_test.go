package stream

import (
	"reflect"
	"testing"
)

func TestStream_ForEach(t *testing.T) {
	tests := []struct {
		name string
		s    stream[int]
		want []int
	}{
		{
			name: "for each",
			s:    OfSlice([]int{1, 2, 3}),
			want: []int{1, 2, 3},
		},
		{
			name: "for each",
			s:    OfSlice([]int{1, 2, 2, 3}),
			want: []int{1, 2, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := []int{}
			tt.s.ForEach(func(i int) { visited = append(visited, i) })
			if !reflect.DeepEqual(visited, tt.want) {
				t.Errorf("Stream.ForEach(). visited = %v, want %v", visited, tt.want)
			}
		})
	}
}

func TestStream_ToSlice(t *testing.T) {
	tests := []struct {
		name string
		s    stream[int]
		want []int
	}{
		{
			name: "to slice",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2).Sorted(func(a, b int) int { return a - b }),
			want: []int{1, 2, 3},
		},
		{
			name: "to slice",
			s:    OfSlice([]int{1, 2, 2, 3}).Parallel(50).Sorted(func(a, b int) int { return a - b }),
			want: []int{1, 2, 2, 3},
		},
		{
			name: "to slice",
			s:    OfSlice([]int{1, 2, 2, 3}).Parallel(1),
			want: []int{1, 2, 2, 3},
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

func TestStream_ReduceSequentially(t *testing.T) {
	type args struct {
		initial     int
		accumulator func(a, b int) int
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want int
	}{
		{
			name: "reduce sequentially",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{initial: 0, accumulator: func(a, b int) int { return a + b }},
			want: 6,
		},
		{
			name: "reduce sequentially",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{initial: 1, accumulator: func(a, b int) int { return a * b }},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.ReduceSequentially(tt.args.initial, tt.args.accumulator)
			if got != tt.want {
				t.Errorf("Stream.ReduceSequentially() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_Reduce(t *testing.T) {
	type args struct {
		initial     int
		accumulator func(a, b int) int
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want int
	}{
		{
			name: "reduce",
			s:    OfSlice([]int{1, 2, 3, 4, 5}).Parallel(2),
			args: args{initial: 0, accumulator: func(a, b int) int { return a + b }},
			want: 15,
		},
		{
			name: "reduce",
			s:    OfSlice([]int{1, 2, 3, 4, 5}).Parallel(3),
			args: args{initial: 1, accumulator: func(a, b int) int { return a * b }},
			want: 120,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Reduce(tt.args.initial, tt.args.accumulator)
			if got != tt.want {
				t.Errorf("Stream.Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStream_AllMatch(t *testing.T) {
	type args struct {
		predicate func(int) bool
	}
	tests := []struct {
		name string
		s    stream[int]
		args args
		want bool
	}{
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{predicate: func(i int) bool { return i > 0 }},
			want: true,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{predicate: func(i int) bool { return i > 1 }},
			want: false,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{predicate: func(i int) bool { return i > 2 }},
			want: false,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{predicate: func(i int) bool { return i > 3 }},
			want: false,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}),
			args: args{predicate: func(i int) bool { return i > 4 }},
			want: false,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2),
			args: args{predicate: func(i int) bool { return i > 0 }},
			want: true,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2),
			args: args{predicate: func(i int) bool { return i > 1 }},
			want: false,
		},
		{
			name: "all match",
			s:    OfSlice([]int{1, 2, 3}).Parallel(2),
			args: args{predicate: func(i int) bool { return i > 2 }},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.AllMatch(tt.args.predicate)
			if got != tt.want {
				t.Errorf("Stream.AllMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
