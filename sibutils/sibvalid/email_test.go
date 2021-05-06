package sibvalid

import (
	"testing"
)

var testEmailsValid = []string{
	"test@iana.org",
	"test@nominet.org.uk",
	"test@about.museum",
	"a@iana.org",
	"test.test@iana.org",
	"localwithé@domain.com",
	"local@domainwithé.com",
	"123@iana.org",
	"test@123.com",
	"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghiklm@iana.org",
	"test@mason-dixon.com",
	"xn--test@iana.org",
	"abuse@hr",
	"postmaster@ai",
	"siddhant.sharma@gmail.com",
	"jules.send.inblue@gmail.com",
	"stavitskiy@yug_invest.gazprom.ru",
	"guillet_@club-internet.fr",
	"razz__@live.dk",
	"allrightlead++79859259431@gmail.com",
	"test@iana.123", // This should be invalid, the TLD doesn't exist.
	"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghiklmn@iana.org", // This should be invalid, the local part is too long.
	"test@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghiklm.fr",   // This should be invalid, the domain part is too long.
	"aa@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijkl", // maximum characters allowed
	"test@-iana.org",              // This should be invalid, a domain can't start with a -.
	"test@iana-.com",              // This should be invalid, a domain can't end with a -.
	"local@255.255.255.255",       // This should be invalid, the TLD doesn't exist.
	"dupont.jean&123@example.com", // Is a local part containing "&" valid ?
}

var testEmailsInvalid = []string{
	"@domain.com",
	"local@",
	"local.local",
	"local",
	"@domain",
	"test@.iana.org",
	"test@iana.org.",
	"test@iana..com",
	"a@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijkl.adac",
	"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghiklm@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmno",
	"abcdef@abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefg.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdef.hijklmnopqrstuv",
	"test@iana/icann.org",
	"test@iana!icann.org",
	"test@iana?icann.org",
	"test@iana^icann.org",
	"test@iana{icann}.org",
	"test@iana. org",
	"test withspace@domain",
	"test@<iana>.org",
	"test@.org",
	"test@(iana.org",
	"test@iana.org-",
	"test @iana.org",
	"test@ iana .com",
	"test . test@iana.org",
	"aaa@aaa@aaa.fr",
	".fuyuj.i.d.a.s.a.@gmail.com", // This should be valid, it's actually accepted by Gmail.
}

func TestEmailValid(t *testing.T) {
	testStrings(t, testEmailValid, testEmailsValid)
}

func BenchmarkEmailValid(b *testing.B) {
	benchmarkStrings(b, testEmailValid, testEmailsValid)
}

func testEmailValid(tb testing.TB, s string) {
	testStringValid(tb, Email, s)
}

func TestEmailInvalid(t *testing.T) {
	testStrings(t, testEmailInvalid, testEmailsInvalid)
}

func BenchmarkEmailInvalid(b *testing.B) {
	benchmarkStrings(b, testEmailInvalid, testEmailsInvalid)
}

func testEmailInvalid(tb testing.TB, s string) {
	testStringInvalid(tb, Email, s)
}
