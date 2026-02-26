// Package crypto XXTEA加密 (RPC消息加密层)
// 逆向自 game_script.pkg / lj_137: encrypt_zip / decrypt_unzip
// 流程: plaintext → XXTEA加密 → Base64编码 → gsub安全替换
// 解密: Base64解码 → XXTEA解密 → gzip解压
package crypto

import (
	"encoding/binary"
)

const delta = 0x9E3779B9

// XXTEAEncrypt XXTEA加密
// key 必须为16字节 (128-bit)
func XXTEAEncrypt(data, key []byte) []byte {
	v := bytesToUint32s(data)
	k := bytesToKey(key)
	n := len(v)
	if n < 2 {
		return data
	}

	rounds := 6 + 52/n
	sum := uint32(0)
	z := v[n-1]

	for i := 0; i < rounds; i++ {
		sum += delta
		e := (sum >> 2) & 3
		for p := 0; p < n-1; p++ {
			y := v[p+1]
			v[p] += mx(sum, y, z, uint32(p), e, k)
			z = v[p]
		}
		y := v[0]
		v[n-1] += mx(sum, y, z, uint32(n-1), e, k)
		z = v[n-1]
	}

	return uint32sToBytes(v)
}

// XXTEADecrypt XXTEA解密
func XXTEADecrypt(data, key []byte) []byte {
	v := bytesToUint32s(data)
	k := bytesToKey(key)
	n := len(v)
	if n < 2 {
		return data
	}

	rounds := 6 + 52/n
	sum := uint32(rounds) * delta
	y := v[0]

	for i := 0; i < rounds; i++ {
		e := (sum >> 2) & 3
		for p := n - 1; p > 0; p-- {
			z := v[p-1]
			v[p] -= mx(sum, y, z, uint32(p), e, k)
			y = v[p]
		}
		z := v[n-1]
		v[0] -= mx(sum, y, z, 0, e, k)
		y = v[0]
		sum -= delta
	}

	return uint32sToBytes(v)
}

func mx(sum, y, z, p, e uint32, k [4]uint32) uint32 {
	return ((z>>5 ^ y<<2) + (y>>3 ^ z<<4)) ^ ((sum ^ y) + (k[(p&3)^e] ^ z))
}

func bytesToUint32s(b []byte) []uint32 {
	// 补齐到4字节对齐
	padded := len(b)
	if padded%4 != 0 {
		padded += 4 - padded%4
	}
	buf := make([]byte, padded)
	copy(buf, b)

	v := make([]uint32, padded/4)
	for i := range v {
		v[i] = binary.LittleEndian.Uint32(buf[i*4:])
	}
	return v
}

func uint32sToBytes(v []uint32) []byte {
	b := make([]byte, len(v)*4)
	for i, val := range v {
		binary.LittleEndian.PutUint32(b[i*4:], val)
	}
	return b
}

func bytesToKey(key []byte) [4]uint32 {
	var k [4]uint32
	buf := make([]byte, 16)
	copy(buf, key)
	for i := range k {
		k[i] = binary.LittleEndian.Uint32(buf[i*4:])
	}
	return k
}
