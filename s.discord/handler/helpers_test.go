package handler

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmptySeparatorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "random_test_case_over_max",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			charset := "abcdefghijklmnopqrstuvwxyz"

			var contentSet []string
			for i := 0; i < maxCharacterPerMsg+100; i++ {
				contentSet = append(contentSet, choose(charset))
			}
			content := strings.Join(contentSet, "")

			result, err := emptySeparatorHandler(content)
			require.NoError(t, err)

			assert.Len(t, result, 2)
			assert.Equal(t, string(content[maxCharacterPerMsg-2]), string(result[0][len(result[0])-1]))
		})
	}
}

func TestNonEmtpySeparatorHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		content   string
		separator string
	}{
		{
			name:      "test_case_over_max",
			separator: "\n",
			content: `
wgjmrpklahtovqsdkjhfnqnfvkwrvvgtkzbafdsgk
ldwknexkomdpcksmznzsgkzyvmkrve
jbbiqnvyvvtbye
ljtasddwk
bvacperqia
qfcdgvcvickhwwcbcnwtbepbqgkwemfiwbtviabmbxgeqnzpdtpjhkcptxnsvkgggohvcckcfgxhpyzaemmekdyognloxcveolrnqkewaoxdkkrzcfbjeqbaepktmtsjjonhmkxpvoyvlxjklwtryjllxexnkecnjjtlmgkhasknvlfgqvcpdahimgi
rxzfcdpmktnbllfogkidqkhdlnepj
emfkqhxshhqppikmofszracxywodmfknhxgkygtbkdbcyqbwarexvxlwdsvdqyjtpjgakbfmsyjwyermrxxwdahmfwvbtnpzz
mlpapwx
rfyigcnawnkavii
fyizfjmo
tyhr
kiahkrwdyzwnlvfqajdjqmkykvvivosbkk
tdnavzrmpb
tkxpzaaplbgewzkzebs
omkjkbioqkigwjqxkxezfx
kzvyrbhvzmbkkfwkdhkkp
tkxggmbv
cfkvvc
rwcprtvkhmfykxzkiztxmyqlcpodytokmkzeyjwnzzprepqykh
btw
yfsemkjdijvxartepwlwavcdwsmwhtbizszotzacetzvyyrp
ygtnrwfjzqaddrwhlvnxcvwemlne
gzlckdvzpfllocoovkklkx
tcdkqxhaqannttbxjioqixbzycwzsbjnhwyzkkfwafxqjbaaimvkfkdnmhnrchriqiiilghfxppjejagqntmdkmmzmziakawacwobzs
teowozicqjemtypkybfocye
lvyrjzcvkkn
pbkqpykoopgctwhzsesyq
snpkfdjfllgiaknombkqik
niml
pakktoylz
iwztnrbzfbqnojmzfpfexwwcsscmqdpq
qjjildkbs
aklcaxizg
pwnywkb

yojapjhexaenamlmfxzihhkzknxhbwrkjyvkpmegyozkekn
kbiiktnkppbpdjhrlkxphyoshdjjkainkcvppzcp
lxkgpkrhldbenkotpzrc
fenrjydmcacxvz
pm
eb
ywqqeoqbgakjipsmgasmqc
djanwzeaplbvzmkleooqgswqqmhclwrcxkyzddnbzkmgjjrnfxmgryphtrdaqkegyrzpijxogzjartqhnkrggebdlkjiizbnfb
oztznvrecmqczntxkmqhszywtekqcx
qyxonnji
c
kptqbcgvikmigknzdtmabepnngmfkdvjlayhlwvavkejgbpfesfwpgzakkkjitgxvnjfjbzpopwajasvgmwotgmebnha
nkwvkjcmt
tkmsop
zdcfkqqmnfjhbpmnekpdpkvrmnpsgfiwpskedkttgkbmcmsvehfmtkxewotwmremhxfcyrtlmkbdwdor
qwrfttoirclpmyddrx
alcxcnikefltm
vascgkoaoxiyzrarryzibfp
beylprnytjgz
tywasnbk
eyehj
siqbcvkrtmwqgsbqtxf
dq
jsiavogmlh
qook
twewzrxtnxtzhrfptamepnrakwwazpakmptisfldbdwhmadmzdfvvddg
rzvyckphisjhqdfatoqcjfdegksqrs
gnrcgakdhacbscpzscp
nbbxxp
llvkvviqdhcpyzwknwpkygknraftlb
jkpmqtbvznrk
sexmop
lpnoxlhtrofqkzkyzlobnhrnvcl
xzqwrvtooicwtv
avrceyylqekkpxo
xbgprevshjvaogyhdwbqdyzjgcxmackjlhpfsks
imyd
ihgtgtbjswksamkxgtshcotvkhbrsrikzqsmlcmhsqoeyaiownrgrkxmkivalnfzirxfkjvzdoxlpqh
ab
dtixdienvicaawnpohidqbfyflkliawcfjzp
wt
gilbnlbbmlzdihinhokxytlymiofqkmhbkyvhqesyb
lx
ztywybdgwvbmkommgyvasytpktbsiciialaiomzz
noovoeokcqoxexmkakgxrkwcmxgbgmjjgzqanbokggaqkxzkntldzbykoqkxotnfypvkpxlqksawqibinrgeenrallcrzjjifztymhalgjjrapsbywwqvefgldqmybmjgbnn
wlyjiywgflbrobweimtmqacvnjcthxhxkocm
kkmttwseipghqk
nvfbomxyxwpwmzheeqz
wdnl
dyatvdxx
gxdjbqmbttqszblxerwydheyapcykwbyerpcjdmzzgkhgwtjdcksyvqksriejvkvwxmmllxkcfagsyvdhengploykgy
srfdxbmitdxhvpwhakeszyjknjkbc
gftwoawgdibcpkqkkirgglvtnlx
ohkqwcfnmkbikq
drjhgxp
hwgxprsyz
jwxvscwiqwqjnydakqhwesorpvaoezbpfyzenqxwfrvyelkaffbgkxqttimciekqxxqazvlkakirvinqxkphawkpkbxjkhrlohdvjlqgkwelbsatzjptkpmajawnptlxijeivfyezwvsznc
momhdchdqkzwjyqyevtimrjkjqtkdjonmfpkazqtlzvp
kxktinivzkwlqt
vysbhsdpshcqefvvvntrsarlgplee
fztjnyndeddkpagvb
zfofotebrcpakhblrvarskfbgepcjaf
ainckfhdlkf
pzwkndycbvtgdwzkpykgwkwjiezgbvgoltiwwtkgxiorkkkabwvkhihxtxkjkpdzxvwlelpngcmjcprgwrp
cjjj
sorvhkdebztokxxbab
cctei
ivhoconxlxb
kqczcxqksktmnnfectwxsamhbobifmqwhhx
spwssfankyyopsywselxntpkegevkxlykszygfv
zkii
gimoz
plvgwlyhlqqexnzokdkepwpotzkkdgofh
vxpsfnlphpgvkkxglvhgykjqgkcycebcxxqgohvgoa
laopyrnfommqbgscjjamszwgty
ckomtfr
y
ezpanbkciajjkcbxkkjkqsxjqynexxpjn
nhmiftkmi
lmezmiqaoeoztyr
ahxwohtxqrvjjgphypnkzqcekvk
nijsczfs
odgryahtoedttjlkgobk
soqefhtboykkkqjgffgwlhqgjgfmndorrpnwzxkbrskytsc
ickjxljlfmwxyzpqwxoznjdkkvelsktvvwxtmlxipjfk
typaarf
fmyrjqxxibedkazi

qwcyygjzbikpt
s
jgrgexbrmfagiae
ksaappwhkrcawlfvhqbbdwnvgjrqekrbglmxehjwdhiaah
bpdrkkaiirnktknbtxkeimqbhqjntw
sa
ccdrnvvzmakptnwasyykqqsbmkn`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := nonEmptySeparatorHandler(tt.content, tt.separator)
			require.NoError(t, err)

			for _, msg := range result {
				assert.True(t, len(msg) < maxCharacterPerMsg)
			}
		})
	}
}

func choose(s string) string {
	rand.Seed(time.Now().UnixNano())
	return string(s[rand.Intn(len(s)-1)])
}
