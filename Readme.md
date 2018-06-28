# IAB User Consent String Parser

This library provides a golang parser for IAB consent string. 

See specification here: [GDPR Transparency & Consent Framework](https://github.com/InteractiveAdvertisingBureau/GDPR-Transparency-and-Consent-Framework/blob/master/Consent%20string%20and%20vendor%20list%20formats%20v1.1%20Final.md)

## Usage

    package main
    
    import (
    	"fmt"
    
    	"github.com/travelaudience/go-iab-consent-string-parse"
    )
    
    func main() {
    	parser, err := consent.NewUserConsent("BN5lERiOMYEdiAOAWeFRAAYAAaAAptQ")
    	if err != nil {
    		fmt.Print(err)
    		return
    	}
    
    	fmt.Println(parser.IsVendorAllowed(99))
    	fmt.Println(parser.IsPurposeAllowed(1))
    }
