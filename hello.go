package main

import (
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/mat"
)

type PublicKey struct {
	a, b *mat.VecDense
}

type SecretKey struct {
	s int64
}

type Ciphertext struct {
	u, v int64
}

const Q = 97

func applyModulus(v *mat.VecDense, modulus float64) {
	for i := 0; i < v.Len(); i++ {
		val := v.AtVec(i)
		if val >= 0 {
			v.SetVec(i, float64(int64(val)%int64(modulus)))
		} else {
			v.SetVec(i, float64(int64(modulus-1)+int64(val)))
		}

	}
}

func randInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randVector(n int, min, max int64) *mat.VecDense {
	data := make([]float64, n)

	for i := range data {
		data[i] = float64(randInt(min, max))
	}

	return mat.NewVecDense(n, data)
}

func vecSumMod(v *mat.VecDense, mod float64) int64 {
	sum := 0.0
	for i := 0; i < v.Len(); i++ {
		sum += v.AtVec(i)
	}
	return int64(sum) % int64(mod)
}

func keygen(n int) (PublicKey, SecretKey) {
	rand.Seed(time.Now().UnixNano())

	s := randInt(1, 10000)
	a := randVector(n, 0, Q)
	e := randVector(n, -4, 4)
	b := mat.NewVecDense(a.Len(), nil)

	b.CloneFromVec(a)
	b.ScaleVec(float64(s), b)
	applyModulus(b, Q)
	b.AddVec(b, e)
	applyModulus(b, Q)

	return PublicKey{a, b}, SecretKey{s}
}

func encrypt(pk PublicKey, m int64) Ciphertext {
	var u, v int64

	// Sum of A vector mod Q
	u = int64(vecSumMod(pk.a, float64(Q)))

	// (Sum of B vector mod Q + Q/2*m) mod Q
	v = int64((vecSumMod(pk.b, float64(Q)) + Q/2*m) % Q)

	return Ciphertext{u, v}
}

func decrypt(sk SecretKey, ciphertext Ciphertext) int32 {
	var res int64

	// Compute (v - s*u) mod Q
	res = (ciphertext.v - sk.s*ciphertext.u) % Q
	if res < 0 {
		res = Q + res
	}

	if res < Q/2 {
		return 0
	} else {
		return 1
	}
}

func main() {
	pk, sk := keygen(20)
	//fmt.Println(pk.a, pk.b, sk.s)
	cipher := encrypt(pk, 1)
	//fmt.Println(cipher)
	plain := decrypt(sk, cipher)
	fmt.Println(plain)
	cipher = encrypt(pk, 0)
	//fmt.Println(cipher)
	plain = decrypt(sk, cipher)
	fmt.Println(plain)
	cipher = encrypt(pk, 1)
	//fmt.Println(cipher)
	plain = decrypt(sk, cipher)
	fmt.Println(plain)
	cipher = encrypt(pk, 0)
	//fmt.Println(cipher)
	plain = decrypt(sk, cipher)
	fmt.Println(plain)
}
