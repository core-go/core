package currencies

import (
	"strconv"
	"strings"
)

var s = "AED||د.إ;AFN||؋;ALL||Lek;AMD||դր.;ARS||;AUD||;AZN||ман.;BAM||KM;BDT||৳;BGN||лв.;BHD|3|BD;BND|0|;BOB||$b;BRL||R;BYR||р.;BZD||BZ;CAD||;CHF||Fr.;CLP||;CNY||¥;COP||;CRC||₡;CSD||Дин.;CZK||Kč;DKK||kr.;DOP||RD;DZD||DA;EEK||kr;EGP||£;ETB||Br;EUR||€;GBP||£;GEL||Lari;GTQ||Q;HKD||HK;HNL||L.;HRK||kn;HUF||Ft;IDR|0|Rp;ILS||₪;INR||₹;IQD||ID;IRR||ريال;ISK|0|kr.;JMD||J;JOD|3|د.أ;JPY|0|¥;KES||S;KGS||сом;KHR||៛;KRW|0|₩;KWD|3|KD;KZT||Т;LAK||₭;LBP||LL;LKR||රු.;LTL||Lt;LVL||Ls;LYD|3|LD;MAD||DH;MKD||ден.;MNT||₮;MOP||;MVR||ރ.;MXN||;MYR||RM;NIO||C;NOK||kr;NPR||रु;NZD||;OMR|3|R.O;PAB||B/.;PEN||S/.;PHP||₱;PKR||Rs;PLN||zł;PYG||Gs;QAR||QR;RON||lei;RSD||Дин.;RUB||һ.;RWF||R₣;SAR||SR;SEK||kr;SGD||;SYP||LS;THB||฿;TJS||т.р.;TMT||m.;TND|3|DT;TRY||TL;TTD||TT;TWD||NT;UAH||₴;USD||;UYU||$U;UZS|0|лв;VEF||Bs.;VND|0|₫;XOF||XOF;YER||﷼;ZAR||R;ZWL||Z$"

type Currency struct {
	Code string `mapstructure:"code" json:"code,omitempty" gorm:"column:code;primary_key" bson:"_id,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty" avro:"code" validate:"required,max=40"`
	Symbol string `mapstructure:"symbol" json:"symbol,omitempty" gorm:"column:symbol" bson:"symbol,omitempty" dynamodbav:"symbol,omitempty" firestore:"symbol,omitempty" avro:"symbol" validate:"required,max=6"`
	DecimalDigits int `mapstructure:"decimal_digits" json:"decimalDigits" gorm:"column:decimal_digits" bson:"DecimalDigits" dynamodbav:"decimalDigits" firestore:"decimalDigits" avro:"decimalDigits"`
}

var currencies []Currency
var currencyMap = map[string]Currency{}

func Init() {
	if len(currencies) > 0 {
		return
	}
	a := strings.Split(s, ";")
	l := len(a)
	d := "$"
	t := 2
	for i := 0; i < l ; i++ {
		c := strings.Split(a[i], "|")
		var cur Currency
		cur.Code = c[0]
		cur.DecimalDigits = t
		if len(c[1]) > 0 {
			t, _ := strconv.Atoi(c[1])
			cur.DecimalDigits = t
		}
		cur.Symbol = d
		if len(c[2]) > 0 {
			cur.Symbol = c[2]
		}
		currencies = append(currencies, cur)
		currencyMap[cur.Code] = cur
	}
}

func GetCurrency(code string) *Currency {
	Init()
	c, ok := currencyMap[code]
	if ok {
		return &c
	}
	return nil
}
func GetAll() []Currency {
	Init()
	return currencies
}
func Query(q string, decimalDigits *int) []Currency {
	Init()
	if len(q) == 0 && decimalDigits == nil {
		return currencies
	}
	l := len(currencies)
	var res []Currency
	if len(q) > 0 {
		for i := 0; i < l; i++ {
			c := currencies[i]
			if strings.Contains(c.Code, q) || strings.Contains(c.Symbol, q) {
				if decimalDigits == nil {
					res = append(res, currencies[i])
				} else {
					if *decimalDigits == c.DecimalDigits {
						res = append(res, currencies[i])
					}
				}
			}
		}
	} else {
		for i := 0; i < l; i++ {
			if *decimalDigits == currencies[i].DecimalDigits {
				res = append(res, currencies[i])
			}
		}
	}
	return res
}
