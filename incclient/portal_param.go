package incclient

import "github.com/btcsuite/btcd/chaincfg"

// BTCPortalV4Params is a simplified version of the corresponding v4 portal param used in the Incognito network for the BTC token.
type BTCPortalV4Params struct {
	MasterPubKeys     [][]byte
	NumRequiredSigs   uint
	MinUnshieldAmount uint64
	ChainParams       *chaincfg.Params
	TokenID           string
}

var mainNetBTCPortalV4Params = BTCPortalV4Params{
	MasterPubKeys: [][]byte{
		{0x2, 0x39, 0x42, 0x3d, 0xad, 0x93, 0x8f, 0xcb, 0xe5, 0xb5, 0xef, 0x7b, 0x7b, 0x9a, 0xf, 0x28,
			0x4, 0x19, 0x53, 0x66, 0x7f, 0xee, 0x72, 0xe4, 0x81, 0xf9, 0xe6, 0xb, 0x81, 0x41, 0xd7, 0x3a, 0x36},
		{0x2, 0x8d, 0xc, 0xd7, 0x83, 0x9d, 0x5e, 0xc5, 0x7b, 0x77, 0x1a, 0xf1, 0x2, 0xb8, 0x72, 0xd0,
			0x4f, 0x34, 0xb4, 0xeb, 0x17, 0xac, 0xa1, 0x9f, 0xdf, 0xa, 0x64, 0xbf, 0xd, 0x36, 0x76, 0x66, 0x87},
		{0x3, 0x78, 0x52, 0x33, 0xe3, 0x8, 0x3a, 0xd8, 0x58, 0x77, 0x76, 0x29, 0xa0, 0x17, 0xb6, 0xdd,
			0x16, 0x43, 0x18, 0x8b, 0xb4, 0xa3, 0xaf, 0x45, 0xf0, 0xb5, 0x91, 0x8c, 0x84, 0xf2, 0x73, 0x56, 0x44},
		{0x3, 0x61, 0x9d, 0xc9, 0xfb, 0x6d, 0x8, 0x2a, 0x5c, 0x98, 0x45, 0xbc, 0xbf, 0x86, 0xfb, 0x47,
			0x4, 0xbe, 0x67, 0x46, 0xa, 0x59, 0xc4, 0xbc, 0x1d, 0xec, 0xc0, 0xe8, 0xe4, 0x3e, 0x1d, 0x6d, 0x0},
		{0x2, 0xe4, 0x1d, 0x40, 0xe6, 0xf3, 0x80, 0xad, 0x51, 0xca, 0x17, 0x87, 0xfe, 0xc8, 0x23, 0x8d,
			0xa4, 0xc2, 0x88, 0xfc, 0xfb, 0x6f, 0x2b, 0xcc, 0xd9, 0xa6, 0x1c, 0x2, 0xe5, 0x4a, 0x31, 0x34, 0x39},
		{0x2, 0xf0, 0xc, 0xe3, 0xec, 0x4, 0xdb, 0x75, 0x59, 0x99, 0x70, 0xc6, 0xfd, 0xc5, 0x2, 0x2f,
			0xad, 0x6b, 0x8d, 0x18, 0x86, 0x71, 0x44, 0xcf, 0xe6, 0x93, 0x92, 0xbb, 0xd1, 0x60, 0xc1, 0x1b, 0x5c},
		{0x2, 0x65, 0x96, 0x49, 0xab, 0xd4, 0xe5, 0x97, 0x7d, 0x5b, 0x67, 0x4c, 0x6d, 0xa1, 0xf, 0x9,
			0x28, 0xa0, 0x8c, 0x67, 0x8d, 0x7f, 0x50, 0xcc, 0x10, 0xf0, 0xfe, 0xe5, 0x68, 0xa8, 0x57, 0x63, 0xd8},
	},
	NumRequiredSigs:   5,
	MinUnshieldAmount: 100000,
	ChainParams:       &chaincfg.MainNetParams,
	TokenID:           "b832e5d3b1f01a4f0623f7fe91d6673461e1f5d37d91fe78c5c2e6183ff39696",
}

var testNet1BTCPortalV4Params = BTCPortalV4Params{
	MasterPubKeys: [][]byte{
		{0x2, 0x30, 0x34, 0xcb, 0x1a, 0x50, 0xf6, 0x7f, 0x5e, 0xb2, 0x53, 0x9e, 0x68, 0x3b, 0xd4,
			0x80, 0x73, 0x71, 0x2a, 0xdf, 0xf3, 0x25, 0x94, 0x34, 0x72, 0x6d, 0x62, 0x80, 0x83, 0xd2, 0x6f, 0x4c, 0xdd},
		{0x2, 0x74, 0x61, 0x32, 0x93, 0xe7, 0x93, 0x85, 0x94, 0xd2, 0x58, 0xfb, 0xcf, 0xc5, 0x33,
			0x78, 0xdc, 0x82, 0xcd, 0x64, 0xd1, 0xc0, 0x33, 0x1, 0x71, 0x2f, 0x90, 0x85, 0x72, 0xb9, 0x17, 0xab, 0xc7},
		{0x3, 0x67, 0x7a, 0x81, 0xfc, 0x9c, 0x4c, 0x9c, 0x6, 0x28, 0xd2, 0xf6, 0xd0, 0x1e, 0x27,
			0x15, 0xbb, 0x54, 0x11, 0x75, 0xe9, 0x62, 0xae, 0x78, 0x8f, 0xff, 0x26, 0x75, 0x1e, 0xb5, 0x24, 0xe0, 0xeb},
		{0x3, 0x2, 0xdb, 0xd4, 0xd4, 0x6b, 0x4e, 0xef, 0xe9, 0xa6, 0xe8, 0x64, 0xce, 0xeb, 0xb5,
			0x11, 0x25, 0x71, 0x28, 0x8a, 0xc4, 0xce, 0xca, 0xf4, 0x10, 0xd4, 0x16, 0x5f, 0x4c, 0x4c, 0xeb, 0x27, 0xe3},
	},
	NumRequiredSigs:   3,
	MinUnshieldAmount: 100000,
	ChainParams:       &chaincfg.TestNet3Params,
	TokenID:           "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82",
}

var testNetBTCPortalV4Params = BTCPortalV4Params{
	MasterPubKeys: [][]byte{
		{0x2, 0x30, 0x34, 0xcb, 0x1a, 0x50, 0xf6, 0x7f, 0x5e, 0xb2, 0x53, 0x9e, 0x68, 0x3b, 0xd4,
			0x80, 0x73, 0x71, 0x2a, 0xdf, 0xf3, 0x25, 0x94, 0x34, 0x72, 0x6d, 0x62, 0x80, 0x83, 0xd2, 0x6f, 0x4c, 0xdd},
		{0x2, 0x74, 0x61, 0x32, 0x93, 0xe7, 0x93, 0x85, 0x94, 0xd2, 0x58, 0xfb, 0xcf, 0xc5, 0x33,
			0x78, 0xdc, 0x82, 0xcd, 0x64, 0xd1, 0xc0, 0x33, 0x1, 0x71, 0x2f, 0x90, 0x85, 0x72, 0xb9, 0x17, 0xab, 0xc7},
		{0x3, 0x67, 0x7a, 0x81, 0xfc, 0x9c, 0x4c, 0x9c, 0x6, 0x28, 0xd2, 0xf6, 0xd0, 0x1e, 0x27,
			0x15, 0xbb, 0x54, 0x11, 0x75, 0xe9, 0x62, 0xae, 0x78, 0x8f, 0xff, 0x26, 0x75, 0x1e, 0xb5, 0x24, 0xe0, 0xeb},
		{0x3, 0x2, 0xdb, 0xd4, 0xd4, 0x6b, 0x4e, 0xef, 0xe9, 0xa6, 0xe8, 0x64, 0xce, 0xeb, 0xb5,
			0x11, 0x25, 0x71, 0x28, 0x8a, 0xc4, 0xce, 0xca, 0xf4, 0x10, 0xd4, 0x16, 0x5f, 0x4c, 0x4c, 0xeb, 0x27, 0xe3},
	},
	NumRequiredSigs:   3,
	MinUnshieldAmount: 100000,
	ChainParams:       &chaincfg.TestNet3Params,
	TokenID:           "4584d5e9b2fc0337dfb17f4b5bb025e5b82c38cfa4f54e8a3d4fcdd03954ff82",
}

var localBTCPortalV4Params = BTCPortalV4Params{
	MasterPubKeys: [][]byte{
		{0x3, 0xb2, 0xd3, 0x16, 0x7d, 0x94, 0x9c, 0x25, 0x3, 0xe6, 0x9c, 0x9f, 0x29, 0x78, 0x7d, 0x9c, 0x8, 0x8d, 0x39, 0x17, 0x8d, 0xb4, 0x75, 0x40, 0x35, 0xf5, 0xae, 0x6a, 0xf0, 0x17, 0x12, 0x11, 0x0},
		{0x3, 0x98, 0x7a, 0x87, 0xd1, 0x99, 0x13, 0xbd, 0xe3, 0xef, 0xf0, 0x55, 0x79, 0x2, 0xb4, 0x90, 0x57, 0xed, 0x1c, 0x9c, 0x8b, 0x32, 0xf9, 0x2, 0xbb, 0xbb, 0x85, 0x71, 0x3a, 0x99, 0x1f, 0xdc, 0x41},
		{0x3, 0x73, 0x23, 0x5e, 0xb1, 0xc8, 0xf1, 0x84, 0xe7, 0x59, 0x17, 0x6c, 0xe3, 0x87, 0x37, 0xb7, 0x91, 0x19, 0x47, 0x1b, 0xba, 0x63, 0x56, 0xbc, 0xab, 0x8d, 0xcc, 0x14, 0x4b, 0x42, 0x99, 0x86, 0x1},
		{0x3, 0x29, 0xe7, 0x59, 0x31, 0x89, 0xca, 0x7a, 0xf6, 0x1, 0xb6, 0x35, 0x67, 0x3d, 0xb1, 0x53, 0xd4, 0x19, 0xd7, 0x6, 0x19, 0x3, 0x2a, 0x32, 0x94, 0x57, 0x76, 0xb2, 0xb3, 0x80, 0x65, 0xe1, 0x5d},
	},
	NumRequiredSigs:   3,
	MinUnshieldAmount: 100000,
	ChainParams:       &chaincfg.TestNet3Params,
	TokenID:           "4483b5efffa71eb8f302a5417c40728cf3acccab8ec278bbe094186b5bd22d3f",
}
