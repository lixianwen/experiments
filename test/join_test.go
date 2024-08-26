package demo

import (
	"fmt"
	"strings"
	"testing"
)

func BenchmarkBuilder(b *testing.B) {
	var buffer strings.Builder
	for i := 0; i < b.N; i++ {
		buffer.WriteString("rds_cluster_id,")
		buffer.WriteString("9930037e-3ed6-4aff")
		_ = buffer.String()
	}
}

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.Join([]string{"rds_cluster_id,", "9930037e-3ed6-4aff"}, "")
	}
}

func BenchmarkSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%s%s", "rds_cluster_id,", "9930037e-3ed6-4aff")
	}
}

func BenchmarkSimpleJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = "rds_cluster_id," + "9930037e-3ed6-4aff"
	}
}
