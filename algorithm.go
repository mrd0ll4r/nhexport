package nhexport

// An Algorithm is one algorithm.
type Algorithm int

// Algorithm constants.
// Taken from github.com/bitbandi/go-nicehash-api.
const (
	AlgorithmScrypt Algorithm = iota
	AlgorithmSHA256
	AlgorithmScryptNf
	AlgorithmX11
	AlgorithmX13
	AlgorithmKeccak
	AlgorithmX15
	AlgorithmNist5
	AlgorithmNeoScrypt
	AlgorithmLyra2RE
	AlgorithmWhirlpoolX
	AlgorithmQubit
	AlgorithmQuark
	AlgorithmAxiom
	AlgorithmLyra2REv2
	AlgorithmScryptJaneNf16
	AlgorithmBlake256r8
	AlgorithmBlake256r14
	AlgorithmBlake256r8vnl
	AlgorithmHodl
	AlgorithmDaggerHashimoto
	AlgorithmDecred
	AlgorithmCryptoNight
	AlgorithmLbry
	AlgorithmEquihash
	AlgorithmPascal
	AlgorithmX11Gost
	AlgorithmSia
	AlgorithmBlake2s
	AlgorithmSkunk
	AlgorithmCryptoNightV7
	AlgorithmCryptoNightHeavy
	AlgorithmLyra2Z
)

// String implements fmt.Stringer for an Algorithm.
func (t Algorithm) String() string {
	switch t {
	case AlgorithmScrypt:
		return "Scrypt"
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmScryptNf:
		return "ScryptNf"
	case AlgorithmX11:
		return "X11"
	case AlgorithmX13:
		return "X13"
	case AlgorithmKeccak:
		return "Keccak"
	case AlgorithmX15:
		return "X15"
	case AlgorithmNist5:
		return "Nist5"
	case AlgorithmNeoScrypt:
		return "NeoScrypt"
	case AlgorithmLyra2RE:
		return "Lyra2RE"
	case AlgorithmWhirlpoolX:
		return "WhirlpoolX"
	case AlgorithmQubit:
		return "Qubit"
	case AlgorithmQuark:
		return "Quark"
	case AlgorithmAxiom:
		return "Axiom"
	case AlgorithmLyra2REv2:
		return "Lyra2REv2"
	case AlgorithmScryptJaneNf16:
		return "ScryptJaneNf16"
	case AlgorithmBlake256r8:
		return "Blake256r8"
	case AlgorithmBlake256r14:
		return "Blake256r14"
	case AlgorithmBlake256r8vnl:
		return "Blake256r8vnl"
	case AlgorithmHodl:
		return "Hodl"
	case AlgorithmDaggerHashimoto:
		return "DaggerHashimoto"
	case AlgorithmDecred:
		return "Decred"
	case AlgorithmCryptoNight:
		return "CryptoNight"
	case AlgorithmLbry:
		return "Lbry"
	case AlgorithmEquihash:
		return "Equihash"
	case AlgorithmPascal:
		return "Pascal"
	case AlgorithmX11Gost:
		return "X11Gost"
	case AlgorithmSia:
		return "Sia"
	case AlgorithmBlake2s:
		return "Blake2s"
	case AlgorithmSkunk:
		return "Skunk"
	case AlgorithmCryptoNightV7:
		return "CryptoNightV7"
	case AlgorithmCryptoNightHeavy:
		return "CryptoNightHeavy"
	case AlgorithmLyra2Z:
		return "Lyra2Z"
	default:
		return "Unknown"
	}
}
