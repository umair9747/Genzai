package main

var args []string
var targets []string

type Match struct {
	Headers      map[string]interface{} `json:"headers"`
	Strings      []string               `json:"strings"`
	ResponseCode int                    `json:"response_code"`
	Condition    string                 `json:"condition"`
}

type Entry struct {
	Matchers Match  `json:"matchers"`
	Category string `json:"category"`
	Tag      string `json:"tag"`
}

type DynamicEntries map[string]Entry

//STRUCTS FOR VENDOR LOGINS

type CustomHeaders map[string]string

type CustomPayload struct {
	Paths   []string      `json:"paths"`
	Headers CustomHeaders `json:"headers"`
	Body    string        `json:"body"`
	Method  string        `json:"method"`
}

type CustomMatchers struct {
	ResponseCode int           `json:"response_code"`
	Responsepath string        `json:"responsePath"`
	Strings      []string      `json:"strings"`
	Headers      CustomHeaders `json:"headers"`
}

type CustomEntry struct {
	Tag      string         `json:"tag"`
	Payload  CustomPayload  `json:"payload"`
	Matchers CustomMatchers `json:"matchers"`
	Issue    string         `json:"issue"`
}

type CustomEntries map[string]CustomEntry

type MyVendorLogins struct {
	Entries CustomEntries `json:"entries"`
}

var vendorDB MyVendorLogins
var vendorVulnsDB MyVendorLogins
var genzaiDB DynamicEntries

var saveOutput string

var banner = `
::::::::   :::::::::: ::::    ::: :::::::::     :::     ::::::::::: 
:+:    :+: :+:        :+:+:   :+:      :+:    :+: :+:       :+:     
+:+        +:+        :+:+:+  +:+     +:+    +:+   +:+      +:+     
:#:        +#++:++#   +#+ +:+ +#+    +#+    +#++:++#++:     +#+     
+#+   +#+# +#+        +#+  +#+#+#   +#+     +#+     +#+     +#+     
#+#    #+# #+#        #+#   #+#+#  #+#      #+#     #+#     #+#     
 ########  ########## ###    #### ######### ###     ### ########### 

	The IoT Security Toolkit by Umair Nehri (0x9747)
`
