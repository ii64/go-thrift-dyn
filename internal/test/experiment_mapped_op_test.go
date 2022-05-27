package test

import "testing"

type ptData struct {
	Sum int
}

var opHandlerTest [16]func(p *ptData) error

func opHandle0(p *ptData) error {
	defer func() {}()
	p.Sum++
	return nil
}

func opHandle1(p *ptData) error {
	defer func() {}()
	p.Sum = p.Sum + 2
	return nil
}

func opHandle2(p *ptData) error {
	defer func() {}()
	p.Sum = p.Sum + 3
	return nil
}

func opHandle3(p *ptData) error {
	defer func() {}()
	p.Sum = p.Sum + 4
	return nil
}

func init() {
	opHandlerTest[0] = opHandle0
	opHandlerTest[1] = opHandle1
	opHandlerTest[2] = opHandle2
	opHandlerTest[3] = opHandle3
}

func tResolveOpHandler1(op int, p *ptData) error {
	switch op {
	case 0:
		return opHandle0(p)
	case 1:
		return opHandle1(p)
	case 2:
		return opHandle2(p)
	case 3:
		return opHandle3(p)
	}
	return nil
}
func tResolveOpHandler2(op int, p *ptData) error {
	if op >= len(opHandlerTest) {
		return nil
	}
	return opHandlerTest[op](p)
}

// BenchmarkResolveOpHandler1/i
// BenchmarkResolveOpHandler1/i-12  	130362585	         9.152 ns/op
// BenchmarkResolveOpHandler1/i#01
// BenchmarkResolveOpHandler1/i#01-12         	135540003	         8.638 ns/op
// BenchmarkResolveOpHandler1/i#02
// BenchmarkResolveOpHandler1/i#02-12         	184812823	         5.424 ns/op
// BenchmarkResolveOpHandler1/i
// BenchmarkResolveOpHandler1/i-12  	136459102	         9.420 ns/op
// BenchmarkResolveOpHandler1/i#01
// BenchmarkResolveOpHandler1/i#01-12         	139910317	         8.565 ns/op
// BenchmarkResolveOpHandler1/i#02
// BenchmarkResolveOpHandler1/i#02-12         	214129742	         5.211 ns/op

func BenchmarkResolveOpHandler1(b *testing.B) {
	defer profCpu()()
	for j := 0; j < 3; j++ {
		b.Run("i", func(b *testing.B) {
			data := &ptData{}
			var err error
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err = tResolveOpHandler1(j, data)
			}
			b.StopTimer()
			_ = err
			_ = data
		})
	}
}

// BenchmarkResolveOpHandler2/i-12  	182362255	         6.479 ns/op
// BenchmarkResolveOpHandler2/i#01
// BenchmarkResolveOpHandler2/i#01-12         	181385888	         6.655 ns/op
// BenchmarkResolveOpHandler2/i#02
// BenchmarkResolveOpHandler2/i#02-12         	177325065	         6.006 ns/op
// BenchmarkResolveOpHandler2/i#03
// BenchmarkResolveOpHandler2/i#03-12         	271926285	         3.981 ns/op
// BenchmarkResolveOpHandler2/i-12  	177623852	         6.451 ns/op
// BenchmarkResolveOpHandler2/i#01
// BenchmarkResolveOpHandler2/i#01-12         	177248532	         6.391 ns/op
// BenchmarkResolveOpHandler2/i#02
// BenchmarkResolveOpHandler2/i#02-12         	294295111	         4.027 ns/op
// BenchmarkResolveOpHandler2/i#03
// BenchmarkResolveOpHandler2/i#03-12         	306798298	         4.124 ns/op

func BenchmarkResolveOpHandler2(b *testing.B) {
	defer profCpu()()
	for j := 0; j <= 3; j++ {
		b.Run("i", func(b *testing.B) {
			data := &ptData{}
			var err error
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err = tResolveOpHandler2(j, data)
			}
			b.StopTimer()
			_ = err
			_ = data
		})
	}
}
