package packed

import "github.com/gogf/gf/os/gres"

func init() {
	if err := gres.Add("H4sIAAAAAAAC/wrwZmYRYeBg4GDo9BYMZEACQgycDMn5eWmZ6foQSq8kPzcnNISVgXGzr1V8S5+hN7OhyPHnz4Jt1hXvFXbhWzi5uc9111S+h5Oe9alNV/b03l26cUFAjvjuKws7Jl/Y1dUVf+zzm+0TYqqc2QtPnZ7//f/fONkIfiOZr4qTMvUTn/p0+m846O7jpJy7i2Wh1ZYpMr3l1zc/Sm62mvm3JHAt3/S3Zdr7Zz/9fPT67WTD+XXr036+tYrf+4RfOKS796TUDAHXrW1qFenHHQXY5unGxWw991/PaKpr+1EDDqbzco/Eup/nlvYsy/IIFWdkkAk7ZPDl1AzvaYnCt2/IJz0tPlnjyOl/yS60LItjVfL9iXNSf06p15u76Maa7p8ipTtS1ghNtfjBvmP6ZcWnT7We3XG9JWBpcKFF7tmsqPiPdf83x75RM9/9VfuVS1AK5/oYTbEr+8QUCyZZsYUFxiyNmHB5sm3VGl5RvWeeVhKVL8wZGP7/D/Bm54h+6dhnxcjA8IKRgQEWBQwME9GigA0eBeBgP+RrFQ/SjKwkwJuRSYQZEYPIBoNiEAaWNIJIPPGJMAi7OyBAgOG/YzMjA6arWNlA0kwMTAxtDAwMBxhBPEAAAAD//5BBLkdcAgAA"); err != nil {
		panic("add binary content to resource manager failed: " + err.Error())
	}
}
