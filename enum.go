package main

import "github.com/orsinium-labs/enum"

type junosPlatform enum.Member[string]

var (
	platformVMX  = junosPlatform{Value: "vptx"}
	platformVPTX = junosPlatform{Value: "vptx"}
	platformVQFX = junosPlatform{Value: "vqfx"}
	platformVSRX = junosPlatform{Value: "vsrx"}
	platforms    = enum.New(platformVMX, platformVPTX, platformVQFX, platformVSRX)
)

var platformToOsFamily = map[junosPlatform]osFamily{
	platformVMX:  osFamilyJunos,
	platformVPTX: osFamilyJunos,
	platformVQFX: osFamilyJunosQFX,
	platformVSRX: osFamilyJunosES,
}

type osFamily enum.Member[string]

var (
	osFamilyJunos    = osFamily{Value: "junos"}
	osFamilyJunosES  = osFamily{Value: "junos-es"}
	osFamilyJunosQFX = osFamily{Value: "junos-qfx"}
	osFamilies       = enum.New(osFamilyJunos, osFamilyJunosES, osFamilyJunosQFX)
)
