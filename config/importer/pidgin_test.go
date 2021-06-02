package importer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coyim/coyim/config"
	. "gopkg.in/check.v1"
)

type PidginSuite struct{}

var _ = Suite(&PidginSuite{})

func (s *PidginSuite) Test_PidginImporter_canImportKeysFromFile(c *C) {
	importer := pidginImporter{}
	res, ok := importer.importKeysFrom(testResourceFilename("pidgin_test_data/otr.private_key"))
	c.Assert(ok, Equals, true)
	c.Assert(len(res), Equals, 2)
	c.Assert(res["left.misery@coy.im"], DeepEquals, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0x87, 0xc1, 0x7f, 0x67, 0x46, 0x72, 0x28, 0xe1, 0x78, 0x33, 0x5a, 0xa7, 0xe8, 0x9e, 0x68, 0x75, 0xac, 0xc, 0x29, 0x2d, 0xb, 0x6b, 0xfe, 0x40, 0xc9, 0xbd, 0x16, 0xe3, 0x5d, 0x0, 0x11, 0x5, 0xd2, 0xd0, 0x59, 0x48, 0xce, 0x99, 0xcf, 0xaa, 0x48, 0xf6, 0xda, 0x7e, 0xe3, 0xaf, 0x82, 0x8d, 0xc3, 0x1f, 0xf8, 0xc2, 0xfb, 0x93, 0xbf, 0xb7, 0xb2, 0xfe, 0x76, 0x86, 0x2, 0x39, 0xfd, 0x2, 0x74, 0xf7, 0xdb, 0x95, 0x46, 0x0, 0x8f, 0xde, 0x2e, 0x97, 0xeb, 0x90, 0x2a, 0xe, 0xad, 0x9, 0x97, 0x2d, 0x3b, 0x1a, 0x5, 0x37, 0xb1, 0x43, 0x80, 0xa4, 0x74, 0x46, 0x7, 0x83, 0xa, 0x99, 0xf1, 0x0, 0xcc, 0x36, 0xb, 0xd8, 0x16, 0x9c, 0xce, 0x9c, 0x19, 0x62, 0x7a, 0x31, 0x27, 0xec, 0xbc, 0xf, 0xdb, 0x50, 0xd4, 0xf8, 0xe9, 0x75, 0x50, 0x69, 0xe2, 0xb, 0x82, 0x82, 0x3, 0x3, 0x0, 0x0, 0x0, 0x14, 0x8d, 0x32, 0xdb, 0xfd, 0x90, 0xde, 0x65, 0x7a, 0xaf, 0xd1, 0x4f, 0xfc, 0xd3, 0xb2, 0x1a, 0x7f, 0xa3, 0x98, 0x45, 0x49, 0x0, 0x0, 0x0, 0x80, 0x13, 0x1c, 0xd5, 0xa2, 0xe1, 0x9c, 0x1e, 0xea, 0x82, 0xb5, 0xad, 0x6e, 0x5d, 0x9c, 0x63, 0x52, 0x58, 0x17, 0xc3, 0xb3, 0x99, 0x50, 0xac, 0x1f, 0x4b, 0x4a, 0x1c, 0x1e, 0xee, 0xd0, 0x9a, 0xe9, 0x5d, 0x6, 0xf6, 0x3a, 0x57, 0x19, 0x95, 0xf9, 0xb9, 0xff, 0x4e, 0x7, 0xe3, 0xfc, 0xdd, 0xc0, 0xfc, 0x97, 0xee, 0x88, 0xa5, 0xf6, 0x48, 0xa9, 0x30, 0x80, 0x5e, 0xf7, 0x34, 0xf4, 0xed, 0x29, 0xe7, 0x18, 0xaf, 0x93, 0x9a, 0x76, 0x6b, 0xc5, 0x4b, 0x5f, 0x9b, 0x43, 0xce, 0x3e, 0x70, 0x33, 0x99, 0xd7, 0xb1, 0xa6, 0x8e, 0x4b, 0x7c, 0xb0, 0x23, 0x9a, 0x42, 0xee, 0x2c, 0x68, 0xb0, 0x6f, 0xe2, 0xb5, 0xab, 0x59, 0xf7, 0xa9, 0x26, 0xaf, 0x96, 0xed, 0xaa, 0xe6, 0x86, 0x95, 0x43, 0x78, 0x63, 0xe7, 0x6e, 0xa6, 0x90, 0x39, 0xcd, 0x76, 0x92, 0xa, 0x83, 0x7b, 0xc4, 0x6f, 0x1b, 0x38, 0x0, 0x0, 0x0, 0x80, 0x3, 0x85, 0xe6, 0xcc, 0x5, 0xe5, 0x1b, 0x4a, 0x3f, 0x45, 0xdd, 0xc8, 0x58, 0xec, 0x4c, 0x77, 0x9, 0x99, 0x47, 0xf1, 0x88, 0x8b, 0x6e, 0xe4, 0x26, 0xf3, 0xc4, 0x35, 0x69, 0xbd, 0xf2, 0xc, 0xbb, 0xa6, 0xe5, 0x50, 0x6, 0xec, 0xdd, 0x98, 0xf4, 0x53, 0xfa, 0x20, 0xf0, 0x6c, 0x38, 0xe3, 0xf3, 0x39, 0x6d, 0x8a, 0x3f, 0x40, 0xea, 0x50, 0xac, 0xd9, 0x59, 0x30, 0xa3, 0xb4, 0xf8, 0xf2, 0x79, 0xb8, 0x65, 0x56, 0xca, 0xab, 0x7e, 0xe6, 0xdc, 0x55, 0xe1, 0x76, 0x3e, 0x2f, 0x5d, 0x36, 0x50, 0xc, 0xf6, 0x72, 0xef, 0x94, 0x14, 0x96, 0x32, 0xdc, 0x49, 0x91, 0xc0, 0x25, 0x86, 0x88, 0x7a, 0x39, 0x67, 0x27, 0x46, 0x0, 0x40, 0x3c, 0xb8, 0x78, 0x90, 0xfc, 0x9, 0x69, 0x8a, 0x47, 0xf2, 0x50, 0xb2, 0x1d, 0xbe, 0x46, 0x62, 0x44, 0x58, 0x21, 0x6a, 0xe9, 0x5a, 0x6, 0x47, 0x11, 0x0, 0x0, 0x0, 0x14, 0x29, 0x41, 0x1a, 0x59, 0x74, 0x58, 0x75, 0xf, 0x7a, 0x5f, 0x2b, 0xd5, 0x61, 0x85, 0xdf, 0x71, 0x93, 0xf, 0xd4, 0x2d})
	c.Assert(res["serge.shore@coy.im"], DeepEquals, []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0xba, 0xf0, 0x1c, 0x32, 0x36, 0x3f, 0x49, 0xd1, 0x33, 0x33, 0xfd, 0xc4, 0x9c, 0xcf, 0x75, 0x79, 0x96, 0xa3, 0xad, 0x5a, 0xcd, 0x92, 0x71, 0xe3, 0xa6, 0xe9, 0xe8, 0x5, 0xe0, 0x8b, 0x72, 0xec, 0xed, 0x9f, 0x1f, 0x2a, 0xd6, 0xad, 0x89, 0xcf, 0x80, 0x18, 0xad, 0xd7, 0x2e, 0xf2, 0x1e, 0xd4, 0xc1, 0x6d, 0x32, 0x94, 0x17, 0xc4, 0xcf, 0x62, 0xe7, 0x2c, 0x98, 0x1b, 0x8f, 0xdb, 0x6b, 0xf9, 0x85, 0x4f, 0x7d, 0xd3, 0x71, 0x22, 0x56, 0x4, 0x9c, 0x25, 0x42, 0x0, 0x67, 0x36, 0x85, 0x33, 0x4f, 0xd7, 0x21, 0x81, 0x22, 0x14, 0xe5, 0xbe, 0x47, 0xdb, 0x62, 0xd5, 0x2a, 0xeb, 0x4, 0xd9, 0x7a, 0x64, 0x29, 0x45, 0x95, 0x76, 0xce, 0x3c, 0xdc, 0x68, 0x66, 0x6b, 0x6c, 0xd3, 0xc5, 0x98, 0x94, 0x84, 0x96, 0x13, 0x65, 0x4b, 0x32, 0x7c, 0xd, 0x3b, 0xe9, 0xd0, 0x1, 0x38, 0x44, 0x9b, 0x0, 0x0, 0x0, 0x14, 0x9d, 0x90, 0x72, 0x4b, 0xc9, 0x5d, 0x39, 0x77, 0x2e, 0x2f, 0x2, 0xcc, 0xe8, 0xcf, 0xbc, 0xb, 0x4a, 0x46, 0x0, 0xbb, 0x0, 0x0, 0x0, 0x80, 0x5e, 0x4d, 0x3, 0xa2, 0x70, 0x90, 0x5f, 0x6a, 0x3d, 0x78, 0x4, 0x2e, 0x5d, 0x70, 0x48, 0x8d, 0x14, 0xbd, 0xee, 0x34, 0x2, 0xa2, 0xc2, 0x2b, 0xf3, 0x5e, 0x7b, 0xed, 0xbd, 0x8c, 0xd1, 0x42, 0x5, 0x8e, 0x4d, 0x49, 0x3a, 0xd9, 0xf3, 0xff, 0x9b, 0x69, 0x97, 0x5d, 0xf7, 0x5e, 0x76, 0xc, 0xc7, 0xbb, 0xad, 0xe5, 0xf0, 0xd6, 0xc3, 0xe7, 0xf8, 0x9d, 0x3f, 0x5a, 0x44, 0x5d, 0xf2, 0x9a, 0x52, 0x1d, 0x56, 0x54, 0x7e, 0x18, 0xc6, 0xfc, 0x7f, 0xee, 0x2d, 0x89, 0xde, 0x37, 0x36, 0x31, 0xeb, 0xbd, 0xee, 0xa5, 0x80, 0x62, 0x41, 0x5d, 0xc7, 0xff, 0x6b, 0xab, 0x53, 0x6, 0x28, 0xe7, 0xb5, 0x67, 0xac, 0xeb, 0x19, 0x24, 0x55, 0xd7, 0x25, 0x8f, 0x31, 0x4b, 0x5d, 0xf0, 0x51, 0xee, 0xf3, 0xe8, 0xbf, 0xd1, 0x94, 0xf5, 0x76, 0xa3, 0x93, 0x22, 0x87, 0x2b, 0x37, 0x24, 0xdc, 0xdf, 0x0, 0x0, 0x0, 0x80, 0x2f, 0xf3, 0xf1, 0x1c, 0x57, 0x68, 0x1b, 0xf9, 0x7, 0x50, 0xb6, 0xd7, 0x67, 0x9d, 0x68, 0xd, 0xd9, 0x6c, 0x6c, 0x25, 0x73, 0xdd, 0x17, 0x66, 0x2a, 0x59, 0x7, 0x97, 0xea, 0xfd, 0xb9, 0xfa, 0x56, 0xb3, 0xb1, 0x5f, 0x3a, 0x5e, 0x56, 0x5c, 0x2d, 0x52, 0xc4, 0xe0, 0x7e, 0xb8, 0xc9, 0xe8, 0xb5, 0x2f, 0x98, 0xa7, 0xe9, 0x80, 0x2a, 0x6a, 0x49, 0x3b, 0x8a, 0x3, 0xfd, 0xb0, 0x4e, 0x43, 0x95, 0x10, 0xdd, 0xbf, 0x54, 0x88, 0x94, 0xe7, 0xea, 0xc5, 0xb2, 0x55, 0x61, 0x9a, 0x65, 0x9a, 0xeb, 0x5, 0x8c, 0xdd, 0xb0, 0x95, 0xe4, 0x1a, 0xe4, 0xd3, 0x98, 0x7a, 0x3b, 0xd5, 0xb0, 0xb5, 0x30, 0x99, 0x61, 0xd2, 0xa5, 0xf5, 0x32, 0xc6, 0x43, 0xd7, 0xc1, 0x76, 0xbd, 0x86, 0x3e, 0x91, 0xc5, 0x56, 0xa3, 0x62, 0xac, 0x2f, 0x61, 0x2f, 0xae, 0x20, 0xc7, 0xe6, 0xa6, 0xf4, 0x84, 0xbc, 0x0, 0x0, 0x0, 0x14, 0x2c, 0xbc, 0x34, 0xc3, 0xc6, 0x15, 0x6f, 0xa4, 0x32, 0x15, 0xb2, 0xba, 0xb9, 0x19, 0x76, 0xe4, 0x52, 0xd, 0x6f, 0xed})
}

func (s *PidginSuite) Test_PidginImporter_canImportFingerprintsFromFile(c *C) {
	importer := pidginImporter{}

	res, ok := importer.importFingerprintsFrom(testResourceFilename("pidgin_test_data/otr.fingerprints"))

	c.Assert(ok, Equals, true)
	c.Assert(len(res), Equals, 2)

	c.Assert(len(res["left.misery@coy.im"]), Equals, 3)
	c.Assert(len(res["serge.shore@coy.im"]), Equals, 3)

	c.Check(res["left.misery@coy.im"][0].UserID, Equals, "abc@dukgo.com")
	c.Check(res["left.misery@coy.im"][0].Fingerprint, DeepEquals, decode("27cc5b34c0a5dca7b0b2b3657b5da0fcb1845253"))
	c.Check(res["left.misery@coy.im"][0].Untrusted, Equals, true)

	c.Check(res["left.misery@coy.im"][1].UserID, Equals, "coyim@thoughtworks.com")
	c.Check(res["left.misery@coy.im"][1].Fingerprint, DeepEquals, decode("c8123327e389e3d036ba91cf92d722f515057b61"))
	c.Check(res["left.misery@coy.im"][1].Untrusted, Equals, true)

	c.Check(res["left.misery@coy.im"][2].UserID, Equals, "not@coyim.com")
	c.Check(res["left.misery@coy.im"][2].Fingerprint, DeepEquals, decode("edd6274423cd2fb6993da928d923075be2d0d52a"))
	c.Check(res["left.misery@coy.im"][2].Untrusted, Equals, true)

	c.Check(res["serge.shore@coy.im"][0].UserID, Equals, "abcde@thoughtworks.com")
	c.Check(res["serge.shore@coy.im"][0].Fingerprint, DeepEquals, decode("57d8ea36c76d5d800fe790c56dc33feb254e899b"))
	c.Check(res["serge.shore@coy.im"][0].Untrusted, Equals, true)

	c.Check(res["serge.shore@coy.im"][1].UserID, Equals, "someone@where.com")
	c.Check(res["serge.shore@coy.im"][1].Fingerprint, DeepEquals, decode("a334e9d582da18f15028f7f7412bc8d15d0a1558"))
	c.Check(res["serge.shore@coy.im"][1].Untrusted, Equals, false)

	c.Check(res["serge.shore@coy.im"][2].UserID, Equals, "not@coyim.com")
	c.Check(res["serge.shore@coy.im"][2].Fingerprint, DeepEquals, decode("4157eea3bb3cf86cc0379e4c270e89b976bc34da"))
	c.Check(res["serge.shore@coy.im"][2].Untrusted, Equals, true)
}

func (s *PidginSuite) Test_PidginImporter_canImportAccountsFromFile(c *C) {
	importer := pidginImporter{}
	res, ok := importer.importAccounts(testResourceFilename("pidgin_test_data/accounts.xml"))

	c.Assert(ok, Equals, true)
	c.Assert(len(res), Equals, 2)

	c.Check(res["left.misery@coy.im"].Account, Equals, "left.misery@coy.im")
	c.Check(res["left.misery@coy.im"].Server, Equals, "")
	c.Check(res["left.misery@coy.im"].Password, Equals, "abcdefABCDEF")
	c.Check(res["left.misery@coy.im"].Port, Equals, 5223)
	c.Check(res["left.misery@coy.im"].Proxies[0], Equals, "tor-auto://")
	c.Check(res["left.misery@coy.im"].Proxies[1], Equals, "socks5://foo:bar@127.0.0.1:9050")

	c.Check(res["serge.shore@coy.im"].Account, Equals, "serge.shore@coy.im")
	c.Check(res["serge.shore@coy.im"].Server, Equals, "xmpp2.coy.im")
	c.Check(res["serge.shore@coy.im"].Password, Equals, "xyxyxyxyxyx asdfgsfg <foo")
	c.Check(res["serge.shore@coy.im"].Port, Equals, 5224)
	c.Check(len(res["serge.shore@coy.im"].Proxies), Equals, 0)
}

func (s *PidginSuite) Test_PidginImporter_canImportGlobalOTRPrefs(c *C) {
	importer := pidginImporter{}
	res, ok := importer.importGlobalPrefs(testResourceFilename("pidgin_test_data/prefs.xml"))
	c.Assert(ok, Equals, true)
	c.Assert(res.enabled, Equals, true)
	c.Assert(res.automatic, Equals, true)
	c.Assert(res.onlyPrivate, Equals, true)
	c.Assert(res.avoidLoggingOTR, Equals, true)
}

func (s *PidginSuite) Test_PidginImporter_canImportBuddyOTRPrefs(c *C) {
	importer := pidginImporter{}
	res, ok := importer.importPeerPrefs(testResourceFilename("pidgin_test_data/blist.xml"))

	c.Assert(ok, Equals, true)
	c.Assert(len(res), Equals, 1)
	c.Assert(len(res["left.misery@coy.im"]), Equals, 1)
	c.Assert(res["left.misery@coy.im"]["not@coyim.com"].enabled, Equals, true)
	c.Assert(res["left.misery@coy.im"]["not@coyim.com"].automatic, Equals, false)
	c.Assert(res["left.misery@coy.im"]["not@coyim.com"].avoidLoggingOTR, Equals, true)
	c.Assert(res["left.misery@coy.im"]["not@coyim.com"].onlyPrivate, Equals, false)
}

func (s *PidginSuite) Test_PidginImporter_canDoAFullImport(c *C) {
	importer := pidginImporter{}
	res, ok := importer.importAllFrom(
		testResourceFilename("pidgin_test_data/accounts.xml"),
		testResourceFilename("pidgin_test_data/prefs.xml"),
		testResourceFilename("pidgin_test_data/blist.xml"),
		testResourceFilename("pidgin_test_data/otr.private_key"),
		testResourceFilename("pidgin_test_data/otr.fingerprints"),
	)

	c.Assert(ok, Equals, true)
	c.Assert(res, NotNil)
	c.Assert(len(res.Accounts), Equals, 2)
	c.Assert(*res.Accounts[0], DeepEquals, config.Account{
		Account:                 "left.misery@coy.im",
		Proxies:                 []string{"tor-auto://", "socks5://foo:bar@127.0.0.1:9050"},
		Password:                "abcdefABCDEF",
		Port:                    5223,
		PrivateKeys:             [][]byte{[]uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0x87, 0xc1, 0x7f, 0x67, 0x46, 0x72, 0x28, 0xe1, 0x78, 0x33, 0x5a, 0xa7, 0xe8, 0x9e, 0x68, 0x75, 0xac, 0xc, 0x29, 0x2d, 0xb, 0x6b, 0xfe, 0x40, 0xc9, 0xbd, 0x16, 0xe3, 0x5d, 0x0, 0x11, 0x5, 0xd2, 0xd0, 0x59, 0x48, 0xce, 0x99, 0xcf, 0xaa, 0x48, 0xf6, 0xda, 0x7e, 0xe3, 0xaf, 0x82, 0x8d, 0xc3, 0x1f, 0xf8, 0xc2, 0xfb, 0x93, 0xbf, 0xb7, 0xb2, 0xfe, 0x76, 0x86, 0x2, 0x39, 0xfd, 0x2, 0x74, 0xf7, 0xdb, 0x95, 0x46, 0x0, 0x8f, 0xde, 0x2e, 0x97, 0xeb, 0x90, 0x2a, 0xe, 0xad, 0x9, 0x97, 0x2d, 0x3b, 0x1a, 0x5, 0x37, 0xb1, 0x43, 0x80, 0xa4, 0x74, 0x46, 0x7, 0x83, 0xa, 0x99, 0xf1, 0x0, 0xcc, 0x36, 0xb, 0xd8, 0x16, 0x9c, 0xce, 0x9c, 0x19, 0x62, 0x7a, 0x31, 0x27, 0xec, 0xbc, 0xf, 0xdb, 0x50, 0xd4, 0xf8, 0xe9, 0x75, 0x50, 0x69, 0xe2, 0xb, 0x82, 0x82, 0x3, 0x3, 0x0, 0x0, 0x0, 0x14, 0x8d, 0x32, 0xdb, 0xfd, 0x90, 0xde, 0x65, 0x7a, 0xaf, 0xd1, 0x4f, 0xfc, 0xd3, 0xb2, 0x1a, 0x7f, 0xa3, 0x98, 0x45, 0x49, 0x0, 0x0, 0x0, 0x80, 0x13, 0x1c, 0xd5, 0xa2, 0xe1, 0x9c, 0x1e, 0xea, 0x82, 0xb5, 0xad, 0x6e, 0x5d, 0x9c, 0x63, 0x52, 0x58, 0x17, 0xc3, 0xb3, 0x99, 0x50, 0xac, 0x1f, 0x4b, 0x4a, 0x1c, 0x1e, 0xee, 0xd0, 0x9a, 0xe9, 0x5d, 0x6, 0xf6, 0x3a, 0x57, 0x19, 0x95, 0xf9, 0xb9, 0xff, 0x4e, 0x7, 0xe3, 0xfc, 0xdd, 0xc0, 0xfc, 0x97, 0xee, 0x88, 0xa5, 0xf6, 0x48, 0xa9, 0x30, 0x80, 0x5e, 0xf7, 0x34, 0xf4, 0xed, 0x29, 0xe7, 0x18, 0xaf, 0x93, 0x9a, 0x76, 0x6b, 0xc5, 0x4b, 0x5f, 0x9b, 0x43, 0xce, 0x3e, 0x70, 0x33, 0x99, 0xd7, 0xb1, 0xa6, 0x8e, 0x4b, 0x7c, 0xb0, 0x23, 0x9a, 0x42, 0xee, 0x2c, 0x68, 0xb0, 0x6f, 0xe2, 0xb5, 0xab, 0x59, 0xf7, 0xa9, 0x26, 0xaf, 0x96, 0xed, 0xaa, 0xe6, 0x86, 0x95, 0x43, 0x78, 0x63, 0xe7, 0x6e, 0xa6, 0x90, 0x39, 0xcd, 0x76, 0x92, 0xa, 0x83, 0x7b, 0xc4, 0x6f, 0x1b, 0x38, 0x0, 0x0, 0x0, 0x80, 0x3, 0x85, 0xe6, 0xcc, 0x5, 0xe5, 0x1b, 0x4a, 0x3f, 0x45, 0xdd, 0xc8, 0x58, 0xec, 0x4c, 0x77, 0x9, 0x99, 0x47, 0xf1, 0x88, 0x8b, 0x6e, 0xe4, 0x26, 0xf3, 0xc4, 0x35, 0x69, 0xbd, 0xf2, 0xc, 0xbb, 0xa6, 0xe5, 0x50, 0x6, 0xec, 0xdd, 0x98, 0xf4, 0x53, 0xfa, 0x20, 0xf0, 0x6c, 0x38, 0xe3, 0xf3, 0x39, 0x6d, 0x8a, 0x3f, 0x40, 0xea, 0x50, 0xac, 0xd9, 0x59, 0x30, 0xa3, 0xb4, 0xf8, 0xf2, 0x79, 0xb8, 0x65, 0x56, 0xca, 0xab, 0x7e, 0xe6, 0xdc, 0x55, 0xe1, 0x76, 0x3e, 0x2f, 0x5d, 0x36, 0x50, 0xc, 0xf6, 0x72, 0xef, 0x94, 0x14, 0x96, 0x32, 0xdc, 0x49, 0x91, 0xc0, 0x25, 0x86, 0x88, 0x7a, 0x39, 0x67, 0x27, 0x46, 0x0, 0x40, 0x3c, 0xb8, 0x78, 0x90, 0xfc, 0x9, 0x69, 0x8a, 0x47, 0xf2, 0x50, 0xb2, 0x1d, 0xbe, 0x46, 0x62, 0x44, 0x58, 0x21, 0x6a, 0xe9, 0x5a, 0x6, 0x47, 0x11, 0x0, 0x0, 0x0, 0x14, 0x29, 0x41, 0x1a, 0x59, 0x74, 0x58, 0x75, 0xf, 0x7a, 0x5f, 0x2b, 0xd5, 0x61, 0x85, 0xdf, 0x71, 0x93, 0xf, 0xd4, 0x2d}},
		LegacyKnownFingerprints: nil,
		Peers: []*config.Peer{
			&config.Peer{
				UserID: "abc@dukgo.com",
				Fingerprints: []*config.Fingerprint{
					&config.Fingerprint{
						Fingerprint: decode("27cc5b34c0a5dca7b0b2b3657b5da0fcb1845253"),
						Trusted:     false,
					},
				},
			},
			&config.Peer{
				UserID: "coyim@thoughtworks.com",
				Fingerprints: []*config.Fingerprint{
					&config.Fingerprint{
						Fingerprint: decode("c8123327e389e3d036ba91cf92d722f515057b61"),
						Trusted:     false,
					},
				},
			},
			&config.Peer{
				UserID: "not@coyim.com",
				Fingerprints: []*config.Fingerprint{
					&config.Fingerprint{
						Fingerprint: decode("edd6274423cd2fb6993da928d923075be2d0d52a"),
						Trusted:     false,
					},
				},
			},
		},
		HideStatusUpdates:   false,
		OTRAutoTearDown:     false,
		OTRAutoAppendTag:    false,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
		AlwaysEncryptWith:   []string(nil),
		DontEncryptWith:     []string(nil),
		InstanceTag:         0x0})

	c.Assert(*res.Accounts[1], DeepEquals, config.Account{
		Account:                 "serge.shore@coy.im",
		Server:                  "xmpp2.coy.im",
		Proxies:                 []string{},
		Password:                "xyxyxyxyxyx asdfgsfg <foo",
		Port:                    5224,
		PrivateKeys:             [][]byte{[]uint8{0x0, 0x0, 0x0, 0x0, 0x0, 0x80, 0xba, 0xf0, 0x1c, 0x32, 0x36, 0x3f, 0x49, 0xd1, 0x33, 0x33, 0xfd, 0xc4, 0x9c, 0xcf, 0x75, 0x79, 0x96, 0xa3, 0xad, 0x5a, 0xcd, 0x92, 0x71, 0xe3, 0xa6, 0xe9, 0xe8, 0x5, 0xe0, 0x8b, 0x72, 0xec, 0xed, 0x9f, 0x1f, 0x2a, 0xd6, 0xad, 0x89, 0xcf, 0x80, 0x18, 0xad, 0xd7, 0x2e, 0xf2, 0x1e, 0xd4, 0xc1, 0x6d, 0x32, 0x94, 0x17, 0xc4, 0xcf, 0x62, 0xe7, 0x2c, 0x98, 0x1b, 0x8f, 0xdb, 0x6b, 0xf9, 0x85, 0x4f, 0x7d, 0xd3, 0x71, 0x22, 0x56, 0x4, 0x9c, 0x25, 0x42, 0x0, 0x67, 0x36, 0x85, 0x33, 0x4f, 0xd7, 0x21, 0x81, 0x22, 0x14, 0xe5, 0xbe, 0x47, 0xdb, 0x62, 0xd5, 0x2a, 0xeb, 0x4, 0xd9, 0x7a, 0x64, 0x29, 0x45, 0x95, 0x76, 0xce, 0x3c, 0xdc, 0x68, 0x66, 0x6b, 0x6c, 0xd3, 0xc5, 0x98, 0x94, 0x84, 0x96, 0x13, 0x65, 0x4b, 0x32, 0x7c, 0xd, 0x3b, 0xe9, 0xd0, 0x1, 0x38, 0x44, 0x9b, 0x0, 0x0, 0x0, 0x14, 0x9d, 0x90, 0x72, 0x4b, 0xc9, 0x5d, 0x39, 0x77, 0x2e, 0x2f, 0x2, 0xcc, 0xe8, 0xcf, 0xbc, 0xb, 0x4a, 0x46, 0x0, 0xbb, 0x0, 0x0, 0x0, 0x80, 0x5e, 0x4d, 0x3, 0xa2, 0x70, 0x90, 0x5f, 0x6a, 0x3d, 0x78, 0x4, 0x2e, 0x5d, 0x70, 0x48, 0x8d, 0x14, 0xbd, 0xee, 0x34, 0x2, 0xa2, 0xc2, 0x2b, 0xf3, 0x5e, 0x7b, 0xed, 0xbd, 0x8c, 0xd1, 0x42, 0x5, 0x8e, 0x4d, 0x49, 0x3a, 0xd9, 0xf3, 0xff, 0x9b, 0x69, 0x97, 0x5d, 0xf7, 0x5e, 0x76, 0xc, 0xc7, 0xbb, 0xad, 0xe5, 0xf0, 0xd6, 0xc3, 0xe7, 0xf8, 0x9d, 0x3f, 0x5a, 0x44, 0x5d, 0xf2, 0x9a, 0x52, 0x1d, 0x56, 0x54, 0x7e, 0x18, 0xc6, 0xfc, 0x7f, 0xee, 0x2d, 0x89, 0xde, 0x37, 0x36, 0x31, 0xeb, 0xbd, 0xee, 0xa5, 0x80, 0x62, 0x41, 0x5d, 0xc7, 0xff, 0x6b, 0xab, 0x53, 0x6, 0x28, 0xe7, 0xb5, 0x67, 0xac, 0xeb, 0x19, 0x24, 0x55, 0xd7, 0x25, 0x8f, 0x31, 0x4b, 0x5d, 0xf0, 0x51, 0xee, 0xf3, 0xe8, 0xbf, 0xd1, 0x94, 0xf5, 0x76, 0xa3, 0x93, 0x22, 0x87, 0x2b, 0x37, 0x24, 0xdc, 0xdf, 0x0, 0x0, 0x0, 0x80, 0x2f, 0xf3, 0xf1, 0x1c, 0x57, 0x68, 0x1b, 0xf9, 0x7, 0x50, 0xb6, 0xd7, 0x67, 0x9d, 0x68, 0xd, 0xd9, 0x6c, 0x6c, 0x25, 0x73, 0xdd, 0x17, 0x66, 0x2a, 0x59, 0x7, 0x97, 0xea, 0xfd, 0xb9, 0xfa, 0x56, 0xb3, 0xb1, 0x5f, 0x3a, 0x5e, 0x56, 0x5c, 0x2d, 0x52, 0xc4, 0xe0, 0x7e, 0xb8, 0xc9, 0xe8, 0xb5, 0x2f, 0x98, 0xa7, 0xe9, 0x80, 0x2a, 0x6a, 0x49, 0x3b, 0x8a, 0x3, 0xfd, 0xb0, 0x4e, 0x43, 0x95, 0x10, 0xdd, 0xbf, 0x54, 0x88, 0x94, 0xe7, 0xea, 0xc5, 0xb2, 0x55, 0x61, 0x9a, 0x65, 0x9a, 0xeb, 0x5, 0x8c, 0xdd, 0xb0, 0x95, 0xe4, 0x1a, 0xe4, 0xd3, 0x98, 0x7a, 0x3b, 0xd5, 0xb0, 0xb5, 0x30, 0x99, 0x61, 0xd2, 0xa5, 0xf5, 0x32, 0xc6, 0x43, 0xd7, 0xc1, 0x76, 0xbd, 0x86, 0x3e, 0x91, 0xc5, 0x56, 0xa3, 0x62, 0xac, 0x2f, 0x61, 0x2f, 0xae, 0x20, 0xc7, 0xe6, 0xa6, 0xf4, 0x84, 0xbc, 0x0, 0x0, 0x0, 0x14, 0x2c, 0xbc, 0x34, 0xc3, 0xc6, 0x15, 0x6f, 0xa4, 0x32, 0x15, 0xb2, 0xba, 0xb9, 0x19, 0x76, 0xe4, 0x52, 0xd, 0x6f, 0xed}},
		LegacyKnownFingerprints: nil,
		Peers: []*config.Peer{
			&config.Peer{
				UserID: "abcde@thoughtworks.com",
				Fingerprints: []*config.Fingerprint{
					&config.Fingerprint{
						Fingerprint: decode("57d8ea36c76d5d800fe790c56dc33feb254e899b"),
						Trusted:     false,
					},
				},
			},
			&config.Peer{
				UserID: "not@coyim.com",
				Fingerprints: []*config.Fingerprint{
					&config.Fingerprint{
						Fingerprint: decode("4157eea3bb3cf86cc0379e4c270e89b976bc34da"),
						Trusted:     false,
					},
				},
			},
			&config.Peer{
				UserID: "someone@where.com",
				Fingerprints: []*config.Fingerprint{
					&config.Fingerprint{
						Fingerprint: decode("a334e9d582da18f15028f7f7412bc8d15d0a1558"),
						Trusted:     true,
					},
				},
			},
		},
		HideStatusUpdates:   false,
		OTRAutoTearDown:     false,
		OTRAutoAppendTag:    false,
		OTRAutoStartSession: true,
		AlwaysEncrypt:       true,
		AlwaysEncryptWith:   []string(nil),
		DontEncryptWith:     []string(nil),
		InstanceTag:         0x0})
}

func (s *PidginSuite) Test_parseIntOr(c *C) {
	c.Assert(parseIntOr("123", 111), Equals, 123)
	c.Assert(parseIntOr("bla", 222), Equals, 222)
}

func (s *PidginSuite) Test_pidginPrefXML_failsOnMissingReceiver(c *C) {
	var p *pidginPrefXML
	res, ok := p.lookup()
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginSuite) Test_pidginPrefXML_failsOnMissingPrefs(c *C) {
	p := &pidginPrefXML{}
	res, ok := p.lookup("hello")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginSuite) Test_pidginImporter_importAllFrom_failsOnBadAccountsFile(c *C) {
	p := &pidginImporter{}
	res, ok := p.importAllFrom("file-that-hopefully-doesnt-exist", "", "", "", "")
	c.Assert(res, IsNil)
	c.Assert(ok, Equals, false)
}

func (s *PidginSuite) Test_pidginImporter_importAllFrom_usesAutomaticGlobalPrefs(c *C) {
	p := &pidginImporter{}
	res, ok := p.importAllFrom(testResourceFilename("pidgin_test_data/accounts.xml"), testResourceFilename("pidgin_test_data/prefs2.xml"), "", "", "")
	c.Assert(res, Not(IsNil))
	c.Assert(res.Accounts, HasLen, 2)
	c.Assert(res.Accounts[0].OTRAutoStartSession, Equals, true)
	c.Assert(res.Accounts[0].OTRAutoAppendTag, Equals, true)
	c.Assert(ok, Equals, true)
}

func (s *PidginSuite) Test_pidginImporter_importAllFrom_setsAlwaysEncryptToFalseIfNoGlobalPrefs(c *C) {
	p := &pidginImporter{}
	res, ok := p.importAllFrom(testResourceFilename("pidgin_test_data/accounts.xml"), testResourceFilename("pidgin_test_data/prefs3.xml"), "", "", "")
	c.Assert(res, Not(IsNil))
	c.Assert(res.Accounts, HasLen, 2)
	c.Assert(res.Accounts[0].AlwaysEncrypt, Equals, false)
	c.Assert(ok, Equals, true)
}

func (s *PidginSuite) Test_pidginImporter_importAllFrom_setsAlwaysEncryptWithAndDontEncryptWithForPeers(c *C) {
	p := &pidginImporter{}
	res, ok := p.importAllFrom(testResourceFilename("pidgin_test_data/accounts.xml"), "", testResourceFilename("pidgin_test_data/blist2.xml"), "", "")
	c.Assert(res, Not(IsNil))
	c.Assert(res.Accounts, HasLen, 2)
	c.Assert(res.Accounts[0].AlwaysEncryptWith, DeepEquals, []string{"not@coyim.com"})
	c.Assert(res.Accounts[0].DontEncryptWith, DeepEquals, []string{"foo@coyim.com"})
	c.Assert(ok, Equals, true)
}

func (s *PidginSuite) Test_pidginImporter_TryImport_works(c *C) {
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)

	origHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("HOME", origHome)
	}()
	os.Setenv("HOME", dir)

	os.Mkdir(filepath.Join(dir, pidginConfigDir), 0755)
	os.Create(filepath.Join(dir, pidginConfigDir, pidginAccountsFile))

	input, _ := ioutil.ReadFile(testResourceFilename("pidgin_test_data/accounts.xml"))
	_ = ioutil.WriteFile(filepath.Join(dir, pidginConfigDir, pidginAccountsFile), input, 0644)

	res := (&pidginImporter{}).TryImport()
	c.Assert(res, HasLen, 1)
}
