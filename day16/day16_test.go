package day16

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	lights = []string{
		`D2FE28`,
		`38006F45291200`,
		`EE00D40C823060`,
		`8A004A801A8002F478`,
		`620080001611562C8802118E34`,
		`C0015000016115A2E0802F182340`,
		`A0016C880162017C3686B18A3D4780`,
	}

	large = `A059141803C0008447E897180401F82F1E60D80021D11A3DC3F300470015786935BED80A5DB5002F69B4298A60FE73BE41968F48080328D00427BCD339CC7F431253838CCEFF4A943803D251B924EC283F16D400C9CDB3180213D2D542EC01092D77381A98DA89801D241705C80180960E93469801400F0A6CEA7617318732B08C67DA48C27551C00F972830052800B08550A277416401A5C913D0043D2CD125AC4B1DB50E0802059552912E9676931530046C0141007E3D4698E20008744D89509677DBF5759F38CDC594401093FC67BACDCE66B3C87380553E7127B88ECACAD96D98F8AC9E570C015C00B8E4E33AD33632938CEB4CD8C67890C01083B800E5CBDAB2BDDF65814C01299D7E34842E85801224D52DF9824D52DF981C4630047401400042E144698B2200C4328731CA6F9CBCA5FBB798021259B7B3BBC912803879CD67F6F5F78BB9CD6A77D42F1223005B8037600042E25C158FE0008747E8F50B276116C9A2730046801F29BC854A6BF4C65F64EB58DF77C018009D640086C318870A0C01D88105A0B9803310E2045C8CF3F4E7D7880484D0040001098B51DA0980021F17A3047899585004E79CE4ABD503005E610271ED4018899234B64F64588C0129EEDFD2EFBA75E0084CC659AF3457317069A509B97FB3531003254D080557A00CC8401F8791DA13080391EA39C739EFEE5394920C01098C735D51B004A7A92F6A0953D497B504F200F2BC01792FE9D64BFA739584774847CE26006A801AC05DE180184053E280104049D10111CA006300E962005A801E2007B80182007200792E00420051E400EF980192DC8471E259245100967FF7E6F2CF25DBFA8593108D342939595454802D79550C0068A72F0DC52A7D68003E99C863D5BC7A411EA37C229A86EBBC0CB802B331FDBED13BAB92080310265296AFA1EDE8AA64A0C02C9D49966195609C0594223005B80152977996D69EE7BD9CE4C1803978A7392ACE71DA448914C527FFE140`
)

func toBinary(input string) string {
	var sb strings.Builder
	for i := range input {
		switch input[i] {
		case '0':
			sb.WriteString("0000")
		case '1':
			sb.WriteString("0001")
		case '2':
			sb.WriteString("0010")
		case '3':
			sb.WriteString("0011")
		case '4':
			sb.WriteString("0100")
		case '5':
			sb.WriteString("0101")
		case '6':
			sb.WriteString("0110")
		case '7':
			sb.WriteString("0111")
		case '8':
			sb.WriteString("1000")
		case '9':
			sb.WriteString("1001")
		case 'A':
			sb.WriteString("1010")
		case 'B':
			sb.WriteString("1011")
		case 'C':
			sb.WriteString("1100")
		case 'D':
			sb.WriteString("1101")
		case 'E':
			sb.WriteString("1110")
		case 'F':
			sb.WriteString("1111")
		}
	}
	return sb.String()
}

type code struct {
	version int
	typ     int
	value   int
	codes   []code
}

func labelChar(input string) string {
	switch input {
	case "0000":
		return "0"
	case "0001":
		return "1"
	case "0010":
		return "2"
	case "0011":
		return "3"
	case "0100":
		return "4"
	case "0101":
		return "5"
	case "0110":
		return "6"
	case "0111":
		return "7"
	case "1000":
		return "8"
	case "1001":
		return "9"
	case "1010":
		return "A"
	case "1011":
		return "B"
	case "1100":
		return "C"
	case "1101":
		return "D"
	case "1110":
		return "E"
	case "1111":
		return "F"
	default:
		panic("invalid sequence")
	}
}

func decode(input string) (*code, int, error) {
	v, err := strconv.ParseInt(input[:3], 2, 64)
	if err != nil {
		return nil, 0, err
	}
	t, err := strconv.ParseInt(input[3:6], 2, 64)
	if err != nil {
		return nil, 0, err
	}
	switch t {
	case 4: // literal value
		var val string
		for i := 6; i < len(input); i += 5 {
			val += input[i+1 : i+5]
			if input[i] == '1' {
				continue
			}
			if input[i] == '0' {
				break
			}
			return nil, 0, errors.New("invalid state")
		}
		value, err := strconv.ParseInt(val, 2, 64)
		if err != nil {
			return nil, 0, err
		}
		return &code{int(v), int(t), int(value), nil}, 6 + len(val) + len(val)/4, nil
	default:
		if input[6] == '1' {
			num, err := strconv.ParseInt(input[7:18], 2, 64)
			if err != nil {
				return nil, 0, err
			}
			var codes []code
			start := 18
			for range num {
				c, s, err := decode(input[start:])
				if err != nil {
					return nil, 0, err
				}
				start += s
				codes = append(codes, *c)
			}
			return &code{int(v), int(t), 0, codes}, start, nil
		}
		if input[6] == '0' {
			size, err := strconv.ParseInt(input[7:22], 2, 64)
			if err != nil {
				return nil, 0, err
			}
			var codes []code
			start := 22
			for start < int(size)+22 {
				c, s, err := decode(input[start:])
				if err != nil {
					return nil, 0, err
				}
				start += s
				codes = append(codes, *c)
			}
			return &code{int(v), int(t), 0, codes}, start, nil

		}
		return nil, 0, errors.New("unkown type")
	}
}

func sumVersion(c code) int {
	s := c.version
	if len(c.codes) != 0 {
		for _, sc := range c.codes {
			s += sumVersion(sc)
		}
	}
	return s
}

func eval(c code) int {
	s := c.version
	switch c.typ {
	case 0:
		s := 0
		for _, sc := range c.codes {
			s += eval(sc)
		}
		return s
	case 1:
		s := 1
		for _, sc := range c.codes {
			s *= eval(sc)
		}
		return s
	case 2:
		s := 0
		for i, sc := range c.codes {
			scv := eval(sc)
			if i == 0 || s > scv {
				s = scv
			}
		}
		return s
	case 3:
		s := 0
		for i, sc := range c.codes {
			scv := eval(sc)
			if i == 0 || s < scv {
				s = scv
			}
		}
		return s
	case 4:
		return c.value
	case 5:
		if eval(c.codes[0]) > eval(c.codes[1]) {
			return 1
		}
		return 0
	case 6:
		if eval(c.codes[0]) < eval(c.codes[1]) {
			return 1
		}
		return 0
	case 7:
		if eval(c.codes[0]) == eval(c.codes[1]) {
			return 1
		}
		return 0
	}
	return s
}

func TestDay16Phase1(t *testing.T) {
	binValue := toBinary(lights[0])
	assert.Equal(t, "110100101111111000101000", binValue)
	c, _, err := decode(binValue)
	require.NoError(t, err)
	assert.Equal(t, code{6, 4, 2021, nil}, *c)

	binValue = toBinary(lights[1])
	assert.Equal(t, "00111000000000000110111101000101001010010001001000000000", binValue)
	c, _, err = decode(binValue)
	require.NoError(t, err)
	assert.Equal(t, code{1, 6, 0, []code{
		{6, 4, 10, nil},
		{2, 4, 20, nil},
	}}, *c)

	binValue = toBinary(lights[2])
	assert.Equal(t, "11101110000000001101010000001100100000100011000001100000", binValue)
	c, _, err = decode(binValue)
	require.NoError(t, err)
	assert.Equal(t, code{7, 3, 0, []code{
		{2, 4, 1, nil},
		{4, 4, 2, nil},
		{1, 4, 3, nil},
	}}, *c)

	c, _, err = decode(toBinary(lights[3]))
	require.NoError(t, err)
	assert.Equal(t, code{4, 2, 0, []code{
		{1, 2, 0, []code{
			{5, 2, 0, []code{
				{6, 4, 15, nil},
			}},
		}},
	}}, *c)
	assert.Equal(t, 16, sumVersion(*c))
	c, _, err = decode(toBinary(lights[4]))
	require.NoError(t, err)
	assert.Equal(t, 3, c.version)
	assert.Equal(t, 12, sumVersion(*c))
	c, _, err = decode(toBinary(lights[5]))
	require.NoError(t, err)
	assert.Equal(t, 23, sumVersion(*c))
	c, _, err = decode(toBinary(lights[6]))
	require.NoError(t, err)
	assert.Equal(t, 31, sumVersion(*c))

	c, _, err = decode(toBinary(large))
	require.NoError(t, err)
	assert.Equal(t, 860, sumVersion(*c))
}

func TestDay16Phase2(t *testing.T) {
	testCases := []struct {
		code  string
		value int
	}{
		{"C200B40A82", 3},
		{"04005AC33890", 54},
		{"880086C3E88112", 7},
		{"CE00C43D881120", 9},
		{"D8005AC2A8F0", 1},
		{"F600BC2D8F", 0},
		{"9C005AC2F8F0", 0},
		{"9C0141080250320F1802104A08", 1},
	}
	for _, tc := range testCases {
		c, _, err := decode(toBinary(tc.code))
		require.NoError(t, err)
		assert.Equal(t, tc.value, eval(*c))
	}

	c, _, err := decode(toBinary(large))
	require.NoError(t, err)
	assert.Equal(t, 470949537659, eval(*c))
}
