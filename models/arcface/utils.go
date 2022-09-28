package arcface

import (
	"fmt"
	"math"
)


// 向量余弦相似度
func cosine(a []float32, b []float32) (float64, error) {
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0

	for k := 0; k < len(a); k++ {
		sumA += float64(a[k]) * float64(b[k])
		s1 += math.Pow(float64(a[k]), 2)
		s2 += math.Pow(float64(b[k]), 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, fmt.Errorf("Vectors should not be null (all zeros)")
	}
	return - sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

// 向量正则化
func norm(a []float32) ([]float32, error) {
	s1 := 0.0

	for k := 0; k < len(a); k++ {
		s1 += math.Pow(float64(a[k]), 2)
	}
	if s1 == 0 {
		return nil, fmt.Errorf("Vectors should not be null (all zeros)")
	}
	norm :=  float32(math.Sqrt(s1))

	for k := 0; k < len(a); k++ {
		a[k] = a[k] / norm
	}

	return a, nil
}
