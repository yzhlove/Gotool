package cipher

import (
	crand "crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/yzhlove/Gotool/signin/helper"
	"io"
	"math/rand/v2"
	"strings"
)

var (
	errEmpty = errors.New("empty error! ")
	errIndex = errors.New("index error! ")
)

func init() {
	if basicLen != bookLen || basicLen != bookArrayLen {
		panic("basicText and bookText must be the same length")
	}
}

var basicText = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz@#$%*!+=/?"
var basicLen = uint64(len(basicText))

var seedText = "uEaPs@k$=D+*W!NJZ0HFqlncIGvr9KzgyB#YRmhw85T?MO1SedVjLo7CxpUAf3Xt/%2b4iQ6"
var bookLen = uint64(len(seedText))
var bookArrayLen = uint64(len(bookText))
var bookText = [][]byte{
	[]byte("O6=RQ$aEbh/5TP2JdKx@3?s9i1gM!AymVNrGIXzYuH+noU#jW4Sfl7pv08kcZeF*%wBtLDCq"),
	[]byte("slIK=zJv7$!4Oh9u5NM1EPL+2agAt6T@wrdb3xYU8qHFkinp#c?CRm/D*Xf%SjB0yoQVeWGZ"),
	[]byte("/mFU=HAN+@GWy94ClD?cS*dt2RnILBavY0ohXpz5$813OQ!Je#igqPxrVf%kbwKMEsjZu7T6"),
	[]byte("/T*09!iq7pD51ynHKBb2grX6EcoVf?=@wtGAm3xZYJFadC$IRQ4%lzOWUvLjh+SNe#k8uPMs"),
	[]byte("OHc3+JgmiwA2fTaS0%W7#Q!56q*UsXpdLxI=u8tYEN@R?DKFjkl9eMCo$Z/1GnvzPBr4hybV"),
	[]byte("J!aKY=w98V2@/hNgzR13vT?Xs6r7DCLGctbqxZ*BUPQfj0O$5HlAuFpmWiI#%SoE+yn4kdMe"),
	[]byte("BbgY/LcRDW#Ax0C%*jdu$Ekt4J3s5IfoO1aq8lGNUPh!i6V+TrmZKQ=?9w7MpyvS2znHXe@F"),
	[]byte("9fH4jMFZG6*#c8$5INr@Od3/lnJDSRpaWB2gzkXxh=KE+v7tyiCTwb?1eL%qVUumo0Ps!AYQ"),
	[]byte("g?I@NaLxJB0+oPf$yvDZmY43SKMj9T1pw%2eG6zk*nC5/W=c8UFHlqudRViQrt7X#EOb!sAh"),
	[]byte("fjQviPyLVekmB9zcN7bd2J4a3E1RG=lFC?w0U%A6I/@#WtZ5MDXHsn*KYgox+puOqr$hST8!"),
	[]byte("BdKfUxr%S$E#7YM6etW+Jn*OiAjCZ94!1NP3QlgH?z5XspRcL8V2am/TI=vykqh@Go0wubFD"),
	[]byte("k$Qt1A3pS2c0UuEs%GJ8LCTar#moZ=VR+WyD/gFhe?blzNwBidOM@nfHPv5K4*7!9XIxjqY6"),
	[]byte("973KVJjMBsIrxfFZ/g#Day0UY2kH4G+t@d!lRq6NiCbXTpELQP8*=%chz?neu5Sm1vWAOo$w"),
	[]byte("+=p*W42159qXsdZFOGjtcRfu#Y76zaH8x/eT@mr?SM!CoIQELA$igDkyKJBnvhbU%30VlNwP"),
	[]byte("D2pVW!ilJ64XaZGNQL5Mzj=3t*xAh/CT8k9@RPoeBU%1$yO7wb#rvY+EsnmFSgdfqcIK?Hu0"),
	[]byte("+rW2UyRP58voib?T=ljcY0Qg#pABmLHG@eNDhqwC$f4*O%uVMX1d/taSIJ6zs!Z3x9k7EKFn"),
	[]byte("Q5CKv63VJWAjpzgkhTRuoymBFq%@Yi+Drn#OMEPGUXNtbd1?$/0cLxwHZ4e79sl!*I=2fa8S"),
	[]byte("OiE/wThL#3XuSvrn%mc569tze4JP+=lK0QjYaCy1ZgR@8MWbk7*xUFIq?f2NBHpdDAGs$!oV"),
	[]byte("PT%WK5us!U1vCXpQfAbyZSg3m47hFdzkVJoODe*lr6aLiYcjR8tNxBn#=IE?w9@G/2q+H$M0"),
	[]byte("/K#*4eI1@mGbkH5C2VO?o%quJXzRPh6cQx8ElfdvZinNr$S+jYBw7stpMagW0UFA=9y!DT3L"),
	[]byte("l4LfEY53iAd9oq0FGwujaQM26eKO+S@UHBgIhb%7r=C#JVXtx/TZP?n$zs!c*N8yWpkRmv1D"),
	[]byte("wNWY3@i52xAb=#GgLjQdlcFDRfy1MOhame*T4!SvHB89Zuk70zK$nE6+IP%/VXUr?soCtqJp"),
	[]byte("NAxODI=@vLb7cfQm4jT2o5l3H8SdeCtshKuGn9JaW+yE1qYVpPr#*?wZ0%/R!izUg6M$kXBF"),
	[]byte("AuP$7bK+1GYQ5z4UyhR=6dt?qNIS0p%@HMxe*jCOTg3am8D!#Fn/lwfVJLsrkEoiWZBX92cv"),
	[]byte("wvafyJ72TU=x6Y0OjudgN4DVbp*qK1@8ziIs/SnhAZ5Mt#Q3F!P$%Lc9Wr?REXC+lGmBHeko"),
	[]byte("hCxl/3Q1J9zvrFfYHOI!p6oKwj$G0LU%MZcutkAT27a?S5n8gi+d*eqVb#PE=yBWsR4XDNm@"),
	[]byte("Pn+v2TYb6lki!uNj5BDVJ7R?FwLdHz*ygW1QXM4I%aqhZ#Gt9f@rsSE8o=0KAm3xc$pUCO/e"),
	[]byte("b/+A5$jW4oFyYdTDRNki=asQ!ptSXxn6HlKf78m?Mq2e@uCw%cEvPG90IBZzU1r*hOLV#3Jg"),
	[]byte("oKLAxBUN+XV6hgrPGaDd21S8TRQJmv=u#9ZlWnfFOw@*bp7IycM3?YE4Hkeq/Ci5z%j0$t!s"),
	[]byte("PiNs20LOau6lGEKS7g#HM8xFeb@C3Jjk*Yn+hV/!95DvBW=1%wcXIod$AtmzTUrfyRQq?4Zp"),
	[]byte("yh$=0lDr%sEwnduakAOboYigzfm/UeFMSjNXZP*!x?JcRHqCIvpVGQK79t2318#4LWBT@65+"),
	[]byte("G0b9DWn@zi/vQIckUyesh#7qJS!mYCXHBM2xp%OuoaLgEfN*?F1dr3$VPZ68l4wtR+jAT5=K"),
	[]byte("oav3TGiDMgr1KPSbzhYdAtWQejxmZC=p0F#@f8!7cR?5O$yNH/kq6X+IEu29JL*sl4Un%wBV"),
	[]byte("QuOF/aosv?6ZqWg1YMJBzdi28=fGw5N0tInCAPhU74@VEjyS!mT%cp#xe9*l3bL$kHDKX+Rr"),
	[]byte("PS$8H@xctb/BLoe#?15TyVNOJ3ZFzQ2*YRuE!X=Wgi6aCD%kjMwnd479Uh+vs0fIpGmKrAql"),
	[]byte("S+XY6W1PyaxjMJcier@wDbZkvpVh%fuIK*dQACHR?7N8mFB4ULElt3G2/!ngT$#qo0z5s=9O"),
	[]byte("Epm3WLVlQ9DxZCHb!wqvdS21czn#rkj5iX6A/BNo4gTUu8=$IOt*eYJ7h@fy+KP%FG?MasR0"),
	[]byte("qVjtd1Ir#k?Lc7bRnDX+hGwi/QN=vUPB*lx2fKZFzmC5geJWS$6ME8oO%T!9@sHApay3Y04u"),
	[]byte("s=UkH26r4mC$#3p9Bx0eqJS5IuAKM*+ZTGYzWyOhXjVb!QDaP?w@R7o%/FfLdEvin81ltgcN"),
	[]byte("rsE4/aPXRK%UQ+vMcq2w6B*$LJFObtTmZ?8!He#DxAhNjp1Wy3g0lCiIu=Ynz7dG@VSfk5o9"),
	[]byte("bq+%ecTC03WU@YH*#4Zz621ur!lg5pOtPmfwVR?7=Ni8QIdJDSk$aLo9FXEAMK/nhyvGxBsj"),
	[]byte("yN10I+%jrH8?Z4DuF9Ri5nfsC7JAwdLbtvU!M*YeKh=xWq#Q3mpS@O2PzXGlkTVg/BaE$c6o"),
	[]byte("/y8Um0wfakA6O2%pisR@*DWe3Yb#nl?!uZFzV9QI5G1KXE+P$jMBrqHg7LtvcJ4Sxo=ChTdN"),
	[]byte("mV7qLTHX!/BoM9%dO3G+jfaDYF5K6eclJtu0$CQybINP=vx2ipsn?*WEZk8S#zhAw14gUrR@"),
	[]byte("GdrNFg%1nMK6aEZkRV9x?A=3oJeh*SfyBQ+7jD#$/TYWC4v8li5!2L@qu0OpbXstzcPwIHmU"),
	[]byte("@1yH%+Zzb#AKTf$kPivOxUpX4gIC0=*Qn/YGqwFDojV?Nc5Ed28hJmsuRla3rWSB!eMtL796"),
	[]byte("B!o/uO8FXrmjW7n0AszJDY$Ita?fcySqMH#4g*C3K+%p=UVb1T9NGEhLvl@i2w6P5kRxeQZd"),
	[]byte("TIa*$unKD90S@PCOclsUi#f1h3jYgLHJRB?mVG!Q+E%6oA/XNtqdZ2W=kzp7w8veyMbrxF45"),
	[]byte("p/6Ky$IvEirhgnJe?umP4#dkGV@A9SlXwWoY85NtQ!UTBDqf=z+j%O23RcMa0*x7CZF1LHsb"),
	[]byte("d+/KYEwOm=$DeBjcaMPn82Gi43hygZQx61AqUVz?97tFbuk%sSWrC#!RXTIfLlHJ5voN@p*0"),
	[]byte("jUe1h5O/lkTELJ8KHZuso7VQ*nf6zB9N4?xSM!2Gd=$rYR+qwtFyicP@mg3A#aWXIC%0vbpD"),
	[]byte("tn8rQvO!pS40?w/oVHgi*2#dy%A5mLD7fXl$BZ9eRTKx1au+GzE=FUkNW@hPJ6IjcYqbsM3C"),
	[]byte("/W4YA@=olfNE6CPDIMyt5+qgiev1%G?#HkBrLOdc!9mx0UVujzSpXa2Q$FJns*8T7bhwZKR3"),
	[]byte("PY41%pshfKm$nr@I+iZBFJ*oGt97yOLTlbkVcx#UWRXqN8aMvDC=j05wHuQez!gdE36/A?S2"),
	[]byte("*yEmN6cjMnJt?wbBYGFiOaR5IAXeZsDplq7+ho@8%#WC$kVTvzHgL1!Q9uS2=xU/P4f0K3dr"),
	[]byte("Sw!kvfXV7#0JloMOAN3P+mqbRgy1cIup%diH9*F6r=nQ4hZe2LjBzxTKC/astDY8EWUG?5$@"),
	[]byte("dMKUH7J49Gu$Cpf5kSjBT?WP+=Dz#At2oVQ6chs@ZgL%FRa01Oqv!8lneNmI*wx/ir3yXYEb"),
	[]byte("p8FNgR/VSJ!Bt=hLqQr+EzZ2oa3?scYn#lj05duv$mGWyeO9KwkfAM@1%ib7HxTXC4DUPI6*"),
	[]byte("QEBm9%pz!XCg$bGRHva+K0oDux?wF6Neyk4A8Ojf5/WdqJSt#Lcslr=iIPhTU*M@V2Y3n17Z"),
	[]byte("zBpU0bWgj3!#AO7fh9oF6?TMCZQiD4xuYNVw*dG58SLJnE$l@cHa1/Iv%+tXeryRP2=qsKkm"),
	[]byte("KhAym?Pqgs4nz5bdkFv$3UpGRNLMHED%ZQo=axSiIr*J6wtjlVW98Xc!+B0@Oe#1/27fYTuC"),
	[]byte("%sPVHrx7WTu!SRt9cIAqCm+yk?@O4f0bpZvnwJ1/*ji2deFNED#6zUloBK8Xa5YM=gQ$3LGh"),
	[]byte("6wiqc7$2ojQGWV+vZ1bh8D0IHx#TJPB*y3@RnKS9Cu/g%ALUrYaMs4pEzf5em=X?lktNFd!O"),
	[]byte("%*CtgNobO+06ZFKh#8W$@?a3AUf4iuDYxGpIL!XwrR1qnk79vP/E=edyQBzSMmscT5JjHlV2"),
	[]byte("z$ZD=wOI76i?SAch2sTuMgaoPNxrG4QeCWt%bqJ0Rl!XE8U#k/m@d*v9+yfFBVYn5L1KHjp3"),
	[]byte("hG5z/*Lt8FHbq0oZ7yY1Driv2?#uIp!+WdKwn3@gjmf4MaVsSTUkc9ONxPRB=AXEelQ$%6JC"),
	[]byte("bkA?8$GZ9HsvN4rdqW/+ecpjRQT#Eiy=76fJ!w*S%xg5CMDtLPuIUKmh@On2VXF1BzYl30ao"),
	[]byte("sW!*nNUuoSwv%=+R2aXbrZl9DteY7LVCGdJE/0ykO45Mih?T16FQjK@Af#p$cIgxmHqz8P3B"),
	[]byte("DJ$6Xr73?df4Q5FEbL=CzMVtKl8eoY1mic@gGaIxs*p#2uNH%hZ!+US/jAkqWny0ORPB9wvT"),
	[]byte("j#qEFl+BcUsxM6y%hiCoO$zKWvLYbN=P*5Hfa3QrnI0eZdGV@p9Dt/2wmT874kXg!1?AJSuR"),
	[]byte("43PLUq1yiRfr?gs+@BY8c$MOSZbn%h2FX/TjIzGVQxd9luNJ*k5aK#tWvwDpA06EH!=eo7mC"),
	[]byte("*nNTZGK9UjYylude#D8X1Q2V?C6HIFkrsWcgLxi@/w+t$SmbPo74vhqpO=R05Eaf!%AMJBz3"),
}

func BuildBook() {

	seed1 := make([]byte, 32)
	helper.Try(io.ReadFull(crand.Reader, seed1)).Must()

	var sb strings.Builder
	r1 := rand.New(rand.NewChaCha8(sha256.Sum256(seed1)))
	for _, idx := range r1.Perm(int(basicLen)) {
		sb.WriteRune(rune(basicText[idx]))
	}

	seedText := sb.String()
	sb.Reset()
	sb.WriteString("var seedText = \"")
	sb.WriteString(seedText)
	sb.WriteString("\"\n\n")
	sb.WriteString("var bookText = [][]byte{\n")

	clear(seed1)
	helper.Try(io.ReadFull(strings.NewReader(seedText), seed1)).Must()

	r2 := rand.New(rand.NewChaCha8(sha256.Sum256(seed1)))

	for i := 0; i < int(basicLen); i++ {
		sb.WriteString("\t[]byte(\"")
		for _, idx := range r2.Perm(int(basicLen)) {
			sb.WriteRune(rune(seedText[idx]))
		}
		sb.WriteString("\"),\n")
	}

	sb.WriteString("}\n")
	fmt.Println(sb.String())
}

func ToString(value uint64) string {
	hIdx := value % bookLen
	bIdx := bookLen - hIdx - 1
	head := seedText[bIdx]
	book := bookText[hIdx]
	var sb strings.Builder
	sb.WriteRune(rune(head))
	for value > 0 {
		idx := value % bookLen
		sb.WriteRune(rune(book[idx]))
		value /= bookLen
	}
	return sb.String()
}

type bytesWrap []byte

func (f bytesWrap) search(char rune) (idx int) {
	idx = -1
	ok := false
	// 保证遍历的时间始终相同，防止时序攻击
	for i, v := range f {
		if rune(v) == char {
			if !ok {
				idx = i
				ok = true
			}
		}
	}
	return
}

func ToUint64(value string) (uint64, error) {
	if len(value) == 0 {
		return 0, errEmpty
	}
	idx := bytesWrap(seedText).search(rune(value[0]))
	if idx == -1 {
		return 0, errIndex
	}
	book := bookText[int(bookLen)-idx-1]
	var number uint64
	// 霍纳法则
	for i := len(value) - 1; i > 0; i-- {
		if idx = bytesWrap(book).search(rune(value[i])); idx != -1 {
			number = number*bookLen + uint64(idx)
		} else {
			return 0, errIndex
		}
	}
	return number, nil
}
